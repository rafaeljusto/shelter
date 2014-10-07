// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package dao manage the objects persistence layer
package dao

import (
	"fmt"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/errors"
	"github.com/rafaeljusto/shelter/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// List of possible errors that can occur in this DAO. There can be also other errors from
// low level drivers.
var (
	// Programmer must set the Database attribute from DomainDAO with a valid connection
	// before using this object
	ErrDomainDAOUndefinedDatabase = fmt.Errorf("No database defined for DomainDAO")

	// Pagination attribute is mandatory, and it's a pointer only to fill some query
	// informations in it. For the user that wants all records without pagination for a B2B
	// integration need to pass zero in the page size
	ErrDomainDAOPaginationUndefined = fmt.Errorf("Pagination was not defined")
)

const (
	domainDAOCollection  = "domain" // Collection used to store all domain objects in the MongoDB database
	concurrentOperations = 1000     // Number of Go routines that will be used to execute many operations at once
)

func init() {
	// Add index on FQDN to speed up searchs. FQDN will be a unique field in database
	mongodb.RegisterIndexFunction(func(database *mgo.Database) error {
		index := mgo.Index{
			Name:     "fqdn",
			Key:      []string{"fqdn"},
			Unique:   true,
			DropDups: true,
		}

		if err := database.C(domainDAOCollection).EnsureIndex(index); err != nil {
			return errors.NewSystemError(err)
		}

		return nil
	})

	// Add index on nameserver.lastokat to speed up the query that check the domains that
	// need to be notified. We don't use laststatus because the selectivity is low,
	// according to http://docs.mongodb.org/manual/tutorial/create-queries-that-ensure-
	// selectivity/
	mongodb.RegisterIndexFunction(func(database *mgo.Database) error {
		index := mgo.Index{
			Name: "nameservers",
			Key:  []string{"nameservers.lastokat"},
		}

		if err := database.C(domainDAOCollection).EnsureIndex(index); err != nil {
			return errors.NewSystemError(err)
		}

		return nil
	})

	// Add index on dsset.lastokat to speed up the query that check the domains that need to
	// be notified. We don't use laststatus because the selectivity is low, according to
	// http://docs.mongodb.org/manual/tutorial/create-queries-that-ensure- selectivity/
	mongodb.RegisterIndexFunction(func(database *mgo.Database) error {
		index := mgo.Index{
			Name: "dsset",
			Key:  []string{"dsset.lastokat"},
		}

		if err := database.C(domainDAOCollection).EnsureIndex(index); err != nil {
			return errors.NewSystemError(err)
		}

		return nil
	})
}

// DomainDAO is the structure responsable for keeping the database connection to save the
// domain anytime during the their existence
type DomainDAO struct {
	Database *mgo.Database // MongoDB Database
}

// Save the domain object in the database and by consequence will also save the
// nameservers and ds set. On creation the domain object is going to receive the id that
// refers to the entry in the database
func (dao DomainDAO) Save(domain *model.Domain) error {
	// Check if the programmer forgot to set the database in DomainDAO object
	if dao.Database == nil {
		return errors.NewSystemError(ErrDomainDAOUndefinedDatabase)
	}

	// When creating a new domain object, the id will be probably nil (or kind of new
	// according to bson.ObjectId), so we must initialize it
	if len(domain.Id.Hex()) == 0 {
		domain.Id = bson.NewObjectId()
	}

	// Every time we modified a domain object we increase the revision counter to identify
	// changes in high level structures. Maybe a better approach would be doing this on the
	// MongoDB server side, check out the link http://docs.mongodb.org/manual/tutorial
	// /optimize-query-performance-with-indexes-and-projections/ - Use the Increment
	// Operator to Perform Operations Server-Side
	domain.Revision += 1

	// Store the last time that the object was modified
	domain.LastModifiedAt = time.Now().UTC()

	// Upsert try to update the collection entry if exists, if not, it creates a new entry.
	// For all the domain objects we are going to use the collection "domain". We also avoid
	// concurency adding the revision as a paremeter for updating the entry
	_, err := dao.Database.C(domainDAOCollection).Upsert(bson.M{
		"_id":      domain.Id,
		"revision": domain.Revision - 1,
	}, domain)

	if err != nil {
		return errors.NewSystemError(err)
	}

	return nil
}

// Save many domains at once, creating go routines to execute each domain in the classic
// Save method. This happens because there's no method to execute Upsert to many documents
// at once in the mgo API. http://stackoverflow.com/questions/19810176/mongodb-mgo-
// library-golang-multiple-insert-updates
func (dao DomainDAO) SaveMany(domains []*model.Domain) []DomainResult {
	return dao.executeMany(domains, dao.Save)
}

// Retrieve all domains using pagination control. This method is used by an end user to
// see all domains that are already registered in the system. The user will probably wants
// pagination to analyze the data in amounts. When pagination values are not informed,
// default values are adopted. There's also an expand flag that can control if each domain
// object from the list will have only the FQDN, last modification, nameserver and DS
// status or the full information
func (dao DomainDAO) FindAll(pagination *model.DomainPagination, expand bool, filter string) ([]model.Domain, error) {
	// Check if the programmer forgot to set the database in DomainDAO object
	if dao.Database == nil {
		return nil, errors.NewSystemError(ErrDomainDAOUndefinedDatabase)
	}

	// Programmer must always give a pagination, with default values if necessary
	if pagination == nil {
		return nil, errors.NewSystemError(ErrDomainDAOPaginationUndefined)
	}

	var query *mgo.Query

	if len(pagination.OrderBy) == 0 {
		pagination.OrderBy = model.DomainDefaultPaginationOrderBy
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = model.DefaultPaginationPageSize
	}

	if pagination.Page == 0 {
		pagination.Page = model.DefaultPaginationPage
	}

	var sortList []string
	for _, sort := range pagination.OrderBy {
		var sortTmp string

		if sort.Direction == model.OrderByDirectionDescending {
			sortTmp = "-"
		}

		switch sort.Field {
		case model.DomainOrderByFieldFQDN:
			sortTmp += "fqdn"
		case model.DomainOrderByFieldLastModifiedAt:
			sortTmp += "lastModifiedAt"
		}

		sortList = append(sortList, sortTmp)
	}

	if len(filter) == 0 {
		query = dao.Database.C(domainDAOCollection).Find(bson.M{})

	} else {
		query = dao.Database.C(domainDAOCollection).Find(bson.M{
			"fqdn": bson.RegEx{Pattern: filter, Options: "i"},
		})
	}

	// We store the number of items before applying pagination, if we do this after we get only the
	// number of items of a page size
	var err error
	if pagination.NumberOfItems, err = query.Count(); err != nil {
		return nil, errors.NewSystemError(err)
	}

	// Safety check to don't allow to set a page higher than the number of pages
	maxNumberOfPages := pagination.NumberOfItems / pagination.PageSize
	if pagination.NumberOfItems%pagination.PageSize > 0 {
		maxNumberOfPages++
	}

	if maxNumberOfPages == 0 {
		// When there's no item, we should stay on the first page (don't skip)
		pagination.Page = 1

	} else if pagination.Page > maxNumberOfPages {
		pagination.Page = maxNumberOfPages
	}

	query.
		Sort(sortList...).
		Skip(pagination.PageSize * (pagination.Page - 1)).
		Limit(pagination.PageSize)

	var domains []model.Domain
	if err := query.All(&domains); err != nil {
		return nil, errors.NewSystemError(err)
	}

	if pagination.PageSize > 0 {
		pagination.NumberOfPages = pagination.NumberOfItems / pagination.PageSize
		if pagination.NumberOfItems%pagination.PageSize > 0 {
			pagination.NumberOfPages += 1
		}
	}

	// When the expand flag if not defined, we should compress the domain object so the
	// network data isn't too big. For now the compressed object will have the FQDN, last
	// modification and the status of the nameservers and DS set, this is useful to detect
	// quickly the domains that have some issue
	if !expand {
		for i := range domains {
			for j := range domains[i].Nameservers {
				domains[i].Nameservers[j] = model.Nameserver{
					LastStatus: domains[i].Nameservers[j].LastStatus,
				}
			}

			for j := range domains[i].DSSet {
				domains[i].DSSet[j] = model.DS{
					LastStatus: domains[i].DSSet[j].LastStatus,
				}
			}

			domains[i].Owners = []model.Owner{}
		}
	}

	return domains, nil
}

// Retrieve all domains for a scan. This method can take a long time to load all domains,
// so it will return a channel and will send a domain as soon as it is loaded from the
// database. The method ends when it returns a nil domain or an error in the channel
// result
func (dao DomainDAO) FindAllAsync() (chan DomainResult, error) {
	// Check if the programmer forgot to set the database in DomainDAO object
	if dao.Database == nil {
		return nil, errors.NewSystemError(ErrDomainDAOUndefinedDatabase)
	}

	// Channel to be used for returning each retrieved domain
	domainChannel := make(chan DomainResult)

	go func() {
		// Gets the database result iterator
		it := dao.Database.C(domainDAOCollection).Find(bson.M{}).Iter()

		var domainIt model.Domain
		for it.Next(&domainIt) {
			domain := domainIt // Copy the domainIt object to send it to the channel
			domainChannel <- DomainResult{
				Domain: &domain,
				Error:  nil,
			}
		}

		err := it.Close()
		if err != nil {
			err = errors.NewSystemError(err)
		}

		domainChannel <- DomainResult{
			Domain: nil,
			Error:  err,
		}
	}()

	return domainChannel, nil
}

// Return all domains that need to be notified due to the error tolerancy policy. The
// objective is to help the user to configure correctly the nameservers alerting about
// problems. We are going to have different notification tolerances for nameserver, ds and
// the type of errors (timeout and others). In the worst case this method can return all
// the domains from the system, so it will work asynchronously, returning the domain as
// soon as it is selected
func (dao DomainDAO) FindAllAsyncToBeNotified(
	nameserverErrorAlertDays,
	nameserverTimeoutAlertDays,
	dsErrorAlertDays,
	dsTimeoutAlertDays,
	maxExpirationAlertDays int,
) (chan DomainResult, error) {

	// Check if the programmer forgot to set the database in DomainDAO object
	if dao.Database == nil {
		return nil, errors.NewSystemError(ErrDomainDAOUndefinedDatabase)
	}

	// Channel to be used for returning each retrieved domain
	domainChannel := make(chan DomainResult)

	go func() {
		// When using indexes with $or queries, remember that each clause of an $or query will
		// execute in parallel. These clauses can each use their own index. We tried another
		// query with $or operators inside the main $or but if we do that the "explain" show
		// us that MongoDB don't use indexes for that sittuation (so avoid it!)

		it := dao.Database.C(domainDAOCollection).Find(bson.M{
			"$or": []bson.M{
				{
					"nameservers": bson.M{"$elemMatch": bson.M{
						"laststatus": bson.M{"$nin": []model.NameserverStatus{
							model.NameserverStatusNotChecked,
							model.NameserverStatusOK,
							model.NameserverStatusTimeout,
						},
						},
						"lastokat": bson.M{
							"$lte": time.Now().Add(time.Duration(-nameserverErrorAlertDays*24) * time.Hour),
						},
					},
					},
				},
				{
					"nameservers": bson.M{"$elemMatch": bson.M{
						"laststatus": model.NameserverStatusTimeout,
						"lastokat": bson.M{
							"$lte": time.Now().Add(time.Duration(-nameserverTimeoutAlertDays*24) * time.Hour),
						},
					},
					},
				},
				{
					"dsset": bson.M{"$elemMatch": bson.M{
						"laststatus": bson.M{"$nin": []model.DSStatus{
							model.DSStatusNotChecked,
							model.DSStatusOK,
							model.DSStatusTimeout,
						},
						},
						"lastokat": bson.M{
							"$lte": time.Now().Add(time.Duration(-dsErrorAlertDays*24) * time.Hour),
						},
					},
					},
				},
				{
					"dsset": bson.M{"$elemMatch": bson.M{"laststatus": model.DSStatusTimeout,
						"lastokat": bson.M{
							"$lte": time.Now().Add(time.Duration(-dsTimeoutAlertDays*24) * time.Hour),
						},
					},
					},
				},
				{
					"dsset": bson.M{"$elemMatch": bson.M{"expiresat": bson.M{
						"$lte": time.Now().Add(time.Duration(maxExpirationAlertDays*24) * time.Hour),
					},
					},
					},
				},
			},
		}).Iter()

		var domainIt model.Domain
		for it.Next(&domainIt) {
			domain := domainIt // Copy the domainIt object to send it to the channel
			domainChannel <- DomainResult{
				Domain: &domain,
				Error:  nil,
			}
		}

		err := it.Close()
		if err != nil {
			err = errors.NewSystemError(err)
		}

		domainChannel <- DomainResult{
			Domain: nil,
			Error:  err,
		}
	}()

	return domainChannel, nil
}

// Try to find the domain using the FQDN attribute. The system was designed to have an
// unique FQDN. The database should be prepared (with indexes) to search faster when using
// FQDN as condition
func (dao DomainDAO) FindByFQDN(fqdn string) (model.Domain, error) {
	var domain model.Domain

	// Check if the programmer forgot to set the database in DomainDAO object
	if dao.Database == nil {
		return domain, errors.NewSystemError(ErrDomainDAOUndefinedDatabase)
	}

	// We must create a BSON object to be compared with MongoDB database entries
	err := dao.Database.C(domainDAOCollection).Find(bson.M{
		"fqdn": fqdn,
	}).One(&domain)

	if err == mgo.ErrNotFound {
		return domain, errors.NotFound

	} else if err != nil {
		return domain, errors.NewSystemError(err)
	}

	return domain, nil
}

// Remove a database entry based on a given domain id. This method is useful as a
// RemoveMany auxiliar method, because it's faster that RemoveByFQDN
func (dao DomainDAO) Remove(domain *model.Domain) error {
	// Check if the programmer forgot to set the database in DomainDAO object
	if dao.Database == nil {
		return errors.NewSystemError(ErrDomainDAOUndefinedDatabase)
	}

	if err := dao.Database.C(domainDAOCollection).RemoveId(domain.Id); err != nil {
		return errors.NewSystemError(err)
	}

	return nil
}

// Remove a database entry that have a given FQDN. The system was designed to have an
// unique FQDN. The database should be prepared (with indexes) to search faster when using
// FQDN as condition
func (dao DomainDAO) RemoveByFQDN(fqdn string) error {
	// Check if the programmer forgot to set the database in DomainDAO object
	if dao.Database == nil {
		return errors.NewSystemError(ErrDomainDAOUndefinedDatabase)
	}

	// We must create a BSON object to be compared with MongoDB database entries to
	// determinate wich one is going to be removed
	err := dao.Database.C(domainDAOCollection).Remove(bson.M{
		"fqdn": fqdn,
	})

	if err != nil {
		return errors.NewSystemError(err)
	}

	return nil
}

// Remove many domain objects from database at once, is faster than removing each one
// because it use go routines to execute everything concurrently
func (dao DomainDAO) RemoveMany(domains []*model.Domain) []DomainResult {
	return dao.executeMany(domains, dao.Remove)
}

// Remove all domain entries from the database. This is a DANGEROUS method, use with
// caution. For now is used only by the integration test enviroments to clear the database
// before starting a new test. We don't drop the collection because we don't wanna lose
// the indexes. Dropping the collection is much faster, but this method is probably never
// going to be a part of a critical system (I don't known any system that wants to erase
// all your data)
func (dao DomainDAO) RemoveAll() error {
	_, err := dao.Database.C(domainDAOCollection).RemoveAll(bson.M{})

	if err != nil {
		return errors.NewSystemError(err)
	}

	return nil
}

// Method used to execute an operation over domains concurrently. It was created because
// the SaveMany and RemoveMany methods were exactly the same, except for the database
// operation
func (dao DomainDAO) executeMany(domains []*model.Domain,
	operation func(domain *model.Domain) error) []DomainResult {

	domainResultsChannel := make(chan []DomainResult, len(domains)/concurrentOperations)

	// Create a poll of goroutines to split the operations and execute them concurrently
	var domainsGroups [concurrentOperations][]*model.Domain

	// Dispatch the domain's operations to different go routines
	if len(domains) < concurrentOperations {
		// When there's less operations than go routines, give one operation to each go
		// routine
		for i := 0; i < len(domains); i++ {
			domainsGroups[i] = append(domainsGroups[i], domains[i])
		}

	} else {
		// Check how many operations we will have on each go routine
		groupSize := len(domains) / concurrentOperations
		lastIndex := 0

		for i := 0; i < concurrentOperations; i++ {
			if i == concurrentOperations {
				// The last go routine will get all the rest of operations from the domain
				// list because the groupSize*concurrentOperations is not exatctly the
				// size of the domain list (division with rest)
				domainsGroups[i] = domains[lastIndex:]
			} else {
				// Create an slice of the group of operations of the go routine
				domainsGroups[i] = domains[lastIndex:(groupSize * (i + 1))]
				lastIndex += groupSize
			}
		}
	}

	// Start the go routines
	for i := 0; i < concurrentOperations; i++ {
		go func(domains []*model.Domain) {
			var domainResults []DomainResult

			// Execute all operations before returning the value to avoid channel waits
			for _, domain := range domains {
				domainResults = append(domainResults, DomainResult{
					Domain: domain,
					Error:  operation(domain),
				})
			}

			// Send back the result to the main thread
			domainResultsChannel <- domainResults
		}(domainsGroups[i])
	}

	var domainResults []DomainResult
	for i := 0; i < concurrentOperations; i++ {
		domainResults = append(domainResults, <-domainResultsChannel...)
	}
	return domainResults
}

// DomainResult is used when calling methods that execute actions over many domain
// objects. On error the caller need to know witch domain got an error. We are using this
// solution because there's no pair structure in Go language
type DomainResult struct {
	Domain *model.Domain // Domain related to the result
	Error  error         // When different of nil, represents an error of the operation
}

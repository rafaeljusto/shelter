package dao

import (
	"errors"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

// List of possible errors that can occur in this DAO. There can be also other errors from
// low level drivers.
var (
	// Programmer must set the Database attribute from DomainDAO with a valid connection
	// before using this object
	ErrDomainDAOUndefinedDatabase = errors.New("No database defined for DomainDAO")

	// Pagination attribute is mandatory, and it's a pointer only to fill some query
	// informations in it. For the user that wants all records without pagination for a B2B
	// integration need to pass zero in the page size
	ErrDomainDAOPaginationUndefined = errors.New("Pagination was not defined")

	// An invalid order by field was given to be converted in one of the known order by
	// fields of the Domain DAO
	ErrDomainDAOOrderByFieldUnknown = errors.New("Unknown order by field")

	// An invalid order by direction was given to be converted in one of the known order by
	// fields of the Domain DAO
	ErrDomainDAOOrderByDirectionUnknown = errors.New("Unknown order by direction")
)

const (
	domainDAOCollection  = "domain" // Collection used to store all domain objects in the MongoDB database
	concurrentOperations = 1000     // Number of Go routines that will be used to execute many operations at once
)

// List of possible fields that can be used to order a result set
const (
	DomainDAOOrderByFieldFQDN           DomainDAOOrderByField = 0 // Order by domain's FQDN
	DomainDAOOrderByFieldLastModifiedAt DomainDAOOrderByField = 1 // Order by the last modification date of the domain object
)

// Enumerate definition for the OrderBy so that we can limit the fields that the user can
// use in a query
type DomainDAOOrderByField int

// Convert the DomainDAO order by field from string into enum. If the string is unknown an
// error will be returned. The string is case insensitive and spaces around it are ignored
func DomainDAOOrderByFieldFromString(value string) (DomainDAOOrderByField, error) {
	value = strings.ToLower(value)
	value = strings.TrimSpace(value)

	switch value {
	case "fqdn":
		return DomainDAOOrderByFieldFQDN, nil
	case "lastmodified":
		return DomainDAOOrderByFieldLastModifiedAt, nil
	}

	return DomainDAOOrderByFieldFQDN, ErrDomainDAOOrderByFieldUnknown
}

// Convert the DomainDAO order by field from enum into string. If the enum is unknown this
// method will return an empty string
func DomainDAOOrderByFieldToString(value DomainDAOOrderByField) string {
	switch value {
	case DomainDAOOrderByFieldFQDN:
		return "fqdn"

	case DomainDAOOrderByFieldLastModifiedAt:
		return "lastmodified"
	}

	return ""
}

// List of possible directions of each field in an order by query
const (
	DomainDAOOrderByDirectionAscending  DomainDAOOrderByDirection = 1  // From lower to higher
	DomainDAOOrderByDirectionDescending DomainDAOOrderByDirection = -1 // From Higher to lower
)

// Enumerate definition for the OrderBy so that we can make it easy to determinate the
// direction of an order by field. This is good to hide low level database interfaces.
// Maybe we can move this enum to a more generic DAO because it can be used by other DAOs
// too
type DomainDAOOrderByDirection int

// Convert the DomainDAO order by direction from string into enum. If the string is
// unknown an error will be returned. The string is case insensitive and spaces around it
// are ignored
func DomainDAOOrderByDirectionFromString(value string) (DomainDAOOrderByDirection, error) {
	value = strings.ToLower(value)
	value = strings.TrimSpace(value)

	switch value {
	case "asc":
		return DomainDAOOrderByDirectionAscending, nil
	case "desc":
		return DomainDAOOrderByDirectionDescending, nil
	}

	return DomainDAOOrderByDirectionAscending, ErrDomainDAOOrderByDirectionUnknown
}

// Convert the DomainDAO order by direction from enum into string. If the enum is unknown
// this method will return an empty string
func DomainDAOOrderByDirectionToString(value DomainDAOOrderByDirection) string {
	switch value {
	case DomainDAOOrderByDirectionAscending:
		return "asc"

	case DomainDAOOrderByDirectionDescending:
		return "desc"
	}

	return ""
}

func init() {
	// Add index on FQDN to speed up searchs. FQDN will be a unique field in database
	mongodb.RegisterIndexFunction(func(database *mgo.Database) error {
		index := mgo.Index{
			Name:     "fqdn",
			Key:      []string{"fqdn"},
			Unique:   true,
			DropDups: true,
		}

		return database.C(domainDAOCollection).EnsureIndex(index)
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
		return ErrDomainDAOUndefinedDatabase
	}

	// When creating a new domain object, the id will be probably nil (or kind of new
	// according to bson.ObjectId), so we must initialize it
	if len(domain.Id.Hex()) == 0 {
		domain.Id = bson.NewObjectId()
	}

	// Every time we modified a domain object we increase the revision counter to identify
	// changes in high level structures
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

	return err
}

// Save many domains at once, creating go routines to execute each domain in the classic
// Save method. This happens because there's no method to execute Upsert to many documents
// at once in the mgo API. http://stackoverflow.com/questions/19810176/mongodb-mgo-
// library-golang-multiple-insert-updates
func (dao DomainDAO) SaveMany(domains []*model.Domain) []DomainResult {
	return dao.executeMany(domains, dao.Save)
}

// Retrieve all domains using pagination control. This method is used by an end user to
// see all domains that are alredy registered in the system, except for a B2B integration
// systems the user will probably want pagination so that it can analyze the data in
// amounts. For the user that wants all records without pagination for a B2B
// integration need to pass zero in the page size
func (dao DomainDAO) FindAll(pagination *DomainDAOPagination) ([]model.Domain, error) {
	// Check if the programmer forgot to set the database in DomainDAO object
	if dao.Database == nil {
		return nil, ErrDomainDAOUndefinedDatabase
	}

	if pagination == nil {
		return nil, ErrDomainDAOPaginationUndefined
	}

	var query *mgo.Query

	if pagination.PageSize == 0 {
		query = dao.Database.C(domainDAOCollection).Find(bson.M{})

		var err error
		if pagination.NumberOfItems, err = query.Count(); err != nil {
			return nil, err
		}

	} else {
		var sortList []string
		for _, sort := range pagination.OrderBy {
			var sortTmp string

			if sort.Direction == DomainDAOOrderByDirectionDescending {
				sortTmp = "-"
			}

			switch sort.Field {
			case DomainDAOOrderByFieldFQDN:
				sortTmp += "fqdn"
			case DomainDAOOrderByFieldLastModifiedAt:
				sortTmp += "lastModifiedAt"
			}

			sortList = append(sortList, sortTmp)
		}

		query = dao.Database.C(domainDAOCollection).Find(bson.M{})

		// We store the number of items before applying pagination, if we do this after we get only the
		// number of items of a page size
		var err error
		if pagination.NumberOfItems, err = query.Count(); err != nil {
			return nil, err
		}

		query.
			Sort(sortList...).
			Skip(pagination.PageSize * (pagination.Page - 1)).
			Limit(pagination.PageSize)
	}

	var domains []model.Domain
	if err := query.All(&domains); err != nil {
		return nil, err
	}

	if pagination.PageSize > 0 {
		pagination.NumberOfPages = pagination.NumberOfItems / pagination.PageSize
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
		return nil, ErrDomainDAOUndefinedDatabase
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
		return domain, ErrDomainDAOUndefinedDatabase
	}

	// We must create a BSON object to be compared with MongoDB database entries
	err := dao.Database.C(domainDAOCollection).Find(bson.M{
		"fqdn": fqdn,
	}).One(&domain)

	return domain, err
}

// Remove a database entry based on a given domain id. This method is useful as a
// RemoveMany auxiliar method, because it's faster that RemoveByFQDN
func (dao DomainDAO) Remove(domain *model.Domain) error {
	// Check if the programmer forgot to set the database in DomainDAO object
	if dao.Database == nil {
		return ErrDomainDAOUndefinedDatabase
	}

	return dao.Database.C(domainDAOCollection).RemoveId(domain.Id)
}

// Remove a database entry that have a given FQDN. The system was designed to have an
// unique FQDN. The database should be prepared (with indexes) to search faster when using
// FQDN as condition
func (dao DomainDAO) RemoveByFQDN(fqdn string) error {
	// Check if the programmer forgot to set the database in DomainDAO object
	if dao.Database == nil {
		return ErrDomainDAOUndefinedDatabase
	}

	// We must create a BSON object to be compared with MongoDB database entries to
	// determinate wich one is going to be removed
	return dao.Database.C(domainDAOCollection).Remove(bson.M{
		"fqdn": fqdn,
	})
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
	return err
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

// DomainDAOPagination was created as a necessity for big result sets that needs to be
// sent for an end-user. With pagination we can control the size of the data and make it
// faster for the user to interact with it in a web interface as example
type DomainDAOPagination struct {
	OrderBy       []DomainDAOSort // Sort the list before the pagination
	PageSize      int             // Number of items that are going to be considered in one page
	Page          int             // Current page that will be returned
	NumberOfItems int             // Total number of items in the result set
	NumberOfPages int             // Total number of pages calculated for the current result set
}

// DomainDAOSort is an object responsable to relate the order by field and direction. Each
// field used for sort, can be sorted in both directions
type DomainDAOSort struct {
	Field     DomainDAOOrderByField     // Field to be sorted
	Direction DomainDAOOrderByDirection // Direction used in the sort
}

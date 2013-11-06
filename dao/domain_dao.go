package dao

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"shelter/database/mongodb"
	"shelter/model"
)

// List of possible errors that can occur in this DAO. There can be also other errors from
// low level drivers.
var (
	// Programmer must set the Database attribute from DomainDAO with a valid connection
	// before using this object
	ErrDomainDAOUndefinedDatabase = errors.New("No database defined for DomainDAO")
)

const (
	// Collection used to store all domain objects in the MongoDB database
	domainDAOCollection = "domain"
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
	// according to bson.ObjectId), so we musr initialize it
	if len(domain.Id.Hex()) == 0 {
		domain.Id = bson.NewObjectId()
	}

	// Upsert try to update the collection entry if exists, if not, it creates a new
	// entry. For all the domain objects we are going to use the collection "domain"
	_, err := dao.Database.C(domainDAOCollection).UpsertId(domain.Id, domain)

	return err
}

// Save many domains at once, creating go routines to execute each domain in the classic
// Save method. This happens because there's no method to execute Upsert to many documents
// at once in the mgo API. http://stackoverflow.com/questions/19810176/mongodb-mgo-
// library- golang-multiple-insert-updates
func (dao DomainDAO) SaveMany(domains []*model.Domain) []DomainResult {
	domainResultChannel := make(chan DomainResult, len(domains))

	for _, domain := range domains {
		go func(domain *model.Domain) {
			domainResultChannel <- DomainResult{
				Domain: domain,
				Error:  dao.Save(domain),
			}
		}(domain)
	}

	var domainResults []DomainResult
	for i := 0; i < len(domains); i++ {
		domainResults = append(domainResults, <-domainResultChannel)
	}

	return domainResults
}

// Retrieve all domains for a scan. This method can take a long time to load all domains,
// so it will return a channel and will send a domain as soon as it is loaded from the
// database
func (dao DomainDAO) FindAll() (chan model.Domain, error) {
	// Check if the programmer forgot to set the database in DomainDAO object
	if dao.Database == nil {
		return nil, ErrDomainDAOUndefinedDatabase
	}

	// Channel to be used for returning each retrieved domain
	domainChannel := make(chan model.Domain)

	go func() {
		// Gets the database result iterator
		it := dao.Database.C(domainDAOCollection).Find(bson.M{}).Iter()

		var domain model.Domain
		if it.Next(domain) {
			domainChannel <- domain
		} else {
			// TODO: How to alert about an error that occurred here?
			// err := it.Err()
			return
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
	domainResultChannel := make(chan DomainResult)

	for _, domain := range domains {
		go func(domain *model.Domain) {
			domainResultChannel <- DomainResult{
				Domain: domain,
				Error:  dao.Remove(domain),
			}
		}(domain)
	}

	var domainResults []DomainResult
	for i := 0; i < len(domains); i++ {
		domainResults = append(domainResults, <-domainResultChannel)
	}

	return domainResults
}

// Remove all domain entries from the database. This is a DANGEROUS method, use with
// caution. For now is used only by the integration test enviroments to clear the database
// before starting a new test
func (dao DomainDAO) RemoveAll() error {
	return dao.Database.C(domainDAOCollection).DropCollection()
}

// DomainResult is used when calling methods that execute actions over many domain
// objects. On error the caller need to know witch domain got an error. We are using this
// solution because there's no pair structure in Go language
type DomainResult struct {
	Domain *model.Domain // Domain related to the result
	Error  error         // When different of nil, represents an error of the operation
}

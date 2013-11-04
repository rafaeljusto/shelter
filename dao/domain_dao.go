package dao

import (
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
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
	DomainDAOCollection = "domain"
)

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
	_, err := dao.Database.C(DomainDAOCollection).UpsertId(domain.Id, domain)

	return err
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
	err := dao.Database.C(DomainDAOCollection).Find(bson.M{
		"fqdn": fqdn,
	}).One(&domain)

	return domain, err
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
	return dao.Database.C(DomainDAOCollection).Remove(bson.M{
		"fqdn": fqdn,
	})
}

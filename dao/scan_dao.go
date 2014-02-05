package dao

import (
	"errors"
	"github.com/rafaeljusto/shelter/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// List of possible errors that can occur in this DAO. There can be also other errors from
// low level drivers.
var (
	// Programmer must set the Database attribute from ScanDAO with a valid connection before using
	// this object
	ErrScanDAOUndefinedDatabase = errors.New("No database defined for ScanDAO")
)

const (
	scanDAOCollection = "scan" // Collection used to store all scan objects in the MongoDB database
)

// ScanDAO is the structure responsible for keeping the database connection to save a new scan
// object after every execution
type ScanDAO struct {
	Database *mgo.Database // MongoDB Database
}

// Save the scan object in the database. On creation the scan object is going to receive the id that
// refers to the entry in the database
func (dao ScanDAO) Save(scan *model.Scan) error {
	// Check if the programmer forgot to set the database in scanDAO object
	if dao.Database == nil {
		return ErrScanDAOUndefinedDatabase
	}

	// When creating a new scan object, the id will be probably nil (or kind of new
	// according to bson.ObjectId), so we must initialize it
	if len(scan.Id.Hex()) == 0 {
		scan.Id = bson.NewObjectId()
	}

	// Every time we modified a scan object we increase the revision counter to identify
	// changes in high level structures
	scan.Revision += 1

	// Store the last time that the object was modified
	scan.LastModifiedAt = time.Now().UTC()

	// Upsert try to update the collection entry if exists, if not, it creates a new entry. For all
	// the scan objects we are going to use the collection "scan". We also avoid concurency adding the
	// revision as a paremeter for updating the entry
	_, err := dao.Database.C(scanDAOCollection).Upsert(bson.M{
		"_id":      scan.Id,
		"revision": scan.Revision - 1,
	}, scan)

	return err
}

// Try to find the scan using the Id attribute
func (dao ScanDAO) FindById(id bson.ObjectId) (model.Scan, error) {
	scan := model.Scan{
		NameserverStatistics: make(map[string]uint64),
		DSStatistics:         make(map[string]uint64),
	}

	// Check if the programmer forgot to set the database in ScanDAO object
	if dao.Database == nil {
		return scan, ErrScanDAOUndefinedDatabase
	}

	err := dao.Database.C(scanDAOCollection).FindId(id).One(&scan)

	return scan, err
}

// Remove a database entry that have a given id
func (dao ScanDAO) RemoveById(id bson.ObjectId) error {
	// Check if the programmer forgot to set the database in ScanDAO object
	if dao.Database == nil {
		return ErrScanDAOUndefinedDatabase
	}

	// We must create a BSON object to be compared with MongoDB database entries to
	// determinate wich one is going to be removed
	return dao.Database.C(scanDAOCollection).RemoveId(id)
}

// Remove all scan entries from the database. This is a DANGEROUS method, use with
// caution. For now is used only by the integration test enviroments to clear the database
// before starting a new test. We don't drop the collection because we don't wanna lose
// the indexes. Dropping the collection is much faster, but this method is probably never
// going to be a part of a critical system (I don't known any system that wants to erase
// all your data)
func (dao ScanDAO) RemoveAll() error {
	_, err := dao.Database.C(scanDAOCollection).RemoveAll(bson.M{})
	return err
}

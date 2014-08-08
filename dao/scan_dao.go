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
	"strings"
	"time"
)

// List of possible errors that can occur in this DAO. There can be also other errors from
// low level drivers.
var (
	// Programmer must set the Database attribute from ScanDAO with a valid connection before using
	// this object
	ErrScanDAOUndefinedDatabase = fmt.Errorf("No database defined for ScanDAO")

	// Pagination attribute is mandatory, and it's a pointer only to fill some query
	// informations in it. For the user that wants all records without pagination for a B2B
	// integration need to pass zero in the page size
	ErrScanDAOPaginationUndefined = fmt.Errorf("Pagination was not defined")
)

const (
	scanDAOCollection = "scan" // Collection used to store all scan objects in the MongoDB database
)

// List of possible fields that can be used to order a result set
const (
	ScanDAOOrderByFieldStartedAt                ScanDAOOrderByField = 0 // Order by scan's begin time
	ScanDAOOrderByFieldDomainsScanned           ScanDAOOrderByField = 1 // Order by the number of domains scanned
	ScanDAOOrderByFieldDomainsWithDNSSECScanned ScanDAOOrderByField = 2 // Order by the number of domains with DNSSEC scanned
)

// Enumerate definition for the OrderBy so that we can limit the fields that the user can
// use in a query
type ScanDAOOrderByField int

// Convert the ScanDAO order by field from string into enum. If the string is unknown an
// error will be returned. The string is case insensitive and spaces around it are ignored
func ScanDAOOrderByFieldFromString(value string) (ScanDAOOrderByField, error) {
	value = strings.ToLower(value)
	value = strings.TrimSpace(value)

	switch value {
	case "startedat":
		return ScanDAOOrderByFieldStartedAt, nil
	case "domainsscanned":
		return ScanDAOOrderByFieldDomainsScanned, nil
	case "domainswithdnssecscanned":
		return ScanDAOOrderByFieldDomainsWithDNSSECScanned, nil
	}

	return ScanDAOOrderByFieldStartedAt,
		errors.NewInputError(errors.ErrorCodeInvalidQueryOrderBy, "orderby", value)
}

// Convert the ScanDAO order by field from enum into string. If the enum is unknown this
// method will return an empty string
func ScanDAOOrderByFieldToString(value ScanDAOOrderByField) string {
	switch value {
	case ScanDAOOrderByFieldStartedAt:
		return "startedat"

	case ScanDAOOrderByFieldDomainsScanned:
		return "domainsscanned"

	case ScanDAOOrderByFieldDomainsWithDNSSECScanned:
		return "domainswithdnssecscanned"
	}

	return ""
}

// Default values when the user don't define pagination. After watching a presentation
// from layer7 at http://www.layer7tech.com/tutorials/api-pagination-tutorial I agree that
// when the user don't define the pagination we shouldn't return all the result set,
// instead we assume default pagination values
var (
	scanDAODefaultPaginationOrderBy = []ScanDAOSort{
		{
			Field:     ScanDAOOrderByFieldStartedAt, // Default ordering is by begin time
			Direction: DAOOrderByDirectionAscending, // Default ordering is ascending
		},
	}
)

func init() {
	// Add index on StartedAt to speed up searchs. StartedAt will be a unique field in
	// database, so we cannot have two scans starting at the same time. One problem is that
	// the scan is only going to realize the problem after scanning all domains
	mongodb.RegisterIndexFunction(func(database *mgo.Database) error {
		index := mgo.Index{
			Name:     "startedat",
			Key:      []string{"startedat"},
			Unique:   true,
			DropDups: true,
		}

		if err := database.C(scanDAOCollection).EnsureIndex(index); err != nil {
			return errors.NewSystemError(err)
		}

		return nil
	})
}

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
		return errors.NewSystemError(ErrScanDAOUndefinedDatabase)
	}

	// When creating a new scan object, the id will be probably nil (or kind of new
	// according to bson.ObjectId), so we must initialize it
	if len(scan.Id.Hex()) == 0 {
		scan.Id = bson.NewObjectId()
	}

	// Every time we modified a scan object we increase the revision counter to identify
	// changes in high level structures. Maybe a better approach would be doing this on the
	// MongoDB server side, check out the link http://docs.mongodb.org/manual/tutorial
	// /optimize-query-performance-with-indexes-and-projections/ - Use the Increment
	// Operator to Perform Operations Server-Side
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

	if err != nil {
		return errors.NewSystemError(err)
	}

	return nil
}

// Try to find the scan using the startedAt time attribute
func (dao ScanDAO) FindByStartedAt(startedAt time.Time) (model.Scan, error) {
	scan := model.Scan{
		NameserverStatistics: make(map[string]uint64),
		DSStatistics:         make(map[string]uint64),
	}

	// Check if the programmer forgot to set the database in ScanDAO object
	if dao.Database == nil {
		return scan, errors.NewSystemError(ErrScanDAOUndefinedDatabase)
	}

	err := dao.Database.C(scanDAOCollection).Find(bson.M{
		"startedat": startedAt,
	}).One(&scan)

	if err == mgo.ErrNotFound {
		return scan, errors.NotFound

	} else if err != nil {
		return scan, errors.NewSystemError(err)
	}

	return scan, nil
}

// Retrieve all scans using pagination control. This method is used by an end user to see
// all scans that were executed in the system. The user will probably wants pagination to
// analyze the data in amounts. When pagination values are not informed, default values
// are adopted. There's also an expand flag that can control if each scan object from the
// list will have only the started date and the last modification date or the full
// information
func (dao ScanDAO) FindAll(pagination *ScanDAOPagination, expand bool) ([]model.Scan, error) {
	// Check if the programmer forgot to set the database in ScanDAO object
	if dao.Database == nil {
		return nil, errors.NewSystemError(ErrScanDAOUndefinedDatabase)
	}

	// Programmer must always give a pagination, with default values if necessary
	if pagination == nil {
		return nil, errors.NewSystemError(ErrScanDAOPaginationUndefined)
	}

	var query *mgo.Query

	if len(pagination.OrderBy) == 0 {
		pagination.OrderBy = scanDAODefaultPaginationOrderBy
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = defaultPaginationPageSize
	}

	if pagination.Page == 0 {
		pagination.Page = defaultPaginationPage
	}

	var sortList []string
	for _, sort := range pagination.OrderBy {
		var sortTmp string

		if sort.Direction == DAOOrderByDirectionDescending {
			sortTmp = "-"
		}

		switch sort.Field {
		case ScanDAOOrderByFieldStartedAt:
			sortTmp += "startedAt"
		case ScanDAOOrderByFieldDomainsScanned:
			sortTmp += "domainsScanned"
		case ScanDAOOrderByFieldDomainsWithDNSSECScanned:
			sortTmp += "domainsWithDNSSECScanned"
		}

		sortList = append(sortList, sortTmp)
	}

	query = dao.Database.C(scanDAOCollection).Find(bson.M{})

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

	var scans []model.Scan
	if err := query.All(&scans); err != nil {
		return nil, errors.NewSystemError(err)
	}

	if pagination.PageSize > 0 {
		pagination.NumberOfPages = pagination.NumberOfItems / pagination.PageSize
		if pagination.NumberOfItems%pagination.PageSize > 0 {
			pagination.NumberOfPages += 1
		}
	}

	// When the expand flag if not defined, we should compress the scan object so the
	// network data isn't too big. For now the compressed object will have the start date
	// and the last modification date
	if !expand {
		for i := range scans {
			scans[i] = model.Scan{
				StartedAt:      scans[i].StartedAt,
				LastModifiedAt: scans[i].LastModifiedAt,
			}
		}
	}

	return scans, nil
}

// Remove a database entry that have a given startedAt time
func (dao ScanDAO) RemoveByStartedAt(startedAt time.Time) error {
	// Check if the programmer forgot to set the database in ScanDAO object
	if dao.Database == nil {
		return errors.NewSystemError(ErrScanDAOUndefinedDatabase)
	}

	// We must create a BSON object to be compared with MongoDB database entries to
	// determinate wich one is going to be removed
	err := dao.Database.C(scanDAOCollection).Remove(bson.M{
		"startedat": startedAt,
	})

	if err != nil {
		return errors.NewSystemError(err)
	}

	return nil
}

// Remove all scan entries from the database. This is a DANGEROUS method, use with
// caution. For now is used only by the integration test enviroments to clear the database
// before starting a new test. We don't drop the collection because we don't wanna lose
// the indexes. Dropping the collection is much faster, but this method is probably never
// going to be a part of a critical system (I don't known any system that wants to erase
// all your data)
func (dao ScanDAO) RemoveAll() error {
	_, err := dao.Database.C(scanDAOCollection).RemoveAll(bson.M{})

	if err != nil {
		return errors.NewSystemError(err)
	}

	return nil
}

// ScanDAOPagination was created as a necessity for big result sets that needs to be
// sent for an end-user. With pagination we can control the size of the data and make it
// faster for the user to interact with it in a web interface as example
type ScanDAOPagination struct {
	OrderBy       []ScanDAOSort // Sort the list before the pagination
	PageSize      int           // Number of items that are going to be considered in one page
	Page          int           // Current page that will be returned
	NumberOfItems int           // Total number of items in the result set
	NumberOfPages int           // Total number of pages calculated for the current result set
}

// ScanDAOSort is an object responsable to relate the order by field and direction. Each
// field used for sort, can be sorted in both directions
type ScanDAOSort struct {
	Field     ScanDAOOrderByField // Field to be sorted
	Direction DAOOrderByDirection // Direction used in the sort
}

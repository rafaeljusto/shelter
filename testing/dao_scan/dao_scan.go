package main

import (
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/testing/utils"
	"time"
)

// This test objective is to verify the scan data persistence. The strategy is to insert
// and search for the information. Check for insert/update consistency (updates don't
// create a new element) and if the object id is set on creation

var (
	configFilePath string // Path for the configuration file with the database connection information
	report         bool   // Flag to generate the scan dao performance report file
)

// ScanDAOTestConfigFile is a structure to store the test configuration file data
type ScanDAOTestConfigFile struct {
	Database struct {
		URI  string
		Name string
	}
}

func init() {
	utils.TestName = "ScanDAO"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for ScanDAO test")
}

func main() {
	flag.Parse()

	var config ScanDAOTestConfigFile
	err := utils.ReadConfigFile(configFilePath, &config)

	if err == utils.ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		utils.Fatalln("Error reading configuration file", err)
	}

	database, databaseSession, err := mongodb.Open(config.Database.URI, config.Database.Name)
	if err != nil {
		utils.Fatalln("Error connecting the database", err)
	}
	defer databaseSession.Close()

	scanDAO := dao.ScanDAO{
		Database: database,
	}

	// If there was some problem in the last test, there could be some data in the
	// database, so let's clear it to don't affect this test. We avoid checking the error,
	// because if the collection does not exist yet, it will be created in the first
	// insert
	scanDAO.RemoveAll()

	scanLifeCycle(scanDAO)
	scanConcurrency(scanDAO)
	scansPagination(scanDAO)
	scanStatistics(scanDAO)

	utils.Println("SUCCESS!")
}

// Test all phases of the scan life cycle
func scanLifeCycle(scanDAO dao.ScanDAO) {
	scan := newScan()

	// Create scan
	if err := scanDAO.Save(&scan); err != nil {
		utils.Fatalln("Couldn't save scan in database", err)
	}

	// Search and compare created scan
	if scanRetrieved, err := scanDAO.FindByStartedAt(scan.StartedAt); err != nil {
		utils.Fatalln("Couldn't find created scan in database", err)

	} else if !compareScans(scan, scanRetrieved) {
		utils.Fatalln("Scan created in being persisted wrongly", nil)
	}

	// Update scan
	scan.DomainsScanned = 100
	if err := scanDAO.Save(&scan); err != nil {
		utils.Fatalln("Couldn't save scan in database", err)
	}

	// Search and compare updated scan
	if scanRetrieved, err := scanDAO.FindByStartedAt(scan.StartedAt); err != nil {
		utils.Fatalln("Couldn't find updated scan in database", err)

	} else if !compareScans(scan, scanRetrieved) {
		utils.Fatalln("Scan updated in being persisted wrongly", nil)
	}

	// Remove scan
	if err := scanDAO.RemoveByStartedAt(scan.StartedAt); err != nil {
		utils.Fatalln("Error while trying to remove a scan", err)
	}

	// Check removal
	if _, err := scanDAO.FindByStartedAt(scan.StartedAt); err == nil {
		utils.Fatalln("Scan was not removed from database", nil)
	}
}

// Check if the revision field avoid data concurrency. Is better to fail than to store the
// wrong state
func scanConcurrency(scanDAO dao.ScanDAO) {
	scan := newScan()

	// Create scan
	if err := scanDAO.Save(&scan); err != nil {
		utils.Fatalln("Couldn't save scan in database", err)
	}

	scan1, err := scanDAO.FindByStartedAt(scan.StartedAt)
	if err != nil {
		utils.Fatalln("Couldn't find created scan in database", err)
	}

	scan2, err := scanDAO.FindByStartedAt(scan.StartedAt)
	if err != nil {
		utils.Fatalln("Couldn't find created scan in database", err)
	}

	if err := scanDAO.Save(&scan1); err != nil {
		utils.Fatalln("Couldn't save scan in database", err)
	}

	if err := scanDAO.Save(&scan2); err == nil {
		utils.Fatalln("Not controlling scan concurrency", nil)
	}

	// Remove scan
	if err := scanDAO.RemoveByStartedAt(scan.StartedAt); err != nil {
		utils.Fatalln("Error while trying to remove a scan", err)
	}
}

func scansPagination(scanDAO dao.ScanDAO) {
	numberOfItems := 1000

	for i := 0; i < numberOfItems; i++ {
		scan := model.Scan{
			StartedAt: time.Now().Add(time.Duration(-i) * time.Minute),
		}

		if err := scanDAO.Save(&scan); err != nil {
			utils.Fatalln("Error saving scan in database", err)
		}
	}

	pagination := dao.ScanDAOPagination{
		PageSize: 10,
		Page:     5,
		OrderBy: []dao.ScanDAOSort{
			{
				Field:     dao.ScanDAOOrderByFieldStartedAt,
				Direction: dao.DAOOrderByDirectionAscending,
			},
		},
	}

	scans, err := scanDAO.FindAll(&pagination)
	if err != nil {
		utils.Fatalln("Error retrieving scans", err)
	}

	if pagination.NumberOfItems != numberOfItems {
		utils.Errorln("Number of items not calculated correctly", nil)
	}

	if pagination.NumberOfPages != numberOfItems/pagination.PageSize {
		utils.Errorln("Number of pages not calculated correctly", nil)
	}

	if len(scans) != pagination.PageSize {
		utils.Errorln("Number of scans not following page size", nil)
	}

	pagination = dao.ScanDAOPagination{
		PageSize: 10000,
		Page:     1,
		OrderBy: []dao.ScanDAOSort{
			{
				Field:     dao.ScanDAOOrderByFieldStartedAt,
				Direction: dao.DAOOrderByDirectionAscending,
			},
		},
	}

	scans, err = scanDAO.FindAll(&pagination)
	if err != nil {
		utils.Fatalln("Error retrieving scans", err)
	}

	if pagination.NumberOfPages != 1 {
		utils.Fatalln("Calculating wrong number of pages when there's only one page", nil)
	}

	if err := scanDAO.RemoveAll(); err != nil {
		utils.Fatalln("Error removing scans from database", err)
	}
}

func scanStatistics(scanDAO dao.ScanDAO) {
	scan := newScan()
	scan.NameserverStatistics = map[string]uint64{
		"OK":      40,
		"TIMEOUT": 10,
	}
	scan.DSStatistics = map[string]uint64{
		"OK":     8,
		"DNSERR": 2,
	}

	// Create scan
	if err := scanDAO.Save(&scan); err != nil {
		utils.Fatalln("Couldn't save scan in database", err)
	}

	var err error
	scan, err = scanDAO.FindByStartedAt(scan.StartedAt)
	if err != nil {
		utils.Fatalln("Couldn't find created scan in database", err)
	}

	if len(scan.NameserverStatistics) != 2 ||
		scan.NameserverStatistics["OK"] != 40 ||
		scan.NameserverStatistics["TIMEOUT"] != 10 {
		utils.Fatalln("Not retrieving nameserver statistics correctly", nil)
	}

	if len(scan.DSStatistics) != 2 ||
		scan.DSStatistics["OK"] != 8 ||
		scan.DSStatistics["DNSERR"] != 2 {
		utils.Fatalln("Not retrieving DS statistics correctly", nil)
	}

	if err := scanDAO.RemoveAll(); err != nil {
		utils.Fatalln("Error removing scans from database", err)
	}
}

// Function to mock a scan object
func newScan() model.Scan {
	return model.Scan{
		Status:                   model.ScanStatusExecuted,
		StartedAt:                time.Now().Add(-10 * time.Minute),
		FinishedAt:               time.Now().Add(-5 * time.Minute),
		LastModifiedAt:           time.Now().Add(-5 * time.Minute),
		DomainsScanned:           50,
		DomainsWithDNSSECScanned: 10,
		NameserverStatistics:     make(map[string]uint64),
		DSStatistics:             make(map[string]uint64),
	}
}

// Function to compare if two scans are equal, cannot use operator == because of the
// slices inside the scan object
func compareScans(s1, s2 model.Scan) bool {
	if s1.Id != s2.Id ||
		s1.Revision != s2.Revision ||
		s1.Status != s2.Status ||
		s1.StartedAt.Unix() != s2.StartedAt.Unix() ||
		s1.FinishedAt.Unix() != s2.FinishedAt.Unix() ||
		s1.DomainsScanned != s2.DomainsScanned ||
		s1.DomainsWithDNSSECScanned != s2.DomainsWithDNSSECScanned {
		return false
	}

	for key, value := range s1.NameserverStatistics {
		if otherValue, ok := s2.NameserverStatistics[key]; !ok || value != otherValue {
			return false
		}
	}

	for key, value := range s1.DSStatistics {
		if otherValue, ok := s2.DSStatistics[key]; !ok || value != otherValue {
			return false
		}
	}

	return true
}

// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

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

	database, databaseSession, err := mongodb.Open(
		[]string{config.Database.URI},
		config.Database.Name,
		false, "", "",
	)

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
	scanStatistics(scanDAO)
	scansPagination(scanDAO)
	scansExpand(scanDAO)

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

	} else if !utils.CompareScan(scan, scanRetrieved) {
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

	} else if !utils.CompareScan(scan, scanRetrieved) {
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

	pagination := model.ScanPagination{
		PageSize: 10,
		Page:     5,
		OrderBy: []model.ScanSort{
			{
				Field:     model.ScanOrderByFieldStartedAt,
				Direction: model.OrderByDirectionAscending,
			},
		},
	}

	scans, err := scanDAO.FindAll(&pagination, true)
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

	pagination = model.ScanPagination{
		PageSize: 10000,
		Page:     1,
		OrderBy: []model.ScanSort{
			{
				Field:     model.ScanOrderByFieldStartedAt,
				Direction: model.OrderByDirectionAscending,
			},
		},
	}

	scans, err = scanDAO.FindAll(&pagination, true)
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

func scansExpand(scanDAO dao.ScanDAO) {
	newScans := newScans()
	for _, scan := range newScans {
		if err := scanDAO.Save(&scan); err != nil {
			utils.Fatalln("Error creating scans", err)
		}
	}

	pagination := model.ScanPagination{}
	scans, err := scanDAO.FindAll(&pagination, false)

	if err != nil {
		utils.Fatalln("Error retrieving scans", err)
	}

	for _, scan := range scans {
		if scan.DomainsScanned != 0 ||
			scan.DomainsWithDNSSECScanned != 0 ||
			!scan.FinishedAt.Equal(time.Time{}) ||
			len(scan.NameserverStatistics) > 0 ||
			len(scan.DSStatistics) > 0 {
			utils.Fatalln("Not compressing scan in results", nil)
		}
	}

	scans, err = scanDAO.FindAll(&pagination, true)

	if err != nil {
		utils.Fatalln("Error retrieving scans", err)
	}

	for _, scan := range scans {
		if scan.DomainsScanned == 0 ||
			scan.DomainsWithDNSSECScanned == 0 ||
			scan.FinishedAt.Equal(time.Time{}) ||
			scan.StartedAt.Equal(time.Time{}) ||
			len(scan.NameserverStatistics) == 0 ||
			len(scan.DSStatistics) == 0 {
			fmt.Println(scan.DomainsScanned == 0,
				scan.DomainsWithDNSSECScanned == 0,
				scan.FinishedAt.Equal(time.Time{}),
				scan.StartedAt.Equal(time.Time{}),
				len(scan.NameserverStatistics) == 0,
				len(scan.DSStatistics) == 0)
			utils.Fatalln("Compressing scan in results when it shouldn't", nil)
		}
	}

	if err := scanDAO.RemoveAll(); err != nil {
		utils.Fatalln("Error removing scans", err)
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
		NameserverStatistics: map[string]uint64{
			"OK":      48,
			"TIMEOUT": 2,
		},
		DSStatistics: map[string]uint64{
			"OK": 10,
		},
	}
}

func newScans() []model.Scan {
	return []model.Scan{
		{
			Status:                   model.ScanStatusExecuted,
			StartedAt:                time.Now().Add(-10 * time.Minute),
			FinishedAt:               time.Now().Add(-5 * time.Minute),
			LastModifiedAt:           time.Now().Add(-5 * time.Minute),
			DomainsScanned:           50,
			DomainsWithDNSSECScanned: 10,
			NameserverStatistics: map[string]uint64{
				"OK":      48,
				"TIMEOUT": 2,
			},
			DSStatistics: map[string]uint64{
				"OK": 10,
			},
		},
		{
			Status:                   model.ScanStatusExecuted,
			StartedAt:                time.Now().Add(-20 * time.Minute),
			FinishedAt:               time.Now().Add(-15 * time.Minute),
			LastModifiedAt:           time.Now().Add(-15 * time.Minute),
			DomainsScanned:           20,
			DomainsWithDNSSECScanned: 12,
			NameserverStatistics: map[string]uint64{
				"OK":      10,
				"TIMEOUT": 10,
			},
			DSStatistics: map[string]uint64{
				"OK":     10,
				"EXPSIG": 2,
			},
		},
	}
}

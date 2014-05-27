// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"github.com/rafaeljusto/shelter/net/http/rest/handler"
	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
	"github.com/rafaeljusto/shelter/testing/utils"
	"labix.org/v2/mgo"
	"net/http"
	"time"
)

var (
	configFilePath string // Path for the config file with the connection information
)

// RESTHandlerScanTestConfigFile is a structure to store the test configuration file data
type RESTHandlerScanTestConfigFile struct {
	Database struct {
		URI  string
		Name string
	}
}

func init() {
	utils.TestName = "RESTHandlerScan"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for RESTHandlerScan test")
}

func main() {
	flag.Parse()

	var config RESTHandlerScanTestConfigFile
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

	// If there was some problem in the last test, there could be some data in the
	// database, so let's clear it to don't affect this test. We avoid checking the error,
	// because if the collection does not exist yet, it will be created in the first
	// insert
	utils.ClearDatabase(database)

	startedAt := createScan(database)
	retrieveScan(database, startedAt)
	retrieveScanMetadata(database, startedAt)
	retrieveUnknownScan(database, startedAt)
	createCurrentScan()
	retrieveCurrentScan(database)
	deleteScan(database, startedAt)

	utils.Println("SUCCESS!")
}

func createScan(database *mgo.Database) time.Time {
	scan := model.Scan{
		Status:                   model.ScanStatusExecuted,
		StartedAt:                time.Now().Add(-60 * time.Minute),
		FinishedAt:               time.Now().Add(-10 * time.Minute),
		DomainsScanned:           5000000,
		DomainsWithDNSSECScanned: 250000,
		NameserverStatistics: map[string]uint64{
			"OK":      4800000,
			"TIMEOUT": 200000,
		},
		DSStatistics: map[string]uint64{
			"OK":     220000,
			"EXPSIG": 30000,
		},
	}

	scanDAO := dao.ScanDAO{
		Database: database,
	}

	if err := scanDAO.Save(&scan); err != nil {
		utils.Fatalln("Error creating scan", err)
	}

	return scan.StartedAt
}

func retrieveScan(database *mgo.Database, startedAt time.Time) {
	r, err := http.NewRequest("GET", fmt.Sprintf("/scan/%s", startedAt.Format(time.RFC3339Nano)), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleScan(r, &context)

	if context.ResponseHTTPStatus != http.StatusOK {
		utils.Fatalln("Error retrieving scan",
			errors.New(string(context.ResponseContent)))
	}

	if context.HTTPHeader["ETag"] != "1" {
		utils.Fatalln("Not setting ETag in scan retrieve response", nil)
	}

	if len(context.HTTPHeader["Last-Modified"]) == 0 {
		utils.Fatalln("Not setting Last-Modified in scan retrieve response", nil)
	}

	var scanResponse protocol.ScanResponse
	json.Unmarshal(context.ResponseContent, &scanResponse)

	if scanResponse.Status != "EXECUTED" {
		utils.Fatalln("Scan's status was not persisted correctly", nil)
	}

	if scanResponse.DomainsScanned != 5000000 {
		utils.Fatalln("Scan's domains scanned information is wrong", nil)
	}

	if scanResponse.DomainsWithDNSSECScanned != 250000 {
		utils.Fatalln("Scan's domains with DNSSEC scanned information is wrong", nil)
	}

	if scanResponse.NameserverStatistics["OK"] != 4800000 ||
		scanResponse.NameserverStatistics["TIMEOUT"] != 200000 {
		utils.Fatalln("Scan's nameserver statistics are wrong", nil)
	}

	if scanResponse.DSStatistics["OK"] != 220000 ||
		scanResponse.DSStatistics["EXPSIG"] != 30000 {
		utils.Fatalln("Scan's nameserver statistics are wrong", nil)
	}
}

func retrieveScanMetadata(database *mgo.Database, startedAt time.Time) {
	r, err := http.NewRequest("GET", fmt.Sprintf("/scan/%s", startedAt.Format(time.RFC3339Nano)), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	context1, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleScan(r, &context1)

	if context1.ResponseHTTPStatus != http.StatusOK {
		utils.Fatalln("Error retrieving scan",
			errors.New(string(context1.ResponseContent)))
	}

	r, err = http.NewRequest("HEAD", fmt.Sprintf("/scan/%s", startedAt.Format(time.RFC3339Nano)), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	context2, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleScan(r, &context2)

	if context2.ResponseHTTPStatus != http.StatusOK {
		utils.Fatalln("Error retrieving scan",
			errors.New(string(context2.ResponseContent)))
	}

	if string(context1.ResponseContent) != string(context2.ResponseContent) {
		utils.Fatalln("At this point the GET and HEAD method should "+
			"return the same body content", nil)
	}
}

func retrieveUnknownScan(database *mgo.Database, startedAt time.Time) {
	r, err := http.NewRequest("GET", fmt.Sprintf("/scan/%s", startedAt.Add(-1*time.Millisecond).Format(time.RFC3339Nano)), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleScan(r, &context)

	if context.ResponseHTTPStatus != http.StatusNotFound {
		utils.Fatalln("Error retrieving unknown scan",
			errors.New(string(context.ResponseContent)))
	}
}

func createCurrentScan() {
	// We are not going to call the function FinishAnalyzingDomainForScan because it will reset the
	// current scan information
	model.StartNewScan()
	model.LoadedDomainForScan()
	model.LoadedDomainForScan()
	model.LoadedDomainForScan()
	model.FinishAnalyzingDomainForScan(false)
	model.LoadedDomainForScan()
	model.LoadedDomainForScan()
	model.FinishLoadingDomainsForScan()
	model.FinishAnalyzingDomainForScan(false)
	model.FinishAnalyzingDomainForScan(false)
	model.FinishAnalyzingDomainForScan(true)
}

func retrieveCurrentScan(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/scan/current", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleScan(r, &context)

	if context.ResponseHTTPStatus != http.StatusOK {
		utils.Fatalln("Error retrieving scan",
			errors.New(string(context.ResponseContent)))
	}

	if len(context.HTTPHeader["Last-Modified"]) > 0 {
		utils.Fatalln("Current scan is not persisted to return Last-Modified header", nil)
	}

	var scanResponse protocol.ScanResponse
	json.Unmarshal(context.ResponseContent, &scanResponse)

	if scanResponse.Status != "RUNNING" ||
		scanResponse.DomainsToBeScanned != 5 ||
		scanResponse.DomainsScanned != 4 ||
		scanResponse.DomainsWithDNSSECScanned != 1 {
		utils.Fatalln("Not retrieving current scan information correctly", nil)
	}
}

func deleteScan(database *mgo.Database, startedAt time.Time) {
	scanDAO := dao.ScanDAO{
		Database: database,
	}

	if err := scanDAO.RemoveByStartedAt(startedAt); err != nil {
		utils.Fatalln("Error removing scan", err)
	}
}

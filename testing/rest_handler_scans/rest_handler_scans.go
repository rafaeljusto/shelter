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

// RESTHandlerScansTestConfigFile is a structure to store the test configuration file data
type RESTHandlerScansTestConfigFile struct {
	Database struct {
		URI  string
		Name string
	}
}

type ScansCacheTestData struct {
	HeaderValue        string
	ExpectedHTTPStatus int
}

func init() {
	utils.TestName = "RESTHandlerScans"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for RESTHandlerScans test")
}

func main() {
	flag.Parse()

	var config RESTHandlerScansTestConfigFile
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

	createScans(database)
	retrieveScans(database)
	retrieveScansMetadata(database)
	retrieveScansIfModifiedSince(database)
	retrieveScansIfUnmodifiedSince(database)
	retrieveScansIfMatch(database)
	retrieveScansIfNoneMatch(database)
	deleteScans(database)

	utils.Println("SUCCESS!")
}

func createScans(database *mgo.Database) {
	scanDAO := dao.ScanDAO{
		Database: database,
	}

	for i := 0; i < 100; i++ {
		scan := model.Scan{
			Status:                   model.ScanStatusExecuted,
			StartedAt:                time.Now().Add(time.Duration(-i*2) * time.Minute),
			FinishedAt:               time.Now().Add(time.Duration(-i) * time.Minute),
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

		if err := scanDAO.Save(&scan); err != nil {
			utils.Fatalln("Error creating scan", err)
		}
	}
}

func retrieveScans(database *mgo.Database) {
	data := []struct {
		URI                string
		ExpectedHTTPStatus int
		ContentCheck       func([]byte)
	}{
		{
			URI:                "/scans/?orderby=xxx:desc&pagesize=10&page=1",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			URI:                "/scans/?orderby=startedat:xxx&pagesize=10&page=1",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			URI:                "/scans/?orderby=startedat:desc&pagesize=xxx&page=1",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			URI:                "/scans/?orderby=startedat:desc&pagesize=10&page=xxx",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			URI:                "/scans/?orderby=startedat:desc&pagesize=10&page=1",
			ExpectedHTTPStatus: http.StatusOK,
			ContentCheck: func(content []byte) {
				var scansResponse protocol.ScansResponse
				json.Unmarshal(content, &scansResponse)

				if len(scansResponse.Scans) != 10 {
					utils.Fatalln("Error retrieving the wrong number of scans", nil)
				}
			},
		},
		{
			URI:                "/scans",
			ExpectedHTTPStatus: http.StatusOK,
			ContentCheck: func(content []byte) {
				var scansResponse protocol.ScansResponse
				json.Unmarshal(content, &scansResponse)

				if len(scansResponse.Scans) == 0 {
					utils.Fatalln("Error retrieving scans", nil)
				}

				if scansResponse.Scans[0].DomainsScanned > 0 {
					utils.Fatalln("Not compressing scans result", nil)
				}
			},
		},
		{
			URI:                "/scans?expand",
			ExpectedHTTPStatus: http.StatusOK,
			ContentCheck: func(content []byte) {
				var scansResponse protocol.ScansResponse
				json.Unmarshal(content, &scansResponse)

				if len(scansResponse.Scans) == 0 {
					utils.Fatalln("Error retrieving scans", nil)
				}

				if scansResponse.Scans[0].DomainsScanned == 0 {
					utils.Fatalln("Compressing scans result when it shouldn't", nil)
				}
			},
		},
	}

	for _, item := range data {
		r, err := http.NewRequest("GET", item.URI, nil)
		if err != nil {
			utils.Fatalln("Error creating the HTTP request", err)
		}

		context, err := context.NewContext(r, database)
		if err != nil {
			utils.Fatalln("Error creating context", err)
		}

		handler.HandleScans(r, &context)

		if context.ResponseHTTPStatus != item.ExpectedHTTPStatus {
			utils.Fatalln(fmt.Sprintf("Error when requesting scans using the URI [%s]. "+
				"Expected HTTP status code %d but got %d", item.URI,
				item.ExpectedHTTPStatus, context.ResponseHTTPStatus), nil)
		}

		if item.ContentCheck != nil {
			item.ContentCheck(context.ResponseContent)
		}
	}
}

func retrieveScansMetadata(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/scans/?orderby=startedat:desc&pagesize=10&page=1", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	context1, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleScans(r, &context1)

	if context1.ResponseHTTPStatus != http.StatusOK {
		utils.Fatalln("Error retrieving scans",
			errors.New(string(context1.ResponseContent)))
	}

	r, err = http.NewRequest("HEAD", "/scans/?orderby=startedat:desc&pagesize=10&page=1", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	context2, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleScans(r, &context2)

	if context2.ResponseHTTPStatus != http.StatusOK {
		utils.Fatalln("Error retrieving scans",
			errors.New(string(context2.ResponseContent)))
	}

	if string(context1.ResponseContent) != string(context2.ResponseContent) {
		utils.Fatalln("At this point the GET and HEAD method should "+
			"return the same body content", nil)
	}
}

func retrieveScansIfModifiedSince(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/scans", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	scansCacheTest(database, r, "If-Modified-Since", []ScansCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			HeaderValue:        time.Now().Add(1 * time.Second).UTC().Format(time.RFC1123),
			ExpectedHTTPStatus: http.StatusNotModified,
		},
	})
}

func retrieveScansIfUnmodifiedSince(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/scans", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	scansCacheTest(database, r, "If-Unmodified-Since", []ScansCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			HeaderValue:        time.Now().Add(-10 * time.Second).UTC().Format(time.RFC1123),
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
	})
}

func retrieveScansIfMatch(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/scans", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleScans(r, &context)

	scansCacheTest(database, r, "If-Match", []ScansCacheTestData{
		{
			HeaderValue:        context.HTTPHeader["ETag"] + "x",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
		{
			HeaderValue:        context.HTTPHeader["ETag"],
			ExpectedHTTPStatus: http.StatusOK,
		},
	})
}

func retrieveScansIfNoneMatch(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/scans", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleScans(r, &context)

	scansCacheTest(database, r, "If-None-Match", []ScansCacheTestData{
		{
			HeaderValue:        context.HTTPHeader["ETag"],
			ExpectedHTTPStatus: http.StatusNotModified,
		},
		{
			HeaderValue:        context.HTTPHeader["ETag"] + "x",
			ExpectedHTTPStatus: http.StatusOK,
		},
	})
}

func deleteScans(database *mgo.Database) {
	scanDAO := dao.ScanDAO{
		Database: database,
	}

	if err := scanDAO.RemoveAll(); err != nil {
		utils.Fatalln("Error removing scans", err)
	}
}

func scansCacheTest(database *mgo.Database, r *http.Request,
	header string, scansCacheTestData []ScansCacheTestData) {

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	for _, item := range scansCacheTestData {
		r.Header.Set(header, item.HeaderValue)

		handler.HandleScans(r, &context)

		if context.ResponseHTTPStatus != item.ExpectedHTTPStatus {
			utils.Fatalln(fmt.Sprintf("Error in %s test using %s [%s] "+
				"HTTP header. Expected HTTP status code %d and got %d",
				r.Method, header, item.HeaderValue, item.ExpectedHTTPStatus,
				context.ResponseHTTPStatus), errors.New(string(context.ResponseContent)))
		}
	}
}

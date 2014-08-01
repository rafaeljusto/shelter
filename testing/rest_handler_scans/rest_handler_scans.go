// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/rafaeljusto/handy"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/handler"
	"github.com/rafaeljusto/shelter/protocol"
	"github.com/rafaeljusto/shelter/testing/utils"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var (
	configFilePath string // Path for the config file with the connection information
)

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

	err := utils.ReadConfigFile(configFilePath, &config.ShelterConfig)

	if err == utils.ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		utils.Fatalln("Error reading configuration file", err)
	}

	database, databaseSession, err := mongodb.Open(
		config.ShelterConfig.Database.URIs,
		config.ShelterConfig.Database.Name,
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
	createCurrentScan()
	retrieveCurrentScan(database)

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
		ContentCheck       func(*protocol.ScansResponse)
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
			ContentCheck: func(scansResponse *protocol.ScansResponse) {
				if len(scansResponse.Scans) != 10 {
					utils.Fatalln("Error retrieving the wrong number of scans", nil)
				}
			},
		},
		{
			URI:                "/scans",
			ExpectedHTTPStatus: http.StatusOK,
			ContentCheck: func(scansResponse *protocol.ScansResponse) {
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
			ContentCheck: func(scansResponse *protocol.ScansResponse) {
				if len(scansResponse.Scans) == 0 {
					utils.Fatalln("Error retrieving scans", nil)
				}

				if scansResponse.Scans[0].DomainsScanned == 0 {
					utils.Fatalln("Compressing scans result when it shouldn't", nil)
				}
			},
		},
		{
			URI:                "/scans?current",
			ExpectedHTTPStatus: http.StatusOK,
			ContentCheck: func(scansResponse *protocol.ScansResponse) {
				if len(scansResponse.Scans) == 0 {
					utils.Fatalln("Error retrieving current scan", nil)
				}

				if scansResponse.Scans[0].Status != "WAITINGEXECUTION" {
					utils.Fatalln("Current scan with wrong status", nil)
				}
			},
		},
	}

	mux := handy.NewHandy()

	h := new(handler.ScansHandler)
	mux.Handle("/scans", func() handy.Handler {
		return h
	})

	for _, item := range data {
		r, err := http.NewRequest("GET", item.URI, nil)
		if err != nil {
			utils.Fatalln("Error creating the HTTP request", err)
		}
		utils.BuildHTTPHeader(r, nil)

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)

		responseContent, err := ioutil.ReadAll(w.Body)
		if err != nil {
			utils.Fatalln("Error reading response body", err)
		}

		if w.Code != item.ExpectedHTTPStatus {
			utils.Fatalln(fmt.Sprintf("Error when requesting scans using the URI [%s]. "+
				"Expected HTTP status code %d but got %d",
				item.URI, item.ExpectedHTTPStatus, w.Code), errors.New(string(responseContent)))
		}

		if item.ContentCheck != nil {
			item.ContentCheck(h.Response)
		}
	}
}

func retrieveScansMetadata(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.ScansHandler)
	mux.Handle("/scans", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/scans/?orderby=startedat:desc&pagesize=10&page=1", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, nil)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	response1 := *h.Response
	responseContent, err := ioutil.ReadAll(w.Body)
	if err != nil {
		utils.Fatalln("Error reading response body", err)
	}

	if w.Code != http.StatusOK {
		utils.Fatalln("Error retrieving scans", errors.New(string(responseContent)))
	}

	r, err = http.NewRequest("HEAD", "/scans/?orderby=startedat:desc&pagesize=10&page=1", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, nil)

	mux.ServeHTTP(w, r)

	response2 := *h.Response
	responseContent, err = ioutil.ReadAll(w.Body)
	if err != nil {
		utils.Fatalln("Error reading response body", err)
	}

	if w.Code != http.StatusOK {
		utils.Fatalln("Error retrieving scans", errors.New(string(responseContent)))
	}

	if !utils.CompareProtocolScans(response1, response2) {
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
	mux := handy.NewHandy()

	h := new(handler.ScansHandler)
	mux.Handle("/scans", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/scans", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, nil)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	scansCacheTest(database, r, "If-Match", []ScansCacheTestData{
		{
			HeaderValue:        w.Header().Get("ETag") + "x",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
		{
			HeaderValue:        w.Header().Get("ETag"),
			ExpectedHTTPStatus: http.StatusOK,
		},
	})
}

func retrieveScansIfNoneMatch(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.ScansHandler)
	mux.Handle("/scans", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/scans", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, nil)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	scansCacheTest(database, r, "If-None-Match", []ScansCacheTestData{
		{
			HeaderValue:        w.Header().Get("ETag"),
			ExpectedHTTPStatus: http.StatusNotModified,
		},
		{
			HeaderValue:        w.Header().Get("ETag") + "x",
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

	mux := handy.NewHandy()

	h := new(handler.ScansHandler)
	mux.Handle("/scans", func() handy.Handler {
		return h
	})

	for _, item := range scansCacheTestData {
		r.Header.Set(header, item.HeaderValue)
		utils.BuildHTTPHeader(r, nil)

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)

		responseContent, err := ioutil.ReadAll(w.Body)
		if err != nil {
			utils.Fatalln("Error reading response body", err)
		}

		if w.Code != item.ExpectedHTTPStatus {
			utils.Fatalln(fmt.Sprintf("Error in %s test using %s [%s] "+
				"HTTP header. Expected HTTP status code %d and got %d",
				r.Method, header, item.HeaderValue, item.ExpectedHTTPStatus,
				w.Code), errors.New(string(responseContent)))
		}
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
	mux := handy.NewHandy()

	h := new(handler.ScansHandler)
	mux.Handle("/scans", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/scans?current", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, nil)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	responseContent, err := ioutil.ReadAll(w.Body)
	if err != nil {
		utils.Fatalln("Error reading response body", err)
	}

	if w.Code != http.StatusOK {
		utils.Fatalln("Error retrieving scan", errors.New(string(responseContent)))
	}

	if len(w.Header().Get("Last-Modified")) == 0 {
		utils.Fatalln("Current scan should return the Last-Modified header", nil)
	}

	if len(h.Response.Scans) != 1 {
		utils.Fatalln("Current scan not returned when it should", nil)
	}

	if h.Response.Scans[0].Status != "RUNNING" ||
		h.Response.Scans[0].DomainsToBeScanned != 5 ||
		h.Response.Scans[0].DomainsScanned != 4 ||
		h.Response.Scans[0].DomainsWithDNSSECScanned != 1 {

		utils.Fatalln("Not retrieving current scan information correctly", nil)
	}
}

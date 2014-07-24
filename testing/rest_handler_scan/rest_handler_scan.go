// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/handler"
	"github.com/rafaeljusto/shelter/testing/utils"
	"github.com/trajber/handy"
	"io/ioutil"
	"gopkg.in/mgo.v2"
	"net/http"
	"net/http/httptest"
	"time"
)

var (
	configFilePath string // Path for the config file with the connection information
)

func init() {
	utils.TestName = "RESTHandlerScan"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for RESTHandlerScan test")
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

	startedAt := createScan(database)
	retrieveScan(database, startedAt)
	retrieveScanMetadata(database, startedAt)
	retrieveUnknownScan(database, startedAt)
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
	mux := handy.NewHandy()

	h := new(handler.ScanHandler)
	mux.Handle("/scan/{started-at}", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", fmt.Sprintf("/scan/%s", startedAt.Format(time.RFC3339Nano)), nil)
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

	if w.Header().Get("ETag") != "1" {
		utils.Fatalln("Not setting ETag in scan retrieve response", nil)
	}

	if len(w.Header().Get("Last-Modified")) == 0 {
		utils.Fatalln("Not setting Last-Modified in scan retrieve response", nil)
	}

	if h.Response.Status != "EXECUTED" {
		utils.Fatalln("Scan's status was not persisted correctly", nil)
	}

	if h.Response.DomainsScanned != 5000000 {
		utils.Fatalln("Scan's domains scanned information is wrong", nil)
	}

	if h.Response.DomainsWithDNSSECScanned != 250000 {
		utils.Fatalln("Scan's domains with DNSSEC scanned information is wrong", nil)
	}

	if h.Response.NameserverStatistics["OK"] != 4800000 ||
		h.Response.NameserverStatistics["TIMEOUT"] != 200000 {
		utils.Fatalln("Scan's nameserver statistics are wrong", nil)
	}

	if h.Response.DSStatistics["OK"] != 220000 ||
		h.Response.DSStatistics["EXPSIG"] != 30000 {
		utils.Fatalln("Scan's nameserver statistics are wrong", nil)
	}
}

func retrieveScanMetadata(database *mgo.Database, startedAt time.Time) {
	mux := handy.NewHandy()

	h := new(handler.ScanHandler)
	mux.Handle("/scan/{started-at}", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", fmt.Sprintf("/scan/%s", startedAt.Format(time.RFC3339Nano)), nil)
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
		utils.Fatalln("Error retrieving scan", errors.New(string(responseContent)))
	}

	r, err = http.NewRequest("HEAD", fmt.Sprintf("/scan/%s", startedAt.Format(time.RFC3339Nano)), nil)
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
		utils.Fatalln("Error retrieving scan", errors.New(string(responseContent)))
	}

	if !utils.CompareProtocolScan(response1, response2) {
		utils.Fatalln("At this point the GET and HEAD method should "+
			"return the same body content", nil)
	}
}

func retrieveUnknownScan(database *mgo.Database, startedAt time.Time) {
	mux := handy.NewHandy()

	h := new(handler.ScanHandler)
	mux.Handle("/scan/{started-at}", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", fmt.Sprintf("/scan/%s", startedAt.Add(-1*time.Millisecond).Format(time.RFC3339Nano)), nil)
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

	if w.Code != http.StatusNotFound {
		utils.Fatalln("Error retrieving unknown scan", errors.New(string(responseContent)))
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

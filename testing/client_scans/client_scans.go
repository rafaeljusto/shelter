// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
	"github.com/rafaeljusto/shelter/testing/utils"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	configFilePath string // Path for the config file with the connection information
)

func init() {
	utils.TestName = "ClientScans"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for client scans test")
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

	finishRESTServer := utils.StartRESTServer()
	defer finishRESTServer()
	utils.StartWebClient()

	database, databaseSession, err := mongodb.Open(
		config.ShelterConfig.Database.URIs,
		config.ShelterConfig.Database.Name,
		config.ShelterConfig.Database.Auth.Enabled,
		config.ShelterConfig.Database.Auth.Username,
		config.ShelterConfig.Database.Auth.Password,
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
	retrieveScans()
	retrieveScansWithPagination()
	retrieveScansWithCache()
	retrieveCurrentScan()

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

func retrieveScans() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	r, err := http.NewRequest("GET", fmt.Sprintf("%s%s", url, "/scans"), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	utils.BuildHTTPHeader(r, nil)

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Fatalln("Error reading response content", err)
	}

	if response.StatusCode != http.StatusOK {
		utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method GET "+
			"and URI /scans",
			http.StatusOK, response.StatusCode),
			errors.New(string(responseContent)),
		)
	}

	var scansResponse protocol.ScansResponse
	if err := json.Unmarshal(responseContent, &scansResponse); err != nil {
		utils.Fatalln("Error decoding scan response", err)
	}

	if scansResponse.PageSize != 20 {
		utils.Fatalln("Default page size isn't 20", nil)
	}

	if scansResponse.Page != 1 {
		utils.Fatalln("Default page isn't the first one", nil)
	}

	if scansResponse.NumberOfItems != 100 {
		utils.Fatalln("Not calculating the correct number of scans", nil)
	}

	if scansResponse.NumberOfPages != 5 {
		utils.Fatalln("Not calculating the correct number of pages", nil)
	}

	if len(scansResponse.Scans) != scansResponse.PageSize {
		utils.Fatalln("Not returning all desired scans", nil)
	}
}

func retrieveScansWithPagination() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	r, err := http.NewRequest("GET", fmt.Sprintf("%s%s", url, "/scans?page=5&pagesize=20"), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	utils.BuildHTTPHeader(r, nil)

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Fatalln("Error reading response content", err)
	}

	if response.StatusCode != http.StatusOK {
		utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method GET "+
			"and URI /scans",
			http.StatusOK, response.StatusCode),
			errors.New(string(responseContent)),
		)
	}

	var scansResponse protocol.ScansResponse
	if err := json.Unmarshal(responseContent, &scansResponse); err != nil {
		utils.Fatalln("Error decoding scan response", err)
	}

	if scansResponse.PageSize != 20 {
		utils.Fatalln("Default page size isn't 20", nil)
	}

	if scansResponse.Page != 5 {
		utils.Fatalln("Default page isn't the first one", nil)
	}

	if scansResponse.NumberOfItems != 100 {
		utils.Fatalln("Not calculating the correct number of scans", nil)
	}

	if scansResponse.NumberOfPages != 5 {
		utils.Fatalln("Not calculating the correct number of pages", nil)
	}

	if len(scansResponse.Scans) != scansResponse.PageSize {
		utils.Fatalln("Not returning all desired scans", nil)
	}
}

func retrieveScansWithCache() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	r, err := http.NewRequest("GET", fmt.Sprintf("%s%s", url, "/scans"), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	utils.BuildHTTPHeader(r, nil)

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	r, err = http.NewRequest("GET", fmt.Sprintf("%s%s", url, "/scans"), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	r.Header.Add("If-None-Match", response.Header.Get("ETag"))
	utils.BuildHTTPHeader(r, nil)

	response, err = client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	if response.StatusCode != http.StatusNotModified {
		utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method GET and URI /scans",
			http.StatusNotModified, response.StatusCode), nil)
	}
}

func retrieveCurrentScan() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	r, err := http.NewRequest("GET", fmt.Sprintf("%s%s", url, "/scans/?current"), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	utils.BuildHTTPHeader(r, nil)

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Fatalln("Error reading response content", err)
	}

	if response.StatusCode != http.StatusOK {
		utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method GET "+
			"and URI /scan/currentscan",
			http.StatusOK, response.StatusCode),
			errors.New(string(responseContent)),
		)
	}

	var scansResponse protocol.ScansResponse
	if err := json.Unmarshal(responseContent, &scansResponse); err != nil {
		utils.Fatalln("Error decoding scan response", err)
	}

	if len(scansResponse.Scans) != 1 {
		utils.Fatalln("Not returning the current scan in the result set or "+
			"more items than expected were found", nil)
	}

	if scansResponse.Scans[0].Status != model.ScanStatusToString(model.ScanStatusWaitingExecution) {
		utils.Fatalln(fmt.Sprintf("Invalid status returned by current scan. Expected %s and got %s",
			model.ScanStatusToString(model.ScanStatusWaitingExecution), scansResponse.Scans[0].Status), nil)
	}
}

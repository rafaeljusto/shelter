// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
	"github.com/rafaeljusto/shelter/testing/utils"
	"io/ioutil"
	"labix.org/v2/mgo"
	"net/http"
	"time"
)

var (
	configFilePath string // Path for the config file with the connection information
)

func init() {
	utils.TestName = "ClientScan"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for client scan test")
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

	startedAt := createScan(database)
	retrieveScan(startedAt)

	utils.Println("SUCCESS!")
}

func createScan(database *mgo.Database) time.Time {
	scanDAO := dao.ScanDAO{
		Database: database,
	}

	scan := model.Scan{
		Status:                   model.ScanStatusExecuted,
		StartedAt:                time.Now().Add(time.Duration(-2) * time.Minute),
		FinishedAt:               time.Now().Add(time.Duration(-1) * time.Minute),
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

	return scan.StartedAt
}

func retrieveScan(startedAt time.Time) {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	uri := fmt.Sprintf("/scan/%s", startedAt.Format(time.RFC3339Nano))
	r, err := http.NewRequest("GET", fmt.Sprintf("%s%s", url, uri), nil)
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

	var scanResponse protocol.ScanResponse
	if err := json.Unmarshal(responseContent, &scanResponse); err != nil {
		utils.Fatalln("Error decoding scan response", err)
	}

	if scanResponse.Status != model.ScanStatusToString(model.ScanStatusExecuted) {
		utils.Fatalln("Invalid status returned by scan", nil)
	}

	// Tolerance of one second
	if scanResponse.StartedAt.Sub(startedAt) > (1*time.Second) ||
		scanResponse.StartedAt.Sub(startedAt) < (-1*time.Second) {
		utils.Fatalln("Invalid start date returned by scan", nil)
	}

	if scanResponse.DomainsScanned != 5000000 {
		utils.Fatalln("Invalid number of domains scanned returned by scan", nil)
	}

	if scanResponse.DomainsWithDNSSECScanned != 250000 {
		utils.Fatalln("Invalid number of domains with DNSSEC scanned returned by scan", nil)
	}

	if len(scanResponse.NameserverStatistics) != 2 {
		utils.Fatalln("Invalid nameserver statistics returned by scan", nil)
	}

	if len(scanResponse.DSStatistics) != 2 {
		utils.Fatalln("Invalid dsset statistics returned by scan", nil)
	}
}

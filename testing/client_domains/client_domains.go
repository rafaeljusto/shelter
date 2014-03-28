// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
	"github.com/rafaeljusto/shelter/testing/utils"
	"io/ioutil"
	"net/http"
)

var (
	configFilePath string // Path for the config file with the connection information
)

func init() {
	utils.TestName = "ClientDomains"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for client domain test")
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
		config.ShelterConfig.Database.URI,
		config.ShelterConfig.Database.Name,
	)

	if err != nil {
		utils.Fatalln("Error connecting the database", err)
	}
	defer databaseSession.Close()

	// If there was some problem in the last test, there could be some data in the database,
	// so let's clear it to don't affect this test. We avoid checking the error, because if
	// the collection does not exist yet, it will be created in the first insert
	utils.ClearDatabase(database)

	finishRESTServer := utils.StartRESTServer()
	defer finishRESTServer()
	utils.StartWebClient()

	createDomains()
	retrieveDomains()
	retrieveDomainsWithPagination()
	removeDomains()

	utils.Println("SUCCESS!")
}

func createDomains() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	var r *http.Request
	var err error

	for i := 0; i < 100; i++ {
		uri := fmt.Sprintf("/domain/example%d.com.br.", i)
		content := `{
  "Nameservers": [
    { "host": "a.dns.br." },
    { "host": "b.dns.br." }
  ],
  "Owners": [
    { "email": "admin@gmail.com.", "language": "pt-br" }
  ]
}`

		r, err = http.NewRequest("PUT", fmt.Sprintf("%s%s", url, uri),
			bytes.NewReader([]byte(content)))

		if err != nil {
			utils.Fatalln("Error creating the HTTP request", err)
		}

		utils.BuildHTTPHeader(r, []byte(content))

		response, err := client.Do(r)
		if err != nil {
			utils.Fatalln("Error sending request", err)
		}

		if response.StatusCode != http.StatusCreated {
			responseContent, err := ioutil.ReadAll(response.Body)

			if err == nil {
				utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method PUT and URI %s",
					http.StatusCreated, response.StatusCode, uri),
					errors.New(string(responseContent)),
				)
			} else {
				utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method PUT and URI %s",
					http.StatusCreated, response.StatusCode, uri), nil)
			}
		}
	}
}

func retrieveDomains() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	var r *http.Request
	var err error

	r, err = http.NewRequest("GET", fmt.Sprintf("%s%s", url, "/domains"), nil)
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
		utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method GET and URI /domains",
			http.StatusCreated, response.StatusCode),
			errors.New(string(responseContent)),
		)
	}

	var domainsResponse protocol.DomainsResponse
	if err := json.Unmarshal(responseContent, &domainsResponse); err != nil {
		utils.Fatalln("Error parsing response JSON", err)
	}

	if domainsResponse.PageSize != 20 {
		utils.Fatalln("Default page size isn't 20", nil)
	}

	if domainsResponse.Page != 1 {
		utils.Fatalln("Default page isn't the first one", nil)
	}

	if domainsResponse.NumberOfItems != 100 {
		utils.Fatalln("Not calculating the correct number of domains", nil)
	}

	if domainsResponse.NumberOfPages != 5 {
		utils.Fatalln("Not calculating the correct number of pages", nil)
	}

	if len(domainsResponse.Domains) != domainsResponse.PageSize {
		utils.Fatalln("Not returning all desired domains", nil)
	}
}

func retrieveDomainsWithPagination() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	var r *http.Request
	var err error

	r, err = http.NewRequest("GET", fmt.Sprintf("%s%s", url, "/domains?page=5&pagesize=20"), nil)
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
		utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method GET and URI /domains",
			http.StatusCreated, response.StatusCode),
			errors.New(string(responseContent)),
		)
	}

	var domainsResponse protocol.DomainsResponse
	if err := json.Unmarshal(responseContent, &domainsResponse); err != nil {
		utils.Fatalln("Error parsing response JSON", err)
	}

	if domainsResponse.PageSize != 20 {
		utils.Fatalln("Default page size isn't 20", nil)
	}

	if domainsResponse.Page != 5 {
		utils.Fatalln("Default page isn't the first one", nil)
	}

	if domainsResponse.NumberOfItems != 100 {
		utils.Fatalln("Not calculating the correct number of domains", nil)
	}

	if domainsResponse.NumberOfPages != 5 {
		utils.Fatalln("Not calculating the correct number of pages", nil)
	}

	if len(domainsResponse.Domains) != domainsResponse.PageSize {
		utils.Fatalln("Not returning all desired domains", nil)
	}
}

func removeDomains() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	var r *http.Request
	var err error

	for i := 0; i < 100; i++ {
		uri := fmt.Sprintf("/domain/example%d.com.br.", i)

		r, err = http.NewRequest("DELETE", fmt.Sprintf("%s%s", url, uri), nil)

		if err != nil {
			utils.Fatalln("Error creating the HTTP request", err)
		}

		utils.BuildHTTPHeader(r, nil)

		response, err := client.Do(r)
		if err != nil {
			utils.Fatalln("Error sending request", err)
		}

		if response.StatusCode != http.StatusNoContent {
			responseContent, err := ioutil.ReadAll(response.Body)

			if err == nil {
				utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method DELETE and URI %s",
					http.StatusCreated, response.StatusCode, uri),
					errors.New(string(responseContent)),
				)
			} else {
				utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method DELETE and URI %s",
					http.StatusCreated, response.StatusCode, uri), nil)
			}
		}
	}
}

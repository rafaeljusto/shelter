// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"github.com/rafaeljusto/shelter/net/http/rest/interceptor"
	"github.com/rafaeljusto/shelter/net/http/rest/messages"
	"github.com/rafaeljusto/shelter/testing/utils"
)

var (
	configFilePath string // Path for the config file with the connection information
)

func init() {
	utils.TestName = "RESTMux"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for RESTMux test")
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

	finishRESTServer := utils.StartRESTServer()
	defer finishRESTServer()

	checkHeaders()
	createDomain()
	checkWrongACL()
	retrieveDomainMetadata()
	retrieveDomainsMetadata()
	deleteDomain()

	utils.Println("SUCCESS!")
}

func checkHeaders() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.RESTServer.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.RESTServer.Listeners[0].IP,
			config.ShelterConfig.RESTServer.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	r, err := http.NewRequest("GET", fmt.Sprintf("%s/xxx/example.com.br.", url), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	if response.StatusCode != http.StatusNotFound {
		utils.Fatalln("Not returning HTTP Not Found when the URI is not registered", nil)
	}

	data := []struct {
		Header             string
		HeaderValue        string
		ExpectedHTTPStatus int
	}{
		{
			Header:             "Accept-Language",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusNotAcceptable,
		},
		{
			Header:             "Accept",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusNotAcceptable,
		},
		{
			Header:             "Accept-Charset",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusNotAcceptable,
		},
		{
			Header:             "Content-Type",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Content-Type",
			HeaderValue:        "",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Content-MD5",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Content-MD5",
			HeaderValue:        "",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Date",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Date",
			HeaderValue:        "",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Authorization",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Authorization",
			HeaderValue:        fmt.Sprintf("%s %d:%s", check.SupportedNamespace, 999, "0PN5J17HBGZHT7JJ3X82"),
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Authorization",
			HeaderValue:        fmt.Sprintf("%s %d:%s", check.SupportedNamespace, 1, "0PN5J17HBGZHT7JJ3X82"),
			ExpectedHTTPStatus: http.StatusUnauthorized,
		},
	}

	content := []byte(`{
      "Nameservers": [
        { "host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "host": "ns2.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        { "email": "admin@example.com.br", "language": "pt-br" }
      ]
    }`)

	for _, item := range data {
		r, err = http.NewRequest("PUT", fmt.Sprintf("%s/domain/example.com.br.", url),
			bytes.NewReader(content))
		if err != nil {
			utils.Fatalln("Error creating the HTTP request", err)
		}

		utils.BuildHTTPHeader(r, content)
		r.Header.Set(item.Header, item.HeaderValue)

		response, err := client.Do(r)
		if err != nil {
			utils.Fatalln("Error sending request", err)
		}

		if response.StatusCode != item.ExpectedHTTPStatus {
			utils.Fatalln(fmt.Sprintf("For HTTP header %s with the value '%s' we were "+
				"expecting the HTTP status %d but got %d",
				item.Header, item.HeaderValue, item.ExpectedHTTPStatus, response.StatusCode), nil)
		}
	}
}

func createDomain() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.RESTServer.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.RESTServer.Listeners[0].IP,
			config.ShelterConfig.RESTServer.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	content := []byte(`{
      "Nameservers": [
        { "host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "host": "ns2.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        { "email": "admin@example.com.br", "language": "en-us" }
      ]
    }`)

	r, err := http.NewRequest("PUT", fmt.Sprintf("%s/domain/example.com.br.", url),
		bytes.NewReader(content))
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	utils.BuildHTTPHeader(r, content)

	messages.ShelterRESTLanguagePacks = messages.LanguagePacks{
		Default: "en-us",
		Packs: []messages.LanguagePack{
			{
				GenericName:  "en",
				SpecificName: "en-us",
			},
			{
				GenericName:  "pt",
				SpecificName: "pt-br",
			},
		},
	}

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	if response.StatusCode != http.StatusCreated {
		utils.Fatalln(fmt.Sprintf("Error creting a domain object. Expected HTTP status %d and got %d",
			http.StatusCreated, response.StatusCode), nil)
	}

	if len(response.Header.Get("Accept")) == 0 {
		utils.Fatalln("Not setting Accept HTTP header in response", nil)
	}

	if len(response.Header.Get("Accept-Language")) == 0 {
		utils.Fatalln("Not setting Accept-Language HTTP header in response", nil)
	}

	if len(response.Header.Get("Accept-Charset")) == 0 {
		utils.Fatalln("Not setting Accept-Charset HTTP header in response", nil)
	}

	if len(response.Header.Get("Date")) == 0 {
		utils.Fatalln("Not setting Date HTTP header in response", nil)
	}
}

func checkWrongACL() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.RESTServer.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.RESTServer.Listeners[0].IP,
			config.ShelterConfig.RESTServer.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	r, err := http.NewRequest("GET", fmt.Sprintf("%s/domain/example.com.br.", url), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	r.RemoteAddr = "127.0.0.1:1234"

	utils.BuildHTTPHeader(r, nil)

	_, cidr, err := net.ParseCIDR("10.0.0.0/8")
	if err != nil {
		utils.Fatalln("Error parsing CIDR", err)
	}
	interceptor.ACL = []*net.IPNet{cidr}

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	if response.StatusCode != http.StatusForbidden {
		utils.Fatalln("Not checking ACL", nil)
	}

	_, cidr, err = net.ParseCIDR("127.0.0.0/8")
	if err != nil {
		utils.Fatalln("Error parsing CIDR", err)
	}
	interceptor.ACL = []*net.IPNet{cidr}

	response, err = client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	if response.StatusCode != http.StatusOK {
		utils.Fatalln("Not allowing a valid ACL", nil)
	}
}

func retrieveDomainMetadata() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.RESTServer.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.RESTServer.Listeners[0].IP,
			config.ShelterConfig.RESTServer.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	r, err := http.NewRequest("HEAD", fmt.Sprintf("%s/domain/example.com.br.", url), nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}
	r.RemoteAddr = "127.0.0.1:1234"

	utils.BuildHTTPHeader(r, nil)

	_, cidr, err := net.ParseCIDR("127.0.0.0/8")
	if err != nil {
		utils.Fatalln("Error parsing CIDR", err)
	}
	interceptor.ACL = []*net.IPNet{cidr}

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	if response.StatusCode != http.StatusOK {
		utils.Fatalln("Not retrieving a valid domain", nil)
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Fatalln("Error reading response body", err)
	}

	if len(responseContent) > 0 {
		utils.Fatalln("HEAD method should not return body", nil)
	}
}

func retrieveDomainsMetadata() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.RESTServer.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.RESTServer.Listeners[0].IP,
			config.ShelterConfig.RESTServer.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	r, err := http.NewRequest("HEAD", fmt.Sprintf("%s/domains", url), nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}
	r.RemoteAddr = "127.0.0.1:1234"

	utils.BuildHTTPHeader(r, nil)

	_, cidr, err := net.ParseCIDR("127.0.0.0/8")
	if err != nil {
		utils.Fatalln("Error parsing CIDR", err)
	}
	interceptor.ACL = []*net.IPNet{cidr}

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	if response.StatusCode != http.StatusOK {
		utils.Fatalln("Not retrieving valid domains", nil)
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Fatalln("Error reading response body", err)
	}

	if len(responseContent) > 0 {
		utils.Fatalln("HEAD method should not return body", nil)
	}
}

func deleteDomain() {
	var client http.Client

	url := ""
	if len(config.ShelterConfig.RESTServer.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.RESTServer.Listeners[0].IP,
			config.ShelterConfig.RESTServer.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	r, err := http.NewRequest("DELETE", fmt.Sprintf("%s/domain/example.com.br.", url), nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	utils.BuildHTTPHeader(r, nil)

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	if response.StatusCode != http.StatusNoContent {
		utils.Fatalln("Error removing a domain object", nil)
	}
}

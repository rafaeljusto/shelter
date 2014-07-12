// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/net/http/rest/handler"
	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
	"github.com/rafaeljusto/shelter/testing/utils"
	"github.com/trajber/handy"
	"io/ioutil"
	"labix.org/v2/mgo"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var (
	configFilePath string // Path for the config file with the connection information
)

type DomainsCacheTestData struct {
	HeaderValue        string
	ExpectedHTTPStatus int
}

func init() {
	utils.TestName = "RESTHandlerDomains"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for RESTHandlerDomains test")
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

	createDomains(database)
	retrieveDomains(database)
	retrieveDomainsMetadata(database)
	retrieveDomainsIfModifiedSince(database)
	retrieveDomainsIfUnmodifiedSince(database)
	retrieveDomainsIfMatch(database)
	retrieveDomainsIfNoneMatch(database)
	deleteDomains(database)

	utils.Println("SUCCESS!")
}

func createDomains(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainHandler)
	mux.Handle("/domain/{fqdn}", func() handy.Handler {
		return h
	})

	for i := 0; i < 100; i++ {
		requestContent := `{
      "Nameservers": [
        { "host": "ns1.example.com.br." },
        { "host": "ns2.example.com.br." }
      ],
      "Owners": [
        { "email": "admin@example.com.br.", "language": "pt-br" }
      ]
    }`

		r, err := http.NewRequest("PUT", fmt.Sprintf("/domain/example%d.com.br.", i),
			strings.NewReader(requestContent))
		if err != nil {
			utils.Fatalln("Error creating the HTTP request", err)
		}
		utils.BuildHTTPHeader(r, []byte(requestContent))

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)

		responseContent, err := ioutil.ReadAll(w.Body)
		if err != nil {
			utils.Fatalln("Error reading response body", err)
		}

		if w.Code != http.StatusCreated {
			utils.Fatalln("Error creating domain", errors.New(string(responseContent)))
		}
	}
}

func retrieveDomains(database *mgo.Database) {
	data := []struct {
		URI                string
		ExpectedHTTPStatus int
		ContentCheck       func(*protocol.DomainsResponse)
	}{
		{
			URI:                "/domains/?orderby=xxx:desc&pagesize=10&page=1",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			URI:                "/domains/?orderby=fqdn:xxx&pagesize=10&page=1",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			URI:                "/domains/?orderby=fqdn:desc&pagesize=xxx&page=1",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			URI:                "/domains/?orderby=fqdn:desc&pagesize=10&page=xxx",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			URI:                "/domains/?orderby=fqdn:desc&pagesize=10&page=1",
			ExpectedHTTPStatus: http.StatusOK,
			ContentCheck: func(domainsResponse *protocol.DomainsResponse) {
				if len(domainsResponse.Domains) != 10 {
					utils.Fatalln("Error retrieving the wrong number of domains", nil)
				}
			},
		},
		{
			URI:                "/domains",
			ExpectedHTTPStatus: http.StatusOK,
			ContentCheck: func(domainsResponse *protocol.DomainsResponse) {
				if len(domainsResponse.Domains) == 0 {
					utils.Fatalln("Error retrieving domains", nil)
				}

				if len(domainsResponse.Domains[0].Nameservers) == 0 {
					utils.Fatalln("Error retrieving domains, no nameservers", nil)
				}

				if len(domainsResponse.Domains[0].Nameservers[0].Host) > 0 {
					utils.Fatalln("Not compressing domains result", nil)
				}
			},
		},
		{
			URI:                "/domains?expand",
			ExpectedHTTPStatus: http.StatusOK,
			ContentCheck: func(domainsResponse *protocol.DomainsResponse) {
				if len(domainsResponse.Domains) == 0 {
					utils.Fatalln("Error retrieving domains", nil)
				}

				if len(domainsResponse.Domains[0].Nameservers) == 0 {
					utils.Fatalln("Error retrieving domains, no nameservers", nil)
				}

				if len(domainsResponse.Domains[0].Nameservers[0].Host) == 0 {
					utils.Fatalln("Compressing domains result when it's not necessary", nil)
				}
			},
		},
	}

	mux := handy.NewHandy()

	h := new(handler.DomainsHandler)
	mux.Handle("/domains", func() handy.Handler {
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
			utils.Fatalln(fmt.Sprintf("Error when requesting domains using the URI [%s]. "+
				"Expected HTTP status code %d but got %d", item.URI,
				item.ExpectedHTTPStatus, w.Code), errors.New(string(responseContent)))
		}

		if item.ContentCheck != nil {
			item.ContentCheck(h.Response)
		}
	}
}

func retrieveDomainsMetadata(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainsHandler)
	mux.Handle("/domains", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/domains/?orderby=fqdn:desc&pagesize=10&page=1", nil)
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
		utils.Fatalln("Error retrieving domains", errors.New(string(responseContent)))
	}

	r, err = http.NewRequest("HEAD", "/domains/?orderby=fqdn:desc&pagesize=10&page=1", nil)
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
		utils.Fatalln("Error retrieving domains", errors.New(string(responseContent)))
	}

	if !utils.CompareProtocolDomains(response1, response2) {
		utils.Fatalln("At this point the GET and HEAD method should "+
			"return the same body content", nil)
	}
}

func retrieveDomainsIfModifiedSince(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/domains", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainsCacheTest(database, r, "If-Modified-Since", []DomainsCacheTestData{
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

func retrieveDomainsIfUnmodifiedSince(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/domains", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainsCacheTest(database, r, "If-Unmodified-Since", []DomainsCacheTestData{
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

func retrieveDomainsIfMatch(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainsHandler)
	mux.Handle("/domains", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/domains", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, nil)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	domainsCacheTest(database, r, "If-Match", []DomainsCacheTestData{
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

func retrieveDomainsIfNoneMatch(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainsHandler)
	mux.Handle("/domains", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/domains", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, nil)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	domainsCacheTest(database, r, "If-None-Match", []DomainsCacheTestData{
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

func deleteDomains(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainHandler)
	mux.Handle("/domain/{fqdn}", func() handy.Handler {
		return h
	})

	for i := 0; i < 100; i++ {
		r, err := http.NewRequest("DELETE", fmt.Sprintf("/domain/example%d.com.br.", i), nil)
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

		if w.Code != http.StatusNoContent {
			utils.Fatalln("Error removing domain", errors.New(string(responseContent)))
		}
	}
}

func domainsCacheTest(database *mgo.Database, r *http.Request,
	header string, domainsCacheTestData []DomainsCacheTestData) {

	mux := handy.NewHandy()

	h := new(handler.DomainsHandler)
	mux.Handle("/domains", func() handy.Handler {
		return h
	})

	for _, item := range domainsCacheTestData {
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

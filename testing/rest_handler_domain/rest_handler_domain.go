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
	"github.com/rafaeljusto/shelter/testing/utils"
	"github.com/trajber/handy"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var (
	configFilePath string // Path for the config file with the connection information
)

type DomainCacheTestData struct {
	HeaderValue        string
	ExpectedHTTPStatus int
}

func init() {
	utils.TestName = "RESTHandlerDomain"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for RESTHandlerDomain test")
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

	invalidFQDN(database)
	createDomain(database)
	updateDomain(database)
	retrieveDomain(database)
	retrieveDomainMetadata(database)
	retrieveDomainIfModifiedSince(database)
	retrieveDomainIfUnmodifiedSince(database)
	retrieveDomainIfMatch(database)
	retrieveDomainIfNoneMatch(database)
	updateDomainIfModifiedSince(database)
	updateDomainIfUnmodifiedSince(database)
	updateDomainIfMatch(database)
	updateDomainIfNoneMatch(database)
	deleteDomainIfModifiedSince(database)
	deleteDomainIfUnmodifiedSince(database)
	deleteDomainIfMatch(database)
	deleteDomainIfNoneMatch(database)
	deleteDomain(database)
	retrieveUnknownDomain(database)

	utils.Println("SUCCESS!")
}

func invalidFQDN(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainHandler)
	mux.Handle("/domain/{fqdn}", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/domain/!!!", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, nil)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		utils.Fatalln(fmt.Sprintf("Not verifying if FQDN exists in the URI. "+
			"Expected status %d and got %d", http.StatusBadRequest, w.Code), nil)
	}
}

func createDomain(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainHandler)
	mux.Handle("/domain/{fqdn}", func() handy.Handler {
		return h
	})

	requestContent := `{
      "Nameservers": [
        { "host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "host": "ns2.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        { "email": "admin@example.com.br.", "language": "en-us" }
      ]
    }`

	r, err := http.NewRequest("PUT", "/domain/example.com.br.", strings.NewReader(requestContent))
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

	if w.Header().Get("ETag") != "1" {
		utils.Fatalln("Not setting ETag in domain creation response", nil)
	}

	if len(w.Header().Get("Last-Modified")) == 0 {
		utils.Fatalln("Not setting Last-Modified in domain creation response", nil)
	}

	if len(w.Header().Get("Location")) == 0 {
		utils.Fatalln("Not setting Location in domain creation response", nil)
	}
}

func updateDomain(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainHandler)
	mux.Handle("/domain/{fqdn}", func() handy.Handler {
		return h
	})

	requestContent := `{
      "Nameservers": [
        { "host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "host": "ns3.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        { "email": "administrator@example.com.br.", "language": "pt-br" }
      ]
    }`

	r, err := http.NewRequest("PUT", "/domain/example.com.br.", strings.NewReader(requestContent))
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

	if w.Code != http.StatusNoContent {
		utils.Fatalln(fmt.Sprintf("Error updating domain. "+
			"Expecting %d and got %d", http.StatusNoContent, w.Code),
			errors.New(string(responseContent)))
	}

	if w.Header().Get("ETag") != "2" {
		utils.Fatalln("Not setting ETag in domain update response", nil)
	}

	if len(w.Header().Get("Last-Modified")) == 0 {
		utils.Fatalln("Not setting Last-Modified in domain update response", nil)
	}

	if len(w.Header().Get("Location")) > 0 {
		utils.Fatalln("Setting Location in domain update response", nil)
	}
}

func retrieveDomain(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainHandler)
	mux.Handle("/domain/{fqdn}", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
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
		utils.Fatalln("Error retrieving domain", errors.New(string(responseContent)))
	}

	if w.Header().Get("ETag") != "2" {
		utils.Fatalln("Not setting ETag in domain retrieve response", nil)
	}

	if len(w.Header().Get("Last-Modified")) == 0 {
		utils.Fatalln("Not setting Last-Modified in domain retrieve response", nil)
	}

	if h.Response.FQDN != "example.com.br." {
		utils.Fatalln("Domain's FQDN was not persisted correctly", nil)
	}

	if len(h.Response.Nameservers) != 2 ||
		h.Response.Nameservers[0].Host != "ns1.example.com.br." ||
		h.Response.Nameservers[0].IPv4 != "127.0.0.1" ||
		h.Response.Nameservers[1].Host != "ns3.example.com.br." ||
		h.Response.Nameservers[1].IPv6 != "::1" {
		utils.Fatalln("Domain's nameservers were not persisted correctly", nil)
	}

	if len(h.Response.Owners) != 1 ||
		h.Response.Owners[0].Email != "administrator@example.com.br." {

		utils.Fatalln("Domain's owners were not persisted correctly", nil)
	}
}

func retrieveDomainMetadata(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainHandler)
	mux.Handle("/domain/{fqdn}", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, nil)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	response1 := *h.Response
	responseContent, err := ioutil.ReadAll(w.Body)
	if err != nil {
		utils.Fatalln("Error reading response body from GET", err)
	}

	if w.Code != http.StatusOK {
		utils.Fatalln("Error retrieving domain", errors.New(string(responseContent)))
	}

	r, err = http.NewRequest("HEAD", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, nil)

	mux.ServeHTTP(w, r)

	response2 := *h.Response
	responseContent, err = ioutil.ReadAll(w.Body)
	if err != nil {
		utils.Fatalln("Error reading response body from HEAD", err)
	}

	if w.Code != http.StatusOK {
		utils.Fatalln("Error retrieving domain", errors.New(string(responseContent)))
	}

	if !utils.CompareProtocolDomain(response1, response2) {
		utils.Fatalln("At this point the GET and HEAD method should "+
			"return the same body content", nil)
	}
}

func retrieveDomainIfModifiedSince(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, "", "If-Modified-Since", []DomainCacheTestData{
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

func retrieveDomainIfUnmodifiedSince(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, "", "If-Unmodified-Since", []DomainCacheTestData{
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

func retrieveDomainIfMatch(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, "", "If-Match", []DomainCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
		{
			HeaderValue:        "3",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
		{
			HeaderValue:        "2",
			ExpectedHTTPStatus: http.StatusOK,
		},
	})
}

func retrieveDomainIfNoneMatch(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, "", "If-None-Match", []DomainCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusOK,
		},
		{
			HeaderValue:        "2",
			ExpectedHTTPStatus: http.StatusNotModified,
		},
	})
}

func updateDomainIfModifiedSince(database *mgo.Database) {
	requestContent := `{
      "Nameservers": [
        { "host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "host": "ns3.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        { "email": "administrator@example.com.br.", "language": "pt-br" }
      ]
    }`

	r, err := http.NewRequest("PUT", "/domain/example.com.br.", strings.NewReader(requestContent))
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, requestContent, "If-Modified-Since", []DomainCacheTestData{
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

func updateDomainIfUnmodifiedSince(database *mgo.Database) {
	requestContent := `{
      "Nameservers": [
        { "host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "host": "ns3.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        { "email": "administrator@example.com.br.", "language": "pt-br" }
      ]
    }`

	r, err := http.NewRequest("PUT", "/domain/example.com.br.", strings.NewReader(requestContent))
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, requestContent, "If-Unmodified-Since", []DomainCacheTestData{
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

func updateDomainIfMatch(database *mgo.Database) {
	requestContent := `{
      "Nameservers": [
        { "host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "host": "ns3.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        { "email": "administrator@example.com.br.", "language": "en-us" }
      ]
    }`

	r, err := http.NewRequest("PUT", "/domain/example.com.br.", strings.NewReader(requestContent))
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, requestContent, "If-Match", []DomainCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
		{
			HeaderValue:        "3",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
		{
			HeaderValue:        "2",
			ExpectedHTTPStatus: http.StatusNoContent,
		},
	})
}

func updateDomainIfNoneMatch(database *mgo.Database) {
	requestContent := `{
      "Nameservers": [
        { "host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "host": "ns3.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        { "email": "administrator@example.com.br.", "language": "en-us" }
      ]
    }`

	r, err := http.NewRequest("PUT", "/domain/example.com.br.", strings.NewReader(requestContent))
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, requestContent, "If-None-Match", []DomainCacheTestData{
		{
			HeaderValue:        "3",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusNoContent,
		},
	})
}

func deleteDomainIfModifiedSince(database *mgo.Database) {
	r, err := http.NewRequest("DELETE", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, "", "If-Modified-Since", []DomainCacheTestData{
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

func deleteDomainIfUnmodifiedSince(database *mgo.Database) {
	r, err := http.NewRequest("DELETE", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, "", "If-Unmodified-Since", []DomainCacheTestData{
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

func deleteDomainIfMatch(database *mgo.Database) {
	r, err := http.NewRequest("DELETE", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, "", "If-Match", []DomainCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
		{
			HeaderValue:        "2",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
	})
}

func deleteDomainIfNoneMatch(database *mgo.Database) {
	r, err := http.NewRequest("DELETE", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	domainCacheTest(database, r, "", "If-None-Match", []DomainCacheTestData{
		{
			HeaderValue:        "4",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
	})
}

func deleteDomain(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainHandler)
	mux.Handle("/domain/{fqdn}", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("DELETE", "/domain/example.com.br.", nil)
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
		utils.Fatalln("Error deleting domain", errors.New(string(responseContent)))
	}
}

func retrieveUnknownDomain(database *mgo.Database) {
	mux := handy.NewHandy()

	h := new(handler.DomainHandler)
	mux.Handle("/domain/{fqdn}", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
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
		utils.Fatalln("Error retrieving unknown domain", errors.New(string(responseContent)))
	}
}

func domainCacheTest(database *mgo.Database, r *http.Request,
	requestContent, header string, domainCacheTestData []DomainCacheTestData) {

	mux := handy.NewHandy()

	h := new(handler.DomainHandler)
	mux.Handle("/domain/{fqdn}", func() handy.Handler {
		return h
	})

	for _, item := range domainCacheTestData {
		r.Header.Set(header, item.HeaderValue)

		if len(requestContent) > 0 {
			utils.BuildHTTPHeader(r, []byte(requestContent))
		} else {
			utils.BuildHTTPHeader(r, nil)
		}

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

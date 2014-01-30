package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"github.com/rafaeljusto/shelter/net/http/rest/handler"
	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
	"github.com/rafaeljusto/shelter/testing/utils"
	"labix.org/v2/mgo"
	"net/http"
	"strings"
	"time"
)

var (
	configFilePath string // Path for the config file with the connection information
)

// RESTHandlerDomainTestConfigFile is a structure to store the test configuration file data
type RESTHandlerDomainTestConfigFile struct {
	Database struct {
		URI  string
		Name string
	}
}

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

	var config RESTHandlerDomainTestConfigFile
	err := utils.ReadConfigFile(configFilePath, &config)

	if err == utils.ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		utils.Fatalln("Error reading configuration file", err)
	}

	database, databaseSession, err := mongodb.Open(config.Database.URI, config.Database.Name)
	if err != nil {
		utils.Fatalln("Error connecting the database", err)
	}
	defer databaseSession.Close()

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	// If there was some problem in the last test, there could be some data in the
	// database, so let's clear it to don't affect this test. We avoid checking the error,
	// because if the collection does not exist yet, it will be created in the first
	// insert
	domainDAO.RemoveAll()

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
	r, err := http.NewRequest("GET", "/domain/", nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleDomain(r, &context)

	if context.ResponseHTTPStatus != http.StatusBadRequest {
		utils.Fatalln("Not verifying if FQDN exists in the URI", nil)
	}
}

func createDomain(database *mgo.Database) {
	r, err := http.NewRequest("PUT", "/domain/example.com.br.",
		strings.NewReader(`{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns2.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        "admin@example.com.br."
      ]
    }`))
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleDomain(r, &context)

	if context.ResponseHTTPStatus != http.StatusCreated {
		utils.Fatalln("Error creating domain",
			errors.New(string(context.ResponseContent)))
	}

	if context.HTTPHeader["ETag"] != "1" {
		utils.Fatalln("Not setting ETag in domain creation response", nil)
	}

	if len(context.HTTPHeader["Last-Modified"]) == 0 {
		utils.Fatalln("Not setting Last-Modified in domain creation response", nil)
	}

	if len(context.HTTPHeader["Location"]) == 0 {
		utils.Fatalln("Not setting Location in domain creation response", nil)
	}
}

func updateDomain(database *mgo.Database) {
	r, err := http.NewRequest("PUT", "/domain/example.com.br.",
		strings.NewReader(`{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns3.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        "administrator@example.com.br."
      ]
    }`))
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleDomain(r, &context)

	if context.ResponseHTTPStatus != http.StatusNoContent {
		utils.Fatalln("Error updating domain",
			errors.New(string(context.ResponseContent)))
	}

	if context.HTTPHeader["ETag"] != "2" {
		utils.Fatalln("Not setting ETag in domain update response", nil)
	}

	if len(context.HTTPHeader["Last-Modified"]) == 0 {
		utils.Fatalln("Not setting Last-Modified in domain update response", nil)
	}

	if len(context.HTTPHeader["Location"]) > 0 {
		utils.Fatalln("Setting Location in domain update response", nil)
	}
}

func retrieveDomain(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleDomain(r, &context)

	if context.ResponseHTTPStatus != http.StatusOK {
		utils.Fatalln("Error retrieving domain",
			errors.New(string(context.ResponseContent)))
	}

	if context.HTTPHeader["ETag"] != "2" {
		utils.Fatalln("Not setting ETag in domain retrieve response", nil)
	}

	if len(context.HTTPHeader["Last-Modified"]) == 0 {
		utils.Fatalln("Not setting Last-Modified in domain retrieve response", nil)
	}

	var domainResponse protocol.DomainResponse
	json.Unmarshal(context.ResponseContent, &domainResponse)

	if domainResponse.FQDN != "example.com.br." {
		utils.Fatalln("Domain's FQDN was not persisted correctly", nil)
	}

	if len(domainResponse.Nameservers) != 2 ||
		domainResponse.Nameservers[0].Host != "ns1.example.com.br." ||
		domainResponse.Nameservers[0].IPv4 != "127.0.0.1" ||
		domainResponse.Nameservers[1].Host != "ns3.example.com.br." ||
		domainResponse.Nameservers[1].IPv6 != "::1" {
		utils.Fatalln("Domain's nameservers were not persisted correctly", nil)
	}

	if len(domainResponse.Owners) != 1 ||
		domainResponse.Owners[0] != "administrator@example.com.br." {

		utils.Fatalln("Domain's owners were not persisted correctly", nil)
	}
}

func retrieveDomainMetadata(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	context1, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleDomain(r, &context1)

	if context1.ResponseHTTPStatus != http.StatusOK {
		utils.Fatalln("Error retrieving domain",
			errors.New(string(context1.ResponseContent)))
	}

	r, err = http.NewRequest("HEAD", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	context2, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleDomain(r, &context2)

	if context2.ResponseHTTPStatus != http.StatusOK {
		utils.Fatalln("Error retrieving domain",
			errors.New(string(context2.ResponseContent)))
	}

	if string(context1.ResponseContent) != string(context2.ResponseContent) {
		utils.Fatalln("At this point the GET and HEAD method should "+
			"return the same body content", nil)
	}
}

func retrieveDomainIfModifiedSince(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-Modified-Since", []DomainCacheTestData{
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
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-Unmodified-Since", []DomainCacheTestData{
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
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-Match", []DomainCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			HeaderValue:        "3",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
	})
}

func retrieveDomainIfNoneMatch(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-None-Match", []DomainCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			HeaderValue:        "2",
			ExpectedHTTPStatus: http.StatusNotModified,
		},
	})
}

func updateDomainIfModifiedSince(database *mgo.Database) {
	r, err := http.NewRequest("PUT", "/domain/example.com.br.",
		strings.NewReader(`{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns3.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        "administrator@example.com.br."
      ]
    }`))
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-Modified-Since", []DomainCacheTestData{
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
	r, err := http.NewRequest("PUT", "/domain/example.com.br.",
		strings.NewReader(`{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns3.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        "administrator@example.com.br."
      ]
    }`))
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-Unmodified-Since", []DomainCacheTestData{
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
	r, err := http.NewRequest("PUT", "/domain/example.com.br.",
		strings.NewReader(`{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns3.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        "administrator@example.com.br."
      ]
    }`))
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-Match", []DomainCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			HeaderValue:        "3",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
	})
}

func updateDomainIfNoneMatch(database *mgo.Database) {
	r, err := http.NewRequest("PUT", "/domain/example.com.br.",
		strings.NewReader(`{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns3.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        "administrator@example.com.br."
      ]
    }`))
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-None-Match", []DomainCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			HeaderValue:        "2",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
	})
}

func deleteDomainIfModifiedSince(database *mgo.Database) {
	r, err := http.NewRequest("DELETE", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-Modified-Since", []DomainCacheTestData{
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
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-Unmodified-Since", []DomainCacheTestData{
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
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-Match", []DomainCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			HeaderValue:        "3",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
	})
}

func deleteDomainIfNoneMatch(database *mgo.Database) {
	r, err := http.NewRequest("DELETE", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	domainCacheTest(database, r, "If-None-Match", []DomainCacheTestData{
		{
			HeaderValue:        "abcdef",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			HeaderValue:        "2",
			ExpectedHTTPStatus: http.StatusPreconditionFailed,
		},
	})
}

func deleteDomain(database *mgo.Database) {
	r, err := http.NewRequest("DELETE", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleDomain(r, &context)

	if context.ResponseHTTPStatus != http.StatusNoContent {
		utils.Fatalln("Error deleting domain",
			errors.New(string(context.ResponseContent)))
	}
}

func retrieveUnknownDomain(database *mgo.Database) {
	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
	}

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	handler.HandleDomain(r, &context)

	if context.ResponseHTTPStatus != http.StatusNotFound {
		utils.Fatalln("Error retrieving unknown domain",
			errors.New(string(context.ResponseContent)))
	}
}

func domainCacheTest(database *mgo.Database, r *http.Request,
	header string, domainCacheTestData []DomainCacheTestData) {

	context, err := context.NewContext(r, database)
	if err != nil {
		utils.Fatalln("Error creating context", err)
	}

	for _, item := range domainCacheTestData {
		r.Header.Set(header, item.HeaderValue)

		handler.HandleDomain(r, &context)

		if context.ResponseHTTPStatus != item.ExpectedHTTPStatus {
			utils.Fatalln(fmt.Sprintf("Error in %s test using %s [%s] "+
				"HTTP header. Expected HTTP status code %d and got %d",
				r.Method, header, item.HeaderValue, item.ExpectedHTTPStatus,
				context.ResponseHTTPStatus), errors.New(string(context.ResponseContent)))
		}
	}
}

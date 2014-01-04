package main

import (
	"errors"
	"flag"
	"fmt"
	"labix.org/v2/mgo"
	"net/http"
	"shelter/dao"
	"shelter/database/mongodb"
	"shelter/net/http/rest/context"
	"shelter/net/http/rest/handler"
	"shelter/testing/utils"
	"strings"
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

	database, err := mongodb.Open(config.Database.URI, config.Database.Name)
	if err != nil {
		utils.Fatalln("Error connecting the database", err)
	}

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	// If there was some problem in the last test, there could be some data in the
	// database, so let's clear it to don't affect this test. We avoid checking the error,
	// because if the collection does not exist yet, it will be created in the first
	// insert
	domainDAO.RemoveAll()

	createDomain(database)
	updateDomain(database)
	retrieveDomain(database)
	deleteDomain(database)

	utils.Println("SUCCESS!")
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

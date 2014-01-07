package main

import (
	"flag"
	"fmt"
	"labix.org/v2/mgo"
	"shelter/dao"
	"shelter/database/mongodb"
	"shelter/testing/utils"
)

var (
	configFilePath string // Path for the config file with the connection information
)

// RESTHandlerDomainsTestConfigFile is a structure to store the test configuration file data
type RESTHandlerDomainsTestConfigFile struct {
	Database struct {
		URI  string
		Name string
	}
}

func init() {
	utils.TestName = "RESTHandlerDomains"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for RESTHandlerDomains test")
}

func main() {
	flag.Parse()

	var config RESTHandlerDomainsTestConfigFile
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

	createDomains(database)
	retrieveDomains(database)
	retrieveDomainsMetadata(database)

	utils.Println("SUCCESS!")
}

func createDomains(database *mgo.Database) {

}

func retrieveDomains(database *mgo.Database) {

}

func retrieveDomainsMetadata(database *mgo.Database) {

}

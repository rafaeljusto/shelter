package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"shelter/config"
	"shelter/dao"
	"shelter/database/mongodb"
	"shelter/model"
	"shelter/service"
	"time"
)

// List of possible errors in this test. There can be also other errors from low level
// structures
var (
	// Config file path is a mandatory parameter
	ErrConfigFileUndefined = errors.New("Config file path undefined")
)

var (
	configFilePath string // Path for the configuration file with all the query parameters
)

func init() {
	flag.StringVar(&configFilePath, "config", "", "Configuration file for ScanQuerier test")
}

func main() {
	flag.Parse()

	if err := readAndLoadConfigFile(); err == ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		fatalln("Error reading configuration file", err)
	}

	database, err := mongodb.Open(config.ShelterConfig.Database.URI,
		config.ShelterConfig.Database.Name)
	if err != nil {
		fatalln("Error connecting the database", err)
	}

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	// Remove all data before starting the test. This is necessary because maybe in the last
	// test there was an error and the data wasn't removed from the database
	domainDAO.RemoveAll()

	domainWithNoErrors(domainDAO)

	println("SUCCESS!")
}

func domainWithNoErrors(domainDAO dao.DomainDAO) {
	domain := newDomain()
	lastCheckAt := time.Now().Add(-72 * time.Hour)
	lastOKAt := lastCheckAt.Add(-24 * time.Hour)

	// Set all nameservers with error and the last check equal of the error check interval,
	// this will force the domain to be checked
	for index, _ := range domain.Nameservers {
		domain.Nameservers[index].LastCheckAt = lastCheckAt
		domain.Nameservers[index].LastOKAt = lastOKAt
		domain.Nameservers[index].LastStatus = model.NameserverStatusServerFailure
	}

	// Set all DS records with error and the last check equal of the error check interval,
	// this will force the domain to be checked
	for index, _ := range domain.DSSet {
		domain.DSSet[index].LastCheckAt = lastCheckAt
		domain.DSSet[index].LastOKAt = lastOKAt
		domain.DSSet[index].LastStatus = model.DSStatusTimeout
	}

	if err := domainDAO.Save(&domain); err != nil {
		fatalln("Error saving domain for scan scenario", err)
	}

	service.ScanDomains()

	domain, err := domainDAO.FindByFQDN(domain.FQDN)
	if err != nil {
		fatalln("Didn't find scanned domain", err)
	}

	for _, nameserver := range domain.Nameservers {
		if nameserver.LastStatus != model.NameserverStatusOK {
			fatalln(fmt.Sprintf("Fail to validate a supposedly well configured nameserver '%s'. Found status: %s",
				nameserver.Host, model.NameserverStatusToString(nameserver.LastStatus)), err)
		}

		if nameserver.LastCheckAt.Before(lastCheckAt) ||
			nameserver.LastCheckAt.Equal(lastCheckAt) {
			fatalln(fmt.Sprintf("Last check date was not updated in nameserver '%s'",
				nameserver.Host), nil)
		}

		if nameserver.LastOKAt.Before(lastOKAt) || nameserver.LastOKAt.Equal(lastOKAt) {
			fatalln(fmt.Sprintf("Last OK date was not updated in nameserver '%s'",
				nameserver.Host), nil)
		}
	}

	for _, ds := range domain.DSSet {
		if ds.LastStatus != model.DSStatusOK {
			fatalln(fmt.Sprintf("Fail to validate a supposedly well configured DS %d. "+
				"Found status: %s", ds.Keytag, model.DSStatusToString(ds.LastStatus)), err)
		}

		if ds.LastCheckAt.Before(lastCheckAt) || ds.LastCheckAt.Equal(lastCheckAt) {
			fatalln(fmt.Sprintf("Last check date was not updated in DS %d",
				ds.Keytag), nil)
		}

		if ds.LastOKAt.Before(lastOKAt) || ds.LastOKAt.Equal(lastOKAt) {
			fatalln(fmt.Sprintf("Last OK date was not updated in DS %d",
				ds.Keytag), nil)
		}
	}
}

// Function to read the configuration file
func readAndLoadConfigFile() error {
	// Config file path is a mandatory program parameter
	if len(configFilePath) == 0 {
		return ErrConfigFileUndefined
	}

	confBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(confBytes, &config.ShelterConfig); err != nil {
		return err
	}

	return nil
}

// Function to mock a domain object
func newDomain() model.Domain {
	var domain model.Domain
	domain.FQDN = "br."

	domain.Nameservers = []model.Nameserver{
		{
			Host: "a.dns.br.",
			IPv4: net.ParseIP("200.160.0.10"),
			IPv6: net.ParseIP("2001:12ff::10"),
		},
		{
			Host: "b.dns.br.",
			IPv4: net.ParseIP("200.189.41.10"),
		},
		{
			Host: "c.dns.br.",
			IPv4: net.ParseIP("200.192.233.10"),
		},
		{
			Host: "e.dns.br.",
			IPv4: net.ParseIP("200.229.248.10"),
			IPv6: net.ParseIP("2001:12f8:1::10"),
		},
		{
			Host: "f.dns.br.",
			IPv4: net.ParseIP("200.219.159.10"),
		},
	}

	// Caution! The .br DNSKEY will change periodically, so this test will fail sometime
	// beucase of this, when this occurs we need to update the DS information for the new
	// .br key
	domain.DSSet = []model.DS{
		{
			Keytag:     41674,
			Algorithm:  model.DSAlgorithmRSASHA1,
			DigestType: model.DSDigestTypeSHA1,
			Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
		},
	}

	owner, _ := mail.ParseAddress("test@rafael.net.br")
	domain.Owners = []*mail.Address{owner}

	return domain
}

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints a simple message without ending the test
func println(message string) {
	message = fmt.Sprintf("Scan integration test: %s", message)
	log.Println(message)
}

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints an error message and ends the test
func fatalln(message string, err error) {
	message = fmt.Sprintf("Scan integration test: %s", message)
	if err != nil {
		message = fmt.Sprintf("%s. Details: %s", message, err.Error())
	}

	log.Fatalln(message)
}

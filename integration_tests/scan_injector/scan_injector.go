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
	"shelter/dao"
	"shelter/database/mongodb"
	"shelter/model"
	"shelter/net/scan"
	"strconv"
	"sync"
	"time"
)

// This test objective is to verify the domain selection rules in the injector scanner
// component. As the injector depends on a database scenario, all this checks are going to
// be made in an integration test enviroment

// List of possible errors in this test. There can be also other errors from low level
// structures
var (
	// Config file path is a mandatory parameter
	ErrConfigFileUndefined = errors.New("Config file path undefined")
)

var (
	configFilePath string // Path for the configuration file with the database connection information
)

// ScanInjectorTestConfigFile is a structure to store the test configuration file data
type ScanInjectorTestConfigFile struct {
	Database struct {
		URI  string
		Name string
	}

	// Indicates if a domain is going to be selected for scan or not, based on the last check,
	// if has errors or if the DNSSEC expiration date is near
	Scan struct {
		DomainsBufferSize int // Size of the channel

		VerificationIntervals struct {
			MaxOKDays              int
			MaxErrorDays           int
			MaxExpirationAlertDays int
		}
	}
}

func init() {
	flag.StringVar(&configFilePath, "config", "", "Configuration file for ScanInjector test")
}

func main() {
	flag.Parse()

	config, err := readConfigFile()
	if err == ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		fatalln("Error reading configuration file", err)
	}

	database, err := mongodb.Open(config.Database.URI, config.Database.Name)
	if err != nil {
		fatalln("Error connecting the database", err)
	}

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	// Remove all data before starting the test. This is necessary because maybe in the last
	// test there was an error and the data wasn't removed from the database
	domainDAO.RemoveAll()

	domainWithDNSErrors(config, domainDAO)
	domainWithDNSSECErrors(config, domainDAO)
	domainWithNoErrors(config, domainDAO)

	println("SUCCESS!")
}

func domainWithDNSErrors(config ScanInjectorTestConfigFile, domainDAO dao.DomainDAO) {
	domain := newDomain()

	// Set all nameservers with error and the last check equal of the error check interval,
	// this will force the domain to be checked
	for index, _ := range domain.Nameservers {
		maxErrorHours := config.Scan.VerificationIntervals.MaxErrorDays * 24
		lessThreeDays, _ := time.ParseDuration("-" + strconv.Itoa(maxErrorHours) + "h")

		domain.Nameservers[index].LastCheckAt = time.Now().Add(lessThreeDays)
		domain.Nameservers[index].LastStatus = model.NameserverStatusServerFailure
	}

	if err := domainDAO.Save(&domain); err != nil {
		fatalln("Error saving domain for scan scenario", err)
	}

	if domains := runScan(config, domainDAO); len(domains) != 1 {
		fatalln(fmt.Sprintf("Couldn't load a domain with DNS errors for scan. "+
			"Expected 1 got %d", len(domains)), nil)
	}

	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		fatalln("Error removing domain", err)
	}
}

func domainWithDNSSECErrors(config ScanInjectorTestConfigFile, domainDAO dao.DomainDAO) {
	domain := newDomain()

	// Set all DS records with error and the last check equal of the error check interval,
	// this will force the domain to be checked
	for index, _ := range domain.DSSet {
		maxErrorHours := config.Scan.VerificationIntervals.MaxErrorDays * 24
		lessThreeDays, _ := time.ParseDuration("-" + strconv.Itoa(maxErrorHours) + "h")

		domain.DSSet[index].LastCheckAt = time.Now().Add(lessThreeDays)
		domain.DSSet[index].LastStatus = model.DSStatusTimeout
	}

	if err := domainDAO.Save(&domain); err != nil {
		fatalln("Error saving domain for scan scenario", err)
	}

	if domains := runScan(config, domainDAO); len(domains) != 1 {
		fatalln(fmt.Sprintf("Couldn't load a domain with DNSSEC errors for scan. "+
			"Expected 1 got %d", len(domains)), nil)
	}

	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		fatalln("Error removing domain", err)
	}
}

func domainWithNoErrors(config ScanInjectorTestConfigFile, domainDAO dao.DomainDAO) {
	domain := newDomain()

	// Set all nameservers as configured correctly and the last check as now, this domain is
	// unlikely to be selected
	for index, _ := range domain.Nameservers {
		domain.Nameservers[index].LastCheckAt = time.Now()
		domain.Nameservers[index].LastStatus = model.NameserverStatusOK
	}

	// Set all DS records as configured correctly and the last check as now, this domain is
	// unlikely to be selected
	for index, _ := range domain.DSSet {
		domain.DSSet[index].LastCheckAt = time.Now()
		domain.DSSet[index].LastStatus = model.DSStatusOK
	}

	if err := domainDAO.Save(&domain); err != nil {
		fatalln("Error saving domain for scan scenario", err)
	}

	if domains := runScan(config, domainDAO); len(domains) > 0 {
		fatalln(fmt.Sprintf("Selected a domain configured correctly for the scan. "+
			"Expected 0 got %d", len(domains)), nil)
	}

	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		fatalln("Error removing domain", err)
	}
}

// Method responsable to configure and start scan injector for tests
func runScan(config ScanInjectorTestConfigFile, domainDAO dao.DomainDAO) []*model.Domain {
	scanInjector := scan.NewInjector(
		domainDAO.Database,
		config.Scan.DomainsBufferSize,
		config.Scan.VerificationIntervals.MaxOKDays,
		config.Scan.VerificationIntervals.MaxErrorDays,
		config.Scan.VerificationIntervals.MaxExpirationAlertDays,
	)

	// Go routines group control created, but not used for this tests, as we are simulating
	// a querier receiver
	var scanGroup sync.WaitGroup

	errorsChannel := make(chan error)
	domainsToQueryChannel := scanInjector.Start(&scanGroup, errorsChannel)

	var domains []*model.Domain

	for {
		exit := false

		select {
		case domain := <-domainsToQueryChannel:
			// Detect the poison pills
			if domain == nil {
				exit = true

			} else {
				domains = append(domains, domain)
			}

		case err := <-errorsChannel:
			fatalln("Error selecting domain", err)
		}

		if exit {
			break
		}
	}

	return domains
}

// Function to read the configuration file
func readConfigFile() (ScanInjectorTestConfigFile, error) {
	var configFile ScanInjectorTestConfigFile

	// Config file path is a mandatory program parameter
	if len(configFilePath) == 0 {
		return configFile, ErrConfigFileUndefined
	}

	confBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return configFile, err
	}

	if err := json.Unmarshal(confBytes, &configFile); err != nil {
		return configFile, err
	}

	return configFile, nil
}

// Function to mock a domain object
func newDomain() model.Domain {
	var domain model.Domain
	domain.FQDN = "rafael.net.br"

	domain.Nameservers = []model.Nameserver{
		{
			Host: "ns1.rafael.net.br",
			IPv4: net.ParseIP("127.0.0.1"),
			IPv6: net.ParseIP("::1"),
		},
		{
			Host: "ns2.rafael.net.br",
			IPv4: net.ParseIP("127.0.0.2"),
		},
	}

	domain.DSSet = []model.DS{
		{
			Keytag:    1234,
			Algorithm: model.DSAlgorithmRSASHA1,
			Digest:    "A790A11EA430A85DA77245F091891F73AA740483",
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
	message = fmt.Sprintf("ScanInjector integration test: %s", message)
	log.Println(message)
}

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints an error message and ends the test
func fatalln(message string, err error) {
	message = fmt.Sprintf("ScanInjector integration test: %s", message)
	if err != nil {
		message = fmt.Sprintf("%s. Details: %s", message, err.Error())
	}

	log.Fatalln(message)
}

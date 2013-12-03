package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"net"
	"net/mail"
	"shelter/dao"
	"shelter/database/mongodb"
	"shelter/model"
	"shelter/net/scan"
	"sync"
)

// List of possible errors in this test. There can be also other errors from low level
// structures
var (
	// Config file path is a mandatory parameter
	ErrConfigFileUndefined = errors.New("Config file path undefined")
)

var (
	configFilePath string // Path for the config file with the connection information
)

// Define some important parameter for the test enviroment
const (
	domainsBufferSize = 100 // Size of the communication channels between the querier and the collector
	saveAtOnce        = 100 // Number of domains to acumulate before saving all togheter
)

// ScanCollectorTestConfigFile is a structure to store the test configuration file data
type ScanCollectorTestConfigFile struct {
	Database struct {
		URI  string
		Name string
	}
}

func init() {
	flag.StringVar(&configFilePath, "config", "", "Configuration file for ScanInjector test")
}

func main() {
	flag.Parse()

	configFile, err := readConfigFile()
	if err == ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		fatalln("Error reading configuration file", err)
	}

	database, err := mongodb.Open(configFile.Database.URI, configFile.Database.Name)
	if err != nil {
		fatalln("Error connecting the database", err)
	}

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	// Remove all data before starting the test. This is necessary because maybe in the last
	// test there was an error and the data wasn't removed from the database
	domainDAO.RemoveAll()

	domainWithErrors(database)
	domainWithNoErrors(database)

	println("SUCCESS!")
}

func domainWithErrors(database *mgo.Database) {
	domainsToSave := make(chan *model.Domain, domainsBufferSize)
	domainsToSave <- &model.Domain{
		FQDN: "br.",
		Nameservers: []model.Nameserver{
			{
				Host:       "ns1.br",
				IPv4:       net.ParseIP("127.0.0.1"),
				LastStatus: model.NameserverStatusTimeout,
			},
		},
		DSSet: []model.DS{
			{
				Keytag:     1234,
				Algorithm:  model.DSAlgorithmRSASHA1NSEC3,
				DigestType: model.DSDigestTypeSHA1,
				Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
				LastStatus: model.DSStatusExpiredSignature,
			},
		},
	}
	domainsToSave <- nil

	runScan(database, domainsToSave)

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	domain, err := domainDAO.FindByFQDN("br.")
	if err != nil {
		fatalln("Error loading domain with problems", err)
	}

	if len(domain.Nameservers) == 0 {
		fatalln("Error saving nameservers", nil)
	}

	if domain.Nameservers[0].LastStatus != model.NameserverStatusTimeout {
		fatalln("Error setting status in the nameserver", nil)
	}

	if len(domain.DSSet) == 0 {
		fatalln("Error saving the DS set", nil)
	}

	if domain.DSSet[0].LastStatus != model.DSStatusExpiredSignature {
		fatalln("Error setting status in the DS", nil)
	}

	if err := domainDAO.RemoveByFQDN("br."); err != nil {
		fatalln("Error removing test domain", err)
	}
}

func domainWithNoErrors(database *mgo.Database) {
	domainsToSave := make(chan *model.Domain, domainsBufferSize)
	domainsToSave <- &model.Domain{
		FQDN: "br.",
		Nameservers: []model.Nameserver{
			{
				Host:       "ns1.br",
				IPv4:       net.ParseIP("127.0.0.1"),
				LastStatus: model.NameserverStatusOK,
			},
		},
		DSSet: []model.DS{
			{
				Keytag:     1234,
				Algorithm:  model.DSAlgorithmRSASHA1NSEC3,
				DigestType: model.DSDigestTypeSHA1,
				Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
				LastStatus: model.DSStatusOK,
			},
		},
	}
	domainsToSave <- nil

	runScan(database, domainsToSave)

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	domain, err := domainDAO.FindByFQDN("br.")
	if err != nil {
		fatalln("Error loading domain with problems", err)
	}

	if len(domain.Nameservers) == 0 {
		fatalln("Error saving nameservers", nil)
	}

	if domain.Nameservers[0].LastStatus != model.NameserverStatusOK {
		fatalln("Error setting status in the nameserver", nil)
	}

	if len(domain.DSSet) == 0 {
		fatalln("Error saving the DS set", nil)
	}

	if domain.DSSet[0].LastStatus != model.DSStatusOK {
		fatalln("Error setting status in the DS", nil)
	}

	if err := domainDAO.RemoveByFQDN("br."); err != nil {
		fatalln("Error removing test domain", err)
	}
}

// Method responsable to configure and start scan injector for tests
func runScan(database *mgo.Database, domainsToSave chan *model.Domain) {
	scanCollector := scan.NewCollector(database, saveAtOnce)

	var scanGroup sync.WaitGroup
	errorsChannel := make(chan error)
	scanCollector.Start(&scanGroup, domainsToSave, errorsChannel)

	go func() {
		select {
		case err := <-errorsChannel:
			fatalln("Error saving domain", err)
		}
	}()

	scanGroup.Wait()
}

// Function to read the configuration file
func readConfigFile() (ScanCollectorTestConfigFile, error) {
	var configFile ScanCollectorTestConfigFile

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
	message = fmt.Sprintf("ScanCollector integration test: %s", message)
	log.Println(message)
}

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints an error message and ends the test
func fatalln(message string, err error) {
	message = fmt.Sprintf("ScanCollector integration test: %s", message)
	if err != nil {
		message = fmt.Sprintf("%s. Details: %s", message, err.Error())
	}

	log.Fatalln(message)
}

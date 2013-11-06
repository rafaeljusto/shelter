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
	"time"
)

// This test objective is to verify the domain data persistence. The strategy is to insert
// and search for the information. Check for insert/update consistency (updates don't
// create a new element) and if the object id is set on creation

// List of possible errors in this test. There can be also other errors from low level
// structures
var (
	// Config file path is a mandatory parameter
	ErrConfigFileUndefined = errors.New("Config file path undefined")
)

var (
	// Path for the configuration file with the database connection information
	configFilePath string
)

// DomainDAOTestConfigFile is a structure to store the test configuration file data
type DomainDAOTestConfigFile struct {
	Database struct {
		URI  string
		Name string
	}
}

func init() {
	flag.StringVar(&configFilePath, "config", "", "Configuration file for DomainDAO test")
}

func main() {
	flag.Parse()

	configFile, err := readConfigFile()
	if err == ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()

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

	// If there was some problem in the last test, there could be some data in the
	// database, so let's clear it to don't affect this test. We avoid checking the error,
	// because if the collection does not exist yet, it will be created in the first
	// insert
	domainDAO.RemoveAll()

	domainLifeCycle(domainDAO)
	domainDAOPerformance(domainDAO)

	println("SUCCESS!")
}

// Test all phases of the domain life cycle
func domainLifeCycle(domainDAO dao.DomainDAO) {
	domain := newDomain()

	// Create domain
	if err := domainDAO.Save(&domain); err != nil {
		fatalln("Couldn't save domain in database", err)
	}

	// Search and compare created domain
	if domainRetrieved, err := domainDAO.FindByFQDN(domain.FQDN); err != nil {
		fatalln("Couldn't find created domain in database", err)

	} else if !compareDomains(domain, domainRetrieved) {
		fatalln("Domain created in being persisted wrongly", nil)
	}

	// Update domain
	domain.Owners = []*mail.Address{}
	if err := domainDAO.Save(&domain); err != nil {
		fatalln("Couldn't save domain in database", err)
	}

	// Search and compare updated domain
	if domainRetrieved, err := domainDAO.FindByFQDN(domain.FQDN); err != nil {
		fatalln("Couldn't find updated domain in database", err)

	} else if !compareDomains(domain, domainRetrieved) {
		fatalln("Domain updated in being persisted wrongly", nil)
	}

	// Remove domain
	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		fatalln("Error while trying to remove a domain", err)
	}

	// Check removal
	if _, err := domainDAO.FindByFQDN(domain.FQDN); err == nil {
		fatalln("Domain was not removed from database", nil)
	}
}

// Check if the DAO operations are optimezed for big volume of data. After some results,
// with indexes we get 80% better performance, another good improvements was to create and
// remove many objects at once using go routines
func domainDAOPerformance(domainDAO dao.DomainDAO) {
	numberOfItems := 20000
	durationTolerance := 5.0 // seconds

	beginTimer := time.Now()
	sectionTimer := beginTimer

	// Build array to create many at once
	var domains []*model.Domain
	for i := 0; i < numberOfItems; i++ {
		domain := model.Domain{
			FQDN: fmt.Sprintf("test%d.com.br", i),
		}

		domains = append(domains, &domain)
	}

	errorInDomainsCreation := false
	domainResults := domainDAO.SaveMany(domains)

	// Check if there was any error while creating them
	for _, domainResult := range domainResults {
		if domainResult.Error != nil {
			errorInDomainsCreation = true
			errorln(fmt.Sprintf("Couldn't save domain %s in database during the performance test",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	if errorInDomainsCreation {
		fatalln("Due to errors in domain creation, the performance test will be aborted", nil)
	}

	insertDuration := time.Since(sectionTimer)
	sectionTimer = time.Now()

	// Try to find domains from different parts of the whole range to check indexes
	queryRanges := numberOfItems / 4
	fqdn1 := fmt.Sprintf("test%d.com.br", queryRanges)
	fqdn2 := fmt.Sprintf("test%d.com.br", queryRanges*2)
	fqdn3 := fmt.Sprintf("test%d.com.br", queryRanges*3)

	if _, err := domainDAO.FindByFQDN(fqdn1); err != nil {
		fatalln(fmt.Sprintf("Couldn't find domain %s in database during "+
			"the performance test", fqdn1), err)
	}

	if _, err := domainDAO.FindByFQDN(fqdn2); err != nil {
		fatalln(fmt.Sprintf("Couldn't find domain %s in database during "+
			"the performance test", fqdn2), err)
	}

	if _, err := domainDAO.FindByFQDN(fqdn3); err != nil {
		fatalln(fmt.Sprintf("Couldn't find domain %s in database during "+
			"the performance test", fqdn3), err)
	}

	queryDuration := time.Since(sectionTimer)
	sectionTimer = time.Now()

	errorInDomainsRemoval := false
	domainResults = domainDAO.RemoveMany(domains)

	// Check if there was any error while removing them
	for _, domainResult := range domainResults {
		if domainResult.Error != nil {
			errorInDomainsRemoval = true
			errorln(fmt.Sprintf("Error while trying to remove a domain %s during the performance test",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	if errorInDomainsRemoval {
		fatalln("Due to errors in domain removal, the performance test will be aborted", nil)
	}

	removeDuration := time.Since(sectionTimer)
	duration := time.Since(beginTimer)

	if duration.Seconds() > durationTolerance {
		fatalln(fmt.Sprintf("Domain DAO operations are too slow (total: %s, insert: %s, query: %s, remove: %s)",
			duration.String(), insertDuration.String(), queryDuration.String(), removeDuration.String()), nil)
	} else {
		println(fmt.Sprintf("Domain DAO operations took %s", time.Since(beginTimer).String()))
	}
}

// Function to read the configuration file
func readConfigFile() (DomainDAOTestConfigFile, error) {
	var configFile DomainDAOTestConfigFile

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

// Function to compare if two domains are equal, cannot use operator == because of the
// slices inside the domain object
func compareDomains(d1, d2 model.Domain) bool {
	if d1.Id != d2.Id || d1.FQDN != d2.FQDN {
		return false
	}

	if len(d1.Nameservers) != len(d2.Nameservers) {
		return false
	}

	for i := 0; i < len(d1.Nameservers); i++ {
		// Cannot compare the nameservers directly with operator == because of the
		// pointers for IP addresses
		if d1.Nameservers[i].Host != d2.Nameservers[i].Host ||
			d1.Nameservers[i].IPv4.String() != d2.Nameservers[i].IPv4.String() ||
			d1.Nameservers[i].IPv6.String() != d2.Nameservers[i].IPv6.String() ||
			d1.Nameservers[i].LastStatus != d2.Nameservers[i].LastStatus ||
			d1.Nameservers[i].LastCheckAt != d2.Nameservers[i].LastCheckAt ||
			d1.Nameservers[i].LastOKAt != d2.Nameservers[i].LastOKAt {
			return false
		}
	}

	if len(d1.DSSet) != len(d2.DSSet) {
		return false
	}

	for i := 0; i < len(d1.DSSet); i++ {
		if d1.DSSet[i] != d2.DSSet[i] {
			return false
		}
	}

	if len(d1.Owners) != len(d2.Owners) {
		return false
	}

	for i := 0; i < len(d1.Owners); i++ {
		if d1.Owners[i].String() != d2.Owners[i].String() {
			return false
		}
	}

	return true
}

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints a simple message without ending the test
func println(message string) {
	message = fmt.Sprintf("DomainDAO integration test: %s", message)
	log.Println(message)
}

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints an error message without ending the test
func errorln(message string, err error) {
	message = fmt.Sprintf("DomainDAO integration test: %s", message)
	if err != nil {
		message = fmt.Sprintf("%s. Details: %s", message, err.Error())
	}

	log.Println(message)
}

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints an error message and ends the test
func fatalln(message string, err error) {
	message = fmt.Sprintf("DomainDAO integration test: %s", message)
	if err != nil {
		message = fmt.Sprintf("%s. Details: %s", message, err.Error())
	}

	log.Fatalln(message)
}

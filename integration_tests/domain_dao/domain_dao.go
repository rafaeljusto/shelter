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
	configFilePath string // Path for the configuration file with the database connection information
	reportFilePath string // Path to generate the domain dao performance report file
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
	flag.StringVar(&reportFilePath, "report", "", "Report file for DomainDAO performance")
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
	domainsLifeCycle(domainDAO)
	domainUniqueFQDN(domainDAO)
	domainDAOPerformance(domainDAO)

	// Domain DAO performance report is optional and only generated when the report file
	// path parameter is given
	if len(reportFilePath) > 0 {
		domainDAOPerformanceReport(domainDAO)
	}

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

// Test all phases from a domain lyfe cycle, but now working with a group of domains
func domainsLifeCycle(domainDAO dao.DomainDAO) {
	domains := newDomains()

	// Create domains
	domainResults := domainDAO.SaveMany(domains)

	for _, domainResult := range domainResults {
		if domainResult.Error != nil {
			fatalln(fmt.Sprintf("Couldn't save domain %s in database",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	for _, domain := range domains {
		// Search and compare created domains
		if domainRetrieved, err := domainDAO.FindByFQDN(domain.FQDN); err != nil {
			fatalln(fmt.Sprintf("Couldn't find created domain %s in database", domain.FQDN), err)

		} else if !compareDomains(*domain, domainRetrieved) {
			fatalln(fmt.Sprintf("Domain %s created in being persisted wrongly", domain.FQDN), nil)
		}
	}

	// Update domains
	for _, domain := range domains {
		domain.Owners = []*mail.Address{}
	}

	domainResults = domainDAO.SaveMany(domains)

	for _, domainResult := range domainResults {
		if domainResult.Error != nil {
			fatalln(fmt.Sprintf("Couldn't update domain %s in database",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	for _, domain := range domains {
		// Search and compare updated domains
		if domainRetrieved, err := domainDAO.FindByFQDN(domain.FQDN); err != nil {
			fatalln(fmt.Sprintf("Couldn't find updated domain %s in database", domain.FQDN), err)

		} else if !compareDomains(*domain, domainRetrieved) {
			fatalln(fmt.Sprintf("Domain %s updated in being persisted wrongly", domain.FQDN), nil)
		}
	}

	// Check if find all really return all domains
	allDomainsChannel, err := domainDAO.FindAll()
	if err != nil {
		fatalln("Error while retrieving all domains from database", err)
	}

	var allDomains []model.Domain
	for {
		domainRetrieved := <-allDomainsChannel
		if domainRetrieved.Error != nil {
			fatalln("Error while retrieving all domains from database", err)
		} else if domainRetrieved.Domain == nil {
			break
		}

		allDomains = append(allDomains, *domainRetrieved.Domain)
	}

	if len(allDomains) != len(domains) {
		fatalln(fmt.Sprintf("FindAll method is not returning all domains we expected %d but got %d",
			len(domains), len(allDomains)), nil)
	}

	// Remove domains
	domainResults = domainDAO.RemoveMany(domains)

	for _, domainResult := range domainResults {
		if domainResult.Error != nil {
			fatalln(fmt.Sprintf("Error while trying to remove domain %s from database",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	for _, domain := range domains {
		// Check removals
		if _, err := domainDAO.FindByFQDN(domain.FQDN); err == nil {
			fatalln(fmt.Sprintf("Domain %s was not removed from database", domain.FQDN), nil)
		}
	}

	// Let's add and remove the domains again to test the remove all method

	domainResults = domainDAO.SaveMany(domains)
	for _, domainResult := range domainResults {
		if domainResult.Error != nil {
			fatalln(fmt.Sprintf("Couldn't save domain %s in database",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	if err := domainDAO.RemoveAll(); err != nil {
		fatalln("Couldn't remove all domains", err)
	}

	allDomainsChannel, err = domainDAO.FindAll()
	if err != nil {
		fatalln("Error while retrieving all domains from database", err)
	}

	allDomains = []model.Domain{}
	for {
		domainRetrieved := <-allDomainsChannel
		if domainRetrieved.Error != nil {
			fatalln("Error while retrieving all domains from database", err)
		} else if domainRetrieved.Domain == nil {
			break
		}

		allDomains = append(allDomains, *domainRetrieved.Domain)
	}

	if len(allDomains) > 0 {
		fatalln("RemoveAll method is not removing the domains from the database", nil)
	}
}

// FQDN must be unique in the database, today we limit this using an unique index key
func domainUniqueFQDN(domainDAO dao.DomainDAO) {
	domain1 := newDomain()

	// Create domain
	if err := domainDAO.Save(&domain1); err != nil {
		fatalln("Couldn't save domain in database", err)
	}

	domain2 := newDomain()

	// Create another domain with the same FQDN
	if err := domainDAO.Save(&domain2); err == nil {
		fatalln("Allowing more than one object with the same FQDN", nil)
	}

	// Remove domain
	if err := domainDAO.RemoveByFQDN(domain1.FQDN); err != nil {
		fatalln("Error while trying to remove a domain", err)
	}
}

// Check if the DAO operations are optimezed for big volume of data. After some results,
// with indexes we get 80% better performance, another good improvements was to create and
// remove many objects at once using go routines
func domainDAOPerformance(domainDAO dao.DomainDAO) {
	numberOfItems := 40000
	durationTolerance := 5.0 // seconds

	totalDuration, insertDuration, queryDuration, removeDuration :=
		calculateDomainDAODurations(domainDAO, numberOfItems)

	if totalDuration.Seconds() > durationTolerance {
		fatalln(fmt.Sprintf("Domain DAO operations are too slow (total: %s, insert: %s, query: %s, remove: %s)",
			totalDuration.String(), insertDuration.String(), queryDuration.String(), removeDuration.String()), nil)
	} else {
		println(fmt.Sprintf("Domain DAO operations took %s (insert: %s, query: %s, remove: %s)",
			totalDuration.String(), insertDuration.String(), queryDuration.String(), removeDuration.String()))
	}
}

// Generates a report with the amount of time for each operation in the domain DAO. For
// more realistic values it does the same operation for the same amount of data X number
// of times to get the average time of the operation
func domainDAOPerformanceReport(domainDAO dao.DomainDAO) {
	// Report header
	report := " #       | Total           | Insert          | Find            | Remove\n" +
		"------------------------------------------------------------------\n"

	// Report variables
	averageTurns := 5
	scale := []int{10, 50, 100, 500, 1000, 5000,
		10000, 50000, 100000, 500000, 1000000, 5000000}

	for _, numberOfItems := range scale {
		var totalDuration, insertDuration, queryDuration, removeDuration time.Duration

		for i := 0; i < averageTurns; i++ {
			println(fmt.Sprintf("Generating report - scale %d - turn %d", numberOfItems, i+1))
			totalDurationTmp, insertDurationTmp, queryDurationTmp, removeDurationTmp :=
				calculateDomainDAODurations(domainDAO, numberOfItems)

			totalDuration += totalDurationTmp
			insertDuration += insertDurationTmp
			queryDuration += queryDurationTmp
			removeDuration += removeDurationTmp
		}

		report += fmt.Sprintf("% -8d | % 15s | % 15s | % 15s | % 15s\n",
			numberOfItems,
			time.Duration(int64(totalDuration)/int64(averageTurns)).String(),
			time.Duration(int64(insertDuration)/int64(averageTurns)).String(),
			time.Duration(int64(queryDuration)/int64(averageTurns)).String(),
			time.Duration(int64(removeDuration)/int64(averageTurns)).String(),
		)
	}

	ioutil.WriteFile(reportFilePath, []byte(report), 0444)
}

func calculateDomainDAODurations(domainDAO dao.DomainDAO, numberOfItems int) (totalDuration, insertDuration,
	queryDuration, removeDuration time.Duration) {

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

	insertDuration = time.Since(sectionTimer)
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

	queryDuration = time.Since(sectionTimer)
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

	removeDuration = time.Since(sectionTimer)
	totalDuration = time.Since(beginTimer)

	return
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

// Function to mock a group of domains
func newDomains() []*model.Domain {
	owner1, _ := mail.ParseAddress("adm@test1.com.br")
	owner2, _ := mail.ParseAddress("adm@test2.com.br")

	return []*model.Domain{
		{
			FQDN: "test1.com.br",
			Nameservers: []model.Nameserver{
				{
					Host: "ns1.test1.com.br",
					IPv4: net.ParseIP("192.168.0.1"),
					IPv6: net.ParseIP("::1"),
				},
				{
					Host: "ns2.test1.com.br",
					IPv4: net.ParseIP("192.168.1.2"),
				},
			},
			DSSet: []model.DS{
				{
					Keytag:    1324,
					Algorithm: model.DSAlgorithmRSASHA256,
					Digest:    "A790A11EA430A85DA77245F091891F73AA7404AA",
				},
			},
			Owners: []*mail.Address{owner1},
		},
		{
			FQDN: "test2.com.br",
			Nameservers: []model.Nameserver{
				{
					Host: "ns1.test2.com.br",
					IPv4: net.ParseIP("192.168.0.3"),
					IPv6: net.ParseIP("::2"),
				},
				{
					Host: "ns2.test2.com.br",
					IPv4: net.ParseIP("192.168.0.4"),
				},
			},
			DSSet: []model.DS{
				{
					Keytag:    4321,
					Algorithm: model.DSAlgorithmGOST,
					Digest:    "A790A11EA430A85DA77245F091891F73AA7404BB",
				},
			},
			Owners: []*mail.Address{owner2},
		},
	}
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

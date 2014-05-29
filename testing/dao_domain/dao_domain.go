// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/testing/utils"
	"net"
	"net/mail"
	"time"
)

// This test objective is to verify the domain data persistence. The strategy is to insert
// and search for the information. Check for insert/update consistency (updates don't
// create a new element) and if the object id is set on creation

var (
	configFilePath string // Path for the configuration file with the database connection information
	report         bool   // Flag to generate the domain dao performance report file
)

// DomainDAOTestConfigFile is a structure to store the test configuration file data
type DomainDAOTestConfigFile struct {
	Database struct {
		URI  string
		Name string
	}

	Report struct {
		ReportFile string
	}
}

func init() {
	utils.TestName = "DomainDAO"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for DomainDAO test")
	flag.BoolVar(&report, "report", false, "Report flag for DomainDAO performance")
}

func main() {
	flag.Parse()

	var config DomainDAOTestConfigFile
	err := utils.ReadConfigFile(configFilePath, &config)

	if err == utils.ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		utils.Fatalln("Error reading configuration file", err)
	}

	database, databaseSession, err := mongodb.Open(
		[]string{config.Database.URI},
		config.Database.Name,
		false, "", "",
	)

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

	domainLifeCycle(domainDAO)
	domainsLifeCycle(domainDAO)
	domainUniqueFQDN(domainDAO)
	domainConcurrency(domainDAO)
	domainsPagination(domainDAO)
	domainsNotification(domainDAO)
	domainsExpand(domainDAO)

	// Domain DAO performance report is optional and only generated when the report file
	// path parameter is given
	if report {
		domainDAOPerformanceReport(config.Report.ReportFile, domainDAO)
	}

	utils.Println("SUCCESS!")
}

// Test all phases of the domain life cycle
func domainLifeCycle(domainDAO dao.DomainDAO) {
	domain := newDomain()

	// Create domain
	if err := domainDAO.Save(&domain); err != nil {
		utils.Fatalln("Couldn't save domain in database", err)
	}

	// Search and compare created domain
	if domainRetrieved, err := domainDAO.FindByFQDN(domain.FQDN); err != nil {
		utils.Fatalln("Couldn't find created domain in database", err)

	} else if !compareDomains(domain, domainRetrieved) {
		utils.Fatalln("Domain created is being persisted wrongly", nil)
	}

	// Update domain
	domain.Owners = []model.Owner{}
	if err := domainDAO.Save(&domain); err != nil {
		utils.Fatalln("Couldn't save domain in database", err)
	}

	// Search and compare updated domain
	if domainRetrieved, err := domainDAO.FindByFQDN(domain.FQDN); err != nil {
		utils.Fatalln("Couldn't find updated domain in database", err)

	} else if !compareDomains(domain, domainRetrieved) {
		utils.Fatalln("Domain updated is being persisted wrongly", nil)
	}

	// Remove domain
	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		utils.Fatalln("Error while trying to remove a domain", err)
	}

	// Check removal
	if _, err := domainDAO.FindByFQDN(domain.FQDN); err == nil {
		utils.Fatalln("Domain was not removed from database", nil)
	}
}

// Test all phases from a domain lyfe cycle, but now working with a group of domains
func domainsLifeCycle(domainDAO dao.DomainDAO) {
	domains := newDomains()

	// Create domains
	domainResults := domainDAO.SaveMany(domains)

	for _, domainResult := range domainResults {
		if domainResult.Error != nil {
			utils.Fatalln(fmt.Sprintf("Couldn't save domain %s in database",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	for _, domain := range domains {
		// Search and compare created domains
		if domainRetrieved, err := domainDAO.FindByFQDN(domain.FQDN); err != nil {
			utils.Fatalln(fmt.Sprintf("Couldn't find created domain %s in database", domain.FQDN), err)

		} else if !compareDomains(*domain, domainRetrieved) {
			utils.Fatalln(fmt.Sprintf("Domain %s created is being persisted wrongly", domain.FQDN), nil)
		}
	}

	// Update domains
	for _, domain := range domains {
		domain.Owners = []model.Owner{}
	}

	domainResults = domainDAO.SaveMany(domains)

	for _, domainResult := range domainResults {
		if domainResult.Error != nil {
			utils.Fatalln(fmt.Sprintf("Couldn't update domain %s in database",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	for _, domain := range domains {
		// Search and compare updated domains
		if domainRetrieved, err := domainDAO.FindByFQDN(domain.FQDN); err != nil {
			utils.Fatalln(fmt.Sprintf("Couldn't find updated domain %s in database", domain.FQDN), err)

		} else if !compareDomains(*domain, domainRetrieved) {
			utils.Fatalln(fmt.Sprintf("Domain %s updated in being persisted wrongly", domain.FQDN), nil)
		}
	}

	// Check if find all really return all domains
	allDomainsChannel, err := domainDAO.FindAllAsync()
	if err != nil {
		utils.Fatalln("Error while retrieving all domains from database", err)
	}

	var allDomains []model.Domain
	for {
		domainRetrieved := <-allDomainsChannel
		if domainRetrieved.Error != nil {
			utils.Fatalln("Error while retrieving all domains from database", err)
		} else if domainRetrieved.Domain == nil {
			break
		}

		allDomains = append(allDomains, *domainRetrieved.Domain)
	}

	if len(allDomains) != len(domains) {
		utils.Fatalln(fmt.Sprintf("FindAll method is not returning all domains we expected %d but got %d",
			len(domains), len(allDomains)), nil)
	}

	// Detected a problem in FindAsync method on 2014-01-17 where we were returning the same
	// object many times because we were reusing the same pointer. For that reason we are
	// going to add a test to check if the items returned are the same set of the inserted
	// ones
	for _, domain := range domains {
		found := false
		for _, domainRetrieved := range allDomains {
			if domainRetrieved.Id.Hex() == domain.Id.Hex() {
				found = true
				break
			}
		}

		if !found {
			utils.Fatalln("FindAll method is not returning all objects "+
				"that were inserted, apparently there are duplicated objects in the result set", nil)
		}
	}

	// Remove domains
	domainResults = domainDAO.RemoveMany(domains)

	for _, domainResult := range domainResults {
		if domainResult.Error != nil {
			utils.Fatalln(fmt.Sprintf("Error while trying to remove domain %s from database",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	for _, domain := range domains {
		// Check removals
		if _, err := domainDAO.FindByFQDN(domain.FQDN); err == nil {
			utils.Fatalln(fmt.Sprintf("Domain %s was not removed from database", domain.FQDN), nil)
		}
	}

	// Let's add and remove the domains again to test the remove all method

	domainResults = domainDAO.SaveMany(domains)
	for _, domainResult := range domainResults {
		if domainResult.Error != nil {
			utils.Fatalln(fmt.Sprintf("Couldn't save domain %s in database",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	if err := domainDAO.RemoveAll(); err != nil {
		utils.Fatalln("Couldn't remove all domains", err)
	}

	allDomainsChannel, err = domainDAO.FindAllAsync()
	if err != nil {
		utils.Fatalln("Error while retrieving all domains from database", err)
	}

	allDomains = []model.Domain{}
	for {
		domainRetrieved := <-allDomainsChannel
		if domainRetrieved.Error != nil {
			utils.Fatalln("Error while retrieving all domains from database", err)
		} else if domainRetrieved.Domain == nil {
			break
		}

		allDomains = append(allDomains, *domainRetrieved.Domain)
	}

	if len(allDomains) > 0 {
		utils.Fatalln("RemoveAll method is not removing the domains from the database", nil)
	}
}

// FQDN must be unique in the database, today we limit this using an unique index key
func domainUniqueFQDN(domainDAO dao.DomainDAO) {
	domain1 := newDomain()

	// Create domain
	if err := domainDAO.Save(&domain1); err != nil {
		utils.Fatalln("Couldn't save domain in database", err)
	}

	domain2 := newDomain()

	// Create another domain with the same FQDN
	if err := domainDAO.Save(&domain2); err == nil {
		utils.Fatalln("Allowing more than one object with the same FQDN", nil)
	}

	// Remove domain
	if err := domainDAO.RemoveByFQDN(domain1.FQDN); err != nil {
		utils.Fatalln("Error while trying to remove a domain", err)
	}
}

// Check if the revision field avoid data concurrency. Is better to fail than to store the
// wrong state
func domainConcurrency(domainDAO dao.DomainDAO) {
	domain := newDomain()

	// Create domain
	if err := domainDAO.Save(&domain); err != nil {
		utils.Fatalln("Couldn't save domain in database", err)
	}

	domain1, err := domainDAO.FindByFQDN(domain.FQDN)
	if err != nil {
		utils.Fatalln("Couldn't find created domain in database", err)
	}

	domain2, err := domainDAO.FindByFQDN(domain.FQDN)
	if err != nil {
		utils.Fatalln("Couldn't find created domain in database", err)
	}

	if err := domainDAO.Save(&domain1); err != nil {
		utils.Fatalln("Couldn't save domain in database", err)
	}

	if err := domainDAO.Save(&domain2); err == nil {
		utils.Fatalln("Not controlling domain concurrency", nil)
	}

	// Remove domain
	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		utils.Fatalln("Error while trying to remove a domain", err)
	}
}

func domainsPagination(domainDAO dao.DomainDAO) {
	numberOfItems := 1000

	for i := 0; i < numberOfItems; i++ {
		domain := model.Domain{
			FQDN: fmt.Sprintf("example%d.com.br", i),
		}

		if err := domainDAO.Save(&domain); err != nil {
			utils.Fatalln("Error saving domain in database", err)
		}
	}

	pagination := dao.DomainDAOPagination{
		PageSize: 10,
		Page:     5,
		OrderBy: []dao.DomainDAOSort{
			{
				Field:     dao.DomainDAOOrderByFieldFQDN,
				Direction: dao.DAOOrderByDirectionAscending,
			},
		},
	}

	domains, err := domainDAO.FindAll(&pagination, true)
	if err != nil {
		utils.Fatalln("Error retrieving domains", err)
	}

	if pagination.NumberOfItems != numberOfItems {
		utils.Errorln("Number of items not calculated correctly", nil)
	}

	if pagination.NumberOfPages != numberOfItems/pagination.PageSize {
		utils.Errorln("Number of pages not calculated correctly", nil)
	}

	if len(domains) != pagination.PageSize {
		utils.Errorln("Number of domains not following page size", nil)
	}

	pagination = dao.DomainDAOPagination{
		PageSize: 10000,
		Page:     1,
		OrderBy: []dao.DomainDAOSort{
			{
				Field:     dao.DomainDAOOrderByFieldFQDN,
				Direction: dao.DAOOrderByDirectionAscending,
			},
		},
	}

	domains, err = domainDAO.FindAll(&pagination, true)
	if err != nil {
		utils.Fatalln("Error retrieving domains", err)
	}

	if pagination.NumberOfPages != 1 {
		utils.Fatalln("Calculating wrong number of pages when there's only one page", nil)
	}

	for i := 0; i < numberOfItems; i++ {
		fqdn := fmt.Sprintf("example%d.com.br", i)
		if err := domainDAO.RemoveByFQDN(fqdn); err != nil {
			utils.Fatalln("Error removing domain from database", err)
		}
	}
}

// Verify if the method that choose the domains that needs to be verified is correct
func domainsNotification(domainDAO dao.DomainDAO) {
	numberOfItemsToBeVerified := 1000
	numberOfItemsToDontBeVerified := 1000
	nameserverErrorAlertDays := 7
	nameserverTimeoutAlertDays := 30
	dsErrorAlertDays := 1
	dsTimeoutAlertDays := 7
	maxExpirationAlertDays := 5

	data := []struct {
		name                      string
		numberOfItems             int
		nameserverTimeoutLastOKAt time.Time
		nameserverErrorLastOKAt   time.Time
		dsTimeoutLastOkAt         time.Time
		dsErrorLastOkAt           time.Time
		dsExpiresAt               time.Time
	}{
		{
			name:                      "shouldbenotified",
			numberOfItems:             numberOfItemsToBeVerified,
			nameserverTimeoutLastOKAt: time.Now().Add(time.Duration(-nameserverTimeoutAlertDays*24) * time.Hour),
			nameserverErrorLastOKAt:   time.Now().Add(time.Duration(-nameserverErrorAlertDays*24) * time.Hour),
			dsTimeoutLastOkAt:         time.Now().Add(time.Duration(-dsTimeoutAlertDays*24) * time.Hour),
			dsErrorLastOkAt:           time.Now().Add(time.Duration(-dsErrorAlertDays*24) * time.Hour),
			dsExpiresAt:               time.Now().Add(time.Duration((maxExpirationAlertDays)*24) * time.Hour),
		},
		{
			name:                      "shouldnotbenotified",
			numberOfItems:             numberOfItemsToDontBeVerified,
			nameserverTimeoutLastOKAt: time.Now().Add(time.Duration((-nameserverTimeoutAlertDays+1)*24) * time.Hour),
			nameserverErrorLastOKAt:   time.Now().Add(time.Duration((-nameserverErrorAlertDays+1)*24) * time.Hour),
			dsTimeoutLastOkAt:         time.Now().Add(time.Duration((-dsTimeoutAlertDays+1)*24) * time.Hour),
			dsErrorLastOkAt:           time.Now().Add(time.Duration((-dsErrorAlertDays+1)*24) * time.Hour),
			dsExpiresAt:               time.Now().Add(time.Duration((maxExpirationAlertDays+1)*24) * time.Hour),
		},
	}

	for _, item := range data {
		for i := 0; i < item.numberOfItems/5; i++ {
			domain := model.Domain{
				FQDN: fmt.Sprintf("%s%d.com.br", item.name, i),
				Nameservers: []model.Nameserver{
					{
						LastStatus: model.NameserverStatusTimeout,
						LastOKAt:   item.nameserverTimeoutLastOKAt,
					},
				},
			}

			if err := domainDAO.Save(&domain); err != nil {
				utils.Fatalln("Error saving domain in database", err)
			}
		}

		for i := item.numberOfItems / 5; i < item.numberOfItems/5*2; i++ {
			domain := model.Domain{
				FQDN: fmt.Sprintf("%s%d.com.br", item.name, i),
				Nameservers: []model.Nameserver{
					{
						LastStatus: model.NameserverStatusNoAuthority,
						LastOKAt:   item.nameserverErrorLastOKAt,
					},
				},
			}

			if err := domainDAO.Save(&domain); err != nil {
				utils.Fatalln("Error saving domain in database", err)
			}
		}

		for i := item.numberOfItems / 5 * 2; i < item.numberOfItems/5*3; i++ {
			domain := model.Domain{
				FQDN: fmt.Sprintf("%s%d.com.br", item.name, i),
				DSSet: []model.DS{
					{
						LastStatus: model.DSStatusTimeout,
						LastOKAt:   item.dsTimeoutLastOkAt,
						ExpiresAt:  time.Now().Add(time.Duration((maxExpirationAlertDays+1)*24) * time.Hour),
					},
				},
			}

			if err := domainDAO.Save(&domain); err != nil {
				utils.Fatalln("Error saving domain in database", err)
			}
		}

		for i := item.numberOfItems / 5 * 3; i < item.numberOfItems/5*4; i++ {
			domain := model.Domain{
				FQDN: fmt.Sprintf("%s%d.com.br", item.name, i),
				DSSet: []model.DS{
					{
						LastStatus: model.DSStatusExpiredSignature,
						LastOKAt:   item.dsErrorLastOkAt,
						ExpiresAt:  time.Now().Add(time.Duration((maxExpirationAlertDays+1)*24) * time.Hour),
					},
				},
			}

			if err := domainDAO.Save(&domain); err != nil {
				utils.Fatalln("Error saving domain in database", err)
			}
		}

		for i := item.numberOfItems / 5 * 4; i < item.numberOfItems; i++ {
			domain := model.Domain{
				FQDN: fmt.Sprintf("%s%d.com.br", item.name, i),
				DSSet: []model.DS{
					{
						LastStatus: model.DSStatusOK,
						LastOKAt:   time.Now(),
						ExpiresAt:  item.dsExpiresAt,
					},
				},
			}

			if err := domainDAO.Save(&domain); err != nil {
				utils.Fatalln("Error saving domain in database", err)
			}
		}
	}

	domainChannel, err := domainDAO.FindAllAsyncToBeNotified(
		nameserverErrorAlertDays,
		nameserverTimeoutAlertDays,
		dsErrorAlertDays,
		dsTimeoutAlertDays,
		maxExpirationAlertDays,
	)

	if err != nil {
		utils.Fatalln("Error retrieving domains to be notified", err)
	}

	var domains []*model.Domain
	for {
		domainResult := <-domainChannel

		if domainResult.Error != nil {
			utils.Fatalln("Error retrieving domain to be notified", domainResult.Error)
		}

		if domainResult.Error != nil || domainResult.Domain == nil {
			break
		}

		domains = append(domains, domainResult.Domain)
	}

	if len(domains) != numberOfItemsToBeVerified {
		utils.Fatalln(fmt.Sprintf("Did not select all the domains ready for notification. "+
			"Expected %d and got %d", numberOfItemsToBeVerified, len(domains)), nil)
	}

	for _, item := range data {
		for i := 0; i < item.numberOfItems; i++ {
			fqdn := fmt.Sprintf("%s%d.com.br", item.name, i)
			if err := domainDAO.RemoveByFQDN(fqdn); err != nil {
				utils.Fatalln("Error removing domain from database", err)
			}
		}
	}
}

func domainsExpand(domainDAO dao.DomainDAO) {
	newDomains := newDomains()
	domainsResult := domainDAO.SaveMany(newDomains)
	for _, domainResult := range domainsResult {
		if domainResult.Error != nil {
			utils.Fatalln("Error creating domains", domainResult.Error)
		}
	}

	pagination := dao.DomainDAOPagination{}
	domains, err := domainDAO.FindAll(&pagination, false)

	if err != nil {
		utils.Fatalln("Error retrieving domains", err)
	}

	for _, domain := range domains {
		if len(domain.Owners) > 0 {
			utils.Fatalln("Not compressing owners in results", nil)
		}

		for _, nameserver := range domain.Nameservers {
			if len(nameserver.Host) > 0 ||
				nameserver.IPv4 != nil ||
				nameserver.IPv6 != nil ||
				!nameserver.LastCheckAt.Equal(time.Time{}) ||
				!nameserver.LastOKAt.Equal(time.Time{}) {
				utils.Fatalln("Not compressing nameservers in results", nil)
			}
		}

		for _, ds := range domain.DSSet {
			if ds.Algorithm != 0 ||
				len(ds.Digest) > 0 ||
				ds.DigestType != 0 ||
				ds.Keytag != 0 ||
				!ds.ExpiresAt.Equal(time.Time{}) ||
				!ds.LastCheckAt.Equal(time.Time{}) ||
				!ds.LastOKAt.Equal(time.Time{}) {
				utils.Fatalln("Not compressing ds set in results", nil)
			}
		}
	}

	domains, err = domainDAO.FindAll(&pagination, true)

	if err != nil {
		utils.Fatalln("Error retrieving domains", err)
	}

	for _, domain := range domains {
		if len(domain.Owners) == 0 {
			utils.Fatalln("Compressing owners in results when it shouldn't", nil)
		}

		for _, nameserver := range domain.Nameservers {
			if len(nameserver.Host) == 0 ||
				nameserver.IPv4 == nil ||
				nameserver.IPv6 == nil ||
				nameserver.LastCheckAt.Equal(time.Time{}) ||
				nameserver.LastOKAt.Equal(time.Time{}) ||
				nameserver.LastStatus != model.NameserverStatusOK {
				utils.Fatalln("Compressing nameservers in results when it shouldn't", nil)
			}
		}

		for _, ds := range domain.DSSet {
			if ds.Algorithm == 0 ||
				len(ds.Digest) == 0 ||
				ds.DigestType == 0 ||
				ds.Keytag == 0 ||
				ds.ExpiresAt.Equal(time.Time{}) ||
				ds.LastCheckAt.Equal(time.Time{}) ||
				ds.LastOKAt.Equal(time.Time{}) ||
				ds.LastStatus != model.DSStatusOK {
				utils.Fatalln("Compressing ds set in results when it shouldn't", nil)
			}
		}
	}

	domainsResult = domainDAO.RemoveMany(newDomains)
	for _, domainResult := range domainsResult {
		if domainResult.Error != nil {
			utils.Fatalln("Error removing domains", domainResult.Error)
		}
	}
}

// Generates a report with the amount of time for each operation in the domain DAO. For
// more realistic values it does the same operation for the same amount of data X number
// of times to get the average time of the operation. After some results, with indexes we
// get 80% better performance, another good improvements was to create and remove many
// objects at once using go routines
func domainDAOPerformanceReport(reportFile string, domainDAO dao.DomainDAO) {
	// Report header
	report := " #       | Total            | Insert           | Find             | Remove\n" +
		"------------------------------------------------------------------------------------\n"

	// Report variables
	averageTurns := 5
	scale := []int{10, 50, 100, 500, 1000, 5000,
		10000, 50000, 100000, 500000, 1000000, 5000000}

	for _, numberOfItems := range scale {
		var totalDuration, insertDuration, queryDuration, removeDuration time.Duration

		for i := 0; i < averageTurns; i++ {
			utils.Println(fmt.Sprintf("Generating report - scale %d - turn %d", numberOfItems, i+1))
			totalDurationTmp, insertDurationTmp, queryDurationTmp, removeDurationTmp :=
				calculateDomainDAODurations(domainDAO, numberOfItems)

			totalDuration += totalDurationTmp
			insertDuration += insertDurationTmp
			queryDuration += queryDurationTmp
			removeDuration += removeDurationTmp
		}

		report += fmt.Sprintf("% -8d | % 16s | % 16s | % 16s | % 16s\n",
			numberOfItems,
			time.Duration(int64(totalDuration)/int64(averageTurns)).String(),
			time.Duration(int64(insertDuration)/int64(averageTurns)).String(),
			time.Duration(int64(queryDuration)/int64(averageTurns)).String(),
			time.Duration(int64(removeDuration)/int64(averageTurns)).String(),
		)
	}

	utils.WriteReport(reportFile, report)
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
			utils.Errorln(fmt.Sprintf("Couldn't save domain %s in database during the performance test",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	if errorInDomainsCreation {
		utils.Fatalln("Due to errors in domain creation, the performance test will be aborted", nil)
	}

	insertDuration = time.Since(sectionTimer)
	sectionTimer = time.Now()

	// Try to find domains from different parts of the whole range to check indexes
	queryRanges := numberOfItems / 4
	fqdn1 := fmt.Sprintf("test%d.com.br", queryRanges)
	fqdn2 := fmt.Sprintf("test%d.com.br", queryRanges*2)
	fqdn3 := fmt.Sprintf("test%d.com.br", queryRanges*3)

	if _, err := domainDAO.FindByFQDN(fqdn1); err != nil {
		utils.Fatalln(fmt.Sprintf("Couldn't find domain %s in database during "+
			"the performance test", fqdn1), err)
	}

	if _, err := domainDAO.FindByFQDN(fqdn2); err != nil {
		utils.Fatalln(fmt.Sprintf("Couldn't find domain %s in database during "+
			"the performance test", fqdn2), err)
	}

	if _, err := domainDAO.FindByFQDN(fqdn3); err != nil {
		utils.Fatalln(fmt.Sprintf("Couldn't find domain %s in database during "+
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
			utils.Errorln(fmt.Sprintf("Error while trying to remove a domain %s during the performance test",
				domainResult.Domain.FQDN), domainResult.Error)
		}
	}

	if errorInDomainsRemoval {
		utils.Fatalln("Due to errors in domain removal, the performance test will be aborted", nil)
	}

	removeDuration = time.Since(sectionTimer)
	totalDuration = time.Since(beginTimer)

	return
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
	domain.Owners = []model.Owner{
		{
			Email:    owner,
			Language: "pt-BR",
		},
	}

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
					Host:        "ns1.test1.com.br",
					IPv4:        net.ParseIP("192.168.0.1"),
					IPv6:        net.ParseIP("::1"),
					LastCheckAt: time.Now(),
					LastOKAt:    time.Now(),
					LastStatus:  model.NameserverStatusOK,
				},
				{
					Host:        "ns2.test1.com.br",
					IPv4:        net.ParseIP("192.168.1.2"),
					LastCheckAt: time.Now(),
					LastOKAt:    time.Now(),
					LastStatus:  model.NameserverStatusOK,
				},
			},
			DSSet: []model.DS{
				{
					Keytag:      1324,
					Algorithm:   model.DSAlgorithmRSASHA256,
					Digest:      "A790A11EA430A85DA77245F091891F73AA7404AA",
					DigestType:  model.DSDigestTypeSHA1,
					ExpiresAt:   time.Now(),
					LastCheckAt: time.Now(),
					LastOKAt:    time.Now(),
					LastStatus:  model.DSStatusOK,
				},
			},
			Owners: []model.Owner{
				{
					Email:    owner1,
					Language: "pt-BR",
				},
			},
		},
		{
			FQDN: "test2.com.br",
			Nameservers: []model.Nameserver{
				{
					Host:        "ns1.test2.com.br",
					IPv4:        net.ParseIP("192.168.0.3"),
					IPv6:        net.ParseIP("::2"),
					LastCheckAt: time.Now(),
					LastOKAt:    time.Now(),
					LastStatus:  model.NameserverStatusOK,
				},
				{
					Host:        "ns2.test2.com.br",
					IPv4:        net.ParseIP("192.168.0.4"),
					LastCheckAt: time.Now(),
					LastOKAt:    time.Now(),
					LastStatus:  model.NameserverStatusOK,
				},
			},
			DSSet: []model.DS{
				{
					Keytag:      4321,
					Algorithm:   model.DSAlgorithmECCGOST,
					Digest:      "A790A11EA430A85DA77245F091891F73AA7404BB",
					DigestType:  model.DSDigestTypeSHA1,
					ExpiresAt:   time.Now(),
					LastCheckAt: time.Now(),
					LastOKAt:    time.Now(),
					LastStatus:  model.DSStatusOK,
				},
			},
			Owners: []model.Owner{
				{
					Email:    owner2,
					Language: "en-US",
				},
			},
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
		// pointers for IP addresses and dates
		if d1.Nameservers[i].Host != d2.Nameservers[i].Host ||
			d1.Nameservers[i].IPv4.String() != d2.Nameservers[i].IPv4.String() ||
			d1.Nameservers[i].IPv6.String() != d2.Nameservers[i].IPv6.String() ||
			d1.Nameservers[i].LastStatus != d2.Nameservers[i].LastStatus ||
			d1.Nameservers[i].LastCheckAt.Unix() != d2.Nameservers[i].LastCheckAt.Unix() ||
			d1.Nameservers[i].LastOKAt.Unix() != d2.Nameservers[i].LastOKAt.Unix() {
			return false
		}
	}

	if len(d1.DSSet) != len(d2.DSSet) {
		return false
	}

	for i := 0; i < len(d1.DSSet); i++ {
		// Cannot compare the nameservers directly with operator == because of the dates
		if d1.DSSet[i].Algorithm != d2.DSSet[i].Algorithm ||
			d1.DSSet[i].Digest != d2.DSSet[i].Digest ||
			d1.DSSet[i].DigestType != d2.DSSet[i].DigestType ||
			d1.DSSet[i].ExpiresAt.Unix() != d2.DSSet[i].ExpiresAt.Unix() ||
			d1.DSSet[i].Keytag != d2.DSSet[i].Keytag ||
			d1.DSSet[i].LastCheckAt.Unix() != d2.DSSet[i].LastCheckAt.Unix() ||
			d1.DSSet[i].LastOKAt.Unix() != d2.DSSet[i].LastOKAt.Unix() ||
			d1.DSSet[i].LastStatus != d2.DSSet[i].LastStatus {
			return false
		}
	}

	if len(d1.Owners) != len(d2.Owners) {
		return false
	}

	for i := 0; i < len(d1.Owners); i++ {
		if d1.Owners[i].Email.String() != d2.Owners[i].Email.String() ||
			d1.Owners[i].Language != d2.Owners[i].Language {
			return false
		}
	}

	return true
}

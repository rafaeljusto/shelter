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
	"github.com/rafaeljusto/shelter/net/scan"
	"github.com/rafaeljusto/shelter/testing/utils"
	"net"
	"net/mail"
	"strconv"
	"sync"
	"time"
)

// This test objective is to verify the domain selection rules in the injector scanner
// component. As the injector depends on a database scenario, all this checks are going to
// be made in an integration test enviroment

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
	utils.TestName = "ScanInjector"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for ScanInjector test")
}

func main() {
	flag.Parse()

	var config ScanInjectorTestConfigFile
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

	// Remove all data before starting the test. This is necessary because maybe in the last
	// test there was an error and the data wasn't removed from the database
	utils.ClearDatabase(database)

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	domainWithDNSErrors(config, domainDAO)
	domainWithDNSSECErrors(config, domainDAO)
	domainWithNoErrors(config, domainDAO)

	utils.Println("SUCCESS!")
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
		utils.Fatalln("Error saving domain for scan scenario", err)
	}

	model.StartNewScan()
	if domains := runScan(config, domainDAO); len(domains) != 1 {
		utils.Fatalln(fmt.Sprintf("Couldn't load a domain with DNS errors for scan. "+
			"Expected 1 got %d", len(domains)), nil)
	}

	currentScan := model.GetCurrentScan()
	if currentScan.Status != model.ScanStatusRunning {
		utils.Fatalln("Not changing the scan info status with DNS errors", nil)
	}

	if currentScan.DomainsToBeScanned != 1 {
		utils.Fatalln("Not counting the domains to be scanned with DNS errors", nil)
	}

	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		utils.Fatalln("Error removing domain", err)
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
		utils.Fatalln("Error saving domain for scan scenario", err)
	}

	model.StartNewScan()
	if domains := runScan(config, domainDAO); len(domains) != 1 {
		utils.Fatalln(fmt.Sprintf("Couldn't load a domain with DNSSEC errors for scan. "+
			"Expected 1 got %d", len(domains)), nil)
	}

	currentScan := model.GetCurrentScan()
	if currentScan.Status != model.ScanStatusRunning {
		utils.Fatalln("Not changing the scan info status for domain with DNSSEC errors", nil)
	}

	if currentScan.DomainsToBeScanned != 1 {
		utils.Fatalln("Not counting the domains to be scanned for domain with DNSSEC errors", nil)
	}

	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		utils.Fatalln("Error removing domain", err)
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
		utils.Fatalln("Error saving domain for scan scenario", err)
	}

	model.StartNewScan()
	if domains := runScan(config, domainDAO); len(domains) > 0 {
		utils.Fatalln(fmt.Sprintf("Selected a domain configured correctly for the scan. "+
			"Expected 0 got %d", len(domains)), nil)
	}

	currentScan := model.GetCurrentScan()
	if currentScan.Status != model.ScanStatusRunning {
		utils.Fatalln("Not changing the scan info status for domain with no errors", nil)
	}

	if currentScan.DomainsToBeScanned > 0 {
		utils.Fatalln("Not counting the domains to be scanned for domain with no errors", nil)
	}

	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		utils.Fatalln("Error removing domain", err)
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
			utils.Fatalln("Error selecting domain", err)
		}

		if exit {
			break
		}
	}

	return domains
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

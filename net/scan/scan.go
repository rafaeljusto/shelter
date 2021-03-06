// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package scan is the scan service
package scan

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/miekg/dns"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
)

// When converting a DNSKEY into a DS we need to choose wich digest type are we going to
// use, as we don't want to bother the user asking this information we assume a default
// digest type
const (
	DefaultDigestType = model.DSDigestTypeSHA256
)

// Function responsible for running the domain scan system, checking the configuration of each
// domain in the database according to an algorithm. This method is synchronous and will return only
// after the scan proccess is done
func ScanDomains() {
	defer func() {
		// Something went really wrong while scanning the domains. Log the error stacktrace
		// and move out
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("Panic detected while scanning domains. Details: %v\n%s", r, buf)
		}
	}()

	log.Info("Start scan job")
	defer func() {
		log.Info("End scan job")
	}()

	log.Debugf("Initializing database with the parameters: URIS - %v | Name - %s | Auth - %t | Username - %s",
		config.ShelterConfig.Database.URIs,
		config.ShelterConfig.Database.Name,
		config.ShelterConfig.Database.Auth.Enabled,
		config.ShelterConfig.Database.Auth.Username,
	)

	database, databaseSession, err := mongodb.Open(
		config.ShelterConfig.Database.URIs,
		config.ShelterConfig.Database.Name,
		config.ShelterConfig.Database.Auth.Enabled,
		config.ShelterConfig.Database.Auth.Username,
		config.ShelterConfig.Database.Auth.Password,
	)

	if err != nil {
		log.Println("Error while initializing database. Details:", err)
		return
	}
	defer databaseSession.Close()

	injector := NewInjector(
		database,
		config.ShelterConfig.Scan.DomainsBufferSize,
		config.ShelterConfig.Scan.VerificationIntervals.MaxOKDays,
		config.ShelterConfig.Scan.VerificationIntervals.MaxErrorDays,
		config.ShelterConfig.Scan.VerificationIntervals.MaxExpirationAlertDays,
	)

	querierDispatcher := NewQuerierDispatcher(
		config.ShelterConfig.Scan.NumberOfQueriers,
		config.ShelterConfig.Scan.DomainsBufferSize,
		config.ShelterConfig.Scan.UDPMaxSize,
		time.Duration(config.ShelterConfig.Scan.Timeouts.DialSeconds)*time.Second,
		time.Duration(config.ShelterConfig.Scan.Timeouts.ReadSeconds)*time.Second,
		time.Duration(config.ShelterConfig.Scan.Timeouts.WriteSeconds)*time.Second,
		config.ShelterConfig.Scan.ConnectionRetries,
	)

	collector := NewCollector(
		database,
		config.ShelterConfig.Scan.SaveAtOnce,
	)

	// Create a new scan information
	model.StartNewScan()

	var scanGroup sync.WaitGroup
	errorsChannel := make(chan error, config.ShelterConfig.Scan.ErrorsBufferSize)
	domainsToQueryChannel := injector.Start(&scanGroup, errorsChannel)
	domainsToSaveChannel := querierDispatcher.Start(&scanGroup, domainsToQueryChannel)
	collector.Start(&scanGroup, domainsToSaveChannel, errorsChannel)

	// Keep track of errors for the scan information structure
	errorDetected := false

	go func() {
		for {
			select {
			case err := <-errorsChannel:
				// Detect the poison pill to finish the error listener go routine. This poison
				// pill should be sent after all parts of the scan are done and we are sure that
				// we don't have any error to log anymore
				if err == nil {
					return

				} else {
					errorDetected = true
					log.Println("Error detected while executing the scan. Details:", err)
				}
			}
		}
	}()

	// Wait for all parts of the scan to finish their job
	scanGroup.Wait()

	// Finish the error listener sending a poison pill
	errorsChannel <- nil

	scanDAO := dao.ScanDAO{
		Database: database,
	}

	// Save the scan information for future reports
	if err := model.FinishAndSaveScan(errorDetected, scanDAO.Save); err != nil {
		log.Println("Error while saving scan information. Details:", err)
	}
}

// Function created to check a single domain without persisting in database. Useful for online
// domain checking. As we update the same object, we update the parameter pointer and don't return
// nothing
func ScanDomain(domain *model.Domain) {
	querierDispatcher := NewQuerierDispatcher(
		config.ShelterConfig.Scan.NumberOfQueriers,
		config.ShelterConfig.Scan.DomainsBufferSize,
		config.ShelterConfig.Scan.UDPMaxSize,
		time.Duration(config.ShelterConfig.Scan.Timeouts.DialSeconds)*time.Second,
		time.Duration(config.ShelterConfig.Scan.Timeouts.ReadSeconds)*time.Second,
		time.Duration(config.ShelterConfig.Scan.Timeouts.WriteSeconds)*time.Second,
		config.ShelterConfig.Scan.ConnectionRetries,
	)

	var scanGroup sync.WaitGroup
	domainsToQueryChannel := make(chan *model.Domain)
	domainsToSaveChannel := querierDispatcher.Start(&scanGroup, domainsToQueryChannel)

	domainsToQueryChannel <- domain
	domainsToQueryChannel <- nil // Poison pill
	domain = <-domainsToSaveChannel

	// Wait for all parts of the scan to finish their job
	scanGroup.Wait()
}

// Send DNS requests to fill a domain object from the information found on the DNS authoritative
// nameservers. This is very usefull to make it easier for the user to fill forms with the domain
// information. The domain must be already delegated by a registry to this function works, because
// it uses a recursive DNS
func QueryDomain(fqdn string) (model.Domain, error) {
	domain := model.Domain{
		FQDN: fqdn,
	}

	querier := newQuerier(
		config.ShelterConfig.Scan.UDPMaxSize,
		time.Duration(config.ShelterConfig.Scan.Timeouts.DialSeconds)*time.Second,
		time.Duration(config.ShelterConfig.Scan.Timeouts.ReadSeconds)*time.Second,
		time.Duration(config.ShelterConfig.Scan.Timeouts.WriteSeconds)*time.Second,
		config.ShelterConfig.Scan.ConnectionRetries,
	)

	resolver := fmt.Sprintf("%s:%d",
		config.ShelterConfig.Scan.Resolver.Address,
		config.ShelterConfig.Scan.Resolver.Port,
	)

	var dnsRequestMessage dns.Msg
	dnsRequestMessage.SetQuestion(fqdn, dns.TypeNS)
	dnsRequestMessage.RecursionDesired = true

	// Allow retrieving domain information when there's a DNSSEC problem in the chain-of-trust
	dnsRequestMessage.CheckingDisabled = true

	dnsResponseMsg, err := querier.sendDNSRequest(resolver, &dnsRequestMessage)
	if err != nil {
		return domain, err
	}

	for _, answer := range dnsResponseMsg.Answer {
		nsRecord, ok := answer.(*dns.NS)
		if !ok {
			continue
		}

		domain.Nameservers = append(domain.Nameservers, model.Nameserver{
			Host: nsRecord.Ns,
		})
	}

	for index, nameserver := range domain.Nameservers {
		// Don't need to retrieve glue records if not necessary
		if !strings.HasSuffix(nameserver.Host, domain.FQDN) {
			continue
		}

		dnsRequestMessage.SetQuestion(nameserver.Host, dns.TypeA)
		dnsResponseMsg, err = querier.sendDNSRequest(resolver, &dnsRequestMessage)
		if err != nil {
			return domain, err
		}

		for _, answer := range dnsResponseMsg.Answer {
			ipv4Record, ok := answer.(*dns.A)
			if !ok {
				continue
			}

			domain.Nameservers[index].IPv4 = ipv4Record.A
		}

		dnsRequestMessage.SetQuestion(nameserver.Host, dns.TypeAAAA)
		dnsResponseMsg, err = querier.sendDNSRequest(resolver, &dnsRequestMessage)
		if err != nil {
			return domain, err
		}

		for _, answer := range dnsResponseMsg.Answer {
			ipv6Record, ok := answer.(*dns.AAAA)
			if !ok {
				continue
			}

			domain.Nameservers[index].IPv6 = ipv6Record.AAAA
		}
	}

	// We are going to retrieve the DNSKEYs from the user zone, and generate the DS records from it.
	// This is good if the user wants to use the Shelter project as a easy-to-fill domain registration
	// form
	dnsRequestMessage.SetQuestion(fqdn, dns.TypeDNSKEY)

	dnsResponseMsg, err = querier.sendDNSRequest(resolver, &dnsRequestMessage)
	if err != nil {
		return domain, err
	}

	for _, answer := range dnsResponseMsg.Answer {
		dnskeyRecord, ok := answer.(*dns.DNSKEY)
		if !ok {
			continue
		}

		// Only add DNSKEYs with bit SEP on
		if (dnskeyRecord.Flags & dns.SEP) == 0 {
			continue
		}

		dsRecord := dnskeyRecord.ToDS(uint8(DefaultDigestType))

		domain.DSSet = append(domain.DSSet, model.DS{
			Keytag:     dsRecord.KeyTag,
			Algorithm:  model.DSAlgorithm(dsRecord.Algorithm),
			DigestType: model.DSDigestType(dsRecord.DigestType),
			Digest:     dsRecord.Digest,
		})
	}

	return domain, nil
}

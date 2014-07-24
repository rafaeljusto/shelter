// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package scan is the scan service
package scan

import (
	"github.com/miekg/dns"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/scan/dspolicy"
	"github.com/rafaeljusto/shelter/net/scan/nspolicy"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	querierDomainsQueueSize = 10 // Number of domains that can wait in the querier channel
)

var (
	// DNS query port. It's not a constant because in test scenarios we change the DNS port
	// to one that don't need root privilleges
	DNSPort = 53
)

// Querier is responsable for sending the DNS queries to check if the namerservers are
// configured correctly with DNS/DNSSEC.  The UDPMaxSize attribute is used for DNSSEC
// queries to notify the maximum UDP package size supported in the network. This object is
// private for this package and should only be accessed by the querier dispatcher
type querier struct {
	client            dns.Client // Low level DNS client for network checks
	UDPMaxSize        uint16     // UDP max package size to pass over firewalls
	ConnectionRetries int        // Number of retries before setting timeout
}

// Return a new Querier object with the necessary fields for the scan filled
func newQuerier(udpMaxSize uint16, dialTimeout, readTimeout,
	writeTimeout time.Duration, connectionRetries int) *querier {

	return &querier{
		client: dns.Client{
			DialTimeout:  dialTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		UDPMaxSize:        udpMaxSize,
		ConnectionRetries: connectionRetries,
	}
}

// Structure to store domains that will be postponed because of to many queries sent to
// onbly one host
type postponedDomain struct {
	domain *model.Domain // Domain postponed
	index  int           // Host index that exceeded QPS
}

// Fire a querier that will process domains sent via channel until receives a poison pill
// (nil domain), for go routines control this method receives a wait group, so that the
// main thread can wait for everubody finishs. It also receives the channel where the
// querier will put the domains for the collector save them in database
func (q *querier) start(queriers *sync.WaitGroup,
	domainsToSaveChannel chan *model.Domain) chan *model.Domain {
	// Create the communication channel that we are going to listen to retrieve domains, we
	// can store more than one domain in this channel because some queriers can slow down
	// when checking domains with timeouts
	querierChannel := make(chan *model.Domain, querierDomainsQueueSize)

	// Add one more in the group of queriers. This go routine don't need to be in the group
	// of the scan go routines because it is controlled by the querier dispatcher
	queriers.Add(1)

	go func() {
		var postponedDomains []postponedDomain

		for {
			// Retrieve the domain from the channel
			domain := <-querierChannel

			// Detect the poison pill from the dispatcher
			if domain == nil {

				// Check domains that were postponed due to QPS limits for the nameservers. We
				// don't use the foreach feature beacause, according to tests, we cannot push a
				// new value into the slice while we iterate over it
				for i := 0; i < len(postponedDomains); i++ {
					postponed := postponedDomains[i]

					// We also send the list to the method so it can postpone the domain again and
					// again and again...
					if q.checkPostponedDomains(postponedDomains, postponed) {
						domainsToSaveChannel <- postponed.domain
					}
				}

				// Tell everyone that we are done!
				queriers.Done()
				return
			}

			if q.checkDomain(domain, postponedDomains) {
				// Send to collector the domain with the new state
				domainsToSaveChannel <- domain
			}
		}
	}()

	return querierChannel
}

// Main function to check a domain DNS/DNSSEC configuration. Returns true if domain is
// done checking and can be saved or false otherwise, that indicates that the domain was
// postponed
func (q *querier) checkDomain(domain *model.Domain,
	postponedDomains []postponedDomain) bool {

	for index, _ := range domain.Nameservers {
		if !q.checkNameserver(domain, index, postponedDomains) {
			return false
		}

		if !q.checkDS(domain, index, q.UDPMaxSize, postponedDomains) {
			return false
		}
	}

	return true
}

// Verify the DNS configuration on the nameservers. This method will send a SOA request
// for each nameserver and verify the results. Returns true if nameserver is done checking
// and can be saved or false otherwise, that indicates that the domain was postponed
func (q *querier) checkNameserver(domain *model.Domain,
	index int, postponedDomains []postponedDomain) bool {

	nameserver := domain.Nameservers[index]
	domainNSPolicy := nspolicy.NewDomainNSPolicy(domain)

	// Build message to send the request
	var dnsRequestMessage dns.Msg
	dnsRequestMessage.SetQuestion(domain.FQDN, dns.TypeSOA)
	dnsRequestMessage.RecursionDesired = false

	host, err := getHost(domain.FQDN, nameserver)
	if err == ErrHostTimeout {
		domain.Nameservers[index].ChangeStatus(model.NameserverStatusTimeout)
		return true

	} else if err == ErrHostQPSExceeded {
		postponedDomains = append(postponedDomains, postponedDomain{
			domain: domain,
			index:  index,
		})
		return false
	}

	dnsResponseMessage, err := q.sendDNSRequest(host, &dnsRequestMessage)
	querierCache.Query(nameserver.Host)

	if status := domainNSPolicy.CheckNetworkError(err); status != model.NameserverStatusOK {
		if status == model.NameserverStatusTimeout {
			querierCache.Timeout(nameserver.Host)
		}

		domain.Nameservers[index].ChangeStatus(status)

	} else {
		domain.Nameservers[index].ChangeStatus(domainNSPolicy.Run(dnsResponseMessage))
	}

	return true
}

// Check the DS with the domain DNSSEC keys and signatures. You need also to inform the
// UDP max package size supported to pass into firewalls. Many firewalls don't allow
// fragmented UDP packages or UDP packages bigger than 512 bytes. Returns true if DS set
// is done checking and can be saved or false otherwise, that indicates that the domain
// was postponed
func (q *querier) checkDS(domain *model.Domain, index int, udpMaxSize uint16,
	postponedDomains []postponedDomain) bool {

	// Check if the domain has DNSSEC, this system will work with both kinds of domain. So
	// when the domain don't have any DS record we assume that it does not have DNSSEC
	// configured and check only the DNS configuration
	if len(domain.DSSet) == 0 {
		return true
	}

	nameserver := domain.Nameservers[index]
	domainDSPolicy := dspolicy.NewDomainDSPolicy(domain)

	// We are going to request the DNSSEC keys to validate with the DS information that we
	// have from the domain
	var dnsRequestMessage dns.Msg
	dnsRequestMessage.SetQuestion(domain.FQDN, dns.TypeDNSKEY)
	dnsRequestMessage.RecursionDesired = false
	dnsRequestMessage.SetEdns0(udpMaxSize, true)

	host, err := getHost(domain.FQDN, nameserver)
	if err == ErrHostTimeout {
		for index, _ := range domain.DSSet {
			domain.DSSet[index].ChangeStatus(model.DSStatusTimeout)
		}
		return true

	} else if err == ErrHostQPSExceeded {
		postponedDomains = append(postponedDomains, postponedDomain{
			domain: domain,
			index:  index,
		})
		return false
	}

	dnsResponseMessage, err := q.sendDNSRequest(host, &dnsRequestMessage)
	querierCache.Query(nameserver.Host)

	if domainDSPolicy.CheckNetworkError(err) {
		domainDSPolicy.Run(dnsResponseMessage)
	}

	return true
}

// Try to check the postponed domains. Maybe we should have some protection here to avoid
// an almost forever loop when we have a lot of domains with the same nameserver. Returns
// true if domain is done checking and can be saved or false otherwise, that indicates
// that the domain was postponed again
func (q *querier) checkPostponedDomains(postponedDomains []postponedDomain,
	postponed postponedDomain) bool {

	// We only need to check from the nameserver that had a problem (exceeded the QPS), so
	// we are directly calling the checkNameserver method instead of the checkDomain method
	for i := postponed.index; i < len(postponed.domain.Nameservers); i++ {
		if !q.checkNameserver(postponed.domain, i, postponedDomains) {
			return false
		}

		if !q.checkDS(postponed.domain, i, q.UDPMaxSize, postponedDomains) {
			return false
		}
	}

	return true
}

func (q *querier) sendDNSRequest(host string, dnsRequestMessage *dns.Msg) (dnsResponseMessage *dns.Msg, err error) {
	for i := 0; i < q.ConnectionRetries; i++ {
		// For now we ignore the RTT, in the future we can use this for some report
		dnsResponseMessage, _, err = q.client.Exchange(dnsRequestMessage, host)

		// Check if there was a timeout in the connection, if so try again a couple of times
		// just to make it sure that we didn't lose any UDP package
		if err == nil {
			break

		} else if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
			break
		}
	}

	// Message truncated, let's retry using TCP connection. TCP connection will also get the
	// same retries chances of the UDP connection for timeouts because the UDP connection
	// proved in some point that the server is alive
	if err == nil && dnsResponseMessage.Truncated {
		q.client.Net = "tcp"

		// Move back the Net value to empty so that the next package sent by this querier is
		// via UDP connection
		defer func() {
			q.client.Net = ""
		}()

		for i := 0; i < q.ConnectionRetries; i++ {
			// For now we ignore the RTT, in the future we can use this for some report
			dnsResponseMessage, _, err = q.client.Exchange(dnsRequestMessage, host)

			// Check if there was a timeout in the connection, if so try again a couple of times
			// just to make it sure that we didn't lose any UDP package
			if err == nil {
				break

			} else if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
				break
			}
		}
	}

	return
}

// Useful function to retrieve the proper host and port to send the request. The host can
// change because of glue records needs or not. This function alsos resolve hostnames and
// store the addresses in a cache
func getHost(fqdn string, nameserver model.Nameserver) (string, error) {
	// Using cache to store host addresses when there's no glue
	if addresses, err := querierCache.Get(nameserver, fqdn); err == nil || len(addresses) > 0 {
		// Found information in cache, lets use it to speed up the scan. We will try to use an IPv4 from
		// the returned addresses. if we don't find any we will use the first IPv6 address
		for _, address := range addresses {
			if address.To4() != nil {
				return "[" + address.String() + "]:" + strconv.Itoa(DNSPort), nil
			}
		}
		return "[" + addresses[0].String() + "]:" + strconv.Itoa(DNSPort), nil

	} else if err == ErrHostTimeout || err == ErrHostQPSExceeded {
		// Control errors were returned, we need to return them to take an action
		return "", err
	}

	// Error ocurred to retrieve the information from cache. Let's query without using the
	// cache
	return nameserver.Host + ":" + strconv.Itoa(DNSPort), nil
}

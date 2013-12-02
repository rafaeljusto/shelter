package scan

import (
	"github.com/miekg/dns"
	"shelter/model"
	"shelter/scan/dspolicy"
	"shelter/scan/nspolicy"
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
// queries to notify the maximum UDP package size supported in the network
type Querier struct {
	client     dns.Client // Low level DNS client for network checks
	UDPMaxSize uint16     // UDP max package size to pass over firewalls
}

// Return a new Querier object with the necessary fields for the scan filled
func NewQuerier(udpMaxSize uint16, dialTimeout, readTimeout,
	writeTimeout time.Duration) *Querier {

	return &Querier{
		client: dns.Client{
			DialTimeout:  dialTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		UDPMaxSize: udpMaxSize,
	}
}

// Fire a querier that will process domains sent via channel until receives a poison pill
// (nil domain), for go routines control this method receives a wait group, so that the
// main thread can wait for everubody finishs. It also receives the channel where the
// querier will put the domains for the collector save them in database
func (q *Querier) Start(queriers *sync.WaitGroup, domainsToSave chan *model.Domain) chan *model.Domain {
	// Create the communication channel that we are going to listen to retrieve domains, we
	// can store more than one domain in this channel because some queriers can slow down
	// when checking domains with timeouts
	querierChannel := make(chan *model.Domain, querierDomainsQueueSize)

	// Add one more in the group of queriers
	queriers.Add(1)

	go func() {
		for {
			// Retrieve the domain from the channel
			domain := <-querierChannel

			// Detect the poison pill from the dispatcher
			if domain == nil {
				// Tell everyone that we are done!
				queriers.Done()
				return
			}

			q.checkNameserver(domain)
			q.checkDS(domain, q.UDPMaxSize)

			// Send to collector the domain with the new state
			domainsToSave <- domain
		}
	}()

	return querierChannel
}

// Verify the DNS configuration on the nameservers. This method will send a SOA request
// for each nameserver and verify the results
func (q *Querier) checkNameserver(domain *model.Domain) {
	domainNSPolicy := nspolicy.NewDomainNSPolicy(domain)

	// Build message to send the request
	var dnsRequestMessage dns.Msg
	dnsRequestMessage.SetQuestion(domain.FQDN, dns.TypeSOA)
	dnsRequestMessage.RecursionDesired = false

	for index, nameserver := range domain.Nameservers {
		host, err := getHost(domain.FQDN, nameserver)
		if err == HostTimeoutErr {
			domain.Nameservers[index].ChangeStatus(model.NameserverStatusTimeout)
			continue

		} else if err == HostQPSExceededErr {
			// TODO: How are we going to postpone a host query?
			continue
		}

		// For now we ignore the RTT, in the future we can use this for some report
		dnsResponseMessage, _, err := q.client.Exchange(&dnsRequestMessage, host)
		querierCache.Query(nameserver.Host)

		if status := domainNSPolicy.CheckNetworkError(err); status != model.NameserverStatusOK {
			if status == model.NameserverStatusTimeout {
				querierCache.Timeout(nameserver.Host)
			}

			domain.Nameservers[index].ChangeStatus(status)

		} else {
			domain.Nameservers[index].ChangeStatus(domainNSPolicy.Run(dnsResponseMessage))
		}
	}
}

// Check the DS with the domain DNSSEC keys and signatures. You need also to inform the
// UDP max package size supported to pass into firewalls. Many firewalls don't allow
// fragmented UDP packages or UDP packages bigger than 512 bytes
func (q *Querier) checkDS(domain *model.Domain, udpMaxSize uint16) {
	// Check if the domain has DNSSEC, this system will work with both kinds of domain. So
	// when the domain don't have any DS record we assume that it does not have DNSSEC
	// configured and check only the DNS configuration
	if len(domain.DSSet) == 0 {
		return
	}

	domainDSPolicy := dspolicy.NewDomainDSPolicy(domain)

	// We are going to request the DNSSEC keys to validate with the DS information that we
	// have from the domain
	var dnsRequestMessage dns.Msg
	dnsRequestMessage.SetQuestion(domain.FQDN, dns.TypeDNSKEY)
	dnsRequestMessage.RecursionDesired = false
	dnsRequestMessage.SetEdns0(udpMaxSize, true)

	for _, nameserver := range domain.Nameservers {
		host, err := getHost(domain.FQDN, nameserver)
		if err == HostTimeoutErr {
			for index, _ := range domain.DSSet {
				domain.DSSet[index].ChangeStatus(model.DSStatusTimeout)
			}
			continue

		} else if err == HostQPSExceededErr {
			// TODO: How are we going to postpone a host query?
			continue
		}

		// For now we ignore the RTT, in the future we can use this for some report
		dnsResponseMessage, _, err := q.client.Exchange(&dnsRequestMessage, host)

		if !domainDSPolicy.CheckNetworkError(err) || !domainDSPolicy.Run(dnsResponseMessage) {
			break
		}
	}
}

// Useful function to retrieve the proper host and port to send the request. The host can
// change because of glue records needs or not. This function alsos resolve hostnames and
// store the addresses in a cache
func getHost(fqdn string, nameserver model.Nameserver) (string, error) {
	if nameserver.NeedsGlue(fqdn) {
		// Nameserver with glue record. For now we are only checking IPv4 addresses, in the
		// future it would be nice to have an algorithm using both addresses
		return "[" + nameserver.IPv4.String() + "]:" + strconv.Itoa(DNSPort), nil
	}

	// Using cache to store host addresses when there's no glue
	if addresses, err := querierCache.Get(nameserver.Host); err == nil || len(addresses) > 0 {
		// Found information in cache, lets use it to speed up the scan
		return "[" + addresses[0].String() + "]:" + strconv.Itoa(DNSPort), nil

	} else if err == HostTimeoutErr || err == HostQPSExceededErr {
		// Control errors were returned, we need to return them to take an action
		return "", err
	}

	// Error ocurred to retrieve the information from cache. Let's query without using the
	// cache
	return nameserver.Host + ":" + strconv.Itoa(DNSPort), nil
}

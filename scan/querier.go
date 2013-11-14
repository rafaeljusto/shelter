package scan

import (
	"github.com/miekg/dns"
	"net"
	"regexp"
	"shelter/model"
	"strconv"
	"sync"
)

const (
	querierDomainsQueueSize = 10 // Number of domains that can wait in the querier channel
	dnsPort                 = 53 // DNS query port
)

// Querier is responsable for sending the DNS queries to check if the namerservers are
// configured correctly with DNS/DNSSEC
type Querier struct {
	client dns.Client
}

func (q Querier) Start(queriers sync.WaitGroup,
	domainsToSave chan *model.Domain) chan *model.Domain {
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

			q.checkDNS(domain)
			q.checkDNSSEC(domain)

			// Send to collector the domain with the new state
			domainsToSave <- domain
		}
	}()

	return querierChannel
}

// Verify the DNS configuration on the nameservers. This method will send a SOA request
// for each nameserver and verify the results
func (q Querier) checkDNS(domain *model.Domain) {
	// Variable used to check if all nameservers have the same zone
	soaVersion := uint32(0)

	// Build message to send the request
	var dnsRequestMessage dns.Msg
	dnsRequestMessage.SetQuestion(domain.FQDN, dns.TypeSOA)
	dnsRequestMessage.RecursionDesired = false

	for index, nameserver := range domain.Nameservers {
		host := ""
		if nameserver.NeedsGlue(domain.FQDN) {
			// Nameserver with glue record. For now we are only checking IPv4 addresses, in the
			// future it would be nice to have an algorithm using both addresses
			host = nameserver.IPv4.String() + ":" + strconv.Itoa(dnsPort)

		} else {
			// Using cache to store host addresses when there's no glue
			if addresses, err := querierCache.Get(nameserver.Host); err != nil || len(addresses) == 0 {
				// Error ocurred to retrieve the information from cache. Let's query without using
				// the cache
				host = nameserver.Host + ":" + strconv.Itoa(dnsPort)
			} else {
				// Found information in cache, lets use it to speed up the scan
				host = addresses[0].String() + ":" + strconv.Itoa(dnsPort)
			}
		}

		// For now we ignore the RTT, in the future we can use this for some report
		dnsResponseMessage, _, err := q.client.Exchange(&dnsRequestMessage, host)

		if err != nil {
			// Check for transport errors. Some erros we don't have a specific structure to
			// detect it, so we are going to catch the problem analyzing the error message

			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				domain.Nameservers[index].ChangeStatus(model.NameserverStatusTimeout)

			} else {
				if match, _ := regexp.MatchString(".*lookup.*", err.Error()); match {
					domain.Nameservers[index].ChangeStatus(model.NameserverStatusUnknownHost)

				} else if match, _ := regexp.MatchString(".*connection refused.*", err.Error()); match {
					domain.Nameservers[index].ChangeStatus(model.NameserverStatusConnectionRefused)

				} else {
					domain.Nameservers[index].ChangeStatus(model.NameserverStatusError)
				}
			}

			continue
		}

		// Checking if we have CNAME in apex. This is not allowed according to RFC 1033, CNAME
		// cannot exist with other records, and SOA record is mandatory in apex.

		cnameError := false
		for _, rr := range dnsResponseMessage.Answer {
			if rr.Header().Name == domain.FQDN && rr.Header().Rrtype == dns.TypeCNAME {
				domain.Nameservers[index].ChangeStatus(model.NameserverStatusCanonicalName)
				cnameError = true
				break
			}
		}

		if cnameError {
			continue
		}

		// We now verify the DNS response package and according to the DNS header result code
		// we set the properly status

		if dnsResponseMessage.MsgHdr.Rcode == dns.RcodeSuccess {
			var soaRR *dns.SOA
			foundSOA := false

			// Check if the SOA resource record exists in the response
			for _, rr := range dnsResponseMessage.Answer {
				if rr.Header().Rrtype == dns.TypeSOA {
					soaRR, _ = rr.(*dns.SOA)
					foundSOA = true
					break
				}
			}

			if !foundSOA {
				domain.Nameservers[index].ChangeStatus(model.NameserverStatusUnknownDomainName)
				continue
			}

			if !dnsResponseMessage.Authoritative {
				domain.Nameservers[index].ChangeStatus(model.NameserverStatusNoAuthority)
				continue
			}

			if soaRR != nil {
				// Find the SOA record to check the version of the zone. If it's the first
				// nameserver that we are checking (version 0) we don't need to compare

				if soaVersion == 0 {
					soaVersion = soaRR.Serial

				} else if soaVersion != soaRR.Serial {
					domain.Nameservers[index].ChangeStatus(model.NameserverStatusNotSynchronized)
					continue
				}
			}

		} else if dnsResponseMessage.MsgHdr.Rcode == dns.RcodeRefused {
			domain.Nameservers[index].ChangeStatus(model.NameserverStatusQueryRefused)
			continue

		} else if dnsResponseMessage.MsgHdr.Rcode == dns.RcodeServerFailure {
			domain.Nameservers[index].ChangeStatus(model.NameserverStatusServerFailure)
			continue

		} else if dnsResponseMessage.MsgHdr.Rcode == dns.RcodeNameError {
			domain.Nameservers[index].ChangeStatus(model.NameserverStatusUnknownDomainName)
			continue

		} else {
			domain.Nameservers[index].ChangeStatus(model.NameserverStatusError)
			continue
		}

		domain.Nameservers[index].ChangeStatus(model.NameserverStatusOK)
	}
}

func (q Querier) checkDNSSEC(domain *model.Domain) {
	// Check if the domain has DNSSEC
	if len(domain.DSSet) == 0 {
		return
	}

	var dnsRequestMessage dns.Msg
	dnsRequestMessage.SetQuestion(domain.FQDN, dns.TypeDNSKEY)
	dnsRequestMessage.RecursionDesired = false

	for _, nameserver := range domain.Nameservers {
		host := ""
		if nameserver.NeedsGlue(domain.FQDN) {
			host = nameserver.Host + ":" + strconv.Itoa(dnsPort)
		} else {
			host = nameserver.IPv4.String() + ":" + strconv.Itoa(dnsPort)
		}

		_, _, err := q.client.Exchange(&dnsRequestMessage, host)
		if err != nil {
			for index, _ := range domain.DSSet {
				domain.DSSet[index].ChangeStatus(model.DSStatusDNSError)
			}
			break
		}

		// TODO
	}
}

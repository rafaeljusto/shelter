package scan

import (
	"github.com/miekg/dns"
	"net"
	"regexp"
	"shelter/model"
	"strconv"
	"sync"
	"time"
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
		// cannot exist with other records, and SOA record is mandatory in apex

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

			// Check if the SOA resource record exists in the response
			if rr := filterFirstRR(dnsResponseMessage.Answer, dns.TypeSOA); rr != nil {
				soaRR, _ = rr.(*dns.SOA)

			} else {
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
	// Check if the domain has DNSSEC, this system will work with both kinds of domain. So
	// when the domain don't have any DS record we assume that it does not have DNSSEC
	// configured and check only the DNS configuration
	if len(domain.DSSet) == 0 {
		return
	}

	// We are going to request the DNSSEC keys to validate with the DS information that we
	// have from the domain
	var dnsRequestMessage dns.Msg
	dnsRequestMessage.SetQuestion(domain.FQDN, dns.TypeDNSKEY)
	dnsRequestMessage.RecursionDesired = false
	dnsRequestMessage.SetEdns0(4096, true) // TODO: UDP max size must be configurable

	for _, nameserver := range domain.Nameservers {
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
			// We can have timeouts only for DNSSEC requests, because usually the response is
			// bigger and firewalls are not configured for big UDP packages, or for DNS over TCP
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				for index, _ := range domain.DSSet {
					domain.DSSet[index].ChangeStatus(model.DSStatusTimeout)
				}

			} else {
				// Other types of network errors are not a specific problem of DNSSEC
				// configuration, so let's just set a status for the user to fix the DNS
				// configuration to make the DNSSEC configuration check possible
				for index, _ := range domain.DSSet {
					domain.DSSet[index].ChangeStatus(model.DSStatusDNSError)
				}
			}

			break
		}

		if dnsResponseMessage.Rcode != dns.RcodeSuccess ||
			!dnsResponseMessage.MsgHdr.Authoritative {
			// Authority errors are not a specific problem of DNSSEC configuration, so let's
			// just set a status for the user to fix the DNS configuration to make the DNSSEC
			// configuration check possible
			for index, _ := range domain.DSSet {
				domain.DSSet[index].ChangeStatus(model.DSStatusDNSError)
			}

			break
		}

		// Get all DNSSEC public keys
		dnskeys := filterRRs(dnsResponseMessage.Answer, dns.TypeDNSKEY)

		// Get all signatures from the DNS response
		rrsigs := filterRRs(dnsResponseMessage.Answer, dns.TypeRRSIG)

		for index, ds := range domain.DSSet {
			// Find the DNSSEC public key related to the DS
			var selectedDNSKEY *dns.DNSKEY
			for _, rr := range dnskeys {
				dnskey, ok := rr.(*dns.DNSKEY)
				if !ok {
					continue
				}

				if dnskey.KeyTag() == ds.Keytag {
					selectedDNSKEY = dnskey
					break
				}
			}

			// Find the signature of the DNSSEC key that signed the keyset
			var selectedRRSIG *dns.RRSIG
			for _, rr := range rrsigs {
				rrsig, ok := rr.(*dns.RRSIG)
				if !ok {
					continue
				}

				if rrsig.KeyTag == ds.Keytag {
					selectedRRSIG = rrsig
					break
				}
			}

			if selectedDNSKEY == nil {
				domain.DSSet[index].ChangeStatus(model.DSStatusNoKey)
				continue
			}

			// Check if the DNSSEC key related to the DS has the security entry point. Check
			// RFCs 3755 and 4034
			if (selectedDNSKEY.Flags & (1 << 15)) == 0 {
				domain.DSSet[index].ChangeStatus(model.DSStatusNoSEP)
				continue
			}

			if selectedRRSIG == nil {
				domain.DSSet[index].ChangeStatus(model.DSStatusNoSignature)
				continue
			}

			// We store the DS expiration date to alert clients whenever an expiration date is
			// near. There's no status in DS to define a near expiration state, because this
			// isn't a problem
			domain.DSSet[index].ExpiresAt = time.Unix(int64(selectedRRSIG.Expiration), 0)

			// Check signature expiration
			if ds.ExpiresAt.Before(time.Now()) {
				domain.DSSet[index].ChangeStatus(model.DSStatusExpiredSignature)
				continue
			}

			// Check signature consistency
			if err := selectedRRSIG.Verify(selectedDNSKEY, dnskeys); err != nil {
				domain.DSSet[index].ChangeStatus(model.DSStatusSignatureError)
				continue
			}

			// Check DNSKEY hash is the same of the DS digest
			if selectedDNSKEY.ToDS(int(ds.DigestType)).Digest != ds.Digest {
				domain.DSSet[index].ChangeStatus(model.DSStatusNoKey)
				continue
			}

			domain.DSSet[index].ChangeStatus(model.DSStatusOK)
		}
	}
}

// Useful function to retrieve all records of a specific type from the DNS response
// message. We assume that the resource record is an instance of the specific type based
// on the Rrtype attribute
func filterRRs(rrs []dns.RR, rrType uint16) []dns.RR {
	var filtered []dns.RR
	for _, rr := range rrs {
		if rr.Header().Rrtype == rrType {
			filtered = append(filtered, rr)
		}
	}
	return filtered
}

// Useful function to return the first occurence of a resource record of a specific type.
// This method is faster than filterRRs when we are interested in only one record (like
// SOA)
func filterFirstRR(rrs []dns.RR, rrType uint16) dns.RR {
	for _, rr := range rrs {
		if rr.Header().Rrtype == rrType {
			return rr
		}
	}
	return nil
}

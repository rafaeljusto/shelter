package scan

import (
	"github.com/miekg/dns"
	"net"
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

func (q Querier) checkDNS(domain *model.Domain) {
	var dnsRequestMessage dns.Msg
	dnsRequestMessage.SetQuestion(domain.FQDN, dns.TypeSOA)
	dnsRequestMessage.RecursionDesired = false

	for index, nameserver := range domain.Nameservers {
		host := ""
		if nameserver.NeedsGlue(domain.FQDN) {
			host = nameserver.Host + ":" + strconv.Itoa(dnsPort)
		} else {
			host = nameserver.IPv4.String() + ":" + strconv.Itoa(dnsPort)
		}

		dnsResponseMessage, _, err := q.client.Exchange(&dnsRequestMessage, host)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				domain.Nameservers[index].ChangeStatus(model.NameserverStatusTimeout)
				continue

			} else {
				// TODO
			}
		}

		if !dnsResponseMessage.Authoritative {
			domain.Nameservers[index].ChangeStatus(model.NameserverStatusNoAuthority)
			continue
		}

		// TODO
	}
}

func (q Querier) checkDNSSEC(domain *model.Domain) {
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

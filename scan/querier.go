package scan

import (
	"github.com/miekg/dns"
	"shelter/model"
	"sync"
)

const (
	QuerierDomainsQueueSize = 10 // Number of domains that can wait in the querier channel
)

// Querier is responsable for sending the DNS queries to check if the namerservers are
// configured correctly with DNS/DNSSEC
type Querier struct {
}

func (q Querier) Start(queriers sync.WaitGroup) chan *model.Domain {
	// Create the communication channel that we are going to listen to retrieve domains, we
	// can store more than one domain in this channel because some queriers can slow down
	// when checking domains with timeouts
	querierChannel := make(chan *model.Domain, QuerierDomainsQueueSize)

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

			var dnsMessage dns.Msg
			dnsMessage.SetQuestion(domain.FQDN, dns.TypeSOA)
			dnsMessage.RecursionDesired = false

			var client dns.Client

			for _, nameserver := range domain.Nameservers {
				// TODO: Glue!
				client.Exchange(&dnsMessage, nameserver.Host)
			}

			// TODO
		}
	}()

	return querierChannel
}

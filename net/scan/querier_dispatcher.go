package scan

import (
	"shelter/model"
	"sync"
	"time"
)

// QuerierDispatcher is responsable for receiving the domains to query and redistribute
// them to the queriers. This allows parallels DNS queriers and speed up the scan system.
// The object has many attributes to control the scan query performance, allow packages to
// pass network firewall rules, determinate the amount of time that the queriers are going
// to wait in network operations when there's no answer and determinate the number of
// concurrently go routines that will resolve the domains
type QuerierDispatcher struct {
	NumberOfQueriers  int           // Number of queriers to concurrently check the domains
	DomainsBufferSize int           // Size of the domains to save channel
	UDPMaxSize        uint16        // UDP max package size to pass over firewalls
	DialTimeout       time.Duration // Timeout while connecting to a server
	ReadTimeout       time.Duration // Timeout while waiting for a response
	WriteTimeout      time.Duration // Timeout to write a query to the DNS server
	ConnectionRetries int           // Number of retries before setting timeout
}

// Return a new QuerierDispatcher object with the necessary fields for the scan filled
func NewQuerierDispatcher(
	numberOfQueriers,
	domainsBufferSize int,
	udpMaxSize uint16,
	dialTimeout, readTimeout, writeTimeout time.Duration,
	connectionRetries int,
) *QuerierDispatcher {

	return &QuerierDispatcher{
		NumberOfQueriers:  numberOfQueriers,
		DomainsBufferSize: domainsBufferSize,
		UDPMaxSize:        udpMaxSize,
		DialTimeout:       dialTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		ConnectionRetries: connectionRetries,
	}
}

// This is the method that start the querier dispatcher and the queriers. It is
// asynchronous and will ends after receiving the poison pill from the injector. It
// receives a object to sinalize to the main thread the end and a channel that tells the
// domains to query, sent by the injector
func (q *QuerierDispatcher) Start(scanGroup *sync.WaitGroup,
	domainsToQueryChannel chan *model.Domain) chan *model.Domain {

	// Create the output channel used for each querier to add the result for the collector,
	// the poison pill is the nil domain object
	domainsToSaveChannel := make(chan *model.Domain, q.DomainsBufferSize)

	// Allocate the number of queriers
	queriersChannels := make([]chan *model.Domain, q.NumberOfQueriers)

	// Create a sync group to control the end of all queriers. The dispatcher can only ends
	// after all queriers finished their jobs
	var queriers sync.WaitGroup

	// Initialize each querier
	for index, _ := range queriersChannels {
		querier := newQuerier(
			q.UDPMaxSize,
			q.DialTimeout,
			q.ReadTimeout,
			q.WriteTimeout,
			q.ConnectionRetries,
		)

		queriersChannels[index] = querier.start(&queriers, domainsToSaveChannel)
	}

	// Add one more to the group of scan go routines
	scanGroup.Add(1)

	go func() {
		index := 0

		for {
			// Retrieve a domain from the injector
			domain := <-domainsToQueryChannel

			// Detect the poinson pill from the injector
			if domain == nil {
				// Finish all queriers sending a nil domain
				for _, queriersChannels := range queriersChannels {
					queriersChannels <- nil
				}

				// Wait for queriers to finish
				queriers.Wait()

				// Send the poison pill to the collector
				domainsToSaveChannel <- nil

				scanGroup.Done()
				return
			}

			// We are going to use a round robin strategy to distribute the domains for the
			// queriers, so if we reach the last channel, go back to the first one
			if index >= len(queriersChannels) {
				index = 0
			}

			// Send to the querier a domain
			queriersChannels[index] <- domain

			// Move to the next querier
			index += 1
		}
	}()

	return domainsToSaveChannel
}

package scan

import (
	"shelter/dao"
	"shelter/model"
	"sync"
)

// QuerierDispatcher is responsable for receiving the domains to query and redistribute
// them to the queriers. This allows parallels DNS queriers and speed up the scan system
type QuerierDispatcher struct {
}

// This is the method that start the querier dispatcher and the queriers. It is
// asynchronous and will ends after receiving the poison pill from the injector
func (q QuerierDispatcher) Start(domainsToQueryChannel chan dao.DomainResult,
	domainsBufferSize, numberOfQueriers int) chan *model.Domain {

	// Create the output channel used for each querier to add the result for the collector,
	// the poison pill is the nil domain object
	domainsToSave := make(chan *model.Domain, domainsBufferSize)

	// Allocate the number of queriers
	queriersChannels := make([]chan *model.Domain, numberOfQueriers)

	// Create a sync group to control the end of all queriers. The dispatcher can only ends
	// after all queriers finished their jobs
	var queriers sync.WaitGroup

	// Initialize each querier
	for index, _ := range queriersChannels {
		queriersChannels[index] = Querier{}.Start(queriers)
	}

	go func() {
		index := 0

		for {
			// Retrieve a domain from the injector
			domainResult := <-domainsToQueryChannel

			// Detect the poinson pill from the injector
			if domainResult.Domain == nil || domainResult.Error != nil {
				// Finish all queriers sending a nil domain
				for _, queriersChannels := range queriersChannels {
					queriersChannels <- nil
				}

				// Wait for queriers and move out
				queriers.Wait()
				return
			}

			// We are going to use a round robin strategy to distribute the domains for the
			// queriers, so if we reach the last channel, go back to the first one
			if index >= len(queriersChannels) {
				index = 0
			}

			// Send to the querier a domain
			queriersChannels[index] <- domainResult.Domain
		}
	}()

	return domainsToSave
}

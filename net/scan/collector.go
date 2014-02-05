package scan

import (
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/model"
	"labix.org/v2/mgo"
	"sync"
)

// Collector is responsable for persisting all domains with their new status into the
// database. For faster approach the collector waits until it has many domains to save
// them at once in the database
type Collector struct {
	Database   *mgo.Database // Low level database connection
	SaveAtOnce int           // Number of domains to save at once
}

// Return a new Collector object with the necessary fields for the scan filled
func NewCollector(database *mgo.Database, saveAtOnce int) *Collector {
	return &Collector{
		Database:   database,
		SaveAtOnce: saveAtOnce,
	}
}

// This method is the last part of the scan, when the new state of the domain object is
// persisted back to the database. It receives a go routine control group to sinalize to
// the main thread when the scan ends, a domain channel to receive each domain that need
// to be save and an error channel to send back all errors while persisting the data. It
// was created to be asynchronous and finish after receiving a poison pill from querier
// dispatcher
func (c *Collector) Start(scanGroup *sync.WaitGroup,
	domainsToSaveChannel chan *model.Domain, errorsChannel chan error) {

	// Add one more to the group of scan go routines
	scanGroup.Add(1)

	go func() {
		// Initialize Domain DAO using injected database connection
		domainDAO := dao.DomainDAO{
			Database: c.Database,
		}

		// Add a safety check to avoid an infinite loop
		if c.SaveAtOnce == 0 {
			c.SaveAtOnce = 1
		}

		finished := false
		nameserverStatistics := make(map[string]uint64)
		dsStatistics := make(map[string]uint64)

		for {
			// Using make for faster allocation
			domains := make([]*model.Domain, 0, c.SaveAtOnce)

			for i := 0; i < c.SaveAtOnce; i++ {
				domain := <-domainsToSaveChannel

				// Detect poison pill. We don't return from function here because we can still
				// have some domains to save in the domains array
				if domain == nil {
					finished = true
					break
				}

				// Count this domain for the scan information to estimate the scan progress
				model.FinishAnalyzingDomainForScan(len(domain.DSSet) > 0)

				// Keep track of nameservers statistics
				for _, nameserver := range domain.Nameservers {
					status := model.NameserverStatusToString(nameserver.LastStatus)
					nameserverStatistics[status] += 1
				}

				// Keep track of DS statistics
				for _, ds := range domain.DSSet {
					status := model.DSStatusToString(ds.LastStatus)
					dsStatistics[status] += 1
				}

				domains = append(domains, domain)
			}

			domainsResults := domainDAO.SaveMany(domains)
			for _, domainResult := range domainsResults {
				if domainResult.Error != nil {
					// Error channel should have a buffer or this will block the collector until
					// someone check this error. One question here is that we are returning the
					// error, but not telling wich domain got the error, we should improve the error
					// communication system between the go routines
					errorsChannel <- domainResult.Error
				}
			}

			// Now that everything is done, check if we received a poison pill
			if finished {
				model.StoreStatisticsOfTheScan(nameserverStatistics, dsStatistics)
				scanGroup.Done()
				return
			}
		}
	}()
}

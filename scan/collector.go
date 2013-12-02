package scan

import (
	"labix.org/v2/mgo"
	"shelter/dao"
	"shelter/model"
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

// This method is asynchronous and will finish after receiving a poison pill from querier
// dispatcher
func (c *Collector) Start(domainsToSave chan *model.Domain, errorsChannel chan error) {
	go func() {
		// Initialize Domain DAO using injected database connection
		domainDAO := dao.DomainDAO{
			Database: c.Database,
		}

		// Add a safety check to avoid an infinite loop
		if c.SaveAtOnce == 0 {
			c.SaveAtOnce = 1
		}

		for {
			finished := false

			var domains []*model.Domain
			for i := 0; i < c.SaveAtOnce; i++ {
				domain := <-domainsToSave

				// Detect poison pill. We don't return from function here because we can still
				// have some domains to save in the domains array
				if domain == nil {
					finished = true
					break
				}

				domains = append(domains, domain)
			}

			domainsResults := domainDAO.SaveMany(domains)
			for _, domainResult := range domainsResults {
				if domainResult.Error != nil {
					errorsChannel <- domainResult.Error
				}
			}

			if finished {
				break
			}
		}
	}()
}

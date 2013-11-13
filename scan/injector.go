package scan

import (
	"labix.org/v2/mgo"
	"shelter/dao"
)

// Injector is responsable for selecting all domains that are going to be checked. While
// selecting the domains the injector will add to a channel, so that the querier can start
// immediately
type Injector struct {
	Database *mgo.Database
}

// Method that starts the injector job, retrieving the data from the database and adding
// the same data into a channel for a querier start sending DNS requests. There's also
// four more parameters to define the size of the channel to add domains to query, the
// maximum number of days to verify a domain configured correctly with DNS/DNSSEC, the
// maximum number of days to verify a domain with problems, and the number of days to
// alert for DNSSEC signatures that are near from the expiration date. This method is
// asynchronous and will finish sending a poison pill (error or nil domain) to indicate to
// the querier that there are no more domains
func (i Injector) Start(domainsBufferSize, maxOKVerificationDays, maxErrorVerificationDays,
	maxExpirationAlertDays int) (chan dao.DomainResult, error) {

	// Create the output channel where we are going to add the domains retrieved from the
	// database for the querier
	domainsToQueryChannel := make(chan dao.DomainResult, domainsBufferSize)

	// Initialize Domain DAO using injected database connection
	domainDAO := dao.DomainDAO{
		Database: i.Database,
	}

	// Load all domains from database to begin the scan
	domainChannel, err := domainDAO.FindAll()
	if err != nil {
		return nil, err
	}

	go func() {
		// Dispatch the asynchronous part of the method
		for {
			// Get domain from the database (one-by-one)
			domainResult := <-domainChannel

			// Problem detected while retrieving a domain or we don't have domains anymore, send
			// the poison pill to alert the querier
			if domainResult.Error != nil || domainResult.Domain == nil {
				domainsToQueryChannel <- domainResult
				return
			}

			// The logic that decides if a domain is going to be a part of this scan or not is
			// inside the domain object for better unit testing
			if domainResult.Domain.ShouldBeScanned(maxOKVerificationDays,
				maxErrorVerificationDays, maxExpirationAlertDays) {
				// Send to the querier
				domainsToQueryChannel <- domainResult
			}
		}
	}()

	return domainsToQueryChannel, nil
}

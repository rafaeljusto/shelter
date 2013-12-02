package scan

import (
	"labix.org/v2/mgo"
	"shelter/dao"
	"shelter/model"
)

// Injector is responsable for selecting all domains that are going to be checked. While
// selecting the domains the injector will add to a channel, so that the querier can start
// immediately
type Injector struct {
	Database                 *mgo.Database // Low level database connection
	DomainsBufferSize        int           // Size of the domains to query channel
	MaxOKVerificationDays    int           // Maximum number of days to verify a domain configured correctly with DNS/DNSSEC
	MaxErrorVerificationDays int           // Maximum number of days to verify a domain with problems
	MaxExpirationAlertDays   int           // Number of days to alert for DNSSEC signatures that are near from the expiration date
}

// Return a new Injector object with the necessary fields for the scan filled
func NewInjector(database *mgo.Database, domainsBufferSize, maxOKVerificationDays,
	maxErrorVerificationDays, maxExpirationAlertDays int) *Injector {

	return &Injector{
		Database:                 database,
		DomainsBufferSize:        domainsBufferSize,
		MaxOKVerificationDays:    maxOKVerificationDays,
		MaxErrorVerificationDays: maxErrorVerificationDays,
		MaxExpirationAlertDays:   maxExpirationAlertDays,
	}
}

// Method that starts the injector job, retrieving the data from the database and adding
// the same data into a channel for a querier start sending DNS requests. There is only
// one parameter to define a channel to report errors while loading the data. This method
// is asynchronous and will finish sending a poison pill (error or nil domain) to indicate
// to the querier that there are no more domains
func (i *Injector) Start(errorsChannel chan error) chan *model.Domain {

	// Create the output channel where we are going to add the domains retrieved from the
	// database for the querier
	domainsToQueryChannel := make(chan *model.Domain, i.DomainsBufferSize)

	go func() {
		// Initialize Domain DAO using injected database connection
		domainDAO := dao.DomainDAO{
			Database: i.Database,
		}

		// Load all domains from database to begin the scan
		domainChannel, err := domainDAO.FindAll()

		// Low level error was detected. No domain was processed yet, but we still need to
		// shutdown the querier and by consequence the collector, so we send back the error
		// and add the poison pill
		if err != nil {
			errorsChannel <- err
			domainsToQueryChannel <- nil
			return
		}

		// Dispatch the asynchronous part of the method
		for {
			// Get domain from the database (one-by-one)
			domainResult := <-domainChannel

			// Send back the error to the caller thread. We don't log the error here directly
			// into the log interface because sometimes we want to do something when an error
			// occurs, like in a test enviroment
			if domainResult.Error != nil {
				errorsChannel <- domainResult.Error
			}

			// Problem detected while retrieving a domain or we don't have domains anymore, send
			// the poison pill to alert the querier
			if domainResult.Error != nil || domainResult.Domain == nil {
				domainsToQueryChannel <- nil
				return
			}

			// The logic that decides if a domain is going to be a part of this scan or not is
			// inside the domain object for better unit testing
			if domainResult.Domain.ShouldBeScanned(i.MaxOKVerificationDays,
				i.MaxErrorVerificationDays, i.MaxExpirationAlertDays) {
				// Send to the querier
				domainsToQueryChannel <- domainResult.Domain
			}
		}
	}()

	return domainsToQueryChannel
}

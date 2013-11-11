package scan

import (
	"labix.org/v2/mgo"
	"math"
	"math/rand"
	"shelter/dao"
	"time"
)

// Injector is responsable for selecting all domains that are going to be checked. While
// selecting the domains the injector will add to a channel, so that the querier can start
// immediately
type Injector struct {
	Database *mgo.Database
}

// Method that starts the injector job, retrieving the data from the database and adding
// the same data into a channel for a querier start sending DNS requests. There's also
// three more parameters to define the maximum number of days to verify a domain
// configured correctly with DNS/DNSSEC, a domain with problems and to alert for DNSSEC
// signatures that are near from the expiration date. This method is asynchronous and will
// finish sending a poison pill (error or nil domain) to indicate to the querier that
// there are no more domains
func (i Injector) Start(domainsToQueryChannel chan dao.DomainResult,
	maxOKVerificationDays, maxErrorVerificationDays, maxExpirationAlertDays int) error {
	// Initialize Domain DAO using injected database connection
	domainDAO := dao.DomainDAO{
		Database: i.Database,
	}

	// Load all domains from database to begin the scan
	domainChannel, err := domainDAO.FindAll()
	if err != nil {
		return err
	}

	go func() {
		// Start the random with seed only once, we are going to reuse it on every domain
		// check to randomly select a domain to the scan or not. As we are using the current
		// nanosecond, we have the entropy necessary to be really random
		random := rand.New(rand.NewSource(time.Now().UnixNano()))

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

			domain := domainResult.Domain

			var maxDays int

			// When all nameservers are OK the domain has less chance to be selected for the
			// scan, because the random range will be bigger
			if domain.AllNameserversOK() && domain.AllDSSetOK() {
				maxDays = maxOKVerificationDays
			} else {
				maxDays = maxErrorVerificationDays
			}

			// The longer the last check occurred, better are the chances to select the domain
			// for the scan
			daysSinceLastCheck := domain.DaysSinceLastCheck()
			selectedDay := 1 + (random.Int() * maxDays / math.MaxInt64)

			// If the domain is configured with DNSSEC and is near the expiration date, we
			// must check even if it's not selected by the random algorithm
			if !domain.IsNearDNSSECExpirationDate(maxExpirationAlertDays) &&
				selectedDay > daysSinceLastCheck {
				continue
			}

			// Send to the querier
			domainsToQueryChannel <- domainResult
		}
	}()

	return nil
}

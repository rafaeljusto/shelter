package scan

import (
	"labix.org/v2/mgo"
	"math"
	"math/rand"
	"shelter/dao"
	"shelter/model"
	"time"
)

// Injector is responsable for selecting all domains that are going to be checked. While
// selecting the domains the injector will add to a channel, so that the querier can start
// immediately
type Injector struct {
	Database *mgo.Database
}

// Method that starts the injector job, retrieving the data and adding the same data into
// a channel
func (i Injector) Start(queue chan model.Domain) error {
	// Strategy:
	//
	// For domains that have all DNS OK and the DS OK, we check at least once a week (7
	// days), for the other cases we check at least once every 2 days. The longer the last
	// check occurred, better are the chances to select the domain for the scan.
	//
	// TODO: Dates must be configurable

	var domainDAO dao.DomainDAO

	// Load all domains from database to begin the scan
	domainChannel, err := domainDAO.FindAll()
	if err != nil {
		return err
	}

	// Start the random with seed only once, we are going to reuse it on every domain
	// check to randomly a domain to the scan or not. As we are using the current
	// nanosecond we have the entropy necessary to be really random
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		// Get domain from the database (one-by-one)
		domain := <-domainChannel

		var maxDays int

		// When all nameservers are OK the domain has less chance to be selected for the
		// scan, because the random range will be bigger
		if domain.AllNameserversOK() && domain.AllDSSetOK() {
			maxDays = 7
		} else {
			maxDays = 2
		}

		daysSinceLastCheck := domain.DaysSinceLastCheck()
		selectedDay := 1 + (random.Int() * maxDays / math.MaxInt64)

		// If the domain is configured with DNSSEC and is near the expiration date, we
		// must check even if it's not selected by the random algorithm
		if !domain.IsNearDNSSECExpirationDate(10) && selectedDay > daysSinceLastCheck {
			continue
		}

		// Send to the querier
		queue <- domain
	}

	return nil
}

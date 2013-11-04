package model

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"math"
	"net/mail"
	"time"
)

// Domain stores all the necessary information for validating the DNS and DNSSEC. It also
// stores information to alert the domain's owners about the problems.
type Domain struct {
	Id          bson.ObjectId   `bson:"_id"` // Database identification
	FQDN        string          // Actual domain name
	Nameservers []Nameserver    // Nameservers that asnwer with authority for this domain
	DSSet       []DS            // Records for the DNS tree chain of trust
	Owners      []*mail.Address // E-mails that will be alerted on any problem
}

// Check if all nameservers are configured correctly with DNS
func (d Domain) AllNameserversOK() bool {
	for i := 0; i < len(d.Nameservers); i++ {
		if d.Nameservers[i].LastStatus != NameserverStatusOK {
			return false
		}
	}

	return true
}

// Check if all DS set is configured correctly with DNSSEC
func (d Domain) AllDSSetOK() bool {
	for i := 0; i < len(d.DSSet); i++ {
		if d.DSSet[i].LastStatus != DSStatusOK {
			return false
		}
	}

	return true
}

// DaysSinceLastCheck returns the number of days since the last check in the nameservers
// or in the DS records
func (d Domain) DaysSinceLastCheck() int {
	// We need to retrieve the most recent check date from the objects of the domain
	var lastCheckAt time.Time

	// Check within the nameservers
	for i := 0; i < len(d.Nameservers); i++ {
		if d.Nameservers[i].LastCheckAt.After(lastCheckAt) {
			lastCheckAt = d.Nameservers[i].LastCheckAt
		}
	}

	// Check within the DS set
	for i := 0; i < len(d.DSSet); i++ {
		if d.DSSet[i].LastCheckAt.After(lastCheckAt) {
			lastCheckAt = d.DSSet[i].LastCheckAt
		}
	}

	// Now that we have the most recent date, let's see how long it was. For better
	// precision lets round the duration to convert it to integer
	hours := math.Ceil(time.Since(lastCheckAt).Hours())
	return int(hours) / int(24*time.Hour.Hours())
}

// Check the DS set to see if the expiration date of the DNSKEYs signatures are near. The
// alert period is defined by the parameter daysBefore, that is the number of days before
// the expiration date that we will consider near
func (d Domain) IsNearDNSSECExpirationDate(daysBefore int) bool {
	// Lets look for the oldest expiration date of the DS set, the it's probably the most
	// problematic one
	var expiresAt time.Time

	for i := 0; i < len(d.DSSet); i++ {
		if expiresAt.IsZero() || d.DSSet[i].ExpiresAt.Before(expiresAt) {
			expiresAt = d.DSSet[i].ExpiresAt
		}
	}

	// When there's no DS, we don't have an expiration date and shouldn't care about it
	if expiresAt.IsZero() {
		return false
	}

	// Now we can check if this expiration date is already in the alert period. The alert
	// period is defined by the function parameter in days
	daysInHours := fmt.Sprintf("%dh", daysBefore*24)
	expirationMinimumAlertPeriod, _ := time.ParseDuration(daysInHours)
	return time.Now().Add(expirationMinimumAlertPeriod).After(expiresAt)
}

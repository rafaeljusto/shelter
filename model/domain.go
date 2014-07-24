// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"gopkg.in/mgo.v2/bson"
	"math"
	"math/rand"
	"time"
)

var (
	// Start the random with seed only once, we are going to reuse it on every domain
	// check to randomly select a domain to the scan or not. As we are using the current
	// nanosecond, we have the entropy necessary to be really random
	domainRandomSelector = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Domain stores all the necessary information for validating the DNS and DNSSEC. It also
// stores information to alert the domain's owners about the problems
type Domain struct {
	Id             bson.ObjectId `bson:"_id"` // Database identification
	Revision       int           // Version of the object
	LastModifiedAt time.Time     // Last time the object was modified
	FQDN           string        // Actual domain name
	Nameservers    []Nameserver  // Nameservers that asnwer with authority for this domain
	DSSet          []DS          // Records for the DNS tree chain of trust
	Owners         []Owner       // Responsables for the domains that will receive alerts
}

// ShouldBeScanned method is responsable for telling if the domain can be scanned or not
// using some business rules based on the last verification, nameservers and DS status and
// DNSSEC signatures expiration date. For now this method is used by scan injector
func (d Domain) ShouldBeScanned(maxOKVerificationDays,
	maxErrorVerificationDays, maxExpirationAlertDays int) bool {

	var maxDays int

	// When all nameservers are OK the domain has less chance to be selected for the
	// scan, because the random range will be bigger
	if d.allNameserversOK() && d.allDSSetOK() {
		maxDays = maxOKVerificationDays
	} else {
		maxDays = maxErrorVerificationDays
	}

	// The longer the last check occurred, better are the chances to select the domain
	// for the scan
	daysSinceLastCheck := d.daysSinceLastCheck()
	selectedDay := 1 + (domainRandomSelector.Int() * maxDays / math.MaxInt64)

	// If the domain is configured with DNSSEC and is near the expiration date, we
	// must check even if it's not selected by the random algorithm
	if !d.isNearDNSSECExpirationDate(maxExpirationAlertDays) &&
		selectedDay > daysSinceLastCheck {
		return false
	}

	return true
}

// Check if all nameservers are configured correctly with DNS
func (d Domain) allNameserversOK() bool {
	for i := 0; i < len(d.Nameservers); i++ {
		if d.Nameservers[i].LastStatus != NameserverStatusOK {
			return false
		}
	}

	return true
}

// Check if all DS set is configured correctly with DNSSEC
func (d Domain) allDSSetOK() bool {
	for i := 0; i < len(d.DSSet); i++ {
		if d.DSSet[i].LastStatus != DSStatusOK {
			return false
		}
	}

	return true
}

// DaysSinceLastCheck returns the number of days since the last check in the nameservers
// or in the DS records
func (d Domain) daysSinceLastCheck() int {
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
func (d Domain) isNearDNSSECExpirationDate(daysBefore int) bool {
	// Lets look for the oldest expiration date of the DS set, the it's probably the most
	// problematic one
	var expiresAt time.Time

	for i := 0; i < len(d.DSSet); i++ {
		// If is the first iteration we don't compare to expiresAt, because we do compare it
		// will be always less than any other date. And we also check if the expiration date
		// of the DS record was initialized to avoid replacing a real date for a not
		// initialized date
		if expiresAt.IsZero() || (!d.DSSet[i].ExpiresAt.IsZero() &&
			d.DSSet[i].ExpiresAt.Before(expiresAt)) {

			expiresAt = d.DSSet[i].ExpiresAt
		}
	}

	// When there's no DS, we don't have an expiration date and shouldn't care about it
	if expiresAt.IsZero() {
		return false
	}

	// Now we can check if this expiration date is already in the alert period. The alert
	// period is defined by the function parameter in days
	return time.Now().Add(time.Duration(daysBefore*24) * time.Hour).After(expiresAt)
}

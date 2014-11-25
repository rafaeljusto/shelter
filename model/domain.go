// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"code.google.com/p/go.net/idna"
	"fmt"
	"github.com/rafaeljusto/shelter/protocol"
	"gopkg.in/mgo.v2/bson"
	"math"
	"math/rand"
	"strings"
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
	Nameservers    Nameservers   // Nameservers that asnwer with authority for this domain
	DSSet          DSSet         // Records for the DNS tree chain of trust
	Owners         Owners        // Responsables for the domains that will receive alerts
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

// Apply is used to merge a domain request object sent by the user into a domain object of the
// database. It can return errors related to merge problems that are problem caused by data format
// of the user input
func (d *Domain) Apply(domainRequest protocol.DomainRequest) bool {
	if domainRequest.FQDN == nil {
		return false
	}

	// Detect when the domain object is empty, that is the case when we are creating a new
	// domain in the Shelter project
	if len(d.FQDN) == 0 {
		d.FQDN = *domainRequest.FQDN

	} else {
		// Cannot merge domains with different FQDNs
		if d.FQDN != *domainRequest.FQDN {
			return false
		}
	}

	var ok bool

	if d.Nameservers, ok = d.Nameservers.Apply(domainRequest.Nameservers); !ok {
		return false
	}

	if d.DSSet, ok = d.DSSet.Apply(domainRequest.DSSet); !ok {
		return false
	}

	if d.DSSet, ok = d.DSSet.ApplyDNSKEYs(domainRequest.DNSKEYS); !ok {
		return false
	}

	// We can replace the whole structure of the e-mail every time that a new UPDATE arrives
	// because there's no extra information in server side that we need to keep
	if d.Owners, ok = d.Owners.Apply(domainRequest.Owners); !ok {
		return false
	}

	return true
}

// Convert the domain system object to a limited information user format. We have a persisted flag
// to known when the object exists in our database or not to choose when we need to add the object
// links or not
func (d *Domain) Protocol() protocol.DomainResponse {
	var links []protocol.Link

	// We don't add links when the object doesn't exist in the system yet
	if d.Id.Valid() {
		// We should add more links here for system navigation. For example, we could add links
		// for object update, delete, list, etc. But I did not found yet in IANA list the
		// correct link type to be used. Also, the URI is hard coded, I didn't have any idea on
		// how can we do this dynamically yet. We cannot get the URI from the handler because we
		// are going to have a cross-reference problem
		links = append(links, protocol.Link{
			Types: []protocol.LinkType{protocol.LinkTypeSelf},
			HRef:  fmt.Sprintf("/domain/%s", d.FQDN),
		})
	}

	// We will try to show the FQDN always in unicode format. To solve a little bug in IDN library, we
	// will always convert the FQDN to lower case
	fqdn := strings.ToLower(d.FQDN)

	var err error
	fqdn, err = idna.ToUnicode(fqdn)
	if err != nil {
		// On error, keep the ace format
		fqdn = strings.ToLower(d.FQDN)
	}

	return protocol.DomainResponse{
		FQDN:        fqdn,
		Nameservers: d.Nameservers.Protocol(),
		DSSet:       d.DSSet.Protocol(),
		Owners:      d.Owners.Protocol(),
		Links:       links,
	}
}

type Domains []Domain

// Convert a list of domain objects into protocol format with pagination support
func (d *Domains) Protocol(pagination DomainPagination, expand bool, filter string) protocol.DomainsResponse {
	var domains []protocol.DomainResponse
	for _, domain := range *d {
		domains = append(domains, domain.Protocol())
	}

	var orderBy string
	for _, sort := range pagination.OrderBy {
		if len(orderBy) > 0 {
			orderBy += "@"
		}

		orderBy += fmt.Sprintf("%s:%s",
			DomainOrderByFieldToString(sort.Field),
			OrderByDirectionToString(sort.Direction),
		)
	}

	expandParameter := ""
	if expand {
		expandParameter = "&expand"
	}

	// Add pagination managment links to the response. The URI is hard coded, I didn't have
	// any idea on how can we do this dynamically yet. We cannot get the URI from the
	// handler because we are going to have a cross-reference problem
	var links []protocol.Link

	// Only add fast backward if we aren't in the first page
	if pagination.Page > 1 {
		links = append(links, protocol.Link{
			Types: []protocol.LinkType{protocol.LinkTypeFirst},
			HRef: fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s&filter=%s%s",
				pagination.PageSize, 1, orderBy, filter, expandParameter),
		})
	}

	// Only add previous if theres a previous page
	if pagination.Page-1 >= 1 {
		links = append(links, protocol.Link{
			Types: []protocol.LinkType{protocol.LinkTypePrev},
			HRef: fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s&filter=%s%s",
				pagination.PageSize, pagination.Page-1, orderBy, filter, expandParameter),
		})
	}

	// Only add next if there's a next page
	if pagination.Page+1 <= pagination.NumberOfPages {
		links = append(links, protocol.Link{
			Types: []protocol.LinkType{protocol.LinkTypeNext},
			HRef: fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s&filter=%s%s",
				pagination.PageSize, pagination.Page+1, orderBy, filter, expandParameter),
		})
	}

	// Only add the fast forward if we aren't on the last page
	if pagination.Page < pagination.NumberOfPages {
		links = append(links, protocol.Link{
			Types: []protocol.LinkType{protocol.LinkTypeLast},
			HRef: fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s&filter=%s%s",
				pagination.PageSize, pagination.NumberOfPages, orderBy, filter, expandParameter),
		})
	}

	return protocol.DomainsResponse{
		Page:          pagination.Page,
		PageSize:      pagination.PageSize,
		NumberOfPages: pagination.NumberOfPages,
		NumberOfItems: pagination.NumberOfItems,
		Domains:       domains,
		Links:         links,
	}
}

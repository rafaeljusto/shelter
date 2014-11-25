// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"strconv"
	"testing"
	"time"
)

func TestShouldBeScanned(t *testing.T) {
	var (
		maxOKVerificationDays    = 7
		maxErrorVerificationDays = 3
		maxExpirationAlertDays   = 10
	)

	maxErrorVerificationDuration, _ :=
		time.ParseDuration("-" + strconv.Itoa(maxErrorVerificationDays*24) + "h")
	maxOKVerificationDuration, _ :=
		time.ParseDuration("-" + strconv.Itoa(maxOKVerificationDays*24) + "h")
	maxExpirationAlertDuration, _ :=
		time.ParseDuration("-" + strconv.Itoa(maxExpirationAlertDays*24) + "h")

	d := Domain{
		Nameservers: []Nameserver{
			{
				LastStatus:  NameserverStatusServerFailure,
				LastCheckAt: time.Now().Add(maxErrorVerificationDuration),
			},
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now().Add(maxErrorVerificationDuration),
			},
		},
	}

	if !d.ShouldBeScanned(maxOKVerificationDays,
		maxErrorVerificationDays, maxExpirationAlertDays) {

		t.Error("Did not select a domain with DNS errors and in the limit " +
			"of errors max verification period")
	}

	d = Domain{
		DSSet: []DS{
			{
				LastStatus:  DSStatusTimeout,
				LastCheckAt: time.Now().Add(maxErrorVerificationDuration),
			},
			{
				LastStatus:  DSStatusTimeout,
				LastCheckAt: time.Now().Add(maxErrorVerificationDuration),
			},
		},
	}

	if !d.ShouldBeScanned(maxOKVerificationDays,
		maxErrorVerificationDays, maxExpirationAlertDays) {

		t.Error("Did not select a domain with DNSSEC errors and in the limit " +
			"of errors max verification period")
	}

	d = Domain{
		Nameservers: []Nameserver{
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now().Add(maxOKVerificationDuration),
			},
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now().Add(maxOKVerificationDuration),
			},
		},
	}

	if !d.ShouldBeScanned(maxOKVerificationDays,
		maxErrorVerificationDays, maxExpirationAlertDays) {

		t.Error("Did not select a domain configured correctly but in the limit " +
			"of ok max verification period")
	}

	d = Domain{
		Nameservers: []Nameserver{
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now(),
			},
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now(),
			},
		},

		DSSet: []DS{
			{
				LastStatus:  DSStatusOK,
				LastCheckAt: time.Now(),
				ExpiresAt:   time.Now().Add(maxExpirationAlertDuration),
			},
			{
				LastStatus:  DSStatusOK,
				LastCheckAt: time.Now(),
			},
		},
	}

	if !d.ShouldBeScanned(maxOKVerificationDays,
		maxErrorVerificationDays, maxExpirationAlertDays) {

		t.Error("Did not select a domain with DNSSEC signatures near expiration")
	}

	d = Domain{
		Nameservers: []Nameserver{
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now(),
			},
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now(),
			},
		},
	}

	if d.ShouldBeScanned(maxOKVerificationDays,
		maxErrorVerificationDays, maxExpirationAlertDays) {

		t.Error("Selected a domain configured correctly and checked now, " +
			"this shouldn't happen, or this domain is very lucky, we should " +
			"take a look in the algorithm")
	}
}

func TestAllNameserversOK(t *testing.T) {
	d := Domain{
		Nameservers: []Nameserver{
			{
				LastStatus: NameserverStatusServerFailure,
			},
			{
				LastStatus: NameserverStatusOK,
			},
		},
	}

	if d.allNameserversOK() {
		t.Error("Fail to detect a nameserver with error")
	}

	d = Domain{
		Nameservers: []Nameserver{},
	}

	if !d.allNameserversOK() {
		t.Error("Fail to inform all nameservers are OK when there are no nameservers")
	}

	d = Domain{
		Nameservers: []Nameserver{
			{
				LastStatus: NameserverStatusOK,
			},
			{
				LastStatus: NameserverStatusOK,
			},
		},
	}

	if !d.allNameserversOK() {
		t.Error("Fail to inform that all nameservers are OK")
	}
}

func TestAllDSSetOK(t *testing.T) {
	d := Domain{
		DSSet: []DS{
			{
				LastStatus: DSStatusExpiredSignature,
			},
			{
				LastStatus: DSStatusOK,
			},
		},
	}

	if d.allDSSetOK() {
		t.Error("Fail to detect a DS with error")
	}

	d = Domain{
		DSSet: []DS{},
	}

	if !d.allDSSetOK() {
		t.Error("Fail to inform all DS set are OK when there are no DS records")
	}

	d = Domain{
		DSSet: []DS{
			{
				LastStatus: DSStatusOK,
			},
			{
				LastStatus: DSStatusOK,
			},
		},
	}

	if !d.allDSSetOK() {
		t.Error("Fail to inform that all DS set are OK")
	}
}

func TestDaysSinceLastCheck(t *testing.T) {
	twoDays, _ := time.ParseDuration("48h")
	threeDays, _ := time.ParseDuration("72h")
	fourDays, _ := time.ParseDuration("96h")
	fiveDays, _ := time.ParseDuration("120h")

	d := Domain{
		Nameservers: []Nameserver{
			{
				LastCheckAt: time.Now().Add(-fiveDays),
			},
			{
				LastCheckAt: time.Now().Add(-twoDays),
			},
			{
				LastCheckAt: time.Now().Add(-twoDays),
			},
		},
		DSSet: []DS{
			{
				LastCheckAt: time.Now().Add(-fourDays),
			},
			{
				LastCheckAt: time.Now().Add(-threeDays),
			},
		},
	}

	if days := d.daysSinceLastCheck(); days != 2 {
		t.Errorf("Not counting correctly the number of days since last "+
			"check for nameservers. Should be 2, but got %d", days)
	}

	d = Domain{
		Nameservers: []Nameserver{
			{
				LastCheckAt: time.Now().Add(-fiveDays),
			},
			{
				LastCheckAt: time.Now().Add(-fourDays),
			},
			{
				LastCheckAt: time.Now().Add(-fourDays),
			},
		},
		DSSet: []DS{
			{
				LastCheckAt: time.Now().Add(-fourDays),
			},
			{
				LastCheckAt: time.Now().Add(-threeDays),
			},
		},
	}

	if days := d.daysSinceLastCheck(); days != 3 {
		t.Errorf("Not counting correctly the number of days since last "+
			"check for DS set. Should be 3, but got %d", days)
	}

	d = Domain{}

	if d.daysSinceLastCheck() < 365 {
		t.Errorf("Not counting correctly the number of days since last check " +
			"when the object is empty")
	}
}

func TestNearDNSSECExpirationDate(t *testing.T) {
	tenDays, _ := time.ParseDuration("240h")
	elevenDays, _ := time.ParseDuration("264h")

	d := Domain{
		DSSet: []DS{
			{
				ExpiresAt: time.Now().Add(tenDays),
			},
		},
	}

	if !d.isNearDNSSECExpirationDate(11) {
		t.Error("Could not detect when the DNSSEC expiration is near")
	}

	d = Domain{
		DSSet: []DS{
			{
				ExpiresAt: time.Now().Add(elevenDays),
			},
		},
	}

	if d.isNearDNSSECExpirationDate(10) {
		t.Error("Could not detect when the DNSSEC expiration is far")
	}

	d = Domain{
		DSSet: []DS{},
	}

	if d.isNearDNSSECExpirationDate(10) {
		t.Error("Could not detect when there's no expiration date")
	}
}

func TestDomainApply(t *testing.T) {
	domainRequest := DomainRequest{
		FQDN: "example.com.br",
		Nameservers: []NameserverRequest{
			{
				Host: "ns1.example.com.br",
				IPv4: "127.0.0.1",
				IPv6: "::1",
			},
			{
				Host: "ns2.example.com.br",
				IPv4: "127.0.0.2",
				IPv6: "::2",
			},
		},
		DSSet: []DSRequest{
			{
				Keytag:     41674,
				Algorithm:  5,
				Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
				DigestType: 1,
			},
			{
				Keytag:     45966,
				Algorithm:  7,
				Digest:     "B7C0BDE8F3C90E573B956B14A14CAF5001A3E841",
				DigestType: 1,
			},
		},
		DNSKEYS: []DNSKEYRequest{
			{
				Flags:     257,
				Algorithm: uint8(5),
				PublicKey: "AwEAAblaEaapG4inrQASY3HzwXwBaRSy5mkj7mZ30F+huI7zL8g0U7dv 7ufnSEQUlsC57OHoTBza+TQIv/mgQed8Fy4XGCGzYiHSYVYvGO9iWG3O 0voBYy/zv0z7ANfrA7Z3lY51CI6m/qoZUcDlNM0yTcJgilaKwUkLBHMA p9NJPuKVt8A7OHab00r2RDEVjiLWIIuTbz74gCXOVfAmvW07c8c=",
			},
			{
				Flags:     256,
				Algorithm: uint8(5),
				PublicKey: "AwEAAfFQjspE7NgjAPclHrlyVFPRUHrU1p1U6POUXDpuIg8grg/s0lG1 8sjMkpxIvecIePLJw24gx48Ta9g0JJzPy35oGX5rYVJAu9BPqdUEuwIN ScTy3fPUhubvXP2fbyS6LeKNX/ZenihCD4HrViZehJmsKKv5fX8qx+RL 7NXCAAM1Xdet13cqR3LduW6wBzMiaQ==",
			},
		},
		Owners: []OwnerRequest{
			{
				Email:    "example1@example.com.br",
				Language: "pt-br",
			},
			{
				Email:    "example2@example.com.br",
				Language: "en-us",
			},
		},
	}

	email, err := mail.ParseAddress("example0@example.com.br.")
	if err != nil {
		t.Fatal(err)
	}

	domain := Domain{
		FQDN: "example.com.br.",
		Nameservers: []Nameserver{
			{
				Host:       "ns1.example.com.br.",
				LastStatus: NameserverStatusTimeout,
			},
		},
		DSSet: []DS{
			{
				Keytag:     41674,
				Algorithm:  5,
				Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
				DigestType: 2,
				LastStatus: DSStatusTimeout,
			},
		},
		Owners: []Owner{
			{
				Email:    email,
				Language: fmt.Sprintf("%s-%s", LanguageTypePT, RegionTypeBR),
			},
		},
	}

	if !domain.Apply(domainRequest) {
		t.Fatal("Not merging correctly a valid domain request")
	}

	if len(domain.Nameservers) != 2 ||
		domain.Nameservers[0].IPv4.IsUnspecified() ||
		domain.Nameservers[0].IPv6.IsUnspecified() ||
		domain.Nameservers[0].LastStatus != NameserverStatusTimeout {

		t.Error("Fail to merge nameservers correctly")
	}

	if len(domain.DSSet) != 4 ||
		domain.DSSet[0].DigestType != DSDigestTypeSHA1 ||
		domain.DSSet[0].LastStatus != DSStatusTimeout {

		t.Error("Fail to merge DS set correctly")
	}

	if len(domain.Owners) != 2 ||
		domain.Owners[0].Email.Address != "example1@example.com.br" ||
		domain.Owners[0].Language != "pt-BR" {
		t.Error("Fail to replace owners")
	}

	domainRequest = DomainRequest{
		FQDN: nil,
	}

	if domain.Apply(domainRequest) {
		t.Error("Not detecting invalid FQDN")
	}

	domainRequest = DomainRequest{
		FQDN: "exampleX.com.br",
	}

	domain = Domain{
		FQDN: "example.com.br.",
	}

	if domain.Apply(domainRequest) {
		t.Error("Not detecting when we merge two different FQDNs")
	}

	domainRequest = DomainRequest{
		FQDN: "example.com.br",
		Nameservers: []NameserverRequest{
			{
				Host: nil,
				IPv4: "127..0.0.1",
				IPv6: "::1",
			},
			{
				Host: "ns2.example.com.br",
				IPv4: "127.0.0.2",
				IPv6: "::2",
			},
		},
	}

	if domain.Apply(domainRequest) {
		t.Error("Not detecting when there's an error in nameservers")
	}

	domainRequest = DomainRequest{
		FQDN: "example.com.br",
		DSSet: []DSRequest{
			{
				Keytag:     41674,
				Algorithm:  5,
				Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
				DigestType: 6,
			},
			{
				Keytag:     nil,
				Algorithm:  7,
				Digest:     "B7C0BDE8F3C90E573B956B14A14CAF5001A3E841",
				DigestType: 1,
			},
		},
	}

	if domain.Apply(domainRequest) {
		t.Error("Not detecting when there's an error in the DS set")
	}

	domainRequest = DomainRequest{
		FQDN: "example.com.br",
		DNSKEYS: []DNSKEYRequest{
			{
				Flags:     257,
				Algorithm: uint8(5),
				PublicKey: "AwEAAblaEaapG4inrQASY3HzwXwBaRSy5mkj7mZ30F+huI7zL8g0U7dv 7ufnSEQUlsC57OHoTBza+TQIv/mgQed8Fy4XGCGzYiHSYVYvGO9iWG3O 0voBYy/zv0z7ANfrA7Z3lY51CI6m/qoZUcDlNM0yTcJgilaKwUkLBHMA p9NJPuKVt8A7OHab00r2RDEVjiLWIIuTbz74gCXOVfAmvW07c8c=",
			},
			{
				Flags:     256,
				Algorithm: uint8(5),
				PublicKey: nil,
			},
		},
	}

	if domain.Apply(domainRequest) {
		t.Error("Not detecting when there's an error in the DNSKEYs")
	}

	domainRequest = DomainRequest{
		FQDN: "example.com.br",
		Owners: []OwnerRequest{
			{
				Email:    "wrongemail.com.br",
				Language: "pt-br",
			},
		},
	}

	if domain.Apply(domainRequest) {
		t.Error("Not detecting when there's an error in owners")
	}
}

func TestDomainProtocol(t *testing.T) {
	email, err := mail.ParseAddress("example@example.com.br.")
	if err != nil {
		t.Fatal(err)
	}

	domain := Domain{
		FQDN: "xn--exmpl-4qa6c.com.br.",
		Nameservers: []Nameserver{
			{
				Host:       "ns1.example.com.br.",
				LastStatus: NameserverStatusTimeout,
			},
		},
		DSSet: []DS{
			{
				Keytag:     41674,
				Algorithm:  5,
				Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
				DigestType: 2,
				LastStatus: DSStatusTimeout,
			},
		},
		Owners: []Owner{
			{
				Email:    email,
				Language: fmt.Sprintf("%s-%s", LanguageTypePT, RegionTypeBR),
			},
		},
	}

	domainResponse := domain.Protocol()

	if domainResponse.FQDN != "exâmplé.com.br." {
		t.Error("Fail to convert FQDN")
	}

	if len(domainResponse.Nameservers) != 1 {
		t.Error("Fail to convert nameservers")
	}

	if len(domainResponse.DSSet) != 1 {
		t.Error("Fail to convert the DS set")
	}

	if len(domainResponse.Owners) != 1 ||
		domainResponse.Owners[0].Email != "example@example.com.br." ||
		domainResponse.Owners[0].Language != "pt-BR" {

		t.Error("Fail to convert owners")
	}

	if len(domainResponse.Links) != 1 {
		t.Error("Wrong number of links")
	}

	domainResponse = ToDomainResponse(domain, false)

	if len(domainResponse.Links) != 0 {
		t.Error("Shouldn't return links when the object doesn't exist in the system")
	}

	// Testing case problem
	domain.FQDN = "XN--EXMPL-4QA6C.COM.BR."
	domainResponse = ToDomainResponse(domain, false)

	if domainResponse.FQDN != "exâmplé.com.br." {
		t.Errorf("Should convert to unicode even in upper case. "+
			"Expected '%s' and got '%s'", "exâmplé.com.br.", domainResponse.FQDN)
	}

	// Testing an invalid FQDN ACE format
	domain.FQDN = "xn--x1x2x3x4x5.com.br."
	domainResponse = ToDomainResponse(domain, false)

	if domainResponse.FQDN != "xn--x1x2x3x4x5.com.br." {
		t.Errorf("Should keep the ACE format when there's an error converting to unicode. "+
			"Expected '%s' and got '%s'", "xn--x1x2x3x4x5.com.br.", domainResponse.FQDN)
	}
}

func TestDomainsProtocol(t *testing.T) {
	domains := Domains{
		{
			FQDN: "example1.com.br.",
		},
		{
			FQDN: "example2.com.br.",
		},
		{
			FQDN: "example3.com.br.",
		},
		{
			FQDN: "example4.com.br.",
		},
		{
			FQDN: "example5.com.br.",
		},
	}

	pagination := DomainPagination{
		PageSize: 10,
		Page:     1,
		OrderBy: []DomainSort{
			{
				Field:     DomainOrderByFieldFQDN,
				Direction: OrderByDirectionAscending,
			},
			{
				Field:     DomainOrderByFieldLastModifiedAt,
				Direction: OrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: len(domains) / 10,
	}

	domainsResponse := domains.Protocol(pagination, true, "example")

	if len(domainsResponse.Domains) != len(domains) {
		t.Error("Not converting domain model objects properly")
	}

	if domainsResponse.PageSize != 10 {
		t.Error("Pagination not storing the page size properly")
	}

	if domainsResponse.Page != 1 {
		t.Error("Pagination not storing the current page properly")
	}

	if domainsResponse.NumberOfItems != len(domains) {
		t.Error("Pagination not storing number of items properly")
	}

	if domainsResponse.NumberOfPages != len(domains)/10 {
		t.Error("Pagination not storing number of pages properly")
	}

	// When we are on the first page, and there's no other page, don't show any link
	if len(domainsResponse.Links) != 0 {
		t.Error("Response not adding the necessary links when there is only one page")
	}
}

func TestDomainsProtocolWithLinks(t *testing.T) {
	domains := Domains{
		{
			FQDN: "example1.com.br.",
		},
		{
			FQDN: "example2.com.br.",
		},
		{
			FQDN: "example3.com.br.",
		},
		{
			FQDN: "example4.com.br.",
		},
		{
			FQDN: "example5.com.br.",
		},
	}

	pagination := DomainPagination{
		PageSize: 2,
		Page:     2,
		OrderBy: []DomainSort{
			{
				Field:     DomainOrderByFieldFQDN,
				Direction: OrderByDirectionAscending,
			},
			{
				Field:     DomainOrderByFieldLastModifiedAt,
				Direction: OrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: 3,
	}

	domainsResponse := domains.Protocol(pagination, true, "example")

	// Show all actions when navigating in the middle of the pagination
	if len(domainsResponse.Links) != 4 {
		t.Error("Response not adding the necessary links when we are navigating")
	}

	pagination = DomainPagination{
		PageSize: 2,
		Page:     1,
		OrderBy: []DomainSort{
			{
				Field:     DomainOrderByFieldFQDN,
				Direction: OrderByDirectionAscending,
			},
			{
				Field:     DomainOrderByFieldLastModifiedAt,
				Direction: OrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: 3,
	}

	domainsResponse = domains.Protocol(pagination, true, "example")

	// Don't show previous or fast backward when we are in the first page
	if len(domainsResponse.Links) != 2 {
		t.Error("Response not adding the necessary links when we are at the first page")
	}

	pagination = DomainPagination{
		PageSize: 2,
		Page:     3,
		OrderBy: []DomainSort{
			{
				Field:     DomainOrderByFieldFQDN,
				Direction: OrderByDirectionAscending,
			},
			{
				Field:     DomainOrderByFieldLastModifiedAt,
				Direction: OrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: 3,
	}

	domainsResponse = domains.Protocol(pagination, true, "example")

	// Don't show next or fast foward when we are in the last page
	if len(domainsResponse.Links) != 2 {
		t.Error("Response not adding the necessary links when we are in the last page")
	}
}

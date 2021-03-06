// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"fmt"
	"net/mail"
	"strings"
	"testing"

	"github.com/rafaeljusto/shelter/model"
)

func TestMerge(t *testing.T) {
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

	email, err := mail.ParseAddress("example0@example.com.br")
	if err != nil {
		t.Fatal(err)
	}

	domain := model.Domain{
		FQDN: "example.com.br.",
		Nameservers: []model.Nameserver{
			{
				Host:       "ns1.example.com.br.",
				LastStatus: model.NameserverStatusTimeout,
			},
		},
		DSSet: []model.DS{
			{
				Keytag:     41674,
				Algorithm:  5,
				Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
				DigestType: 2,
				LastStatus: model.DSStatusTimeout,
			},
		},
		Owners: []model.Owner{
			{
				Email:    email,
				Language: fmt.Sprintf("%s-%s", model.LanguageTypePT, model.RegionTypeBR),
			},
		},
	}

	domain, err = Merge(domain, domainRequest)
	if err != nil {
		t.Fatal(err)
	}

	if len(domain.Nameservers) != 2 ||
		domain.Nameservers[0].IPv4.IsUnspecified() ||
		domain.Nameservers[0].IPv6.IsUnspecified() ||
		domain.Nameservers[0].LastStatus != model.NameserverStatusTimeout {

		t.Error("Fail to merge nameservers correctly")
	}

	if len(domain.DSSet) != 4 ||
		domain.DSSet[0].DigestType != model.DSDigestTypeSHA1 ||
		domain.DSSet[0].LastStatus != model.DSStatusTimeout {

		t.Error("Fail to merge DS set correctly")
	}

	if len(domain.Owners) != 2 ||
		domain.Owners[0].Email.Address != "example1@example.com.br" ||
		domain.Owners[0].Language != "pt-BR" {
		t.Error("Fail to replace owners")
	}

	domainRequest = DomainRequest{
		FQDN: strings.Repeat("x", 65536) + "\uff00", // int32 overflow
	}

	domain, err = Merge(model.Domain{}, domainRequest)
	if err == nil {
		t.Error("Not detecting invalid FQDN")
	}

	domainRequest = DomainRequest{
		FQDN: "exampleX.com.br",
	}

	domain = model.Domain{
		FQDN: "example.com.br.",
	}

	domain, err = Merge(domain, domainRequest)
	if err == nil {
		t.Error("Not detecting when we merge two different FQDNs")
	}

	domainRequest = DomainRequest{
		FQDN: "example.com.br",
		Nameservers: []NameserverRequest{
			{
				Host: "ns1.example.com.br",
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

	domain, err = Merge(model.Domain{}, domainRequest)
	if err == nil {
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
				Keytag:     45966,
				Algorithm:  7,
				Digest:     "B7C0BDE8F3C90E573B956B14A14CAF5001A3E841",
				DigestType: 1,
			},
		},
	}

	domain, err = Merge(model.Domain{}, domainRequest)
	if err == nil {
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
				PublicKey: strings.Repeat("x", 65536) + "\uff00",
			},
		},
	}

	domain, err = Merge(model.Domain{}, domainRequest)
	if err == nil {
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

	domain, err = Merge(model.Domain{}, domainRequest)
	if err == nil {
		t.Error("Not detecting when there's an error in owners")
	}
}

func TestToDomainResponse(t *testing.T) {
	email, err := mail.ParseAddress("example@example.com.br")
	if err != nil {
		t.Fatal(err)
	}

	domain := model.Domain{
		FQDN: "xn--exmpl-4qa6c.com.br.",
		Nameservers: []model.Nameserver{
			{
				Host:       "ns1.example.com.br.",
				LastStatus: model.NameserverStatusTimeout,
			},
		},
		DSSet: []model.DS{
			{
				Keytag:     41674,
				Algorithm:  5,
				Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
				DigestType: 2,
				LastStatus: model.DSStatusTimeout,
			},
		},
		Owners: []model.Owner{
			{
				Email:    email,
				Language: fmt.Sprintf("%s-%s", model.LanguageTypePT, model.RegionTypeBR),
			},
		},
	}

	domainResponse := ToDomainResponse(domain, true)

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
		domainResponse.Owners[0].Email != "example@example.com.br" ||
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

package protocol

import (
	"net/mail"
	"shelter/model"
	"strings"
	"testing"
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
		Owners: []string{
			"example1@example.com.br",
			"example2@example.com.br",
		},
	}

	email, err := mail.ParseAddress("example0@example.com.br.")
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
		Owners: []*mail.Address{
			email,
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

	if len(domain.DSSet) != 2 ||
		domain.DSSet[0].DigestType != 1 ||
		domain.DSSet[0].LastStatus != model.DSStatusTimeout {

		t.Error("Fail to merge DS set correctly")
	}

	if len(domain.Owners) != 2 {
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
		Owners: []string{
			"wrongemail.com.br",
		},
	}

	domain, err = Merge(model.Domain{}, domainRequest)
	if err == nil {
		t.Error("Not detecting when there's an error in owners")
	}
}

func TestToDomainResponse(t *testing.T) {
	email, err := mail.ParseAddress("example@example.com.br.")
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
		Owners: []*mail.Address{
			email,
		},
	}

	domainResponse := ToDomainResponse(domain)

	if domainResponse.FQDN != "example.com.br." {
		t.Error("Fail to convert FQDN")
	}

	if len(domainResponse.Nameservers) != 1 {
		t.Error("Fail to convert nameservers")
	}

	if len(domainResponse.DSSet) != 1 {
		t.Error("Fail to convert the DS set")
	}

	println(domainResponse.Owners[0])
	if len(domainResponse.Owners) != 1 ||
		domainResponse.Owners[0] != "example@example.com.br." {

		t.Error("Fail to convert owners")
	}
}

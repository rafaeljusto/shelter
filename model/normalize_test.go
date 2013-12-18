package model

import (
	"testing"
)

func TestNormalizeDomainName(t *testing.T) {
	domainName := "NS1.bücher.EXAMPLE.com"
	normalizedDomainName, err := NormalizeDomainName(domainName)

	if err != nil {
		t.Fatal(err)
	}

	if normalizedDomainName != "ns1.xn--bcher-kva.example.com." {
		t.Error("Not normalizing correctly the domain name")
	}

	domainName = ""
	normalizedDomainName, err = NormalizeDomainName(domainName)

	if err != nil {
		t.Fatal(err)
	}

	if normalizedDomainName != "" {
		t.Error("Not normalizing correctly empty domain name")
	}
}

func TestNormalizeDSDigest(t *testing.T) {
	digest := "  B7C0BDE8F3C90E573B9 56B14A14CAF5001A3E841  "
	if NormalizeDSDigest(digest) != "b7c0bde8f3c90e573b9 56b14a14caf5001a3e841" {
		t.Error("Not normalizing correctly the DS digest")
	}
}

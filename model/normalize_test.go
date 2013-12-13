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
}

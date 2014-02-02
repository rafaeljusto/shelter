package handler

import (
	"testing"
)

func TestGetFQDNFromURI(t *testing.T) {
	if getFQDNFromURI("/domain/rafael.net.br.") != "rafael.net.br." {
		t.Error("Not retrieving FQDN from URI correctly in simple case")
	}

	if getFQDNFromURI("/domain/rafael.net.br./verification") != "rafael.net.br." {
		t.Error("Not retrieving FQDN from URI correctly when in the middle of the URI")
	}

	if getFQDNFromURI("/domain/example.com.br./crazyuri/rafael.net.br.") != "rafael.net.br." {
		t.Error("Not retrieving FQDN from URI correctly when there're multiple FQDNs")
	}

	if len(getFQDNFromURI("xxx")) != 0 {
		t.Error("Not returning empty when there's no FQDN in the URI")
	}
}

func BenchmarkGetFQDNFromURI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getFQDNFromURI("/domain/rafael.net.br.")
	}
}

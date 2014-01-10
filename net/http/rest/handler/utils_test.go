package handler

import (
	"testing"
)

func TestGetFQDNFromURI(t *testing.T) {
	if getFQDNFromURI("/domain/rafael.net.br.") != "rafael.net.br." {
		t.Error("Not retrieving FQDN from URI correctly")
	}

	if len(getFQDNFromURI("xxx")) != 0 {
		t.Error("Not returning empty when there's no FQDN in the URI")
	}
}

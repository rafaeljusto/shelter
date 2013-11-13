package model

import (
	"testing"
	"time"
)

func TestNeedsGlue(t *testing.T) {
	fqdn := "test.com.br"
	nameserver := Nameserver{
		Host: "ns1.test.com.br",
	}

	if !nameserver.NeedsGlue(fqdn) {
		t.Error("Not identifying when a nameserver needs a glue record")
	}

	fqdn = "test.com.br."
	nameserver = Nameserver{
		Host: "ns1.test.com.br",
	}

	if !nameserver.NeedsGlue(fqdn) {
		t.Error("Not identifying when a nameserver needs a glue record with " +
			"a final dot in the FQDN")
	}

	fqdn = "test.com.br"
	nameserver = Nameserver{
		Host: "ns1.test.com.br.",
	}

	if !nameserver.NeedsGlue(fqdn) {
		t.Error("Not identifying when a nameserver needs a glue record with " +
			"a final dot in the host")
	}

	fqdn = "test.com.br"
	nameserver = Nameserver{
		Host: "ns1.test.com",
	}

	if nameserver.NeedsGlue(fqdn) {
		t.Error("Not identifying when a nameserver don't need a glue record")
	}
}

func TestNameserverChangeStatus(t *testing.T) {
	timeMark := time.Now()

	nameserver := Nameserver{
		LastStatus:  NameserverStatusFail,
		LastCheckAt: timeMark,
	}

	nameserver.ChangeStatus(NameserverStatusOK)

	if nameserver.LastStatus != NameserverStatusOK {
		t.Error("ChangeStatus method did not change nameserver attribute")
	}

	if nameserver.LastCheckAt.Before(timeMark) || nameserver.LastCheckAt.Equal(timeMark) {
		t.Error("ChangeStatus method did not update the last check date")
	}
}

// model - Description of the objects
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

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
		LastStatus:  NameserverStatusServerFailure,
		LastCheckAt: timeMark,
	}

	// Wait for a while to make a difference between the first and the second
	// LastCheck fields
	time.Sleep(1 * time.Millisecond)

	nameserver.ChangeStatus(NameserverStatusOK)

	if nameserver.LastStatus != NameserverStatusOK {
		t.Error("ChangeStatus method did not change nameserver attribute")
	}

	if nameserver.LastCheckAt.Before(timeMark) || nameserver.LastCheckAt.Equal(timeMark) {
		t.Error("ChangeStatus method did not update the last check date")
	}
}

func TestNameserverStatusToString(t *testing.T) {
	if NameserverStatusToString(NameserverStatusNotChecked) != "NOTCHECKED" {
		t.Error("Nameserver status NOTCHECKED not converting correctly to string")
	}

	if NameserverStatusToString(NameserverStatusOK) != "OK" {
		t.Error("Nameserver status OK not converting correctly to string")
	}

	if NameserverStatusToString(NameserverStatusTimeout) != "TIMEOUT" {
		t.Error("Nameserver status TIMEOUT not converting correctly to string")
	}

	if NameserverStatusToString(NameserverStatusNoAuthority) != "NOAA" {
		t.Error("Nameserver status NOAA not converting correctly to string")
	}

	if NameserverStatusToString(NameserverStatusUnknownDomainName) != "UDN" {
		t.Error("Nameserver status UDN not converting correctly to string")
	}

	if NameserverStatusToString(NameserverStatusUnknownHost) != "UH" {
		t.Error("Nameserver status UH not converting correctly to string")
	}

	if NameserverStatusToString(NameserverStatusServerFailure) != "SERVFAIL" {
		t.Error("Nameserver status SERVFAIL not converting correctly to string")
	}

	if NameserverStatusToString(NameserverStatusQueryRefused) != "QREFUSED" {
		t.Error("Nameserver status QREFUSED not converting correctly to string")
	}

	if NameserverStatusToString(NameserverStatusConnectionRefused) != "CREFUSED" {
		t.Error("Nameserver status CREFUSED not converting correctly to string")
	}

	if NameserverStatusToString(NameserverStatusCanonicalName) != "CNAME" {
		t.Error("Nameserver status CNAME not converting correctly to string")
	}

	if NameserverStatusToString(NameserverStatusNotSynchronized) != "NOTSYNCH" {
		t.Error("Nameserver status NOTSYNCH not converting correctly to string")
	}

	if NameserverStatusToString(NameserverStatusError) != "ERROR" {
		t.Error("Nameserver status ERROR not converting correctly to string")
	}

	if NameserverStatusToString(999999) != "" {
		t.Error("Unknown nameserver status associated to some existing status")
	}
}

// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
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

func TestNameserverApply(t *testing.T) {
	nameserverRequest := NameserverRequest{
		Host: "ns1.example.com.br",
		IPv4: "127.0.0.1",
		IPv6: "::1",
	}

	var nameserver Nameserver
	if !nameserver.Apply(nameserverRequest) {
		t.Fatal("Error applying nameserver")
	}

	if nameserver.Host != "ns1.example.com.br." {
		t.Error("Not normalizing nameserver's host")
	}

	if nameserver.IPv4.String() != "127.0.0.1" {
		t.Error("Not parsing correctly IPv4")
	}

	if nameserver.IPv6.String() != "::1" {
		t.Error("Not parsing correctly IPv6")
	}

	nameserverRequest = NameserverRequest{
		Host: "ns1.example.com.br",
		IPv4: "127..0.0.1",
		IPv6: "::1",
	}

	if nameserver.Apply(nameserverRequest) {
		t.Error("Accepting an invalid IPv4")
	}

	nameserverRequest = NameserverRequest{
		Host: "ns1.example.com.br",
		IPv4: "127.0.0.1",
		IPv6: ":::1",
	}

	if nameserver.Apply(nameserverRequest) {
		t.Error("Accepting an invalid IPv6")
	}
}

func TestNameserverProtocol(t *testing.T) {
	now := time.Now()

	nameserver := model.Nameserver{
		Host:        "ns1.example.com.br.",
		IPv4:        net.ParseIP("127.0.0.1"),
		IPv6:        net.ParseIP("::1"),
		LastStatus:  model.NameserverStatusOK,
		LastCheckAt: now,
		LastOKAt:    now,
	}

	nameserverResponse := nameserver.Protocol()

	if nameserverResponse.Host != "ns1.example.com.br." {
		t.Error("Fail to convert host")
	}

	if nameserverResponse.IPv4 != "127.0.0.1" {
		t.Error("Fail to convert IPv4")
	}

	if nameserverResponse.IPv6 != "::1" {
		t.Error("Fail to convert IPv6")
	}

	if nameserverResponse.LastStatus != NameserverStatusToString(model.NameserverStatusOK) {
		t.Error("Fail to convert last status")
	}

	if nameserverResponse.LastCheckAt.Unix() != now.Unix() ||
		nameserverResponse.LastOKAt.Unix() != now.Unix() {

		t.Error("Fail to convert dates")
	}
}

func TestNameserversApply(t *testing.T) {
	nameserversRequest := []NameserverRequest{
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
	}

	var nameservers Nameservers
	if !nameservers.Apply(nameserversRequest) {
		t.Error("Not accepting a valid nameservers request")
	}

	nameserversRequest = []NameserverRequest{
		{
			Host: "ns1.example.com.br",
			IPv4: "127.ABC.0.1",
			IPv6: "::1",
		},
		{
			Host: "ns2.example.com.br",
			IPv4: "127.0.0.2",
			IPv6: "::2",
		},
	}

	if nameservers.Apply(nameserversRequest) {
		t.Error("Not checking errors from nameservers conversion")
	}
}

func TestNameserversProtocol(t *testing.T) {
	now := time.Now()

	nameservers := Nameservers{
		{
			Host:        "ns1.example.com.br.",
			IPv4:        net.ParseIP("127.0.0.1"),
			IPv6:        net.ParseIP("::1"),
			LastStatus:  model.NameserverStatusOK,
			LastCheckAt: now,
			LastOKAt:    now,
		},
		{
			Host:        "ns2.example.com.br.",
			IPv4:        net.ParseIP("127.0.0.2"),
			IPv6:        net.ParseIP("::2"),
			LastStatus:  model.NameserverStatusError,
			LastCheckAt: now,
		},
	}

	nameserversResponse := nameservers.Protocol()

	if len(nameserversResponse) != 2 {
		t.Error("Fail to convert multiple nameservers")
	}
}

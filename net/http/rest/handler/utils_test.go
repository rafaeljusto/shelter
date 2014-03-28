// handler - REST handler of specific URI
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package handler

import (
	"testing"
	"time"
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

func TestGetScanIdFromURI(t *testing.T) {
	date, current, err := getScanIdFromURI("/scan/2014-02-06T22:16:04.865737-02:00")
	if err != nil || current || date.Format(time.RFC3339) != "2014-02-06T22:16:04-02:00" {
		t.Error("Not retrieving date correctly from URI in simple case")
	}

	_, current, _ = getScanIdFromURI("/scan/cUrREnt")
	if !current {
		t.Error("Not retrieving current scan case")
	}

	date, current, err = getScanIdFromURI("/scan/2014-02-06t22:16:04-02:00")
	if err != nil || current || date.Format(time.RFC3339) != "2014-02-06T22:16:04-02:00" {
		t.Error("Not retrieving date correctly from URI without nanoseconds")
	}

	date, current, err = getScanIdFromURI("/scan/2014-02-06t22:16:04.39847234Z")
	if err != nil || current || date.Format(time.RFC3339) != "2014-02-06T22:16:04Z" {
		t.Error("Not retrieving date correctly from URI in UTC")
	}

	_, current, err = getScanIdFromURI("/scan/abc")
	if current || err == nil {
		t.Error("Not returning error when it receives an invalid date")
	}
}

func BenchmarkGetScanIdFromURI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getScanIdFromURI("/scan/2014-02-06T22:16:04.865737-02:00")
	}
}

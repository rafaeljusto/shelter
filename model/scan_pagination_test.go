// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"testing"
)

func TestScanOrderByFieldFromString(t *testing.T) {
	if _, ok := ScanOrderByFieldFromString("xxx"); ok {
		t.Error("Accepting an invalid order by field")
	}

	if field, ok := ScanOrderByFieldFromString("  STARTEDAT  "); !ok || field != ScanOrderByFieldStartedAt {
		t.Error("Not accepting a valid order by field FQDN")
	}

	if field, ok := ScanOrderByFieldFromString("domainsSCANNED"); !ok || field != ScanOrderByFieldDomainsScanned {
		t.Error("Not accepting a valid order by field DomainsScanned")
	}

	if field, ok := ScanOrderByFieldFromString("domainswithdnssecSCANNED"); !ok || field != ScanOrderByFieldDomainsWithDNSSECScanned {
		t.Error("Not accepting a valid order by field DomainsScanned")
	}
}

func TestScanOrderByFieldToString(t *testing.T) {
	if field := ScanOrderByFieldToString(ScanOrderByField(9999)); len(field) > 0 {
		t.Error("Not returning empty string when is an unknown order by field")
	}

	if field := ScanOrderByFieldToString(ScanOrderByFieldStartedAt); field != "startedat" {
		t.Error("Not returning the correct order by field for FQDN")
	}

	if field := ScanOrderByFieldToString(ScanOrderByFieldDomainsScanned); field != "domainsscanned" {
		t.Error("Not returning the correct order by field for DomainsScanned")
	}

	if field := ScanOrderByFieldToString(ScanOrderByFieldDomainsWithDNSSECScanned); field != "domainswithdnssecscanned" {
		t.Error("Not returning the correct order by field for DomainsWithDNSSECScanned")
	}
}

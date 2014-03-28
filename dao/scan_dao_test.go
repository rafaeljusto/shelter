// dao - Objects persistence layer
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package dao

import (
	"testing"
)

func TestScanDAOOrderByFieldFromString(t *testing.T) {
	if _, err := ScanDAOOrderByFieldFromString("xxx"); err == nil {
		t.Error("Accepting an invalid order by field")
	}

	if field, err := ScanDAOOrderByFieldFromString("  STARTEDAT  "); err != nil || field != ScanDAOOrderByFieldStartedAt {
		t.Error("Not accepting a valid order by field FQDN")
	}

	if field, err := ScanDAOOrderByFieldFromString("domainsSCANNED"); err != nil || field != ScanDAOOrderByFieldDomainsScanned {
		t.Error("Not accepting a valid order by field DomainsScanned")
	}

	if field, err := ScanDAOOrderByFieldFromString("domainswithdnssecSCANNED"); err != nil || field != ScanDAOOrderByFieldDomainsWithDNSSECScanned {
		t.Error("Not accepting a valid order by field DomainsScanned")
	}
}

func TestScanDAOOrderByFieldToString(t *testing.T) {
	if field := ScanDAOOrderByFieldToString(ScanDAOOrderByField(9999)); len(field) > 0 {
		t.Error("Not returning empty string when is an unknown order by field")
	}

	if field := ScanDAOOrderByFieldToString(ScanDAOOrderByFieldStartedAt); field != "startedat" {
		t.Error("Not returning the correct order by field for FQDN")
	}

	if field := ScanDAOOrderByFieldToString(ScanDAOOrderByFieldDomainsScanned); field != "domainsscanned" {
		t.Error("Not returning the correct order by field for DomainsScanned")
	}

	if field := ScanDAOOrderByFieldToString(ScanDAOOrderByFieldDomainsWithDNSSECScanned); field != "domainswithdnssecscanned" {
		t.Error("Not returning the correct order by field for DomainsWithDNSSECScanned")
	}
}

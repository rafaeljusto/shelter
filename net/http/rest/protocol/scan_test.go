// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package protocol

import (
	"fmt"
	"github.com/rafaeljusto/shelter/model"
	"testing"
	"time"
)

func TestScanToScanResponse(t *testing.T) {
	scan := model.Scan{
		Status:                   model.ScanStatusExecuted,
		StartedAt:                time.Now().Add(-1 * time.Hour),
		FinishedAt:               time.Now().Add(-30 * time.Minute),
		DomainsScanned:           10,
		DomainsWithDNSSECScanned: 4,
		NameserverStatistics: map[string]uint64{
			model.NameserverStatusToString(model.NameserverStatusOK):      16,
			model.NameserverStatusToString(model.NameserverStatusTimeout): 4,
		},
		DSStatistics: map[string]uint64{
			model.DSStatusToString(model.DSStatusOK):               3,
			model.DSStatusToString(model.DSStatusExpiredSignature): 1,
		},
	}

	scanResponse := ScanToScanResponse(scan)

	if scanResponse.Status != "EXECUTED" {
		t.Error("Status is not being translated correctly for a scan")
	}

	if scanResponse.DomainsToBeScanned != 0 {
		t.Error("Domains to be scanned field was not converted correctly")
	}

	if scanResponse.DomainsScanned != 10 {
		t.Error("Domains scanned field was not converted correctly")
	}

	if scanResponse.DomainsWithDNSSECScanned != 4 {
		t.Error("Domains with DNSSEC scanned field was not converted correctly")
	}

	if !scanResponse.StartedAt.Equal(scan.StartedAt) {
		t.Error("Started time was not converted correctly")
	}

	if !scanResponse.FinishedAt.Equal(scan.FinishedAt) {
		t.Error("Finished time was not converted correctly")
	}

	if scanResponse.NameserverStatistics["OK"] != 16 ||
		scanResponse.NameserverStatistics["TIMEOUT"] != 4 {
		t.Error("Nameserver statistics weren't converted correctly")
	}

	if scanResponse.DSStatistics["OK"] != 3 ||
		scanResponse.DSStatistics["EXPSIG"] != 1 {
		t.Error("DS statistics weren't converted correctly")
	}

	if len(scanResponse.Links) != 1 ||
		scanResponse.Links[0].HRef != fmt.Sprintf("/scan/%s", scan.StartedAt.Format(time.RFC3339Nano)) {
		t.Error("Links weren't added correctly")
	}
}

func TestCurrentScanToScanResponse(t *testing.T) {
	currentScan := model.CurrentScan{
		DomainsToBeScanned: 4,
		Scan: model.Scan{
			Status:                   model.ScanStatusRunning,
			StartedAt:                time.Now().Add(-1 * time.Hour),
			DomainsScanned:           2,
			DomainsWithDNSSECScanned: 0,
		},
	}

	scanResponse := CurrentScanToScanResponse(currentScan)

	if scanResponse.Status != "RUNNING" {
		t.Error("Status is not being translated correctly for a current scan")
	}

	if scanResponse.DomainsToBeScanned != 4 {
		t.Error("Domains to be scanned field was not converted correctly")
	}

	if scanResponse.DomainsScanned != 2 {
		t.Error("Domains scanned field was not converted correctly")
	}

	if scanResponse.DomainsWithDNSSECScanned != 0 {
		t.Error("Domains with DNSSEC scanned field was not converted correctly")
	}

	if !scanResponse.StartedAt.Equal(currentScan.StartedAt) {
		t.Error("Started time was not converted correctly")
	}

	if !scanResponse.FinishedAt.IsZero() {
		t.Error("Finished time was not converted correctly")
	}

	if len(scanResponse.NameserverStatistics) > 0 {
		t.Error("Nameserver statistics weren't converted correctly")
	}

	if len(scanResponse.DSStatistics) > 0 {
		t.Error("DS statistics weren't converted correctly")
	}

	if len(scanResponse.Links) != 2 ||
		scanResponse.Links[0].HRef != "/scan/current" {
		t.Error("Links weren't added correctly")
	}
}

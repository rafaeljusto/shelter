// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/model"
	"testing"
	"time"
)

func TestScansToScansResponse(t *testing.T) {
	scans := []model.Scan{
		{
			StartedAt: time.Now().Add(-1 * time.Minute),
		},
		{
			StartedAt: time.Now().Add(-2 * time.Minute),
		},
		{
			StartedAt: time.Now().Add(-3 * time.Minute),
		},
		{
			StartedAt: time.Now().Add(-4 * time.Minute),
		},
		{
			StartedAt: time.Now().Add(-5 * time.Minute),
		},
	}

	pagination := dao.ScanDAOPagination{
		PageSize: 10,
		Page:     1,
		OrderBy: []dao.ScanDAOSort{
			{
				Field:     dao.ScanDAOOrderByFieldStartedAt,
				Direction: dao.DAOOrderByDirectionAscending,
			},
			{
				Field:     dao.ScanDAOOrderByFieldDomainsScanned,
				Direction: dao.DAOOrderByDirectionDescending,
			},
			{
				Field:     dao.ScanDAOOrderByFieldDomainsWithDNSSECScanned,
				Direction: dao.DAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(scans),
		NumberOfPages: len(scans) / 10,
	}

	scansResponse := ScansToScansResponse(scans, pagination)

	if len(scansResponse.Scans) != len(scans) {
		t.Error("Not converting scan model objects properly")
	}

	if scansResponse.PageSize != 10 {
		t.Error("Pagination not storing the page size properly")
	}

	if scansResponse.Page != 1 {
		t.Error("Pagination not storing the current page properly")
	}

	if scansResponse.NumberOfItems != len(scans) {
		t.Error("Pagination not storing number of items properly")
	}

	if scansResponse.NumberOfPages != len(scans)/10 {
		t.Error("Pagination not storing number of pages properly")
	}

	// We should show only the current link when there's only one page
	if len(scansResponse.Links) != 1 {
		t.Error("Response not adding the necessary links when there is only one page")
	}
}

func TestScansToScansResponseLinks(t *testing.T) {
	scans := []model.Scan{
		{
			StartedAt: time.Now().Add(-1 * time.Minute),
		},
		{
			StartedAt: time.Now().Add(-2 * time.Minute),
		},
		{
			StartedAt: time.Now().Add(-3 * time.Minute),
		},
		{
			StartedAt: time.Now().Add(-4 * time.Minute),
		},
		{
			StartedAt: time.Now().Add(-5 * time.Minute),
		},
	}

	pagination := dao.ScanDAOPagination{
		PageSize: 2,
		Page:     2,
		OrderBy: []dao.ScanDAOSort{
			{
				Field:     dao.ScanDAOOrderByFieldStartedAt,
				Direction: dao.DAOOrderByDirectionAscending,
			},
			{
				Field:     dao.ScanDAOOrderByFieldDomainsScanned,
				Direction: dao.DAOOrderByDirectionDescending,
			},
			{
				Field:     dao.ScanDAOOrderByFieldDomainsWithDNSSECScanned,
				Direction: dao.DAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(scans),
		NumberOfPages: 3,
	}

	scansResponse := ScansToScansResponse(scans, pagination)

	// Show all actions when navigating in the middle of the pagination
	if len(scansResponse.Links) != 5 {
		t.Error("Response not adding the necessary links when we are navigating")
	}

	pagination = dao.ScanDAOPagination{
		PageSize: 2,
		Page:     1,
		OrderBy: []dao.ScanDAOSort{
			{
				Field:     dao.ScanDAOOrderByFieldStartedAt,
				Direction: dao.DAOOrderByDirectionAscending,
			},
			{
				Field:     dao.ScanDAOOrderByFieldDomainsScanned,
				Direction: dao.DAOOrderByDirectionDescending,
			},
			{
				Field:     dao.ScanDAOOrderByFieldDomainsWithDNSSECScanned,
				Direction: dao.DAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(scans),
		NumberOfPages: 3,
	}

	scansResponse = ScansToScansResponse(scans, pagination)

	// On the first page we show current, next and fast foward links
	if len(scansResponse.Links) != 3 {
		t.Error("Response not adding the necessary links when we are at the first page")
	}

	pagination = dao.ScanDAOPagination{
		PageSize: 2,
		Page:     3,
		OrderBy: []dao.ScanDAOSort{
			{
				Field:     dao.ScanDAOOrderByFieldStartedAt,
				Direction: dao.DAOOrderByDirectionAscending,
			},
			{
				Field:     dao.ScanDAOOrderByFieldDomainsScanned,
				Direction: dao.DAOOrderByDirectionDescending,
			},
			{
				Field:     dao.ScanDAOOrderByFieldDomainsWithDNSSECScanned,
				Direction: dao.DAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(scans),
		NumberOfPages: 3,
	}

	scansResponse = ScansToScansResponse(scans, pagination)

	// Don't show next or fast foward when we are in the last page
	if len(scansResponse.Links) != 3 {
		t.Error("Response not adding the necessary links when we are in the last page")
	}
}

func TestCurrentScanToScansResponse(t *testing.T) {
	currentScan := model.CurrentScan{
		DomainsToBeScanned: 4,
		Scan: model.Scan{
			Status:                   model.ScanStatusRunning,
			StartedAt:                time.Now().Add(-1 * time.Hour),
			DomainsScanned:           2,
			DomainsWithDNSSECScanned: 0,
		},
	}

	pagination := dao.ScanDAOPagination{
		PageSize: 2,
		Page:     2,
		OrderBy: []dao.ScanDAOSort{
			{
				Field:     dao.ScanDAOOrderByFieldStartedAt,
				Direction: dao.DAOOrderByDirectionAscending,
			},
			{
				Field:     dao.ScanDAOOrderByFieldDomainsScanned,
				Direction: dao.DAOOrderByDirectionDescending,
			},
			{
				Field:     dao.ScanDAOOrderByFieldDomainsWithDNSSECScanned,
				Direction: dao.DAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: 6,
		NumberOfPages: 3,
	}

	scansResponse := CurrentScanToScansResponse(currentScan, pagination)

	if len(scansResponse.Scans) != 1 {
		t.Error("When retrieving the current scan, no other scan can appear")
	}

	if scansResponse.Scans[0].Status != "RUNNING" {
		t.Error("Status is not being translated correctly for a current scan")
	}

	if scansResponse.Scans[0].DomainsToBeScanned != 4 {
		t.Error("Domains to be scanned field was not converted correctly")
	}

	if scansResponse.Scans[0].DomainsScanned != 2 {
		t.Error("Domains scanned field was not converted correctly")
	}

	if scansResponse.Scans[0].DomainsWithDNSSECScanned != 0 {
		t.Error("Domains with DNSSEC scanned field was not converted correctly")
	}

	if !scansResponse.Scans[0].StartedAt.Equal(currentScan.StartedAt) {
		t.Error("Started time was not converted correctly")
	}

	if !scansResponse.Scans[0].FinishedAt.IsZero() {
		t.Error("Finished time was not converted correctly")
	}

	if len(scansResponse.Scans[0].NameserverStatistics) > 0 {
		t.Error("Nameserver statistics weren't converted correctly")
	}

	if len(scansResponse.Scans[0].DSStatistics) > 0 {
		t.Error("DS statistics weren't converted correctly")
	}

	if len(scansResponse.Scans[0].Links) != 2 ||
		scansResponse.Scans[0].Links[0].HRef != "/scan/current" {
		t.Error("Links weren't added correctly")
	}
}

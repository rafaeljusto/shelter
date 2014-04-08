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

	// We can't show next when there's only one page, but the previous will be the current page
	if len(scansResponse.Links) != 3 {
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

	// On the first page we show as previous the current scan
	if len(scansResponse.Links) != 4 {
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

	// Don't show next when we are in the last page
	if len(scansResponse.Links) != 4 {
		t.Error("Response not adding the necessary links when we are in the last page")
	}
}

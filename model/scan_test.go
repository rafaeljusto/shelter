// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"errors"
	"github.com/rafaeljusto/shelter/scheduler"
	"sync"
	"testing"
	"time"
)

func TestScanProtocol(t *testing.T) {
	scan := Scan{
		Status:                   ScanStatusExecuted,
		StartedAt:                time.Now().Add(-1 * time.Hour),
		FinishedAt:               time.Now().Add(-30 * time.Minute),
		DomainsScanned:           10,
		DomainsWithDNSSECScanned: 4,
		NameserverStatistics: map[string]uint64{
			NameserverStatusToString(NameserverStatusOK):      16,
			NameserverStatusToString(NameserverStatusTimeout): 4,
		},
		DSStatistics: map[string]uint64{
			DSStatusToString(DSStatusOK):               3,
			DSStatusToString(DSStatusExpiredSignature): 1,
		},
	}

	scanResponse := scan.Protocol()

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

func TestScansProtocol(t *testing.T) {
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

	pagination := model.ScanPagination{
		PageSize: 10,
		Page:     1,
		OrderBy: []model.ScanSort{
			{
				Field:     model.ScanOrderByFieldStartedAt,
				Direction: model.OrderByDirectionAscending,
			},
			{
				Field:     model.ScanOrderByFieldDomainsScanned,
				Direction: model.OrderByDirectionDescending,
			},
			{
				Field:     model.ScanOrderByFieldDomainsWithDNSSECScanned,
				Direction: model.OrderByDirectionDescending,
			},
		},
		NumberOfItems: len(scans),
		NumberOfPages: len(scans) / 10,
	}

	scansResponse := scans.Protocol(pagination)

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

func TestScansProtocolWithLinks(t *testing.T) {
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

	pagination := model.ScanPagination{
		PageSize: 2,
		Page:     2,
		OrderBy: []model.ScanSort{
			{
				Field:     model.ScanOrderByFieldStartedAt,
				Direction: model.OrderByDirectionAscending,
			},
			{
				Field:     model.ScanOrderByFieldDomainsScanned,
				Direction: model.OrderByDirectionDescending,
			},
			{
				Field:     model.ScanOrderByFieldDomainsWithDNSSECScanned,
				Direction: model.OrderByDirectionDescending,
			},
		},
		NumberOfItems: len(scans),
		NumberOfPages: 3,
	}

	scansResponse := scans.Protocol(pagination)

	// Show all actions when navigating in the middle of the pagination
	if len(scansResponse.Links) != 5 {
		t.Error("Response not adding the necessary links when we are navigating")
	}

	pagination = model.ScanPagination{
		PageSize: 2,
		Page:     1,
		OrderBy: []model.ScanSort{
			{
				Field:     model.ScanOrderByFieldStartedAt,
				Direction: model.OrderByDirectionAscending,
			},
			{
				Field:     model.ScanOrderByFieldDomainsScanned,
				Direction: model.OrderByDirectionDescending,
			},
			{
				Field:     model.ScanOrderByFieldDomainsWithDNSSECScanned,
				Direction: model.OrderByDirectionDescending,
			},
		},
		NumberOfItems: len(scans),
		NumberOfPages: 3,
	}

	scansResponse = scans.Protocol(pagination)

	// On the first page we show current, next and fast foward links
	if len(scansResponse.Links) != 3 {
		t.Error("Response not adding the necessary links when we are at the first page")
	}

	pagination = model.ScanPagination{
		PageSize: 2,
		Page:     3,
		OrderBy: []model.ScanSort{
			{
				Field:     model.ScanOrderByFieldStartedAt,
				Direction: model.OrderByDirectionAscending,
			},
			{
				Field:     model.ScanOrderByFieldDomainsScanned,
				Direction: model.OrderByDirectionDescending,
			},
			{
				Field:     model.ScanOrderByFieldDomainsWithDNSSECScanned,
				Direction: model.OrderByDirectionDescending,
			},
		},
		NumberOfItems: len(scans),
		NumberOfPages: 3,
	}

	scansResponse = scans.Protocol(pagination)

	// Don't show next or fast foward when we are in the last page
	if len(scansResponse.Links) != 3 {
		t.Error("Response not adding the necessary links when we are in the last page")
	}
}

func TestCurrentScanProtocol(t *testing.T) {
	currentScan := model.CurrentScan{
		DomainsToBeScanned: 4,
		Scan: model.Scan{
			Status:                   model.ScanStatusRunning,
			StartedAt:                time.Now().Add(-1 * time.Hour),
			DomainsScanned:           2,
			DomainsWithDNSSECScanned: 0,
		},
	}

	pagination := model.ScanPagination{
		PageSize: 2,
		Page:     2,
		OrderBy: []model.ScanSort{
			{
				Field:     model.ScanOrderByFieldStartedAt,
				Direction: model.OrderByDirectionAscending,
			},
			{
				Field:     model.ScanOrderByFieldDomainsScanned,
				Direction: model.OrderByDirectionDescending,
			},
			{
				Field:     model.ScanOrderByFieldDomainsWithDNSSECScanned,
				Direction: model.OrderByDirectionDescending,
			},
		},
		NumberOfItems: 6,
		NumberOfPages: 3,
	}

	scansResponse := currentScan.Protocol(pagination)

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

func TestInitializeCurrentScan(t *testing.T) {
	if err := InitializeCurrentScan(); err == nil {
		t.Fatal("Not alerting about scheduler without scan job")
	}

	nextExecution := time.Now().Add(10 * time.Minute)

	scheduler.Register(scheduler.Job{
		Type:          scheduler.JobTypeScan,
		NextExecution: nextExecution,
		Task:          func() {},
	})

	if err := InitializeCurrentScan(); err != nil {
		t.Fatal(err)
	}

	if shelterCurrentScan.Status != ScanStatusWaitingExecution {
		t.Error("Wrong initial status of the current scan")
	}

	if !shelterCurrentScan.ScheduledAt.Equal(nextExecution) {
		t.Error("Wrong next execution time of the current scan")
	}

	scheduler.Clear()
}

func TestStartNewScan(t *testing.T) {
	shelterCurrentScan.DomainsScanned = 5
	shelterCurrentScan.DomainsWithDNSSECScanned = 2
	shelterCurrentScan.DomainsToBeScanned = 10
	shelterCurrentScan.Status = ScanStatusRunning

	StartNewScan()

	if shelterCurrentScan.DomainsScanned == 0 &&
		shelterCurrentScan.DomainsWithDNSSECScanned == 0 &&
		shelterCurrentScan.DomainsToBeScanned == 0 &&
		shelterCurrentScan.Status != ScanStatusLoadingData {

		t.Error("Not setting start scan information correctly")
	}
}

func TestFinishAndSaveScan(t *testing.T) {
	scheduler.Register(scheduler.Job{
		Type:          scheduler.JobTypeScan,
		NextExecution: time.Now().Add(10 * time.Minute),
		Task:          func() {},
	})

	StartNewScan()

	called := false
	if err := FinishAndSaveScan(false, func(s *Scan) error {
		if s.Status != ScanStatusExecuted {
			t.Error("Not setting finish scan information correctly")
		}
		called = true
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error("Not calling scan save method")
	}

	if shelterCurrentScan.Status != ScanStatusWaitingExecution {
		t.Error("Not setting next scan information correctly")
	}

	called = false
	if err := FinishAndSaveScan(true, func(s *Scan) error {
		if s.Status != ScanStatusExecutedWithErrors {
			t.Error("Not setting finish scan information correctly " +
				"when we had errors")
		}
		called = true
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error("Not calling scan save method when we had errors")
	}

	if err := FinishAndSaveScan(true, func(s *Scan) error {
		return errors.New("Error saving scan!")
	}); err == nil {
		t.Error("Not returning err detected by save method")
	}

	scheduler.Clear()

	if err := FinishAndSaveScan(true, func(s *Scan) error {
		return nil
	}); err == nil {
		t.Error("Not returning err when next scan is not found")
	}
}

func TestLoadedDomainForScan(t *testing.T) {
	StartNewScan()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			LoadedDomainForScan()
		}()
	}

	wg.Wait()

	if shelterCurrentScan.DomainsToBeScanned != 100 {
		t.Error("Not counting correctly the number of domains to scan " +
			"in a concurrent enviroment")
	}
}

func TestFinishLoadingDomainsForScan(t *testing.T) {
	StartNewScan()
	FinishLoadingDomainsForScan()

	if shelterCurrentScan.Status != ScanStatusRunning {
		t.Error("Not setting the correct status when finish " +
			"loading the data")
	}
}

func TestFinishAnalyzingDomainForScan(t *testing.T) {
	StartNewScan()

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				FinishAnalyzingDomainForScan(true)
			} else {
				FinishAnalyzingDomainForScan(false)
			}
		}(i)
	}

	wg.Wait()

	if shelterCurrentScan.DomainsScanned != 1000 {
		t.Error("Not counting correctly the domains scanned in " +
			"a concurrent enviroment")
	}

	if shelterCurrentScan.DomainsWithDNSSECScanned != 500 {
		println(shelterCurrentScan.DomainsWithDNSSECScanned)
		t.Error("Not counting correctly the domains with DNSSEC " +
			"scanned in a concurrent enviroment")
	}
}

func TestStoreStatisticsOfTheScan(t *testing.T) {
	StartNewScan()

	nameserverStatistics := make(map[string]uint64)
	nameserverStatistics[NameserverStatusToString(NameserverStatusOK)] = 534
	nameserverStatistics[NameserverStatusToString(NameserverStatusTimeout)] = 19
	nameserverStatistics[NameserverStatusToString(NameserverStatusUnknownHost)] = 3

	dsStatistics := make(map[string]uint64)
	dsStatistics[DSStatusToString(DSStatusOK)] = 32
	dsStatistics[DSStatusToString(DSStatusExpiredSignature)] = 7

	StoreStatisticsOfTheScan(nameserverStatistics, dsStatistics)

	if len(shelterCurrentScan.NameserverStatistics) != 3 {
		t.Error("Not storing namserver statistics")
	}

	if len(shelterCurrentScan.DSStatistics) != 2 {
		t.Error("Not storing DS statistics")
	}
}

func TestGetCurrentScan(t *testing.T) {
	StartNewScan()

	LoadedDomainForScan()
	LoadedDomainForScan()

	FinishLoadingDomainsForScan()

	FinishAnalyzingDomainForScan(false)
	FinishAnalyzingDomainForScan(true)

	currentScan := GetCurrentScan()
	if &currentScan == &shelterCurrentScan {
		t.Error("Not copying scan information object")
	}

	if currentScan.DomainsScanned != 2 ||
		currentScan.DomainsWithDNSSECScanned != 1 ||
		currentScan.DomainsToBeScanned != 2 ||
		currentScan.Status != ScanStatusRunning {

		t.Error("Scan information retrieved doesn't have the correct information")
	}
}

func TestScanStatusToString(t *testing.T) {
	if ScanStatusToString(ScanStatusExecuted) != "EXECUTED" {
		t.Error("Scan status EXECUTED not converting correctly to string")
	}

	if ScanStatusToString(ScanStatusExecutedWithErrors) != "EXECUTEDWITHERRORS" {
		t.Error("Scan status EXECUTEDWITHERRORS not converting correctly to string")
	}

	if ScanStatusToString(ScanStatusRunning) != "RUNNING" {
		t.Error("Scan status RUNNING not converting correctly to string")
	}

	if ScanStatusToString(ScanStatusLoadingData) != "LOADINGDATA" {
		t.Error("Scan status LOADINGDATA not converting correctly to string")
	}

	if ScanStatusToString(ScanStatusWaitingExecution) != "WAITINGEXECUTION" {
		t.Error("Scan status WAITINGEXECUTION not converting correctly to string")
	}

	if ScanStatusToString(999999) != "" {
		t.Error("Unknown scan status associated to some existing status")
	}
}

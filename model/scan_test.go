package model

import (
	"errors"
	"github.com/rafaeljusto/shelter/scheduler"
	"sync"
	"testing"
	"time"
)

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

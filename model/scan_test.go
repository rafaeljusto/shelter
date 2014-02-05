package model

import (
	"sync"
	"testing"
)

func TestStartNewScan(t *testing.T) {
	shelterCurrentScan.DomainsScanned = 5
	shelterCurrentScan.DomainsWihDNSSECScanned = 2
	shelterCurrentScan.DomainsToBeScanned = 10
	shelterCurrentScan.Status = ScanStatusRunning

	StartNewScan()

	if shelterCurrentScan.DomainsScanned == 0 &&
		shelterCurrentScan.DomainsWihDNSSECScanned == 0 &&
		shelterCurrentScan.DomainsToBeScanned == 0 &&
		shelterCurrentScan.Status != ScanStatusLoadingData {

		t.Error("Not setting start scan information correctly")
	}
}

func TestFinishAndSaveScan(t *testing.T) {
	StartNewScan()

	called := false
	FinishAndSaveScan(false, func(s *Scan) error {
		if s.Status != ScanStatusExecuted {
			t.Error("Not setting finish scan information correctly")
		}
		called = true
		return nil
	})

	if !called {
		t.Error("Not calling scan save method")
	}

	if shelterCurrentScan.Status != ScanStatusWaitingExecution {
		t.Error("Not setting next scan information correctly")
	}

	called = false
	FinishAndSaveScan(true, func(s *Scan) error {
		if s.Status != ScanStatusExecutedWithErrors {
			t.Error("Not setting finish scan information correctly " +
				"when we had errors")
		}
		called = true
		return nil
	})

	if !called {
		t.Error("Not calling scan save method when we had errors")
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

	if shelterCurrentScan.DomainsWihDNSSECScanned != 500 {
		println(shelterCurrentScan.DomainsWihDNSSECScanned)
		t.Error("Not counting correctly the domains with DNSSEC " +
			"scanned in a concurrent enviroment")
	}
}

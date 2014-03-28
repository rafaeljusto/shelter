// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package model

import (
	"github.com/rafaeljusto/shelter/scheduler"
	"labix.org/v2/mgo/bson"
	"sync"
	"sync/atomic"
	"time"
)

// Global variables that will store the scan information and guarantee that it can be
// accessed concurrently by the scan process and by the interface
var (
	shelterCurrentScan     CurrentScan // Store all data of the current scan
	shelterCurrentScanLock sync.Mutex  // Make current scan thread safe
)

// List of possible values of a scan status. The status ScanStatusWaitingExecution,
// ScanStatusLoadingData, ScanStatusRunning will only be visible in a CurrentScan struct
const (
	ScanStatusWaitingExecution   ScanStatus = 0 // The scan is going to be executed in the future
	ScanStatusLoadingData        ScanStatus = 1 // Loading the domains objects for Scan
	ScanStatusRunning            ScanStatus = 2 // The scan is current scanning the system
	ScanStatusExecuted           ScanStatus = 3 // Scan alredy finished succesfully
	ScanStatusExecutedWithErrors ScanStatus = 4 // Scan had problems during the execution
)

// We keep a state from the scan to identify scans that had problem or not, and the current
// executio. This is useful for reports and a richer interface
type ScanStatus int

// Convert the scan status enum to text for printing in reports or debugging
func ScanStatusToString(status ScanStatus) string {
	switch status {
	case ScanStatusWaitingExecution:
		return "WAITINGEXECUTION"
	case ScanStatusLoadingData:
		return "LOADINGDATA"
	case ScanStatusRunning:
		return "RUNNING"
	case ScanStatusExecuted:
		return "EXECUTED"
	case ScanStatusExecutedWithErrors:
		return "EXECUTEDWITHERRORS"
	}

	return ""
}

// Store all data related to a scan executed on the system. The statistics attributes cannot use the
// ENUM format because we cannot have a non-string key in the JSON format when saving into the
// database
type Scan struct {
	Id                       bson.ObjectId     `bson:"_id"` // Database identification
	Revision                 int               // Version of the object
	Status                   ScanStatus        // Status of the scan
	StartedAt                time.Time         // Date and time that the scan started
	FinishedAt               time.Time         // Date and time that the scan finished
	LastModifiedAt           time.Time         // Last time the object was modified
	DomainsScanned           uint64            // Number of domains scanned
	DomainsWithDNSSECScanned uint64            // Number of domains with DS recods scanned
	NameserverStatistics     map[string]uint64 // Statistics from nameserver status (text format) in number of hosts
	DSStatistics             map[string]uint64 // Statistics from DS records' status (text format) in number of DS records
}

// CurrentScan is a Scan that is the next to be executed or is executing at this moment. The data
// from this struct is not stored until the scan is finished and become only a Scan struct. This
// should be used to tell the user (using a service) how is a progress of a scan on-the-fly
type CurrentScan struct {
	Scan                         // CurrentScan is a Scan
	ScheduledAt        time.Time // Initial date and time that the scan was schedule to execute
	DomainsToBeScanned uint64    // Domains selected to be scanned
}

// Function to fill current scan variable for the first time. Should run after the
// scheduler registered the scan job, to determinate the next execution time. Returns an
// error if this function is executed before the scheduler register the scan job
func InitializeCurrentScan() error {
	shelterCurrentScanLock.Lock()
	defer shelterCurrentScanLock.Unlock()

	nextExecution, err := scheduler.NextExecutionByType(scheduler.JobTypeScan)

	shelterCurrentScan = CurrentScan{
		Scan: Scan{
			Status:               ScanStatusWaitingExecution,
			NameserverStatistics: make(map[string]uint64),
			DSStatistics:         make(map[string]uint64),
		},
		ScheduledAt: nextExecution,
	}

	// If err from different from nil we didn't find a scan job in the scheduler! Propably
	// this function was executed before the scheduler registered the scan job
	return err
}

// Function to alert that a new scan is going to be started. This function is necessary to
// garantee concurrency acces for the current scan information
func StartNewScan() {
	shelterCurrentScanLock.Lock()
	defer shelterCurrentScanLock.Unlock()

	shelterCurrentScan = CurrentScan{
		Scan: Scan{
			Status:               ScanStatusLoadingData,
			StartedAt:            time.Now().UTC(),
			NameserverStatistics: make(map[string]uint64),
			DSStatistics:         make(map[string]uint64),
		},
	}
}

// FinishAndSaveScan was created to alert that the scan being executed finished. This
// function is necessary to garantee concurrency acces for the current scan information.
// Will set all necessary information and save the scan into the database for future
// reports. Only a part of the scan is saved into the database, because some information
// is only useful during the execution
func FinishAndSaveScan(hadErrors bool, f func(*Scan) error) error {
	shelterCurrentScanLock.Lock()
	defer shelterCurrentScanLock.Unlock()

	if hadErrors {
		shelterCurrentScan.Status = ScanStatusExecutedWithErrors
	} else {
		shelterCurrentScan.Status = ScanStatusExecuted
	}

	shelterCurrentScan.FinishedAt = time.Now()

	// Save the scan
	err := f(&shelterCurrentScan.Scan)

	// Change the current scan state to prepare for the next scan
	shelterCurrentScan = CurrentScan{
		Scan: Scan{
			Status:               ScanStatusWaitingExecution,
			NameserverStatistics: make(map[string]uint64),
			DSStatistics:         make(map[string]uint64),
		},
	}

	// We only check the err after reseting the shelterCurrentScan variable because we want
	// to change the variable even if we had an error
	if err != nil {
		return err
	}

	// Retrieve the scan next execution. We only do this here because we want to be the last
	// thing from the method, avoiding that an error of this action prevent other actions to
	// run
	nextExecution, err := scheduler.NextExecutionByType(scheduler.JobTypeScan)
	if err != nil {
		// Didn't find a scan job in the scheduler, really strange! Return the error to report
		// the problem (probably by log messages)
		return err
	}
	shelterCurrentScan.ScheduledAt = nextExecution

	return nil
}

// When the injector successfully loaded a domain to be scanned, it will increment the
// number of domains in the scan information. This is important to estimate the progress
// of the scan
func LoadedDomainForScan() {
	atomic.AddUint64(&shelterCurrentScan.DomainsToBeScanned, 1)
}

// When the injector finished loading all domains, it tells the scan information
// structure. Because from this moment we are sure that the variable DomainsToBeScanned is
// not going to increment anymore and we can estimate the progress
func FinishLoadingDomainsForScan() {
	shelterCurrentScanLock.Lock()
	defer shelterCurrentScanLock.Unlock()

	shelterCurrentScan.Status = ScanStatusRunning
}

// When the collector receives a domain it tells the scan information structure to help
// predicting when the scan will ends
func FinishAnalyzingDomainForScan(withDNSSEC bool) {
	atomic.AddUint64(&shelterCurrentScan.DomainsScanned, 1)

	if withDNSSEC {
		atomic.AddUint64(&shelterCurrentScan.DomainsWithDNSSECScanned, 1)
	}
}

// Function to store scan result statistics. It can be accessed concurrently because it
// use a general lock to access the global structure
func StoreStatisticsOfTheScan(nameserverStatistics map[string]uint64,
	dsStatistics map[string]uint64) {

	shelterCurrentScanLock.Lock()
	defer shelterCurrentScanLock.Unlock()

	shelterCurrentScan.NameserverStatistics = nameserverStatistics
	shelterCurrentScan.DSStatistics = dsStatistics
}

// Function to copy the global variable and return it to allow other parts of the system
// to read it. It is necessary because the global variable needs locks for read/write
// access
func GetCurrentScan() CurrentScan {
	shelterCurrentScanLock.Lock()
	defer shelterCurrentScanLock.Unlock()

	// Copying the object
	currentScan := shelterCurrentScan

	return currentScan
}

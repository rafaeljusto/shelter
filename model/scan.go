package model

import (
	"labix.org/v2/mgo/bson"
	"time"
)

// List of possible values of a scan status. The status ScanStatusWaitingExecution,
// ScanStatusLoadingData, ScanStatusRunning will only be visible in a CurrentScan struct
const (
	ScanStatusWaitingExecution ScanStatus = 0 // The scan is going to be executed in the future
	ScanStatusLoadingData      ScanStatus = 1 // Loading the domains objects for Scan
	ScanStatusRunning          ScanStatus = 2 // The scan is current scanning the system
	ScanStatusExecuted         ScanStatus = 3 // Scan alredy finished succesfully
	ScanStatusAborted          ScanStatus = 4 // Scan had problems during the execution
)

// We keep a state from the scan to identify scans that had problem or not, and the current
// executio. This is useful for reports and a richer interface
type ScanStatus int

// Store all data related to a scan executed on the system. The statistics attributes cannot use the
// ENUM format because we cannot have a non-string key in the JSON format when saving into the
// database
type Scan struct {
	Id                      bson.ObjectId  `bson:"_id"` // Database identification
	Revision                int            // Version of the object
	Status                  ScanStatus     // Status of the scan
	StartedAt               time.Time      // Date and time that the scan started
	FinishedAt              time.Time      // Date and time that the scan finished
	LastModifiedAt          time.Time      // Last time the object was modified
	DomainsScanned          int            // Number of domains scanned
	DomainsWihDNSSECScanned int            // Number of domains with DS recods scanned
	NameserverStatistics    map[string]int // Statistics from nameserver status (text format) in percentage
	DSStatistics            map[string]int // Statistics from DS records' status (text format) in percentage
}

// CurrentScan is a Scan that is the next to be executed or is executing at this moment. The data
// from this struct is not stored until the scan is finished and become only a Scan struct. This
// should be used to tell the user (using a service) how is a progress of a scan on-the-fly
type CurrentScan struct {
	Scan                   // CurrentScan is a Scan
	DomainsToBeScanned int // Domains selected to be scanned
	Progress           int // Current progress, in percentage, of the scan execution
}

package protocol

import (
	"github.com/rafaeljusto/shelter/model"
	"time"
)

// ScanResponse structure represents the system Scan object to be returned via protocol. With this
// object the user can retrieve information about executed scans or current progress of a specific
// scan
type ScanResponse struct {
	Id                      int32             `json:"id"`                                 // Unique id to identify this scan
	Status                  string            `json:"status"`                             // Current scan situation
	StartedAt               time.Time         `json:"startedAt,omitempty"`                // Start date and time of the scan
	FinishedAt              time.Time         `json:"finishedAt,omitempty"`               // Finish date and time of the scan
	DomainsToBeScanned      uint64            `json:"domainsTOBeScanned,omitempty"`       // Number of domains to verify (scan is executing)
	DomainsScanned          uint64            `json:"domainsScanned,omitempty"`           // Number of domains already verified
	DomainsWihDNSSECScanned uint64            `json:"domainsWithDNSSECScanned,omitempty"` // Number of domains with DNSSEC already verified
	NameserverStatistics    map[string]uint64 `json:"nameserverStatistics,omitempty"`     // Domains' nameservers statistics (status and quantity)
	DSStatistics            map[string]uint64 `json:"dsStatistics,omitempty"`             // Domains' DS records statistics (status and quantity)
	Links                   []Link            `json:"links,omitempty"`                    // Links to move around the scans
}

// Convert a scan object data of the system into a format easy to interpret by the user
func ScanToScanResponse(scan model.Scan) ScanResponse {
	return ScanResponse{
		Id:                      scan.Id.Counter(),
		Status:                  model.ScanStatusToString(scan.Status),
		StartedAt:               scan.StartedAt,
		FinishedAt:              scan.FinishedAt,
		DomainsToBeScanned:      0,
		DomainsScanned:          scan.DomainsScanned,
		DomainsWihDNSSECScanned: scan.DomainsWihDNSSECScanned,
		NameserverStatistics:    scan.NameserverStatistics,
		DSStatistics:            scan.DSStatistics,
	}

	// TODO: Add links
}

// Convert a current scan object data being executed of the system into a format easy to interpret
// by the user
func CurrentScanToScanResponse(currentScan model.CurrentScan) ScanResponse {
	return ScanResponse{
		Id:                      0,
		Status:                  model.ScanStatusToString(currentScan.Status),
		StartedAt:               currentScan.StartedAt,
		FinishedAt:              currentScan.FinishedAt,
		DomainsToBeScanned:      currentScan.DomainsToBeScanned,
		DomainsScanned:          currentScan.DomainsScanned,
		DomainsWihDNSSECScanned: currentScan.DomainsWihDNSSECScanned,
		NameserverStatistics:    currentScan.NameserverStatistics,
		DSStatistics:            currentScan.DSStatistics,
	}

	// TODO: Add links
}

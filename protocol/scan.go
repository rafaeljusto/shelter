// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

// ScanResponse structure represents the system Scan object to be returned via protocol. With this
// object the user can retrieve information about executed scans or current progress of a specific
// scan
type ScanResponse struct {
	Status                   string            `json:"status"`                             // Current scan situation
	ScheduledAt              PreciseTime       `json:"scheduledAt,omitempty"`              // Scheduled date and time that the scan will be executed
	StartedAt                PreciseTime       `json:"startedAt,omitempty"`                // Start date and time of the scan, is also used to identify the scan
	FinishedAt               PreciseTime       `json:"finishedAt,omitempty"`               // Finish date and time of the scan
	DomainsToBeScanned       uint64            `json:"domainsToBeScanned,omitempty"`       // Number of domains to verify (scan is executing)
	DomainsScanned           uint64            `json:"domainsScanned,omitempty"`           // Number of domains already verified
	DomainsWithDNSSECScanned uint64            `json:"domainsWithDNSSECScanned,omitempty"` // Number of domains with DNSSEC already verified
	NameserverStatistics     map[string]uint64 `json:"nameserverStatistics,omitempty"`     // Domains' nameservers statistics (status and quantity)
	DSStatistics             map[string]uint64 `json:"dsStatistics,omitempty"`             // Domains' DS records statistics (status and quantity)
	Links                    []Link            `json:"links,omitempty"`                    // Links to move around the scans
}

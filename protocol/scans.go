// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

// ScansResponse store multiple scans objects with pagination support
type ScansResponse struct {
	Page          int            `json:"page"`            // Current page selected
	PageSize      int            `json:"pageSize"`        // Number of scans in a page
	NumberOfPages int            `json:"numberOfPages"`   // Total number of pages for the result set
	NumberOfItems int            `json:"numberOfItems"`   // Total number of scans in the result set
	Scans         []ScanResponse `json:"scans,omitempty"` // List of scan objects for the current page
	Links         []Link         `json:"links,omitempty"` // Links for pagination managment
}

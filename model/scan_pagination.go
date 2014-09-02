// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"strings"
)

// List of possible fields that can be used to order a result set
const (
	ScanOrderByFieldStartedAt                ScanOrderByField = 0 // Order by scan's begin time
	ScanOrderByFieldDomainsScanned           ScanOrderByField = 1 // Order by the number of domains scanned
	ScanOrderByFieldDomainsWithDNSSECScanned ScanOrderByField = 2 // Order by the number of domains with DNSSEC scanned
)

// Enumerate definition for the OrderBy so that we can limit the fields that the user can
// use in a query
type ScanOrderByField int

// Default values when the user don't define pagination. After watching a presentation
// from layer7 at http://www.layer7tech.com/tutorials/api-pagination-tutorial I agree that
// when the user don't define the pagination we shouldn't return all the result set,
// instead we assume default pagination values
var (
	ScanDefaultPaginationOrderBy = []ScanSort{
		{
			Field:     ScanOrderByFieldStartedAt, // Default ordering is by begin time
			Direction: OrderByDirectionAscending, // Default ordering is ascending
		},
	}
)

// Convert the Scan order by field from string into enum. If the string is unknown false will be
// returned. The string is case insensitive and spaces around it are ignored
func ScanOrderByFieldFromString(value string) (ScanOrderByField, bool) {
	value = strings.ToLower(value)
	value = strings.TrimSpace(value)

	switch value {
	case "startedat":
		return ScanOrderByFieldStartedAt, true
	case "domainsscanned":
		return ScanOrderByFieldDomainsScanned, true
	case "domainswithdnssecscanned":
		return ScanOrderByFieldDomainsWithDNSSECScanned, true
	}

	return ScanOrderByFieldStartedAt, false
}

// Convert the Scan order by field from enum into string. If the enum is unknown this
// method will return an empty string
func ScanOrderByFieldToString(value ScanOrderByField) string {
	switch value {
	case ScanOrderByFieldStartedAt:
		return "startedat"

	case ScanOrderByFieldDomainsScanned:
		return "domainsscanned"

	case ScanOrderByFieldDomainsWithDNSSECScanned:
		return "domainswithdnssecscanned"
	}

	return ""
}

// ScanPagination was created as a necessity for big result sets that needs to be
// sent for an end-user. With pagination we can control the size of the data and make it
// faster for the user to interact with it in a web interface as example
type ScanPagination struct {
	OrderBy       []ScanSort // Sort the list before the pagination
	PageSize      int        // Number of items that are going to be considered in one page
	Page          int        // Current page that will be returned
	NumberOfItems int        // Total number of items in the result set
	NumberOfPages int        // Total number of pages calculated for the current result set
}

// ScanSort is an object responsable to relate the order by field and direction. Each
// field used for sort, can be sorted in both directions
type ScanSort struct {
	Field     ScanOrderByField // Field to be sorted
	Direction OrderByDirection // Direction used in the sort
}

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
	DomainOrderByFieldFQDN           DomainOrderByField = 0 // Order by domain's FQDN
	DomainOrderByFieldLastModifiedAt DomainOrderByField = 1 // Order by the last modification date of the domain object
)

// Enumerate definition for the OrderBy so that we can limit the fields that the user can
// use in a query
type DomainOrderByField int

// Default values when the user don't define pagination. After watching a presentation
// from layer7 at http://www.layer7tech.com/tutorials/api-pagination-tutorial I agree that
// when the user don't define the pagination we shouldn't return all the result set,
// instead we assume default pagination values
var (
	DomainDefaultPaginationOrderBy = []DomainSort{
		{
			Field:     DomainOrderByFieldFQDN,    // Default ordering is by FQDN
			Direction: OrderByDirectionAscending, // Default ordering is ascending
		},
	}
)

// Convert the Domain order by field from string into enum. If the string is unknown false will be
// returned. The string is case insensitive and spaces around it are ignored
func DomainOrderByFieldFromString(value string) (DomainOrderByField, bool) {
	value = strings.ToLower(value)
	value = strings.TrimSpace(value)

	switch value {
	case "fqdn":
		return DomainOrderByFieldFQDN, true
	case "lastmodified":
		return DomainOrderByFieldLastModifiedAt, true
	}

	return DomainOrderByFieldFQDN, false
}

// Convert the Domain order by field from enum into string. If the enum is unknown this
// method will return an empty string
func DomainOrderByFieldToString(value DomainOrderByField) string {
	switch value {
	case DomainOrderByFieldFQDN:
		return "fqdn"

	case DomainOrderByFieldLastModifiedAt:
		return "lastmodified"
	}

	return ""
}

// DomainPagination was created as a necessity for big result sets that needs to be
// sent for an end-user. With pagination we can control the size of the data and make it
// faster for the user to interact with it in a web interface as example
type DomainPagination struct {
	OrderBy       []DomainSort // Sort the list before the pagination
	PageSize      int          // Number of items that are going to be considered in one page
	Page          int          // Current page that will be returned
	NumberOfItems int          // Total number of items in the result set
	NumberOfPages int          // Total number of pages calculated for the current result set
}

// DomainSort is an object responsable to relate the order by field and direction. Each
// field used for sort, can be sorted in both directions
type DomainSort struct {
	Field     DomainOrderByField // Field to be sorted
	Direction OrderByDirection   // Direction used in the sort
}

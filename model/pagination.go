// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"strings"
)

var (
	DefaultPaginationPageSize = 20 // By default we show 20 items per page
	DefaultPaginationPage     = 1  // By default we show the first page
)

// List of possible directions of each field in an order by query
const (
	OrderByDirectionAscending  OrderByDirection = 1  // From lower to higher
	OrderByDirectionDescending OrderByDirection = -1 // From Higher to lower
)

// Enumerate definition for the OrderBy so that we can make it easy to determinate the direction of
// an order by field
type OrderByDirection int

// Convert the  order by direction from string into enum. If the string is unknown false will be
// returned. The string is case insensitive and spaces around it are ignored
func OrderByDirectionFromString(value string) (OrderByDirection, bool) {
	value = strings.ToLower(value)
	value = strings.TrimSpace(value)

	switch value {
	case "asc":
		return OrderByDirectionAscending, true
	case "desc":
		return OrderByDirectionDescending, true
	}

	return OrderByDirectionAscending, false
}

// Convert the  order by direction from enum into string. If the enum is unknown this method will
// return an empty string
func OrderByDirectionToString(value OrderByDirection) string {
	switch value {
	case OrderByDirectionAscending:
		return "asc"

	case OrderByDirectionDescending:
		return "desc"
	}

	return ""
}

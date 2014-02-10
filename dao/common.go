package dao

import (
	"errors"
	"strings"
)

var (
	// An invalid order by direction was given to be converted in one of the known order by
	// fields of the DAO
	ErrDAOOrderByDirectionUnknown = errors.New("Unknown order by direction")
)

var (
	defaultPaginationPageSize = 20 // By default we show 20 items per page
	defaultPaginationPage     = 1  // By default we show the first page
)

// List of possible directions of each field in an order by query
const (
	DAOOrderByDirectionAscending  DAOOrderByDirection = 1  // From lower to higher
	DAOOrderByDirectionDescending DAOOrderByDirection = -1 // From Higher to lower
)

// Enumerate definition for the OrderBy so that we can make it easy to determinate the direction of
// an order by field
type DAOOrderByDirection int

// Convert the DAO order by direction from string into enum. If the string is unknown an error will
// be returned. The string is case insensitive and spaces around it are ignored
func DAOOrderByDirectionFromString(value string) (DAOOrderByDirection, error) {
	value = strings.ToLower(value)
	value = strings.TrimSpace(value)

	switch value {
	case "asc":
		return DAOOrderByDirectionAscending, nil
	case "desc":
		return DAOOrderByDirectionDescending, nil
	}

	return DAOOrderByDirectionAscending, ErrDAOOrderByDirectionUnknown
}

// Convert the DAO order by direction from enum into string. If the enum is unknown this method will
// return an empty string
func DAOOrderByDirectionToString(value DAOOrderByDirection) string {
	switch value {
	case DAOOrderByDirectionAscending:
		return "asc"

	case DAOOrderByDirectionDescending:
		return "desc"
	}

	return ""
}

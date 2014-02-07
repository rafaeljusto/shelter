package handler

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	// fqdnFlexible is used to find FQDNs in the URI. It is flexible because it does not follow all
	// FQDN rules to identify it, we will check FQDN format on low layers, here we are only interested
	// on running the basic checks for the handler
	fqdnFlexible = regexp.MustCompile(`/([[:alnum:]]|\-|\.)+\.([[:alnum:]]|\-|\.)*`)

	// isCurrentScan identify if the user desires to check only the scan in progress. As the current
	// scan doesn't have an id, we use the special world "current". We use the option i for case
	// insensitive
	isCurrentScan = regexp.MustCompile(`.+/(?i:current)$`)

	// dateFlexible retrieve a date from the URI. It is flexible because it doesn't follow exactly all
	// the RFC3339 with nanoseconds rules. We use the option i for case insensitive on the T and Z
	// letters
	dateFlexible = regexp.MustCompile(`/[0-9]+\-[0-9]+\-[0-9]+(?i:T)[0-9]+\:[0-9]+\:[0-9]+(\.[0-9]+)?(((\-)?[0-9]+\:[0-9]+)|(?i:Z))?`)
)

var (
	// Error throw when we don't find any valid date in the URI
	ErrDateNotFound = errors.New("Did not found a valid RFC 3339 date in the URI")
)

// Retrieve the FQDN from the URI. When there're multiple FQDNs in the URI it will return the last
// one. If there're no FQDNs an empty string is returned
func getFQDNFromURI(uri string) string {
	matches := fqdnFlexible.FindAllString(uri, -1)
	if matches == nil || len(matches) == 0 {
		return ""
	}

	// As the regular expression returns multiple FQDNs matchs in the URI, we are going to return the
	// last FQDN match of the URI, that's probably what the handler wants. Also we cannot return the
	// first character '/' that was used in the regular expression to identify more precisly the FQDN
	return matches[len(matches)-1][1:]
}

// Retrieve the scan identifier from the URI, can be the started time of the scan, that is used as a
// unique field to identify the scan, or the special word "current" that is used to reference the
// scan in progress
func getScanIdFromURI(uri string) (date time.Time, current bool, err error) {
	// Search for the special word "current"
	if isCurrentScan.MatchString(uri) {
		current = true
		return
	}

	matches := dateFlexible.FindAllString(uri, -1)
	if matches == nil || len(matches) == 0 {
		err = ErrDateNotFound
		return
	}

	// As the regular expression returns multiple date matchs in the URI, we are going to return the
	// last date match of the URI, that's probably what the handler wants. Also we cannot return the
	// first character '/' that was used in the regular expression to identify more precisly the date
	date, err = time.Parse(time.RFC3339Nano, strings.ToUpper(matches[len(matches)-1][1:]))

	return
}

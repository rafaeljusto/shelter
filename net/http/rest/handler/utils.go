package handler

import (
	"regexp"
)

var (
	// isFQDNFlexible is used to find FQDNs in the URI. It is flexible because it does not follow all
	// FQDN rules to identify it, we will check FQDN format on low layers, here we are only interested
	// on running the basic checks for the handler
	isFQDNFlexible = regexp.MustCompile(`/([[:alnum:]]|\-|\.)+\.([[:alnum:]]|\-|\.)*`)
)

// Retrieve the FQDN from the URI. When there're multiple FQDNs in the URI it will return the last
// one. If there're no FQDNs an empty string is returned
func getFQDNFromURI(uri string) string {
	matches := isFQDNFlexible.FindAllString(uri, -1)
	if matches == nil {
		return ""
	}

	// As the regular expression returns multiple FQDNs matchs in the URI, we are going to return the
	// last FQDN match of the URI, that's probably what the handler wants. Also we cannot return the
	// first character '/' that was used in the regular expression to identify more precisly the FQDN
	return matches[len(matches)-1][1:]
}

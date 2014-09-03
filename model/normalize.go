// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"code.google.com/p/go.net/idna"
	"regexp"
	"strings"
)

var (
	// RFC 1034 and RFC 1123
	isFQDN = regexp.MustCompile(`^(([[:alnum:]](([[:alnum:]]|\-){0,61}[[:alnum:]])?\.)*[[:alnum:]](([[:alnum:]]|\-){0,61}[[:alnum:]])?)?(\.)?$`)
)

// Normalize the domain name to have always the same mask. The following rules applied
// are, all in lower case, no spaces in the edges (spaces in the middle are going to be
// detected by other validation mechanisms), dot at the end of the name and the ASCII
// format (for IDNA domains)
func NormalizeDomainName(domainName string) (string, bool) {
	domainName = strings.ToLower(domainName)
	domainName = strings.TrimSpace(domainName)

	// We will not add the dot for empty domain name so that become root zone ("."), because
	// maybe the user forgot to inform the domain name and we don't want to make it easy for
	// him to pass in the validation verifications
	if len(domainName) > 0 && !strings.HasSuffix(domainName, ".") {
		domainName += "."
	}

	// We always manage domains in ASCII format, so we convert unicode domains, that use
	// accents according to IDNA rules
	var err error
	domainName, err = idna.ToASCII(domainName)
	if err != nil {
		return domainName, false
	}

	// Check FQDN format
	if !isFQDN.MatchString(domainName) {
		return domainName, false
	}

	return domainName, true
}

// We will store the digest always in lower case. According to RFC 4034 (section 5.3):
// "Digest MUST be represented as a sequence of case-insensitive hexadecimal digits.
// Whitespace is allowed within the hexadecimal text.", so there's no problem changing the
// letter case
func NormalizeDSDigest(digest string) string {
	digest = strings.ToLower(digest)
	digest = strings.TrimSpace(digest)
	return digest
}

// RFC 1766: The ISO 639/ISO 3166 convention is that language names are written in lower
// case, while country codes are written in upper case. This convention is recommended,
// but not enforced; the tags are case insensitive
func NormalizeLanguage(language string) string {
	language = strings.TrimSpace(language)
	languageParts := strings.Split(language, "-")
	languageParts[0] = strings.TrimSpace(languageParts[0])

	// Region or country code can be optional, also empty strings can be caught here
	if len(languageParts) == 1 {
		return strings.ToLower(languageParts[0])
	}

	languageParts[1] = strings.TrimSpace(languageParts[1])
	return strings.ToLower(languageParts[0]) + "-" + strings.ToUpper(languageParts[1])
}

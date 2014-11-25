// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package language keep track of all IANA languages
package language

import (
	"strings"
)

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

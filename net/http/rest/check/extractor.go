// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package check verify REST policies
package check

import (
	"net/http"
	"strings"
)

func getHTTPContentType(r *http.Request) string {
	contentType := r.Header.Get("Content-Type")
	contentType = strings.TrimSpace(contentType)
	contentType = strings.ToLower(contentType)
	return contentType
}

func getHTTPContentMD5(r *http.Request) string {
	contentMD5 := r.Header.Get("Content-MD5")
	contentMD5 = strings.TrimSpace(contentMD5)
	return contentMD5
}

func getHTTPDate(r *http.Request) string {
	dateStr := r.Header.Get("Date")
	dateStr = strings.TrimSpace(dateStr)
	return dateStr
}

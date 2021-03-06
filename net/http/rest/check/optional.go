// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package check verify REST policies
package check

import (
	"net/http"
	"strings"
	"time"
)

func HTTPIfModifiedSince(r *http.Request, lastModifiedAt time.Time) (bool, error) {
	ifModifiedSinceStr := r.Header.Get("If-Modified-Since")
	ifModifiedSinceStr = strings.TrimSpace(ifModifiedSinceStr)

	if len(ifModifiedSinceStr) == 0 {
		return true, nil
	}

	ifModifiedSince, err := time.Parse(time.RFC1123, ifModifiedSinceStr)
	if err != nil {
		return true, err
	}

	if lastModifiedAt.Before(ifModifiedSince) || lastModifiedAt.Equal(ifModifiedSince) {
		return false, nil
	}

	return true, nil
}

func HTTPIfUnmodifiedSince(r *http.Request, lastModifiedAt time.Time) (bool, error) {
	ifUnmodifiedSinceStr := r.Header.Get("If-Unmodified-Since")
	ifUnmodifiedSinceStr = strings.TrimSpace(ifUnmodifiedSinceStr)

	if len(ifUnmodifiedSinceStr) == 0 {
		return true, nil
	}

	ifUnmodifiedSince, err := time.Parse(time.RFC1123, ifUnmodifiedSinceStr)
	if err != nil {
		return true, err
	}

	if lastModifiedAt.After(ifUnmodifiedSince) {
		return false, nil
	}

	return true, nil
}

func HTTPIfMatch(r *http.Request, etag string) bool {
	ifMatch := r.Header.Get("If-Match")
	ifMatch = strings.TrimSpace(ifMatch)

	if len(ifMatch) == 0 {
		return true
	}

	ifMatchParts := strings.Split(ifMatch, ",")

	for _, ifMatchPart := range ifMatchParts {
		ifMatchPart = strings.TrimSpace(ifMatchPart)

		// If "*" is given and no current entity exists, the server MUST NOT perform the
		// requested method, and MUST return a 412 (Precondition Failed) response
		if ifMatchPart == "*" {
			return len(etag) > 0 && etag != "0"
		}

		if ifMatchPart == etag {
			return true
		}
	}

	// If none of the entity tags match the server MUST NOT perform the requested method,
	// and MUST return a 412 (Precondition Failed) response
	return false
}

func HTTPIfNoneMatch(r *http.Request, etag string) bool {
	ifNoneMatch := r.Header.Get("If-None-Match")
	ifNoneMatch = strings.TrimSpace(ifNoneMatch)

	if len(ifNoneMatch) == 0 {
		return true
	}

	ifNoneMatchParts := strings.Split(ifNoneMatch, ",")

	for _, ifNoneMatchPart := range ifNoneMatchParts {
		ifNoneMatchPart = strings.TrimSpace(ifNoneMatchPart)

		// if "*" is given and any current entity exists for that resource, then the server
		// MUST NOT perform the requested method, unless required to do so because the
		// resource's modification date fails to match that supplied in an If-Modified-Since
		// header field in the request
		if ifNoneMatchPart == "*" {
			return len(etag) == 0 || etag == "0"
		}

		if ifNoneMatchPart == etag {
			// Instead, if the request method was GET or HEAD, the server SHOULD respond with a
			// 304 (Not Modified) response, including the cache-related header fields
			// (particularly ETag) of one of the entities that matched. For all other request
			// methods, the server MUST respond with a status of 412 (Precondition Failed)
			return false
		}
	}

	return true
}

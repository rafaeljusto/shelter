// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"net/http"
	"testing"
	"time"
)

func TestHTTPCacheHeaders(t *testing.T) {
	r, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("If-Modified-Since", time.Now().Format(time.RFC1123))
	r.Header.Set("If-Unmodified-Since", time.Now().Add(10*time.Second).Format(time.RFC1123))
	r.Header.Set("If-Match", "2")
	r.Header.Set("If-None-Match", "1")

	c, err := context.NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	lastModified := time.Now().Add(2 * time.Second)
	if !CheckHTTPCacheHeaders(r, &c, lastModified, 2) {
		t.Error("Not detecting a valid scenario")
	}
}

func TestCheckIfModifiedSince(t *testing.T) {
	r, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("If-Modified-Since", time.Now().Format(time.RFC1123))

	c, err := context.NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	lastModified := time.Now().Add(-10 * time.Second)
	if checkIfModifiedSince(r, &c, lastModified) {
		t.Error("Not detecting when last modified is lower than If-Modified-Since")
	}

	if c.ResponseHTTPStatus != http.StatusNotModified {
		t.Error("Not returning correct status when last modified is lower than If-Modified-Since")
	}

	lastModified = time.Now().Add(10 * time.Second)
	if !checkIfModifiedSince(r, &c, lastModified) {
		t.Error("Not detecting when last modified is greater than If-Modified-Since")
	}

	r.Header.Set("If-Modified-Since", "This is not valid!")
	if checkIfModifiedSince(r, &c, lastModified) {
		t.Error("Not detecting invalid header field")
	}

	if c.ResponseHTTPStatus != http.StatusBadRequest {
		t.Error("Not setting bad request status on error")
	}
}

func TestCheckIfUnmodifiedSince(t *testing.T) {
	r, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("If-Unmodified-Since", time.Now().Format(time.RFC1123))

	c, err := context.NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	lastModified := time.Now().Add(-10 * time.Second)
	if !checkIfUnmodifiedSince(r, &c, lastModified) {
		t.Error("Not detecting when last modified is lower than If-Unmodified-Since")
	}

	lastModified = time.Now().Add(10 * time.Second)
	if checkIfUnmodifiedSince(r, &c, lastModified) {
		t.Error("Not detecting when last modified is greater than If-Unmodified-Since")
	}

	if c.ResponseHTTPStatus != http.StatusPreconditionFailed {
		t.Error("Not returning correct status when last modified is lower than If-Unmodified-Since")
	}

	r.Header.Set("If-Unmodified-Since", "This is not valid!")
	if checkIfUnmodifiedSince(r, &c, lastModified) {
		t.Error("Not detecting invalid header field")
	}
}

func TestCheckIfMatch(t *testing.T) {
	r, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("If-Match", "1")

	c, err := context.NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !checkIfMatch(r, &c, 1) {
		t.Error("Not detecting when If-Match matched")
	}

	if checkIfMatch(r, &c, 2) {
		t.Error("Not detecting when If-Match did not match")
	}

	if c.ResponseHTTPStatus != http.StatusPreconditionFailed {
		t.Error("Not returning correct status when If-Match did not match")
	}

	r.Header.Set("If-Match", "This is not an revision identifier")
	if checkIfMatch(r, &c, 1) {
		t.Error("Not detecting invalid header field")
	}
}

func TestCheckIfNoneMatch(t *testing.T) {
	r, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("If-None-Match", "1")

	c, err := context.NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	if checkIfNoneMatch(r, &c, 1) {
		t.Error("Not detecting when If-None-Match matched")
	}

	if c.ResponseHTTPStatus != http.StatusNotModified {
		t.Error("Not returning correct status when If-None-Match matched for GET")
	}

	if !checkIfNoneMatch(r, &c, 2) {
		t.Error("Not detecting when If-None-Match did not match")
	}

	r.Header.Set("If-None-Match", "This is not an revision identifier")
	if checkIfNoneMatch(r, &c, 1) {
		t.Error("Not detecting invalid header field")
	}

	r, err = http.NewRequest("PUT", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("If-None-Match", "1")

	if checkIfNoneMatch(r, &c, 1) {
		t.Error("Not detecting when If-None-Match matched")
	}

	if c.ResponseHTTPStatus != http.StatusPreconditionFailed {
		t.Error("Not returning correct status when If-None-Match matched for PUT")
	}
}

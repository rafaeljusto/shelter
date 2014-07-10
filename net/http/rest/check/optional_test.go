// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package check verify REST policies
package check

import (
	"net/http"
	"testing"
	"time"
)

func TestHTTPIfModifiedSince(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("If-Modified-Since", "")
	if ok, err := HTTPIfModifiedSince(r, time.Now()); !ok || err != nil {
		t.Error("Not accepting when there's no HTTP If-Modified-Since header field")
	}

	r.Header.Set("If-Modified-Since", "2006-01-02T15:04:05Z07:00")
	if ok, err := HTTPIfModifiedSince(r, time.Now()); !ok || err == nil {
		t.Error("Accepting an invalid date format. Should only support RFC 1123")
	}

	r.Header.Set("If-Modified-Since", time.Now().Add(10*time.Second).Format(time.RFC1123))
	if ok, _ := HTTPIfModifiedSince(r, time.Now()); ok {
		t.Error("Not comparing properly when object was modified")
	}

	r.Header.Set("If-Modified-Since", time.Now().Format(time.RFC1123))
	if ok, _ := HTTPIfModifiedSince(r, time.Now().Add(1*time.Second)); !ok {
		t.Error("Not comparing properly when object was modified")
	}
}

func TestHTTPIfUnmodifiedSince(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("If-Unmodified-Since", "")
	if ok, err := HTTPIfUnmodifiedSince(r, time.Now()); !ok || err != nil {
		t.Error("Not accepting when there's no HTTP If-Unmodified-Since header field")
	}

	r.Header.Set("If-Unmodified-Since", "2006-01-02T15:04:05Z07:00")
	if ok, err := HTTPIfUnmodifiedSince(r, time.Now()); !ok || err == nil {
		t.Error("Accepting an invalid date format. Should only support RFC 1123")
	}

	r.Header.Set("If-Unmodified-Since", time.Now().Format(time.RFC1123))
	if ok, _ := HTTPIfUnmodifiedSince(r, time.Now().Add(1*time.Second)); ok {
		t.Error("Not comparing properly when object was modified")
	}

	r.Header.Set("If-Unmodified-Since", time.Now().Add(1*time.Second).Format(time.RFC1123))
	if ok, _ := HTTPIfUnmodifiedSince(r, time.Now()); !ok {
		t.Error("Not comparing properly when object was modified")
	}
}

func TestHTTPIfMatch(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("If-Match", "")
	if !HTTPIfMatch(r, "1") {
		t.Error("Not accepting when there's no HTTP If-Match header field")
	}

	r.Header.Set("If-Match", "2")
	if HTTPIfMatch(r, "1") {
		t.Error("Match with an ETag that does not exists")
	}

	r.Header.Set("If-Match", "1")
	if !HTTPIfMatch(r, "1") {
		t.Error("Not accepting a valid ETag in If-Match header field")
	}

	r.Header.Set("If-Match", "*")
	if !HTTPIfMatch(r, "1") {
		t.Error("Not accepting a valid ETag wildcard in If-Match header field")
	}

	r.Header.Set("If-Match", "*")
	if HTTPIfMatch(r, "0") {
		t.Error("Not returning the correct response when the entity " +
			"does not exist in the system with a valid ETag wildcard " +
			"in If-Match header field")
	}
}

func TestHTTPIfNoneMatch(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("If-None-Match", "")
	if !HTTPIfNoneMatch(r, "1") {
		t.Error("Not accepting when there's no HTTP If-None-Match header field")
	}

	r.Header.Set("If-None-Match", "1")
	if HTTPIfNoneMatch(r, "1") {
		t.Error("Match with an ETag should fail the If-None-Match condition")
	}

	r.Header.Set("If-None-Match", "2")
	if !HTTPIfNoneMatch(r, "1") {
		t.Error("Not accepting a valid ETag in If-None-Match header field")
	}

	r.Header.Set("If-None-Match", "*")
	if HTTPIfNoneMatch(r, "1") {
		t.Error("Not accepting a valid ETag wildcard in If-None-Match header field")
	}

	r.Header.Set("If-None-Match", "*")
	if !HTTPIfNoneMatch(r, "0") {
		t.Error("Not returning the correct response when the entity " +
			"does not exist in the system with a valid ETag wildcard " +
			"in If-None-Match header field")
	}
}

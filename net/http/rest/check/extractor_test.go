// check - Verify REST policies
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package check

import (
	"net/http"
	"testing"
)

func TestGetHTTPContentType(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("Content-Type", "   tExt/pLaIn   ")
	if getHTTPContentType(r) != "text/plain" {
		t.Error("Not retrieving content type HTTP header properly")
	}
}

func TestGetHTTPContentMD5(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("Content-MD5", "   ZDhlOGZjYTJkYzBmODk2ZmQ3Y2I0Y2IwMDMxYmEyNDkgIC0K==   ")
	if getHTTPContentMD5(r) != "ZDhlOGZjYTJkYzBmODk2ZmQ3Y2I0Y2IwMDMxYmEyNDkgIC0K==" {
		t.Error("Not retrieving content MD5 HTTP header properly")
	}
}

func TestGetHTTPDate(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("Date", "   Mon, 02 Jan 2006 15:04:05 MST   ")
	if getHTTPDate(r) != "Mon, 02 Jan 2006 15:04:05 MST" {
		t.Error("Not retrieving date HTTP header properly")
	}
}

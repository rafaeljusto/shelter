// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRemoteAddress(t *testing.T) {
	var handler RemoteAddressCompliant
	i := NewRemoteAddress(&handler)
	if i.handler == nil {
		t.Error("Not defining the handler in the constructor")
	}
}

func TestRemoteAddressBefore(t *testing.T) {
	data := []struct {
		description     string
		header          map[string]string
		remoteAddress   string
		expectedStatus  int
		expectedAddress net.IP
	}{
		{
			description: "X-Forwarded-For defined with one address",
			header: map[string]string{
				"X-Forwarded-For": "192.168.1.1",
			},
			expectedStatus:  http.StatusOK,
			expectedAddress: net.ParseIP("192.168.1.1"),
		},
		{
			description: "X-Forwarded-For defined with multiple addresses with spaces",
			header: map[string]string{
				"X-Forwarded-For": "    192.168.1.1    ,   192.168.1.2   ",
			},
			expectedStatus:  http.StatusOK,
			expectedAddress: net.ParseIP("192.168.1.1"),
		},
		{
			description: "X-Forwarded-For defined with multiple addresses",
			header: map[string]string{
				"X-Forwarded-For": "192.168.1.1,192.168.1.2,192.168.1.3,192.168.1.4",
			},
			expectedStatus:  http.StatusOK,
			expectedAddress: net.ParseIP("192.168.1.3"),
		},
		{
			description: "X-Real-IP defined with spaces",
			header: map[string]string{
				"X-Real-IP": "  192.168.1.5  ",
			},
			expectedStatus:  http.StatusOK,
			expectedAddress: net.ParseIP("192.168.1.5"),
		},
		{
			description:    "Invalid remote address",
			remoteAddress:  "192.168.1.6",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			description:     "Valid remote address",
			remoteAddress:   "192.168.1.6:123",
			expectedStatus:  http.StatusOK,
			expectedAddress: net.ParseIP("192.168.1.6"),
		},
	}

	for i, item := range data {
		r, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.RemoteAddr = item.remoteAddress
		for key, value := range item.header {
			r.Header.Set(key, value)
		}
		w := httptest.NewRecorder()

		var handler RemoteAddressCompliant
		remoteAddress := NewRemoteAddress(&handler)
		remoteAddress.Before(w, r)

		if w.Code != item.expectedStatus {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%d'; found '%d'",
				i,
				item.description,
				item.expectedStatus,
				w.Code,
			)
		}

		if item.expectedAddress != nil && !handler.remoteAddress.Equal(item.expectedAddress) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%s'; found '%s'",
				i,
				item.description,
				item.remoteAddress,
				handler.remoteAddress,
			)
		}
	}
}

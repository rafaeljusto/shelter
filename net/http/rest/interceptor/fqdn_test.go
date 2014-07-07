// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockFQDNHandler struct {
	FQDN        string
	Message     string
	URI         string
	ReturnError error
}

func (m *MockFQDNHandler) SetFQDN(fqdn string) {
	m.FQDN = fqdn
}

func (m *MockFQDNHandler) GetFQDN() string {
	return m.FQDN
}

func (m *MockFQDNHandler) MessageResponse(message, uri string) error {
	m.Message = message
	m.URI = uri
	return m.ReturnError
}

func TestFQDNBefore(t *testing.T) {
	data := []struct {
		FQDN            string
		ReturnError     error
		ExpectedMessage string
		ExpectedFQDN    string
	}{
		{FQDN: "   EXAMPLE.com.br   ", ExpectedMessage: "", ExpectedFQDN: "example.com.br."},
		{FQDN: "1234", ExpectedMessage: "", ExpectedFQDN: "1234."},
		{FQDN: "xn-------------", ExpectedMessage: "invalid-uri", ExpectedFQDN: "xn-------------"},
		{
			FQDN:            "xn-------------",
			ExpectedMessage: "invalid-uri",
			ExpectedFQDN:    "xn-------------",
			ReturnError:     errors.New("Low level error"),
		},
	}

	for _, item := range data {
		fqdnHandler := MockFQDNHandler{
			FQDN:        item.FQDN,
			ReturnError: item.ReturnError,
		}
		interceptor := NewFQDN(&fqdnHandler)

		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", fmt.Sprintf("/domain/%s", item.FQDN), nil)
		if err != nil {
			t.Fatal(err)
		}

		interceptor.Before(w, r)
		interceptor.After(w, r)

		if fqdnHandler.Message != item.ExpectedMessage {
			t.Errorf("Not setting the correct message. Expected '%s' and got '%s'",
				item.ExpectedMessage, fqdnHandler.Message)
		}

		if fqdnHandler.GetFQDN() != item.ExpectedFQDN {
			t.Errorf("Not normalizing FQDN properly. Expected '%s' and got '%s'",
				item.ExpectedFQDN, fqdnHandler.GetFQDN())
		}
	}
}

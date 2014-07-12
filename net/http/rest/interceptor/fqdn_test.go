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
	MessageId   string
	ReturnError error
}

func (h *MockFQDNHandler) SetFQDN(fqdn string) {
	h.FQDN = fqdn
}

func (h *MockFQDNHandler) GetFQDN() string {
	return h.FQDN
}

func (h *MockFQDNHandler) MessageResponse(messageId, roid string) error {
	h.MessageId = messageId
	return h.ReturnError
}

func TestFQDN(t *testing.T) {
	data := []struct {
		FQDN              string
		ReturnError       error
		ExpectedMessageId string
		ExpectedFQDN      string
	}{
		{FQDN: "   EXAMPLE.com.br   ", ExpectedMessageId: "", ExpectedFQDN: "example.com.br."},
		{FQDN: "1234", ExpectedMessageId: "", ExpectedFQDN: "1234."},
		{FQDN: "xn-------------", ExpectedMessageId: "invalid-uri", ExpectedFQDN: "xn-------------"},
		{
			FQDN:              "xn-------------",
			ExpectedMessageId: "invalid-uri",
			ExpectedFQDN:      "xn-------------",
			ReturnError:       errors.New("Low level error"),
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

		if fqdnHandler.MessageId != item.ExpectedMessageId {
			t.Errorf("Not setting the correct message id. Expected '%s' and got '%s'",
				item.ExpectedMessageId, fqdnHandler.MessageId)
		}

		if fqdnHandler.GetFQDN() != item.ExpectedFQDN {
			t.Errorf("Not normalizing FQDN properly. Expected '%s' and got '%s'",
				item.ExpectedFQDN, fqdnHandler.GetFQDN())
		}
	}
}

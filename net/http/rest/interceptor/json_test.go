// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"errors"
	"github.com/rafaeljusto/shelter/net/http/rest/messages"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockCodecHandler struct {
	MessageId     string
	ReturnError   error
	Request       map[string]string `request:"put"`
	Response      map[string]string `response:"get"`
	ErrorResponse *string           `error`
}

func (c *MockCodecHandler) GetLanguage() *messages.LanguagePack {
	return &messages.LanguagePack{}
}

func (c *MockCodecHandler) MessageResponse(messageId, roid string) error {
	c.MessageId = messageId
	return c.ReturnError
}

func TestJSONBefore(t *testing.T) {
	var codecHandler MockCodecHandler
	jsonCodec := NewJSONCodec(&codecHandler)

	data := []struct {
		Request           string
		ReturnError       error
		ExpectedCode      int
		ExpectedMessageId string
	}{
		{
			Request: `{
	"key1": "value1",
  "key2": "value2"		
}`,
			ExpectedCode: http.StatusOK,
		},
		{
			Request: `{
	"key1": {
		"subkey1": "value1"
	}
}`,
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "invalid-json-content",
		},
		{
			Request: `{
	"key1": {
		"subkey1": "value1"
	}
}`,
			ReturnError:       errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "invalid-json-content",
		},
	}

	for _, item := range data {
		r, err := http.NewRequest("PUT", "/keys", strings.NewReader(item.Request))
		if err != nil {
			t.Fatal(err)
		}

		codecHandler.ReturnError = item.ReturnError
		codecHandler.MessageId = "" // Reset

		w := httptest.NewRecorder()
		jsonCodec.Before(w, r)

		if w.Code != item.ExpectedCode {
			t.Errorf("Wrong status code for request '%s'. "+
				"Expected %d and got %d", item.Request, item.ExpectedCode, w.Code)
		}

		if codecHandler.MessageId != item.ExpectedMessageId {
			t.Errorf("Wrong message id for request '%s'. "+
				"Expected '%s' and got '%s'", item.Request, item.ExpectedMessageId, codecHandler.MessageId)
		}
	}
}

func TestJSONAfter(t *testing.T) {
	var codecHandler MockCodecHandler
	jsonCodec := NewJSONCodec(&codecHandler)

	data := []struct {
		Method            string
		Response          map[string]string
		ErrorResponse     string
		ReturnError       error
		ExpectedCode      int
		ExpectedMessageId string
		ExpectedBody      string
	}{
		{
			Method: "GET",
			Response: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: `{"key1":"value1","key2":"value2"}`,
		},
		{
			Method:        "GET",
			ErrorResponse: "Error",
			ExpectedCode:  http.StatusOK,
			ExpectedBody:  `"Error"`,
		},
		{
			Method:       "GET",
			ExpectedCode: http.StatusOK,
		},
		{
			Method:       "DELETE",
			ExpectedCode: http.StatusOK,
		},
	}

	for _, item := range data {
		r, err := http.NewRequest(item.Method, "/keys", nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(item.ErrorResponse) > 0 {
			codecHandler.ErrorResponse = &item.ErrorResponse
		} else {
			codecHandler.ErrorResponse = nil
		}

		codecHandler.Response = item.Response
		codecHandler.ReturnError = item.ReturnError
		codecHandler.MessageId = "" // Reset

		w := httptest.NewRecorder()
		jsonCodec.Before(w, r)
		jsonCodec.After(w, r)

		if w.Code != item.ExpectedCode {
			t.Errorf("Wrong status code for response '%v'. "+
				"Expected %d and got %d", item.Response, item.ExpectedCode, w.Code)
		}

		if w.Body.String() != item.ExpectedBody {
			t.Errorf("Wrong body for response '%v'. "+
				"Expected '%s' and got '%s'", item.Response, item.ExpectedBody, w.Body.String())
		}

		if codecHandler.MessageId != item.ExpectedMessageId {
			t.Errorf("Wrong message id for response '%v'. "+
				"Expected %s and got %s", item.Response, item.ExpectedMessageId, codecHandler.MessageId)
		}
	}
}

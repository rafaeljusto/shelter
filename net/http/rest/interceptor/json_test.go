// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"bytes"
	encjson "encoding/json"
	"github.com/rafaeljusto/shelter/protocol"
	"github.com/rafaeljusto/shelter/testing/types"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type MockJSONHandler interface {
	requestResponser
	setResponse(map[string]string)
}

type MockJSONHandlerWithResponse struct {
	RemoteAddressCompliant
	LogCompliant
	JSONCompliant
	language string
	Req      map[string]string `request:"put"`
	Resp     map[string]string `response:"get,head"`
}

func (h MockJSONHandlerWithResponse) Language() string {
	return h.language
}

func (h *MockJSONHandlerWithResponse) setResponse(response map[string]string) {
	h.Resp = response
}

type MockJSONHandlerWithoutResponse struct {
	RemoteAddressCompliant
	LogCompliant
	JSONCompliant
	language string
}

func (h MockJSONHandlerWithoutResponse) Language() string {
	return h.language
}

func (h *MockJSONHandlerWithoutResponse) setResponse(response map[string]string) {
	// Dummy
}

func TestJSONBefore(t *testing.T) {
	var handler MockJSONHandlerWithResponse
	json := NewJSON(&handler)

	data := []struct {
		Req           string
		ExpectedCode  int
		ExpectedError bool
	}{
		{
			Req: `{
	"key1": "value1",
  "key2": "value2"		
}`,
			ExpectedCode: http.StatusOK,
		},
		{
			Req: `{
	"key1": {
		"subkey1": "value1"
	}
}`,
			ExpectedCode:  http.StatusBadRequest,
			ExpectedError: true,
		},
	}

	for i, item := range data {
		r, err := http.NewRequest("PUT", "/keys", strings.NewReader(item.Req))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		json.Before(w, r)

		if w.Code != item.ExpectedCode {
			t.Errorf("Wrong status code for item '%d'. "+
				"Expected %d and got %d", i, item.ExpectedCode, w.Code)
		}

		if item.ExpectedError && handler.Message() == nil {
			t.Errorf("Expected error for item '%d'", i)

		} else if !item.ExpectedError && handler.Message() != nil {
			t.Errorf("Did not expected error for item '%d': %+v", i, handler.Message())
		}
	}
}

type mockTranslatorWithBrokenMarshal struct {
	types.BrokenMarshaler
}

func (mockTranslatorWithBrokenMarshal) Translate(string) bool {
	return true
}

func TestJSONAfter(t *testing.T) {
	protocol.Translations["en"] = protocol.Translation{
		"corrupted-data": "Test!",
	}

	data := []struct {
		Handler       MockJSONHandler
		Headers       map[string]string
		RequestBody   io.ReadCloser
		Resp          map[string]string
		Response      interface{}
		ExpectingBody bool
		Message       protocol.Translator
		ExpectedCode  int
		HTTPMethod    string
	}{
		{
			Handler: &MockJSONHandlerWithResponse{},
			Resp: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			ExpectingBody: true,
			ExpectedCode:  http.StatusOK,
			HTTPMethod:    "GET",
		},
		{
			Handler: &MockJSONHandlerWithResponse{},
			Message: &protocol.Message{
				Code:  protocol.MsgCodeCorruptedData,
				Value: "test",
			},
			ExpectingBody: true,
			ExpectedCode:  http.StatusOK, // status code is defined by the handler
			HTTPMethod:    "GET",
		},
		{
			Handler: &MockJSONHandlerWithResponse{},
			Headers: map[string]string{
				"Accept-Language": "idontexist",
			},
			Message: &protocol.Message{
				Code:  protocol.MsgCodeCorruptedData,
				Value: "test",
			},
			ExpectedCode: http.StatusNotAcceptable,
			HTTPMethod:   "GET",
		},
		{
			Handler: &MockJSONHandlerWithResponse{},
			Message: protocol.Messages{
				protocol.NewMessage(protocol.MsgCodeCorruptedData, "test1"),
				protocol.NewMessage(protocol.MsgCodeCorruptedData, "test2"),
			},
			ExpectingBody: true,
			ExpectedCode:  http.StatusOK, // status code is defined by the handler
			HTTPMethod:    "GET",
		},
		{
			Handler: &MockJSONHandlerWithResponse{},
			Headers: map[string]string{
				"Accept-Language": "idontexist",
			},
			Message: protocol.Messages{
				protocol.NewMessage(protocol.MsgCodeCorruptedData, "test1"),
				protocol.NewMessage(protocol.MsgCodeCorruptedData, "test2"),
			},
			ExpectedCode: http.StatusNotAcceptable,
			HTTPMethod:   "GET",
		},
		{
			Handler:      &MockJSONHandlerWithResponse{},
			ExpectedCode: http.StatusOK,
			HTTPMethod:   "GET",
		},
		{
			Handler:      &MockJSONHandlerWithoutResponse{},
			ExpectedCode: http.StatusOK,
			HTTPMethod:   "GET",
		},
		{
			Handler:      &MockJSONHandlerWithoutResponse{},
			RequestBody:  types.NewBrokenCloser(),
			ExpectedCode: http.StatusOK,
			HTTPMethod:   "GET",
		},
		{
			Handler:      &MockJSONHandlerWithoutResponse{},
			Message:      mockTranslatorWithBrokenMarshal{},
			ExpectedCode: http.StatusInternalServerError,
			HTTPMethod:   "GET",
		},
		{
			Handler:      &MockJSONHandlerWithoutResponse{},
			ExpectedCode: http.StatusInternalServerError,
			Response:     types.NewBrokenMarshaler(),
			HTTPMethod:   "GET",
		},
		{
			Handler: &MockJSONHandlerWithResponse{},
			Resp: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			ExpectedCode: http.StatusOK,
			HTTPMethod:   "HEAD",
		},
	}

	for i, item := range data {
		json := NewJSON(item.Handler)

		r, err := http.NewRequest(item.HTTPMethod, "/keys", item.RequestBody)
		if err != nil {
			t.Fatal(err)
		}

		for key, value := range item.Headers {
			r.Header.Set(key, value)
		}

		w := httptest.NewRecorder()

		if item.Response != nil {
			item.Handler.(*MockJSONHandlerWithoutResponse).response = reflect.ValueOf(item.Response)
		} else {
			item.Handler.setResponse(item.Resp)
		}
		item.Handler.SetMessage(item.Message)

		json.Before(w, r)
		json.After(w, r)

		if w.Code != item.ExpectedCode {
			t.Errorf("Wrong status code for item '%d'. "+
				"Expected %d and got %d", i, item.ExpectedCode, w.Code)
		}

		var expectedBody []byte

		if item.ExpectingBody {
			if item.Resp != nil {
				expectedBody, err = encjson.Marshal(item.Resp)
				if err != nil {
					t.Fatal(err)
				}

			} else if item.Message != nil {
				expectedBody, err = encjson.Marshal(item.Message)
				if err != nil {
					t.Fatal(err)
				}
			}
		}

		if !bytes.Equal(w.Body.Bytes(), expectedBody) {
			t.Errorf("Wrong body for item '%d'. "+
				"Expected '%s' and got '%s'", i, expectedBody, w.Body.String())
		}
	}
}

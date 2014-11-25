// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"fmt"
	"github.com/rafaeljusto/handy"
	"github.com/rafaeljusto/shelter/protocol"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockHTTPHandler struct {
	LogCompliant
	JSONCompliant
	RemoteAddressCompliant
}

func TestNewHTTPHeaders(t *testing.T) {
	httpHandler := NewHTTPHeaders(&MockHTTPHandler{})
	if httpHandler.handler == nil {
		t.Fatal("Handler is not being defined in the constructor")
	}
}

func TestHTTPHeadersBefore(t *testing.T) {
	data := []struct {
		description  string
		headers      map[string]string
		expectedCode int
	}{
		{
			description: "it should accept valid headers",
			headers: map[string]string{
				"Date":            time.Now().Format(time.RFC1123),
				"Content-Type":    "APPLICATION/vnd.br+json",
				"Accept":          "APPLICATION/vnd.br+json;q=1.0,application/json,text/html",
				"Accept-Language": "pt-BR;q=1.0,en-US;q=0.8,es-ES",
				"Accept-Charset":  "UTF-8;q=1.0,iso-8859-1;q=0.8",
			},
			expectedCode: http.StatusOK,
		},
		{
			description: "it should report error for invalid headers",
			headers: map[string]string{
				"Date":            time.Now().Format(time.RFC1123),
				"Content-Type":    "XXXXXXXXXXXXXXXXX",
				"Accept":          "APPLICATION/vnd.br+json;q=1.0,application/json,text/html",
				"Accept-Language": "pt-BR;q=1.0,en-US;q=0.8,es-ES",
				"Accept-Charset":  "UTF-8;q=1.0,iso-8859-1;q=0.8",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	protocol.Translations["pt"] = protocol.Translation{}
	defer func() {
		protocol.Translations = make(map[string]protocol.Translation)
	}()

	for i, item := range data {
		r, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		for key, value := range item.headers {
			r.Header.Add(key, value)
		}
		w := httptest.NewRecorder()

		handler := NewHTTPHeaders(&MockHTTPHandler{})
		handler.Before(w, r)

		if w.Code != item.expectedCode {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%d'; found '%d'",
				i,
				item.description,
				item.expectedCode,
				w.Code,
			)
		}
	}
}

func TestHTTPHeadersAfter(t *testing.T) {
	data := []struct {
		description     string
		writer          http.ResponseWriter
		responseBody    string
		expectedHeaders map[string]string
		expectedCode    int
	}{
		{
			description:  "it should add all the necessary HTTP headers in the response",
			writer:       handy.NewBufferedResponseWriter(httptest.NewRecorder()),
			responseBody: "This is a test",
			expectedHeaders: map[string]string{
				"Content-Type":    "application/vnd.br+json; charset=utf-8",
				"Content-Length":  "14",
				"Content-Md5":     "zhFORQHS9OLc6j4XtUbzOQ==",
				"Accept":          "application/vnd.br+json",
				"Accept-Language": "pt",
				"Accept-Charset":  "utf-8",
			},
			expectedCode: 0,
		},
		{
			description:  "it should report error when it's not a handy.BufferedResponseWriter",
			writer:       httptest.NewRecorder(),
			expectedCode: http.StatusInternalServerError,
		},
	}

	protocol.Translations["pt"] = protocol.Translation{}
	defer func() {
		protocol.Translations = make(map[string]protocol.Translation)
	}()

	for i, item := range data {
		r, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(item.responseBody) > 0 {
			item.writer.Write([]byte(item.responseBody))
		}

		handler := NewHTTPHeaders(&MockHTTPHandler{})
		handler.After(item.writer, r)

		if w, ok := item.writer.(*httptest.ResponseRecorder); ok {
			if w.Code != item.expectedCode {
				t.Errorf(
					"Item %d, “%s”: mismatch results. Expecting '%d'; found '%d'",
					i,
					item.description,
					item.expectedCode,
					w.Code,
				)
			}
		}

		if w, ok := item.writer.(*handy.BufferedResponseWriter); ok {
			if w.Status() != item.expectedCode {
				t.Errorf(
					"Item %d, “%s”: mismatch results. Expecting '%d'; found '%d'",
					i,
					item.description,
					item.expectedCode,
					w.Status(),
				)
			}
		}

		headers := map[string][]string(item.writer.Header())
		for key, expectedValue := range item.expectedHeaders {
			values, found := headers[key]
			if !found {
				t.Fatalf(
					"Item %d, “%s”: Header '%s' not found in response",
					i,
					item.description,
					key,
				)
			}

			if values[0] != expectedValue {
				t.Errorf(
					"Item %d, “%s”: mismatch results. Expecting '%s'; found '%s'",
					i,
					item.description,
					expectedValue,
					values[0],
				)
			}
		}
	}
}

func TestHTTPHeadersCheckDate(t *testing.T) {
	data := []struct {
		description     string
		date            string
		expectedMessage bool
	}{
		{
			description: "it should ignore an empty HTTP Date header",
			date:        "",
		},
		{
			description: "it should validate a HTTP Date header",
			date:        time.Now().Format(time.RFC1123),
		},
		{
			description:     "it should report error for an invalid HTTP Date header (strange format)",
			date:            "XXXX",
			expectedMessage: true,
		},
		{
			description:     "it should report error for an HTTP Date header outside time frame",
			date:            time.Now().Add((timeFrameDuration + 1) * -1).Format(time.RFC1123),
			expectedMessage: true,
		},
	}

	for i, item := range data {
		r, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Date", item.date)

		handler := NewHTTPHeaders(&MockHTTPHandler{})
		message := handler.checkDate(r)

		if message == nil && item.expectedMessage {
			t.Errorf("Item %d, “%s”: expected message", i, item.description)

		} else if message != nil && !item.expectedMessage {
			t.Errorf("Item %d, “%s”: unexpected message '%#v'", i, item.description, message)
		}
	}
}

func TestHTTPHeadersCheckContentType(t *testing.T) {
	data := []struct {
		description      string
		contentType      string
		expectedMessages bool
	}{
		{
			description: "it should ignore an empty HTTP Content-Type header",
			contentType: "",
		},
		{
			description: "it should validate a HTTP Content-Type header",
			contentType: brContentType,
		},
		{
			description: "it should validate a HTTP Content-Type header with charset",
			contentType: fmt.Sprintf("%s;  charset=%s", brContentType, brCharset),
		},
		{
			description:      "it should report problem for an invalid content-type",
			contentType:      "application/json",
			expectedMessages: true,
		},
		{
			description:      "it should report problem for an invalid charset",
			contentType:      fmt.Sprintf("%s;charset=XXX", brContentType),
			expectedMessages: true,
		},
		{
			description: "it should ignore an invalid option format",
			contentType: fmt.Sprintf("%s;charset=%s;xxxx", brContentType, brCharset),
		},
	}

	for i, item := range data {
		r, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", item.contentType)

		handler := NewHTTPHeaders(&MockHTTPHandler{})
		messages := handler.checkContentType(r)

		if messages == nil && item.expectedMessages {
			t.Errorf("Item %d, “%s”: expected message", i, item.description)

		} else if messages != nil && !item.expectedMessages {
			t.Errorf("Item %d, “%s”: unexpected message '%#v'", i, item.description, messages)
		}
	}
}

func TestHTTPHeadersCheckContentMD5(t *testing.T) {
	data := []struct {
		description     string
		body            string
		contentMD5      string
		expectedMessage bool
	}{
		{
			description: "it should ignore an empty HTTP Content-MD5 header",
			body:        "This is a test!",
			contentMD5:  "",
		},
		{
			description: "it should validate a HTTP Content-MD5 header",
			body:        "This is a test!",
			contentMD5:  "cC7coLIYHBXUV+rKw53jmw==",
		},
		{
			description:     "it should report an error for an invalid HTTP Content-MD5 header",
			body:            "This is a test!",
			contentMD5:      "cC7coLIYHBXUV+rKw53jXX==",
			expectedMessage: true,
		},
		{
			description:     "it should report an error for a HTTP Content-MD5 header withou body",
			body:            "",
			contentMD5:      "cC7coLIYHBXUV+rKw53jmw==",
			expectedMessage: true,
		},
	}

	for i, item := range data {
		var r *http.Request
		var err error

		if len(item.body) > 0 {
			r, err = http.NewRequest("GET", "/test", strings.NewReader(item.body))
		} else {
			r, err = http.NewRequest("GET", "/test", nil)
		}

		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-MD5", item.contentMD5)

		handler := NewHTTPHeaders(&MockHTTPHandler{})
		message := handler.checkContentMD5(r)

		if message == nil && item.expectedMessage {
			t.Errorf("Item %d, “%s”: expected message", i, item.description)

		} else if message != nil && !item.expectedMessage {
			t.Errorf("Item %d, “%s”: unexpected message '%#v'", i, item.description, message)
		}
	}
}

func TestHTTPHeadersCheckAccept(t *testing.T) {
	data := []struct {
		description      string
		accept           string
		expectedMessages bool
	}{
		{
			description: "it should ignore an empty HTTP Accept header",
			accept:      "",
		},
		{
			description: "it should allow a valid HTTP Accept header",
			accept:      "APPLICATION/vnd.br+json;q=1.0,application/json,text/html",
		},
		{
			description: "it should allow a wildcard in HTTP Accept header",
			accept:      "application/json,text/html,*",
		},
		{
			description: "it should allow a wildcard in HTTP Accept header",
			accept:      "application/json,text/html,*/*",
		},
		{
			description:      "it should report a problem when there's no valid entries in HTTP Accept header",
			accept:           "application/json,text/html",
			expectedMessages: true,
		},
	}

	for i, item := range data {
		r, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Accept", item.accept)

		handler := NewHTTPHeaders(&MockHTTPHandler{})
		messages := handler.checkAccept(r)

		if messages == nil && item.expectedMessages {
			t.Errorf("Item %d, “%s”: expected message", i, item.description)

		} else if messages != nil && !item.expectedMessages {
			t.Errorf("Item %d, “%s”: unexpected message '%#v'", i, item.description, messages)
		}
	}
}

func TestHTTPHeadersCheckAcceptLanguage(t *testing.T) {
	data := []struct {
		description      string
		acceptLanguage   string
		expectedMessages bool
	}{
		{
			description:    "it should ignore an empty HTTP Accept Language header",
			acceptLanguage: "",
		},
		{
			description:    "it should allow a valid HTTP Accept Language header",
			acceptLanguage: "pt-BR;q=1.0,en-US;q=0.8,es-ES",
		},
		{
			description:    "it should allow a wildcard in HTTP Accept Language header",
			acceptLanguage: "en-US;q=0.8,es-ES,*",
		},
		{
			description:      "it should report a problem when there's no valid entries in HTTP Accept Language header",
			acceptLanguage:   "en-US;q=0.8,es-ES",
			expectedMessages: true,
		},
	}

	protocol.Translations["pt"] = protocol.Translation{}
	protocol.Translations["zh-ZH"] = protocol.Translation{}
	defer func() {
		protocol.Translations = make(map[string]protocol.Translation)
	}()

	for i, item := range data {
		r, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Accept-Language", item.acceptLanguage)

		handler := NewHTTPHeaders(&MockHTTPHandler{})
		messages := handler.checkAcceptLanguage(r)

		if messages == nil && item.expectedMessages {
			t.Errorf("Item %d, “%s”: expected message", i, item.description)

		} else if messages != nil && !item.expectedMessages {
			t.Errorf("Item %d, “%s”: unexpected message '%#v'", i, item.description, messages)
		}
	}
}

func TestHTTPHeadersCheckAcceptCharset(t *testing.T) {
	data := []struct {
		description      string
		acceptCharset    string
		expectedMessages bool
	}{
		{
			description:   "it should ignore an empty HTTP Accept Charset header",
			acceptCharset: "",
		},
		{
			description:   "it should allow a valid HTTP Accept Charset header",
			acceptCharset: "UTF-8;q=1.0,iso-8859-1;q=0.8",
		},
		{
			description:   "it should allow a wildcard in HTTP Accept Charset header",
			acceptCharset: "iso-8859-1;q=0.8,*",
		},
		{
			description:      "it should report a problem when there's no valid entries in HTTP Accept Charset header",
			acceptCharset:    "iso-8859-1;q=0.8",
			expectedMessages: true,
		},
	}

	for i, item := range data {
		r, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Accept-Charset", item.acceptCharset)

		handler := NewHTTPHeaders(&MockHTTPHandler{})
		messages := handler.checkAcceptCharset(r)

		if messages == nil && item.expectedMessages {
			t.Errorf("Item %d, “%s”: expected message", i, item.description)

		} else if messages != nil && !item.expectedMessages {
			t.Errorf("Item %d, “%s”: unexpected message '%#v'", i, item.description, messages)
		}
	}
}

// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockHTTPCacheHandler struct {
	LastModifiedAt time.Time
	ETag           string
	MessageId      string
	ReturnError    error
}

func (h *MockHTTPCacheHandler) GetLastModifiedAt() time.Time {
	return h.LastModifiedAt
}

func (h *MockHTTPCacheHandler) GetETag() string {
	return h.ETag
}

func (h *MockHTTPCacheHandler) MessageResponse(messageId string, roid string) error {
	h.MessageId = messageId
	return h.ReturnError
}

func (h *MockHTTPCacheHandler) ClearResponse() {
	h.MessageId = ""
}

func TestCheckHTTPCache(t *testing.T) {
	var httpCacheHandler MockHTTPCacheHandler

	data := []struct {
		Method            string
		Headers           map[string]string
		LastModifiedAt    string
		ETag              string
		ReturnError       error
		ExpectedCode      int
		ExpectedMessageId string
	}{
		{
			Method: "GET",
			Headers: map[string]string{
				"If-Modified-Since": "Sat, 29 Oct 1994 19:43:31 GMT",
			},
			LastModifiedAt: "Sat, 28 Oct 1994 19:43:31 GMT",
			ExpectedCode:   http.StatusNotModified,
		},
		{
			Method: "GET",
			Headers: map[string]string{
				"If-Modified-Since": "XXXXX",
			},
			LastModifiedAt:    "Sat, 28 Oct 1994 19:43:31 GMT",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "invalid-header-date",
		},
		{
			Method: "GET",
			Headers: map[string]string{
				"If-Modified-Since": "XXXXX",
			},
			LastModifiedAt:    "Sat, 28 Oct 1994 19:43:31 GMT",
			ReturnError:       errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "invalid-header-date",
		},
		{
			Method: "GET",
			Headers: map[string]string{
				"If-Unmodified-Since": "Sat, 29 Oct 1994 19:43:31 GMT",
			},
			LastModifiedAt: "Sat, 30 Oct 1994 19:43:31 GMT",
			ExpectedCode:   http.StatusPreconditionFailed,
		},
		{
			Method: "GET",
			Headers: map[string]string{
				"If-Unmodified-Since": "XXXXX",
			},
			LastModifiedAt:    "Sat, 30 Oct 1994 19:43:31 GMT",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "invalid-header-date",
		},
		{
			Method: "GET",
			Headers: map[string]string{
				"If-Unmodified-Since": "XXXXX",
			},
			LastModifiedAt:    "Sat, 30 Oct 1994 19:43:31 GMT",
			ReturnError:       errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "invalid-header-date",
		},
		{
			Method: "GET",
			Headers: map[string]string{
				"If-Match": "1",
			},
			ETag:              "2",
			ExpectedCode:      http.StatusPreconditionFailed,
			ExpectedMessageId: "if-match-failed",
		},
		{
			Method: "GET",
			Headers: map[string]string{
				"If-Match": "1",
			},
			ETag:              "2",
			ReturnError:       errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "if-match-failed",
		},
		{
			Method: "GET",
			Headers: map[string]string{
				"If-None-Match": "1",
			},
			ETag:         "1",
			ExpectedCode: http.StatusNotModified,
		},
		{
			Method: "DELETE",
			Headers: map[string]string{
				"If-None-Match": "1",
			},
			ETag:              "1",
			ExpectedCode:      http.StatusPreconditionFailed,
			ExpectedMessageId: "if-none-match-failed",
		},
		{
			Method: "DELETE",
			Headers: map[string]string{
				"If-None-Match": "1",
			},
			ETag:              "1",
			ReturnError:       errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "if-none-match-failed",
		},
		{
			Method: "GET",
			Headers: map[string]string{
				"If-Modified-Since":   "Sat, 28 Oct 1994 19:43:31 GMT",
				"If-Unmodified-Since": "Sat, 30 Oct 1994 19:43:31 GMT",
				"If-Match":            "2",
				"If-None-Match":       "1",
			},
			LastModifiedAt: "Sat, 29 Oct 1994 19:43:31 GMT",
			ETag:           "2",
			ExpectedCode:   http.StatusOK,
		},
	}

	for _, item := range data {
		r, err := http.NewRequest(item.Method, "/test", nil)
		if err != nil {
			t.Fatal(err)
		}

		for key, value := range item.Headers {
			r.Header.Set(key, value)
		}

		if len(item.LastModifiedAt) > 0 {
			lastModifiedAt, err := time.Parse(time.RFC1123, item.LastModifiedAt)
			if err != nil {
				t.Fatal(err)
			}
			httpCacheHandler.LastModifiedAt = lastModifiedAt
		}

		httpCacheHandler.ETag = item.ETag
		httpCacheHandler.ReturnError = item.ReturnError
		httpCacheHandler.ClearResponse()

		w := httptest.NewRecorder()
		CheckHTTPCache(w, r, &httpCacheHandler)

		if w.Code != item.ExpectedCode {
			t.Errorf("Wrong status code for headers '%v'. "+
				"Expected %d and got %d", item.Headers, item.ExpectedCode, w.Code)
		}

		if httpCacheHandler.MessageId != item.ExpectedMessageId {
			t.Errorf("Wrong message id for headers '%v'. "+
				"Expected %s and got %s", item.Headers, item.ExpectedMessageId, httpCacheHandler.MessageId)
		}
	}
}

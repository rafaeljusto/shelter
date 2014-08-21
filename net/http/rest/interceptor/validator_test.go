// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"errors"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"github.com/rafaeljusto/shelter/net/http/rest/messages"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockValidatorHandler struct {
	MessageId   string
	ReturnError error
	Language    *messages.LanguagePack
}

func (h *MockValidatorHandler) SetLanguage(language *messages.LanguagePack) {
	h.Language = language
}

func (h *MockValidatorHandler) MessageResponse(messageId string, roid string) error {
	h.MessageId = messageId
	return h.ReturnError
}

func TestValidatorBefore(t *testing.T) {
	var validatorHandler MockValidatorHandler
	validator := NewValidator(&validatorHandler)

	config.ShelterConfig.RESTServer.Secrets = map[string]string{
		"1": "ohV43/9bKlVNaXeNTqEuHQp57LCPCQ==",
		"2": "xJVO/OLkL35GnFTSDtQaVqZfOE8rtw==",
	}

	data := []struct {
		Headers           map[string]string
		Body              string
		LanguagePacks     messages.LanguagePacks
		Authorize         bool
		SecretId          string
		RetunError        error
		ExpectedCode      int
		ExpectedMessageId string
	}{
		{
			Headers: map[string]string{
				"Accept-Language": "XXXX",
				"Date":            time.Now().Format(time.RFC1123),
			},
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			ExpectedCode:      http.StatusNotAcceptable,
			ExpectedMessageId: "accept-language-error",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "XXXX",
				"Date":            time.Now().Format(time.RFC1123),
			},
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			RetunError:        errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "accept-language-error",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "XXXX",
				"Date":            time.Now().Format(time.RFC1123),
			},
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			ExpectedCode:      http.StatusNotAcceptable,
			ExpectedMessageId: "accept-error",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "XXXX",
				"Date":            time.Now().Format(time.RFC1123),
			},
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			RetunError:        errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "accept-error",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "XXXX",
				"Date":            time.Now().Format(time.RFC1123),
			},
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			ExpectedCode:      http.StatusNotAcceptable,
			ExpectedMessageId: "accept-charset-error",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "XXXX",
				"Date":            time.Now().Format(time.RFC1123),
			},
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			RetunError:        errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "accept-charset-error",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "XXXX",
				"Date":            time.Now().Format(time.RFC1123),
			},
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "invalid-content-type",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "XXXX",
				"Date":            time.Now().Format(time.RFC1123),
			},
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			RetunError:        errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "invalid-content-type",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "XXXX",
				"Date":            time.Now().Format(time.RFC1123),
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "invalid-content-md5",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "XXXX",
				"Date":            time.Now().Format(time.RFC1123),
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			RetunError:        errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "invalid-content-md5",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
				"Date":            "XXXX",
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "invalid-header-date",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
				"Date":            "XXXX",
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			RetunError:        errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "invalid-header-date",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
				"Date":            time.Now().Add(-1 * time.Hour).Format(time.RFC1123),
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "invalid-date-time-frame",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
				"Date":            time.Now().Add(-1 * time.Hour).Format(time.RFC1123),
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			RetunError:        errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "invalid-date-time-frame",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
				"Date":            time.Now().Format(time.RFC1123),
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:    true,
			SecretId:     "2",
			ExpectedCode: http.StatusUnauthorized,
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
				"Date":            time.Now().Format(time.RFC1123),
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "content-type-missing",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Date":            time.Now().Format(time.RFC1123),
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "content-md5-missing",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "1",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "date-missing",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
				"Date":            time.Now().Format(time.RFC1123),
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         false,
			SecretId:          "1",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "authorization-missing",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
				"Date":            time.Now().Format(time.RFC1123),
				"Authorization":   "XXXX",
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         false,
			SecretId:          "1",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "invalid-authorization",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
				"Date":            time.Now().Format(time.RFC1123),
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         true,
			SecretId:          "3",
			ExpectedCode:      http.StatusBadRequest,
			ExpectedMessageId: "secret-not-found",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
				"Date":            time.Now().Format(time.RFC1123),
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:         false,
			SecretId:          "1",
			RetunError:        errors.New("Low level error!"),
			ExpectedCode:      http.StatusInternalServerError,
			ExpectedMessageId: "authorization-missing",
		},
		{
			Headers: map[string]string{
				"Accept-Language": "en-US",
				"Accept":          "application/vnd.shelter+json",
				"Accept-Charset":  "utf-8",
				"Content-Type":    "application/vnd.shelter+json",
				"Content-MD5":     "cC7coLIYHBXUV+rKw53jmw==",
				"Date":            time.Now().Format(time.RFC1123),
			},
			Body: "This is a test!",
			LanguagePacks: messages.LanguagePacks{
				Default: "en-US",
				Packs: []messages.LanguagePack{
					{
						GenericName:  "en",
						SpecificName: "en-US",
					},
				},
			},
			Authorize:    true,
			SecretId:     "1",
			ExpectedCode: http.StatusOK,
		},
	}

	for _, item := range data {
		var r *http.Request
		var err error

		if len(item.Body) > 0 {
			r, err = http.NewRequest("GET", "/test", strings.NewReader(item.Body))
		} else {
			r, err = http.NewRequest("GET", "/test", nil)
		}

		if err != nil {
			t.Fatal(err)
		}

		for key, value := range item.Headers {
			r.Header.Set(key, value)
		}

		if item.Authorize {
			// Don't check BuildStringToSign errors, because it will be checked by each test
			stringToSign, _ := check.BuildStringToSign(r, item.SecretId)
			signature := check.GenerateSignature(stringToSign, "abc123")
			r.Header.Set("Authorization", fmt.Sprintf("shelter %s:%s", item.SecretId, signature))
		}

		validatorHandler.ReturnError = item.RetunError
		validatorHandler.MessageId = ""

		messages.ShelterRESTLanguagePacks = item.LanguagePacks
		messages.ShelterRESTLanguagePack =
			messages.ShelterRESTLanguagePacks.Select(messages.ShelterRESTLanguagePacks.Default)

		w := httptest.NewRecorder()
		validator.Before(w, r)

		if w.Code != item.ExpectedCode {
			t.Errorf("Wrong status code for headers '%v'. "+
				"Expected %d and got %d", item.Headers, item.ExpectedCode, w.Code)
		}

		if validatorHandler.MessageId != item.ExpectedMessageId {
			t.Errorf("Wrong message id for headers '%v'. "+
				"Expected %s and got %s", item.Headers, item.ExpectedMessageId, validatorHandler.MessageId)
		}
	}
}

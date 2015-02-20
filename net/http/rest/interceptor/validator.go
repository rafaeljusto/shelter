// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"bytes"
	"errors"
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/rafaeljusto/handy/interceptor"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"github.com/rafaeljusto/shelter/net/http/rest/messages"
	"github.com/rafaeljusto/shelter/secret"
	"io/ioutil"
	"net/http"
)

// List of possible errors that can occur when calling functions from this file. Other
// erros can also occurs from low level layers
var (
	ErrSecretNotFound = errors.New("Secret related to Authorization's secret id not found")
)

type ValidatorHandler interface {
	SetLanguage(*messages.LanguagePack)
	MessageResponse(string, string) error
}

type Validator struct {
	interceptor.NoAfterInterceptor
	validatorHandler ValidatorHandler
}

func NewValidator(h ValidatorHandler) *Validator {
	return &Validator{validatorHandler: h}
}

// Verify HTTP headers and fill context with user preferences
func (i *Validator) Before(w http.ResponseWriter, r *http.Request) {
	// We first check the language header, because if it's acceptable the next messages are
	// going to be returned in the language choosen by the user
	language, ok := check.HTTPAcceptLanguage(r)

	if ok {
		i.validatorHandler.SetLanguage(language)

	} else {
		if err := i.validatorHandler.MessageResponse("accept-language-error", r.RequestURI); err == nil {
			w.WriteHeader(http.StatusNotAcceptable)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	if !check.HTTPAccept(r) {
		if err := i.validatorHandler.MessageResponse("accept-error", r.RequestURI); err == nil {
			w.WriteHeader(http.StatusNotAcceptable)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	if !check.HTTPAcceptCharset(r) {
		if err := i.validatorHandler.MessageResponse("accept-charset-error", r.RequestURI); err == nil {
			w.WriteHeader(http.StatusNotAcceptable)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	if !check.HTTPContentType(r) {
		if err := i.validatorHandler.MessageResponse("invalid-content-type", r.RequestURI); err == nil {
			w.WriteHeader(http.StatusBadRequest)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading the request. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}()

		if !check.HTTPContentMD5(r, body) {
			if err := i.validatorHandler.MessageResponse("invalid-content-md5", r.RequestURI); err == nil {
				w.WriteHeader(http.StatusBadRequest)

			} else {
				log.Println("Error while writing response. Details:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}
	}

	timeFrameOK, err := check.HTTPDate(r)

	if err != nil {
		if err := i.validatorHandler.MessageResponse("invalid-header-date", r.RequestURI); err == nil {
			w.WriteHeader(http.StatusBadRequest)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return

	} else if !timeFrameOK {
		if err := i.validatorHandler.MessageResponse("invalid-date-time-frame", r.RequestURI); err == nil {
			w.WriteHeader(http.StatusBadRequest)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	authorized, err := check.HTTPAuthorization(r, func(secretId string) (string, error) {
		s, ok := config.ShelterConfig.RESTServer.Secrets[secretId]

		if !ok {
			return "", ErrSecretNotFound
		}

		s, err = secret.Decrypt(s)
		if err != nil {
			return "", err
		}

		// In the near future the secret will be encrypted in the configuration file and the
		// decrypt process can generate problems
		return s, nil
	})

	if err == nil && !authorized {
		w.WriteHeader(http.StatusUnauthorized)
		return

	} else if authorized {
		return
	}

	messageId := ""

	switch err {
	case check.ErrHTTPContentTypeNotFound:
		messageId = "content-type-missing"
	case check.ErrHTTPContentMD5NotFound:
		messageId = "content-md5-missing"
	case check.ErrHTTPDateNotFound:
		messageId = "date-missing"
	case check.ErrHTTPAuthorizationNotFound:
		messageId = "authorization-missing"
	case check.ErrInvalidHTTPAuthorization:
		messageId = "invalid-authorization"
	case ErrSecretNotFound:
		messageId = "secret-not-found"
	}

	if len(messageId) == 0 {
		log.Println("Error checking authorization. Details:", err)
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		if err := i.validatorHandler.MessageResponse(messageId, r.RequestURI); err == nil {
			w.WriteHeader(http.StatusBadRequest)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

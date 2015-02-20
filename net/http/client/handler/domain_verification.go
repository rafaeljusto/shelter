// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package handler store the web client handlers of specific URI
package handler

import (
	"fmt"
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/rafaeljusto/handy"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/net/http/client/interceptor"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func init() {
	HandleFunc("/domain/{fqdn}/verification", func() handy.Handler {
		return new(DomainVerificationHandler)
	})
}

type DomainVerificationHandler struct {
	handy.DefaultHandler
	FQDN string `param:"fqdn"`
}

func (h *DomainVerificationHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.handleDomainVerification(w, r)
}

func (h *DomainVerificationHandler) Head(w http.ResponseWriter, r *http.Request) {
	h.handleDomainVerification(w, r)
}

func (h *DomainVerificationHandler) Put(w http.ResponseWriter, r *http.Request) {
	h.handleDomainVerification(w, r)
}

func (h *DomainVerificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.handleDomainVerification(w, r)
}

func (h *DomainVerificationHandler) handleDomainVerification(w http.ResponseWriter, r *http.Request) {
	restAddress, err := retrieveRESTAddress()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error while retrieving the REST address. Details:", err)
		return
	}

	var content []byte
	if r.ContentLength > 0 {
		content, err = ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error while reading request body in web client. Details:", err)
			return
		}
	}

	var request *http.Request
	if content == nil {
		request, err = http.NewRequest(
			r.Method,
			fmt.Sprintf("%s%s", restAddress, r.RequestURI),
			nil,
		)

	} else {
		request, err = http.NewRequest(
			r.Method,
			fmt.Sprintf("%s%s", restAddress, r.RequestURI),
			strings.NewReader(string(content)),
		)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating a request in web client. Details:", err)
		return
	}

	request.Header.Set("Accept-Language", r.Header.Get("Accept-Language"))

	response, err := signAndSend(request, content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error signing and sending a request in web client. Details:", err)
		return
	}

	if response.StatusCode != http.StatusOK &&
		response.StatusCode != http.StatusBadRequest {

		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Sprintf("Unexepected status code %d from /domain result "+
			"in web client", response.StatusCode))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	if _, err := io.Copy(w, response.Body); err != nil {
		// Here we already set the response code, so the client will receive a OK result
		// without body
		log.Println("Error copying REST response to web client response. Details:", err)
		return
	}
}

func (h *DomainVerificationHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(new(interceptor.Permission))
}

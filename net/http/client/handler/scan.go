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
	"net/http"
)

func init() {
	HandleFunc("/scan/{started-at}", func() handy.Handler {
		return new(ScanHandler)
	})
}

// ScanHandler is responsable for keeping the state of a /scan/{started-at} resource
type ScanHandler struct {
	handy.DefaultHandler        // Inject the HTTP methods that this resource does not implement
	StartedAt            string `param:"started-at"` // Scan start date in the URI
}

func (h *ScanHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.handleScan(w, r)
}

func (h *ScanHandler) Head(w http.ResponseWriter, r *http.Request) {
	h.handleScan(w, r)
}

func (h *ScanHandler) Put(w http.ResponseWriter, r *http.Request) {
	h.handleScan(w, r)
}

func (h *ScanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.handleScan(w, r)
}

func (h *ScanHandler) handleScan(w http.ResponseWriter, r *http.Request) {
	restAddress, err := retrieveRESTAddress()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error while retrieving the REST address. Details:", err)
		return
	}

	request, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s", restAddress, r.RequestURI),
		nil,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating a request in web client. Details:", err)
		return
	}

	request.Header.Set("Accept-Language", r.Header.Get("Accept-Language"))
	request.Header.Set("If-None-Match", r.Header.Get("If-None-Match"))

	response, err := signAndSend(request, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error signing and sending a request in web client. Details:", err)
		return
	}

	if response.StatusCode != http.StatusOK &&
		response.StatusCode != http.StatusNotModified &&
		response.StatusCode != http.StatusBadRequest {

		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Sprintf("Unexepected status code %d from /scan/{started-at} "+
			"result in web client", response.StatusCode))
		return
	}

	if response.ContentLength > 0 {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Etag", response.Header.Get("Etag"))
	}

	w.WriteHeader(response.StatusCode)

	if response.ContentLength > 0 {
		if _, err := io.Copy(w, response.Body); err != nil {
			// Here we already set the response code, so the client will receive a OK result
			// without body
			log.Println("Error copying REST response to web client response. Details:", err)
			return
		}
	}
}

func (h *ScanHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(new(interceptor.Permission))
}

// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
	"github.com/trajber/handy/interceptor"
	"net/http"
)

type FQDNHandler interface {
	SetFQDN(fqdn string)
	GetFQDN() string
	MessageResponse(string, string) error
}

type FQDN struct {
	interceptor.NoAfterInterceptor
	fqdnHandler FQDNHandler
}

func NewFQDN(h FQDNHandler) *FQDN {
	return &FQDN{fqdnHandler: h}
}

func (i *FQDN) Before(w http.ResponseWriter, r *http.Request) {
	fqdn, err := model.NormalizeDomainName(i.fqdnHandler.GetFQDN())
	if err != nil {
		if err := i.fqdnHandler.MessageResponse("invalid-uri", r.URL.RequestURI()); err == nil {
			w.WriteHeader(http.StatusBadRequest)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	i.fqdnHandler.SetFQDN(fqdn)
}

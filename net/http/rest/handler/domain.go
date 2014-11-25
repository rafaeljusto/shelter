// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package handler store the REST handlers of specific URI
package handler

import (
	"github.com/rafaeljusto/handy"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/errors"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/messages"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/interceptor"
	"github.com/rafaeljusto/shelter/protocol"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	HandleFunc("/domain/{fqdn}", func() handy.Handler {
		return new(DomainHandler)
	})
}

// DomainHandler is responsable for keeping the state of a /domain/{fqdn} resource
type DomainHandler struct {
	handy.DefaultHandler // Inject the HTTP methods that this resource does not implement
	interceptor.JSONCompliant
	interceptor.DatabaseCompliant

	domain   model.Domain
	language *messages.LanguagePack    // User preferred language based on HTTP header
	FQDN     string                    `param:"fqdn"`   // FQDN defined in the URI
	Request  protocol.DomainRequest    `request:"put"`  // Domain request sent by the user
	Response *protocol.DomainResponse  `response:"get"` // Domain response sent back to the user
	Message  *protocol.MessageResponse `error`          // Message on error sent to the user
}

func (h *DomainHandler) LastModifiedAt() time.Time {
	return h.domain.LastModifiedAt
}

func (h *DomainHandler) ETag() string {
	return strconv.Itoa(h.domain.Revision)
}

func (h *DomainHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.retrieveDomain(w, r)
}

func (h *DomainHandler) Head(w http.ResponseWriter, r *http.Request) {
	h.retrieveDomain(w, r)
}

// The HEAD method is identical to GET except that the server MUST NOT return a message-
// body in the response. But now the responsability for don't adding the body is from the
// mux while writing the response
func (h *DomainHandler) retrieveDomain(w http.ResponseWriter, r *http.Request) {
	domainDAO := dao.DomainDAO{
		Database: h.DB(),
	}

	var err error
	h.domain, err = domainDAO.FindByFQDN(h.FQDN)
	if err == errors.NotFound {
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		// TODO: Log
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("ETag", h.ETag())
	w.Header().Add("Last-Modified", h.LastModifiedAt().Format(time.RFC1123))
	w.WriteHeader(http.StatusOK)

	domainResponse := h.domain.Protocol()
	h.Response = &domainResponse
}

func (h *DomainHandler) Put(w http.ResponseWriter, r *http.Request) {
	// We need to set the FQDN in the domain request object because it is sent only in the
	// URI and not in the domain request body to avoid information redudancy
	h.Request.FQDN = &h.domain.FQDN

	domainDAO := dao.DomainDAO{
		Database: h.DB(),
	}

	var err error
	h.domain, err = domainDAO.FindByFQDN(h.domain.FQDN)
	if err == errors.NotFound {
		w.WriteHeader(http.StatusNotFound)
		return

	} else if err != nil {
		// TODO: Log
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !h.domain.Apply(h.Request) {
		log.Println("Error while merging domain objects for create or update operation.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := domainDAO.Save(&h.domain); err != nil {
		if strings.Index(err.Error(), "duplicate key error index") != -1 {
			h.SetMessage(protocol.NewMessageResponse(protocol.ErrorCodeConflict, nil))
			w.WriteHeader(http.StatusConflict)

		} else {
			log.Println("Error while saving domain object for create or "+
				"update operation. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.Header().Add("ETag", h.ETag())
	w.Header().Add("Last-Modified", h.LastModifiedAt().Format(time.RFC1123))

	if domain.Revision == 1 {
		w.Header().Add("Location", "/domain/"+domain.FQDN)
		w.WriteHeader(http.StatusCreated)

	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *DomainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	domainDAO := dao.DomainDAO{
		Database: h.DB(),
	}

	if err := domainDAO.Remove(&h.domain); err != nil {
		log.Println("Error while removing domain object. Details:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DomainHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(new(interceptor.Permission)).
		Chain(interceptor.NewFQDN(h)).
		Chain(interceptor.NewValidator(h)).
		Chain(interceptor.NewDatabase(h)).
		Chain(interceptor.NewDomain(h)).
		Chain(interceptor.NewHTTPCacheBefore(h)).
		Chain(interceptor.NewJSONCodec(h))
}

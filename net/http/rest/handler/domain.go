// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package handler store the REST handlers of specific URI
package handler

import (
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/rafaeljusto/handy"
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/interceptor"
	"github.com/rafaeljusto/shelter/net/http/rest/messages"
	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
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
	handy.DefaultHandler                           // Inject the HTTP methods that this resource does not implement
	database             *mgo.Database             // Database connection of the MongoDB session
	databaseSession      *mgo.Session              // MongoDB session
	domain               model.Domain              // Domain object related to the resource
	language             *messages.LanguagePack    // User preferred language based on HTTP header
	FQDN                 string                    `param:"fqdn"`   // FQDN defined in the URI
	Request              protocol.DomainRequest    `request:"put"`  // Domain request sent by the user
	Response             *protocol.DomainResponse  `response:"get"` // Domain response sent back to the user
	Message              *protocol.MessageResponse `error`          // Message on error sent to the user
}

func (h *DomainHandler) SetDatabaseSession(session *mgo.Session) {
	h.databaseSession = session
}

func (h *DomainHandler) GetDatabaseSession() *mgo.Session {
	return h.databaseSession
}

func (h *DomainHandler) SetDatabase(database *mgo.Database) {
	h.database = database
}

func (h *DomainHandler) GetDatabase() *mgo.Database {
	return h.database
}

func (h *DomainHandler) SetFQDN(fqdn string) {
	h.FQDN = fqdn
}

func (h *DomainHandler) GetFQDN() string {
	return h.FQDN
}

func (h *DomainHandler) SetDomain(domain model.Domain) {
	h.domain = domain
}

func (h *DomainHandler) GetLastModifiedAt() time.Time {
	return h.domain.LastModifiedAt
}

func (h *DomainHandler) GetETag() string {
	return strconv.Itoa(h.domain.Revision)
}

func (h *DomainHandler) SetLanguage(language *messages.LanguagePack) {
	h.language = language
}

func (h *DomainHandler) GetLanguage() *messages.LanguagePack {
	return h.language
}

func (h *DomainHandler) MessageResponse(messageId string, roid string) error {
	var err error
	h.Message, err = protocol.NewMessageResponse(messageId, roid, h.language)
	return err
}

func (h *DomainHandler) ClearResponse() {
	h.Response = nil
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
	w.Header().Add("ETag", h.GetETag())
	w.Header().Add("Last-Modified", h.GetLastModifiedAt().Format(time.RFC1123))
	w.WriteHeader(http.StatusOK)

	domainResponse := protocol.ToDomainResponse(h.domain, true)
	h.Response = &domainResponse
}

func (h *DomainHandler) Put(w http.ResponseWriter, r *http.Request) {
	// We need to set the FQDN in the domain request object because it is sent only in the
	// URI and not in the domain request body to avoid information redudancy
	h.Request.FQDN = h.GetFQDN()

	var err error
	if h.domain, err = protocol.Merge(h.domain, h.Request); err != nil {
		messageId := ""

		switch err {
		case model.ErrInvalidFQDN:
			messageId = "invalid-fqdn"
		case protocol.ErrInvalidDNSKEY:
			messageId = "invalid-dnskey"
		case protocol.ErrInvalidDSAlgorithm:
			messageId = "invalid-ds-algorithm"
		case protocol.ErrInvalidDSDigestType:
			messageId = "invalid-ds-digest-type"
		case protocol.ErrInvalidIP:
			messageId = "invalid-ip"
		case protocol.ErrInvalidLanguage:
			messageId = "invalid-language"
		}

		if len(messageId) == 0 {
			log.Println("Error while merging domain objects for create or "+
				"update operation. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			if err := h.MessageResponse(messageId, r.URL.RequestURI()); err == nil {
				w.WriteHeader(http.StatusBadRequest)

			} else {
				log.Println("Error while writing response. Details:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		return
	}

	domainDAO := dao.DomainDAO{
		Database: h.GetDatabase(),
	}

	if err := domainDAO.Save(&h.domain); err != nil {
		if strings.Index(err.Error(), "duplicate key error index") != -1 {
			if err := h.MessageResponse("conflict", r.URL.RequestURI()); err == nil {
				w.WriteHeader(http.StatusConflict)

			} else {
				log.Println("Error while writing response. Details:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}

		} else {
			log.Println("Error while saving domain object for create or "+
				"update operation. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.Header().Add("ETag", h.GetETag())
	w.Header().Add("Last-Modified", h.GetLastModifiedAt().Format(time.RFC1123))

	if h.domain.Revision == 1 {
		w.Header().Add("Location", "/domain/"+h.domain.FQDN)
		w.WriteHeader(http.StatusCreated)

	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *DomainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	domainDAO := dao.DomainDAO{
		Database: h.GetDatabase(),
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

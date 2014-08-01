// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package handler store the REST handlers of specific URI
package handler

import (
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/messages"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/interceptor"
	"github.com/rafaeljusto/shelter/net/scan"
	"github.com/rafaeljusto/shelter/protocol"
	"github.com/trajber/handy"
	"gopkg.in/mgo.v2"
	"net/http"
)

func init() {
	HandleFunc("/domain/{fqdn}/verification", func() handy.Handler {
		return new(DomainVerificationHandler)
	})
}

type DomainVerificationHandler struct {
	handy.DefaultHandler
	database        *mgo.Database
	databaseSession *mgo.Session
	language        *messages.LanguagePack
	FQDN            string                    `param:"fqdn"`
	Request         protocol.DomainRequest    `request:"put"`
	Response        *protocol.DomainResponse  `response:"put,get"`
	Message         *protocol.MessageResponse `error`
}

func (h *DomainVerificationHandler) SetDatabaseSession(session *mgo.Session) {
	h.databaseSession = session
}

func (h *DomainVerificationHandler) GetDatabaseSession() *mgo.Session {
	return h.databaseSession
}

func (h *DomainVerificationHandler) SetDatabase(database *mgo.Database) {
	h.database = database
}

func (h *DomainVerificationHandler) GetDatabase() *mgo.Database {
	return h.database
}

func (h *DomainVerificationHandler) SetFQDN(fqdn string) {
	h.FQDN = fqdn
}

func (h *DomainVerificationHandler) GetFQDN() string {
	return h.FQDN
}

func (h *DomainVerificationHandler) SetLanguage(language *messages.LanguagePack) {
	h.language = language
}

func (h *DomainVerificationHandler) GetLanguage() *messages.LanguagePack {
	return h.language
}

func (h *DomainVerificationHandler) MessageResponse(messageId string, roid string) error {
	var err error
	h.Message, err = protocol.NewMessageResponse(messageId, roid, h.language)
	return err
}

func (h *DomainVerificationHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.queryDomain(w, r)
}

func (h *DomainVerificationHandler) Head(w http.ResponseWriter, r *http.Request) {
	h.queryDomain(w, r)
}

// Build the domain object doing a DNS query. To this function works the domain must be
// registered correctly and delegated in the DNS tree
func (h *DomainVerificationHandler) queryDomain(w http.ResponseWriter, r *http.Request) {
	domain, err := scan.QueryDomain(h.FQDN)
	if err != nil {
		log.Println("Error while resolving FQDN. Details:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	domainResponse := protocol.ToDomainResponse(domain, false)
	h.Response = &domainResponse
}

// Put is responsable for checking a domain object on-the-fly without persisting in
// database, useful for pre-registration validations in the registry
func (h *DomainVerificationHandler) Put(w http.ResponseWriter, r *http.Request) {
	// We need to set the FQDN in the domain request object because it is sent only in the
	// URI and not in the domain request body to avoid information redudancy
	h.Request.FQDN = h.GetFQDN()

	var domain model.Domain
	var err error

	if domain, err = protocol.Merge(domain, h.Request); err != nil {
		messageId := ""

		switch err {
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
			log.Println("Error while merging domain objects for domain verification "+
				"operation. Details:", err)
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

	scan.ScanDomain(&domain)

	// As we alredy did the scan, if the domain is registered in the system, we update it for this
	// results. This also gives a more intuitive design for when the user wants to force a check a
	// specific domain in the Shelter system
	domainDAO := dao.DomainDAO{
		Database: h.GetDatabase(),
	}

	if dbDomain, err := domainDAO.FindByFQDN(domain.FQDN); err == nil {
		update := true

		// Check if we have the same nameservers, and if so update the last status
		if len(dbDomain.Nameservers) == len(domain.Nameservers) {
			for i := range dbDomain.Nameservers {
				dbNameserver := dbDomain.Nameservers[i]
				nameserver := domain.Nameservers[i]

				if dbNameserver.Host == nameserver.Host &&
					dbNameserver.IPv4.Equal(nameserver.IPv4) &&
					dbNameserver.IPv6.Equal(nameserver.IPv6) {

					dbDomain.Nameservers[i].ChangeStatus(nameserver.LastStatus)

				} else {
					update = false
					break
				}
			}

		} else {
			update = false
		}

		// Check if we have the same DS set, and if so update the last status
		if len(dbDomain.DSSet) == len(domain.DSSet) {
			for i := range dbDomain.DSSet {
				dbDS := dbDomain.DSSet[i]
				ds := domain.DSSet[i]

				if dbDS.Keytag == ds.Keytag &&
					dbDS.Algorithm == ds.Algorithm &&
					dbDS.DigestType == ds.DigestType &&
					dbDS.Digest == ds.Digest {

					dbDomain.DSSet[i].ChangeStatus(ds.LastStatus)

				} else {
					update = false
					break
				}
			}

		} else {
			update = false
		}

		if update {
			// We don't care about errors resulted here, because the main idea of this service is to scan
			// a domaion, not persist the results
			domainDAO.Save(&dbDomain)
		}
	}

	w.WriteHeader(http.StatusOK)
	domainResponse := protocol.ToDomainResponse(domain, false)
	h.Response = &domainResponse
}

func (h *DomainVerificationHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(new(interceptor.Permission)).
		Chain(interceptor.NewFQDN(h)).
		Chain(interceptor.NewValidator(h)).
		Chain(interceptor.NewDatabase(h)).
		Chain(interceptor.NewJSONCodec(h))
}

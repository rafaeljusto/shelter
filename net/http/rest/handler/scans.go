// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package handler store the REST handlers of specific URI
package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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
	HandleFunc("/scans", func() handy.Handler {
		return new(ScansHandler)
	})
}

type ScansHandler struct {
	handy.DefaultHandler
	database        *mgo.Database
	databaseSession *mgo.Session
	language        *messages.LanguagePack
	Response        *protocol.ScansResponse   `response:"get"`
	Message         *protocol.MessageResponse `error`
	lastModifiedAt  time.Time
}

func (h *ScansHandler) SetDatabaseSession(session *mgo.Session) {
	h.databaseSession = session
}

func (h *ScansHandler) GetDatabaseSession() *mgo.Session {
	return h.databaseSession
}

func (h *ScansHandler) SetDatabase(database *mgo.Database) {
	h.database = database
}

func (h *ScansHandler) GetDatabase() *mgo.Database {
	return h.database
}

func (h *ScansHandler) GetLastModifiedAt() time.Time {
	return h.lastModifiedAt
}

// The ETag header will be the hash of the content on list services
func (h *ScansHandler) GetETag() string {
	body, err := json.Marshal(h.Response)
	if err != nil {
		return ""
	}

	hash := md5.New()
	if _, err := hash.Write(body); err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func (h *ScansHandler) SetLanguage(language *messages.LanguagePack) {
	h.language = language
}

func (h *ScansHandler) GetLanguage() *messages.LanguagePack {
	return h.language
}

func (h *ScansHandler) MessageResponse(messageId string, roid string) error {
	var err error
	h.Message, err = protocol.NewMessageResponse(messageId, roid, h.language)
	return err
}

func (h *ScansHandler) ClearResponse() {
	h.Response = nil
}

func (h *ScansHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.retrieveScans(w, r)
}

func (h *ScansHandler) Head(w http.ResponseWriter, r *http.Request) {
	h.retrieveScans(w, r)
}

// The HEAD method is identical to GET except that the server MUST NOT return a message-
// body in the response. But now the responsability for don't adding the body is from the
// mux while writing the response
func (h *ScansHandler) retrieveScans(w http.ResponseWriter, r *http.Request) {
	var pagination dao.ScanDAOPagination

	expand := false
	returnCurrent := false

	for key, values := range r.URL.Query() {
		key = strings.TrimSpace(key)
		key = strings.ToLower(key)

		// A key can have multiple values in a query string, we are going to always consider
		// the last one (overwrite strategy)
		for _, value := range values {
			value = strings.TrimSpace(value)
			value = strings.ToLower(value)

			switch key {
			case "orderby":
				// OrderBy parameter will store the fields that the user want to be the keys of the sort
				// algorithm in the result set and the direction that each sort field will have. The format
				// that will be used is:
				//
				// <field1>:<direction1>@<field2>:<direction2>@...@<fieldN>:<directionN>

				orderByParts := strings.Split(value, "@")

				for _, orderByPart := range orderByParts {
					orderByPart = strings.TrimSpace(orderByPart)
					orderByAndDirection := strings.Split(orderByPart, ":")

					var field, direction string

					if len(orderByAndDirection) == 1 {
						field, direction = orderByAndDirection[0], "asc"

					} else if len(orderByAndDirection) == 2 {
						field, direction = orderByAndDirection[0], orderByAndDirection[1]

					} else {
						if err := h.MessageResponse("invalid-query-order-by", ""); err == nil {
							w.WriteHeader(http.StatusBadRequest)

						} else {
							log.Println("Error while writing response. Details:", err)
							w.WriteHeader(http.StatusInternalServerError)
						}

						return
					}

					orderByField, err := dao.ScanDAOOrderByFieldFromString(field)
					if err != nil {
						if err := h.MessageResponse("invalid-query-order-by", ""); err == nil {
							w.WriteHeader(http.StatusBadRequest)

						} else {
							log.Println("Error while writing response. Details:", err)
							w.WriteHeader(http.StatusInternalServerError)
						}

						return
					}

					orderByDirection, err := dao.DAOOrderByDirectionFromString(direction)
					if err != nil {
						if err := h.MessageResponse("invalid-query-order-by", ""); err == nil {
							w.WriteHeader(http.StatusBadRequest)

						} else {
							log.Println("Error while writing response. Details:", err)
							w.WriteHeader(http.StatusInternalServerError)
						}

						return
					}

					pagination.OrderBy = append(pagination.OrderBy, dao.ScanDAOSort{
						Field:     orderByField,
						Direction: orderByDirection,
					})
				}

			case "pagesize":
				var err error
				pagination.PageSize, err = strconv.Atoi(value)
				if err != nil {
					if err := h.MessageResponse("invalid-query-page-size", ""); err == nil {
						w.WriteHeader(http.StatusBadRequest)

					} else {
						log.Println("Error while writing response. Details:", err)
						w.WriteHeader(http.StatusInternalServerError)
					}

					return
				}

			case "page":
				var err error
				pagination.Page, err = strconv.Atoi(value)
				if err != nil {
					if err := h.MessageResponse("invalid-query-page", ""); err == nil {
						w.WriteHeader(http.StatusBadRequest)

					} else {
						log.Println("Error while writing response. Details:", err)
						w.WriteHeader(http.StatusInternalServerError)
					}

					return
				}

			case "expand":
				expand = true

			case "current":
				returnCurrent = true
			}
		}
	}

	scanDAO := dao.ScanDAO{
		Database: h.GetDatabase(),
	}

	// As we need to inform the user about the number of items, we always try to retrieve the scan
	// objects even if is requested only the current object
	scans, err := scanDAO.FindAll(&pagination, expand)
	if err != nil {
		log.Println("Error while searching scans objects. Details:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var scansResponse protocol.ScansResponse
	var current model.CurrentScan

	if returnCurrent {
		// The current page will be page zero to avoid misunderstandment
		pagination.Page = 0
		current = model.GetCurrentScan()
		scansResponse = protocol.CurrentScanToScansResponse(current, pagination)

	} else {
		scansResponse = protocol.ScansToScansResponse(scans, pagination)
	}

	h.Response = &scansResponse

	// Last-Modified is going to be the most recent date of the list
	if returnCurrent {
		h.lastModifiedAt = current.LastModifiedAt

	} else {
		for _, scan := range scans {
			if scan.LastModifiedAt.After(h.lastModifiedAt) {
				h.lastModifiedAt = scan.LastModifiedAt
			}
		}
	}

	w.Header().Add("ETag", h.GetETag())
	w.Header().Add("Last-Modified", h.lastModifiedAt.Format(time.RFC1123))
	w.WriteHeader(http.StatusOK)
}

func (h *ScansHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(new(interceptor.Permission)).
		Chain(interceptor.NewValidator(h)).
		Chain(interceptor.NewDatabase(h)).
		Chain(interceptor.NewJSONCodec(h)).
		Chain(interceptor.NewHTTPCacheAfter(h))
}

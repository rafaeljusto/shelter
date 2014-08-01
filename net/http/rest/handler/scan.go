// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package handler store the REST handlers of specific URI
package handler

import (
	"github.com/rafaeljusto/handy"
	"github.com/rafaeljusto/shelter/messages"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/interceptor"
	"github.com/rafaeljusto/shelter/protocol"
	"github.com/trajber/handy"
	"gopkg.in/mgo.v2"
	"net/http"
	"strconv"
	"time"
)

func init() {
	HandleFunc("/scan/{started-at}", func() handy.Handler {
		return new(ScanHandler)
	})
}

// ScanHandler is responsable for keeping the state of a /scan/{started-at} resource
type ScanHandler struct {
	handy.DefaultHandler                           // Inject the HTTP methods that this resource does not implement
	database             *mgo.Database             // Database connection of the MongoDB session
	databaseSession      *mgo.Session              // MongoDB session
	scan                 model.Scan                // Scan object related to the resource
	language             *messages.LanguagePack    // User preferred language based on HTTP header
	StartedAt            string                    `param:"started-at"` // Scan start date in the URI
	Response             *protocol.ScanResponse    `response:"get"`     // Scan response sent back to the user
	Message              *protocol.MessageResponse `error`              // Message on error sent to the user
}

func (h *ScanHandler) SetDatabaseSession(session *mgo.Session) {
	h.databaseSession = session
}

func (h *ScanHandler) GetDatabaseSession() *mgo.Session {
	return h.databaseSession
}

func (h *ScanHandler) SetDatabase(database *mgo.Database) {
	h.database = database
}

func (h *ScanHandler) GetDatabase() *mgo.Database {
	return h.database
}

func (h *ScanHandler) SetScan(scan model.Scan) {
	h.scan = scan
}

func (h *ScanHandler) GetLastModifiedAt() time.Time {
	return h.scan.LastModifiedAt
}

func (h *ScanHandler) GetETag() string {
	return strconv.Itoa(h.scan.Revision)
}

func (h *ScanHandler) SetLanguage(language *messages.LanguagePack) {
	h.language = language
}

func (h *ScanHandler) GetLanguage() *messages.LanguagePack {
	return h.language
}

func (h *ScanHandler) GetStartedAt() string {
	return h.StartedAt
}

func (h *ScanHandler) MessageResponse(messageId string, roid string) error {
	var err error
	h.Message, err = protocol.NewMessageResponse(messageId, roid, h.language)
	return err
}

func (h *ScanHandler) ClearResponse() {
	h.Response = nil
}

func (h *ScanHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.retrieveScan(w, r)
}

func (h *ScanHandler) Head(w http.ResponseWriter, r *http.Request) {
	h.retrieveScan(w, r)
}

func (h *ScanHandler) retrieveScan(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("ETag", strconv.Itoa(h.scan.Revision))
	w.Header().Add("Last-Modified", h.scan.LastModifiedAt.Format(time.RFC1123))
	w.WriteHeader(http.StatusOK)

	scanResponse := protocol.ScanToScanResponse(h.scan)
	h.Response = &scanResponse
}

func (h *ScanHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(new(interceptor.Permission)).
		Chain(interceptor.NewValidator(h)).
		Chain(interceptor.NewDatabase(h)).
		Chain(interceptor.NewScan(h)).
		Chain(interceptor.NewHTTPCacheBefore(h)).
		Chain(interceptor.NewJSONCodec(h))
}

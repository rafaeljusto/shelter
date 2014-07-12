// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"net/http"
	"time"
)

type HTTPCacheHandler interface {
	GetLastModifiedAt() time.Time
	GetETag() string
	MessageResponse(string, string) error
	ClearResponse()
}

func CheckHTTPCache(w http.ResponseWriter, r *http.Request, handler HTTPCacheHandler) {
	if !checkIfModifiedSince(w, r, handler) {
		return
	}

	if !checkIfUnmodifiedSince(w, r, handler) {
		return
	}

	if !checkIfMatch(w, r, handler) {
		return
	}

	checkIfNoneMatch(w, r, handler)
}

func checkIfModifiedSince(w http.ResponseWriter, r *http.Request, handler HTTPCacheHandler) bool {
	modifiedSince, err := check.HTTPIfModifiedSince(r, handler.GetLastModifiedAt())
	if err != nil {
		handler.ClearResponse()

		if err := handler.MessageResponse("invalid-header-date", r.URL.RequestURI()); err == nil {
			w.WriteHeader(http.StatusBadRequest)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return false

	} else if !modifiedSince {
		// If the requested variant has not been modified since the time specified in this
		// field, an entity will not be returned from the server; instead, a 304 (not
		// modified) response will be returned without any message-body
		handler.ClearResponse()
		w.WriteHeader(http.StatusNotModified)
		return false
	}

	return true
}

func checkIfUnmodifiedSince(w http.ResponseWriter, r *http.Request, handler HTTPCacheHandler) bool {
	unmodifiedSince, err := check.HTTPIfUnmodifiedSince(r, handler.GetLastModifiedAt())
	if err != nil {
		handler.ClearResponse()

		if err := handler.MessageResponse("invalid-header-date", r.URL.RequestURI()); err == nil {
			w.WriteHeader(http.StatusBadRequest)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return false

	} else if !unmodifiedSince {
		// If the requested variant has been modified since the specified time, the server
		// MUST NOT perform the requested operation, and MUST return a 412 (Precondition
		// Failed)
		handler.ClearResponse()
		w.WriteHeader(http.StatusPreconditionFailed)
		return false
	}

	return true
}

func checkIfMatch(w http.ResponseWriter, r *http.Request, handler HTTPCacheHandler) bool {
	match := check.HTTPIfMatch(r, handler.GetETag())
	if !match {
		handler.ClearResponse()

		// If "*" is given and no current entity exists or if none of the entity tags match
		// the server MUST NOT perform the requested method, and MUST return a 412
		// (Precondition Failed) response
		if err := handler.MessageResponse("if-match-failed", r.URL.RequestURI()); err == nil {
			w.WriteHeader(http.StatusPreconditionFailed)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return false
	}

	return true
}

func checkIfNoneMatch(w http.ResponseWriter, r *http.Request, handler HTTPCacheHandler) bool {
	noneMatch := check.HTTPIfNoneMatch(r, handler.GetETag())
	if !noneMatch {
		handler.ClearResponse()

		// Instead, if the request method was GET or HEAD, the server SHOULD respond with a
		// 304 (Not Modified) response, including the cache-related header fields
		// (particularly ETag) of one of the entities that matched. For all other request
		// methods, the server MUST respond with a status of 412 (Precondition Failed)
		if r.Method == "GET" || r.Method == "HEAD" {
			// The 304 response MUST NOT contain a message-body, and thus is always terminated
			// by the first empty line after the header fields.
			w.Header().Add("ETag", handler.GetETag())
			w.WriteHeader(http.StatusNotModified)

		} else {
			if err := handler.MessageResponse("if-none-match-failed", r.URL.RequestURI()); err == nil {
				w.WriteHeader(http.StatusPreconditionFailed)

			} else {
				log.Println("Error while writing response. Details:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}

		}
		return false
	}
	return true
}

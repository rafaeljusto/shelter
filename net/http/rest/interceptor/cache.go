// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"github.com/trajber/handy/interceptor"
	"net/http"
	"time"
)

type CacheHandler interface {
	LastModified() time.Time
	ETag() string
	MessageResponse(string, string) error
}

type Cache struct {
	interceptor.NoAfterInterceptor
	cacheHandler CacheHandler
}

func NewCache(h CacheHandler) *Cache {
	return &Cache{cacheHandler: h}
}

func (i *Cache) Before(w http.ResponseWriter, r *http.Request) {
	if !i.checkIfModifiedSince(w, r) {
		return
	}
	if !i.checkIfUnmodifiedSince(w, r) {
		return
	}
	if !i.checkIfMatch(w, r) {
		return
	}
	i.checkIfNoneMatch(w, r)
}

func (i *Cache) checkIfModifiedSince(w http.ResponseWriter, r *http.Request) bool {
	modifiedSince, err := check.HTTPIfModifiedSince(r, i.cacheHandler.LastModified())
	if err != nil {
		if err := i.cacheHandler.MessageResponse("invalid-header-date", r.URL.RequestURI()); err == nil {
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
		w.WriteHeader(http.StatusNotModified)
		return false
	}

	return true
}

func (i *Cache) checkIfUnmodifiedSince(w http.ResponseWriter, r *http.Request) bool {
	unmodifiedSince, err := check.HTTPIfUnmodifiedSince(r, i.cacheHandler.LastModified())
	if err != nil {
		if err := i.cacheHandler.MessageResponse("invalid-header-date", r.URL.RequestURI()); err == nil {
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
		w.WriteHeader(http.StatusPreconditionFailed)
		return false
	}

	return true
}

func (i *Cache) checkIfMatch(w http.ResponseWriter, r *http.Request) bool {
	match, err := check.HTTPIfMatch(r, i.cacheHandler.ETag())
	if err != nil {
		if err := i.cacheHandler.MessageResponse("invalid-if-match", r.URL.RequestURI()); err == nil {
			w.WriteHeader(http.StatusBadRequest)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return false

	} else if !match {
		// If "*" is given and no current entity exists or if none of the entity tags match
		// the server MUST NOT perform the requested method, and MUST return a 412
		// (Precondition Failed) response
		if err := i.cacheHandler.MessageResponse("if-match-failed", r.URL.RequestURI()); err == nil {
			w.WriteHeader(http.StatusPreconditionFailed)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return false
	}

	return true
}

func (i *Cache) checkIfNoneMatch(w http.ResponseWriter, r *http.Request) bool {
	noneMatch, err := check.HTTPIfNoneMatch(r, i.cacheHandler.ETag())
	if err != nil {
		if err := i.cacheHandler.MessageResponse("invalid-if-none-match", r.URL.RequestURI()); err == nil {
			w.WriteHeader(http.StatusBadRequest)

		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return false

	} else if !noneMatch {
		// Instead, if the request method was GET or HEAD, the server SHOULD respond with a
		// 304 (Not Modified) response, including the cache-related header fields
		// (particularly ETag) of one of the entities that matched. For all other request
		// methods, the server MUST respond with a status of 412 (Precondition Failed)
		if r.Method == "GET" || r.Method == "HEAD" {
			// The 304 response MUST NOT contain a message-body, and thus is always terminated
			// by the first empty line after the header fields.
			w.Header().Add("ETag", i.cacheHandler.ETag())
			w.WriteHeader(http.StatusNotModified)

		} else {
			if err := i.cacheHandler.MessageResponse("if-none-match-failed", r.URL.RequestURI()); err == nil {
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

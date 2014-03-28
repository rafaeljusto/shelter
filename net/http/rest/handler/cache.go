// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"net/http"
	"strconv"
	"time"
)

func CheckHTTPCacheHeaders(
	r *http.Request,
	context *context.Context,
	lastModified time.Time,
	revision int,
) bool {

	return checkIfModifiedSince(r, context, lastModified) &&
		checkIfUnmodifiedSince(r, context, lastModified) &&
		checkIfMatch(r, context, revision) &&
		checkIfNoneMatch(r, context, revision)
}

func checkIfModifiedSince(r *http.Request, context *context.Context, lastModified time.Time) bool {
	modifiedSince, err := check.HTTPIfModifiedSince(r, lastModified)
	if err != nil {
		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-header-date", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return false

	} else if !modifiedSince {
		// If the requested variant has not been modified since the time specified in this
		// field, an entity will not be returned from the server; instead, a 304 (not
		// modified) response will be returned without any message-body
		context.Response(http.StatusNotModified)
		return false
	}

	return true
}

func checkIfUnmodifiedSince(r *http.Request, context *context.Context, lastModified time.Time) bool {
	unmodifiedSince, err := check.HTTPIfUnmodifiedSince(r, lastModified)
	if err != nil {
		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-header-date", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return false

	} else if !unmodifiedSince {
		// If the requested variant has been modified since the specified time, the server
		// MUST NOT perform the requested operation, and MUST return a 412 (Precondition
		// Failed)
		context.Response(http.StatusPreconditionFailed)
		return false
	}

	return true
}

func checkIfMatch(r *http.Request, context *context.Context, revision int) bool {
	match, err := check.HTTPIfMatch(r, revision)
	if err != nil {
		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-if-match", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return false

	} else if !match {
		// If "*" is given and no current entity exists or if none of the entity tags match
		// the server MUST NOT perform the requested method, and MUST return a 412
		// (Precondition Failed) response
		if err := context.MessageResponse(http.StatusPreconditionFailed,
			"if-match-failed", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return false
	}

	return true
}

func checkIfNoneMatch(r *http.Request, context *context.Context, revision int) bool {
	noneMatch, err := check.HTTPIfNoneMatch(r, revision)
	if err != nil {
		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-if-none-match", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return false

	} else if !noneMatch {
		// Instead, if the request method was GET or HEAD, the server SHOULD respond with a
		// 304 (Not Modified) response, including the cache-related header fields
		// (particularly ETag) of one of the entities that matched. For all other request
		// methods, the server MUST respond with a status of 412 (Precondition Failed)
		if r.Method == "GET" || r.Method == "HEAD" {
			context.AddHeader("ETag", strconv.Itoa(revision))
			if err := context.MessageResponse(http.StatusNotModified,
				"if-match-none-failed", r.URL.RequestURI()); err != nil {

				log.Println("Error while writing response. Details:", err)
				context.Response(http.StatusInternalServerError)
			}

		} else {
			if err := context.MessageResponse(http.StatusPreconditionFailed,
				"if-match-none-failed", r.URL.RequestURI()); err != nil {

				log.Println("Error while writing response. Details:", err)
				context.Response(http.StatusInternalServerError)
			}
		}
		return false
	}
	return true
}

package handler

import (
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func init() {
	HandleFunc(regexp.MustCompile(`^/scan/([[:alnum:]]|\-|\.|\:)+$`), HandleScan)
}

func HandleScan(r *http.Request, context *context.Context) {
	date, current, err := getScanIdFromURI(r.URL.Path)
	if err != nil {
		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-uri", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "GET" || r.Method == "HEAD" {
		if current {
			retrieveCurrentScan(r, context)
		} else {
			retrieveScan(r, context, date)
		}

	} else {
		context.Response(http.StatusMethodNotAllowed)
	}
}

func retrieveCurrentScan(r *http.Request, context *context.Context) {
	if err := context.JSONResponse(http.StatusOK,
		protocol.CurrentScanToScanResponse(model.GetCurrentScan())); err != nil {

		log.Println("Error while writing response. Details:", err)
		context.Response(http.StatusInternalServerError)
	}
}

func retrieveScan(r *http.Request, context *context.Context, date time.Time) {
	scanDAO := dao.ScanDAO{
		Database: context.Database,
	}

	scan, err := scanDAO.FindByStartedAt(date)
	if err != nil {
		context.Response(http.StatusNotFound)
		return
	}

	modifiedSince, err := check.HTTPIfModifiedSince(r, scan.LastModifiedAt)
	if err != nil {
		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-header-date", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return

	} else if !modifiedSince {
		// If the requested variant has not been modified since the time specified in this
		// field, an entity will not be returned from the server; instead, a 304 (not
		// modified) response will be returned without any message-body
		context.Response(http.StatusNotModified)
		return
	}

	unmodifiedSince, err := check.HTTPIfUnmodifiedSince(r, scan.LastModifiedAt)
	if err != nil {
		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-header-date", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return

	} else if !unmodifiedSince {
		// If the requested variant has been modified since the specified time, the server
		// MUST NOT perform the requested operation, and MUST return a 412 (Precondition
		// Failed)
		context.Response(http.StatusPreconditionFailed)
		return
	}

	match, err := check.HTTPIfMatch(r, scan.Revision)
	if err != nil {
		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-if-match", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return

	} else if !match {
		// If "*" is given and no current entity exists or if none of the entity tags match
		// the server MUST NOT perform the requested method, and MUST return a 412
		// (Precondition Failed) response
		if err := context.MessageResponse(http.StatusPreconditionFailed,
			"if-match-failed", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return
	}

	noneMatch, err := check.HTTPIfNoneMatch(r, scan.Revision)
	if err != nil {
		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-if-none-match", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return

	} else if !noneMatch {
		// Instead, if the request method was GET or HEAD, the server SHOULD respond with a
		// 304 (Not Modified) response, including the cache-related header fields
		// (particularly ETag) of one of the entities that matched. For all other request
		// methods, the server MUST respond with a status of 412 (Precondition Failed)
		context.AddHeader("ETag", strconv.Itoa(scan.Revision))
		if err := context.MessageResponse(http.StatusNotModified,
			"if-match-none-failed", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return
	}

	context.AddHeader("ETag", strconv.Itoa(scan.Revision))
	context.AddHeader("Last-Modified", scan.LastModifiedAt.Format(time.RFC1123))

	if err := context.JSONResponse(http.StatusOK,
		protocol.ScanToScanResponse(scan)); err != nil {

		log.Println("Error while writing response. Details:", err)
		context.Response(http.StatusInternalServerError)
	}
}

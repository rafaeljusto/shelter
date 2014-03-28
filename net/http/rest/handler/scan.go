// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
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

	if !CheckHTTPCacheHeaders(r, context, scan.LastModifiedAt, scan.Revision) {
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

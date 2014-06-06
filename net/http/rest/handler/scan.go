// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package handler store the REST handlers of specific URI
package handler

// import (
// 	"crypto/md5"
// 	"encoding/hex"
// 	"github.com/rafaeljusto/shelter/dao"
// 	"github.com/rafaeljusto/shelter/log"
// 	"github.com/rafaeljusto/shelter/model"
// 	"github.com/rafaeljusto/shelter/net/http/rest/context"
// 	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
// 	"net/http"
// 	"regexp"
// 	"strconv"
// 	"time"
// )

// func init() {
// 	HandleFunc(regexp.MustCompile(`^/scan/([[:alnum:]]|\-|\.|\:)+$`), HandleScan)
// }

// func HandleScan(r *http.Request, context *context.Context) {
// 	date, current, err := getScanIdFromURI(r.URL.Path)
// 	if err != nil {
// 		if err := context.MessageResponse(http.StatusBadRequest,
// 			"invalid-uri", r.URL.RequestURI()); err != nil {

// 			log.Println("Error while writing response. Details:", err)
// 			context.Response(http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	if r.Method == "GET" || r.Method == "HEAD" {
// 		if current {
// 			retrieveCurrentScan(r, context)
// 		} else {
// 			retrieveScan(r, context, date)
// 		}

// 	} else {
// 		context.Response(http.StatusMethodNotAllowed)
// 	}
// }

// func retrieveCurrentScan(r *http.Request, context *context.Context) {
// 	currentScan := model.GetCurrentScan()

// 	if err := context.JSONResponse(http.StatusOK,
// 		protocol.CurrentScanToScanResponse(currentScan)); err != nil {

// 		log.Println("Error while writing response. Details:", err)
// 		context.Response(http.StatusInternalServerError)
// 		return
// 	}

// 	hash := md5.New()
// 	if _, err := hash.Write(context.ResponseContent); err != nil {
// 		log.Println("Error calculating response ETag. Details:", err)
// 		context.Response(http.StatusInternalServerError)
// 		return
// 	}

// 	// The ETag header will be the hash of the content on list services
// 	etag := hex.EncodeToString(hash.Sum(nil))

// 	// TODO: We don't support Last-Modified related cache conditions, what are we going to
// 	// do? Just remove the headers or give a bad request?
// 	if !CheckHTTPCacheHeaders(r, context, time.Time{}, etag) {
// 		return
// 	}

// 	context.AddHeader("ETag", etag)
// }

// func retrieveScan(r *http.Request, context *context.Context, date time.Time) {
// 	scanDAO := dao.ScanDAO{
// 		Database: context.Database,
// 	}

// 	scan, err := scanDAO.FindByStartedAt(date)
// 	if err != nil {
// 		context.Response(http.StatusNotFound)
// 		return
// 	}

// 	if !CheckHTTPCacheHeaders(r, context, scan.LastModifiedAt, strconv.Itoa(scan.Revision)) {
// 		return
// 	}

// 	context.AddHeader("ETag", strconv.Itoa(scan.Revision))
// 	context.AddHeader("Last-Modified", scan.LastModifiedAt.Format(time.RFC1123))

// 	if err := context.JSONResponse(http.StatusOK,
// 		protocol.ScanToScanResponse(scan)); err != nil {

// 		log.Println("Error while writing response. Details:", err)
// 		context.Response(http.StatusInternalServerError)
// 	}
// }

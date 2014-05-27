// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package handler store the REST handlers of specific URI
package handler

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	HandleFunc(regexp.MustCompile("^/domains(/.*)?$"), HandleDomains)
}

func HandleDomains(r *http.Request, context *context.Context) {
	if r.Method == "GET" || r.Method == "HEAD" {
		retrieveDomains(r, context)

	} else {
		context.Response(http.StatusMethodNotAllowed)
	}
}

// The HEAD method is identical to GET except that the server MUST NOT return a message-
// body in the response. But now the responsability for don't adding the body is from the
// mux while writing the response
func retrieveDomains(r *http.Request, context *context.Context) {
	var pagination dao.DomainDAOPagination

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
						if err := context.MessageResponse(http.StatusBadRequest,
							"invalid-query-order-by", ""); err != nil {

							log.Println("Error while writing response. Details:", err)
							context.Response(http.StatusInternalServerError)
						}
						return
					}

					orderByField, err := dao.DomainDAOOrderByFieldFromString(field)
					if err != nil {
						if err := context.MessageResponse(http.StatusBadRequest,
							"invalid-query-order-by", ""); err != nil {

							log.Println("Error while writing response. Details:", err)
							context.Response(http.StatusInternalServerError)
						}
						return
					}

					orderByDirection, err := dao.DAOOrderByDirectionFromString(direction)
					if err != nil {
						if err := context.MessageResponse(http.StatusBadRequest,
							"invalid-query-order-by", ""); err != nil {

							log.Println("Error while writing response. Details:", err)
							context.Response(http.StatusInternalServerError)
						}
						return
					}

					pagination.OrderBy = append(pagination.OrderBy, dao.DomainDAOSort{
						Field:     orderByField,
						Direction: orderByDirection,
					})
				}

			case "pagesize":
				var err error
				pagination.PageSize, err = strconv.Atoi(value)
				if err != nil {
					if err := context.MessageResponse(http.StatusBadRequest,
						"invalid-query-page-size", ""); err != nil {

						log.Println("Error while writing response. Details:", err)
						context.Response(http.StatusInternalServerError)
					}
					return
				}

			case "page":
				var err error
				pagination.Page, err = strconv.Atoi(value)
				if err != nil {
					if err := context.MessageResponse(http.StatusBadRequest,
						"invalid-query-page", ""); err != nil {

						log.Println("Error while writing response. Details:", err)
						context.Response(http.StatusInternalServerError)
					}
					return
				}
			}
		}
	}

	domainDAO := dao.DomainDAO{
		Database: context.Database,
	}

	domains, err := domainDAO.FindAll(&pagination)
	if err != nil {
		log.Println("Error while searching domains objects. Details:", err)
		context.Response(http.StatusInternalServerError)
		return
	}

	if err := context.JSONResponse(http.StatusOK,
		protocol.ToDomainsResponse(domains, pagination)); err != nil {

		log.Println("Error while writing response. Details:", err)
		context.Response(http.StatusInternalServerError)
		return
	}

	hash := md5.New()
	if _, err := hash.Write(context.ResponseContent); err != nil {
		log.Println("Error calculating response ETag. Details:", err)
		context.Response(http.StatusInternalServerError)
		return
	}

	// The ETag header will be the hash of the content on list services
	etag := hex.EncodeToString(hash.Sum(nil))

	// Last-Modified is going to be the most recent date of the list
	var lastModifiedAt time.Time
	for _, domain := range domains {
		if domain.LastModifiedAt.After(lastModifiedAt) {
			lastModifiedAt = domain.LastModifiedAt
		}
	}

	if !CheckHTTPCacheHeaders(r, context, lastModifiedAt, etag) {
		return
	}

	context.AddHeader("ETag", etag)
	context.AddHeader("Last-Modified", lastModifiedAt.Format(time.RFC1123))
}

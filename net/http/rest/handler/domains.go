package handler

import (
	"net/http"
	"shelter/dao"
	"shelter/net/http/rest"
	"shelter/net/http/rest/context"
	"strconv"
	"strings"
)

func init() {
	rest.HandleFunc("/domains", handleDomains)
}

func handleDomains(r *http.Request, context *context.ShelterRESTContext) {
	if r.Method == "GET" {
		retrieveDomains(r, context)

	} else {
		context.Response(http.StatusMethodNotAllowed)
	}
}

func retrieveDomains(r *http.Request, context *context.ShelterRESTContext) {
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
						context.MessageResponse(http.StatusBadRequest, "invalid-query-order-by")
						return
					}

					var orderByField dao.DomainDAOOrderByField

					if field == "fqdn" {
						orderByField = dao.DomainDAOOrderByFieldFQDN

					} else if field == "lastmodified" {
						orderByField = dao.DomainDAOOrderByFieldLastModifiedAt

					} else {
						context.MessageResponse(http.StatusBadRequest, "invalid-query-order-by")
						return
					}

					var orderByDirection dao.DomainDAOOrderByDirection

					if direction == "asc" {
						orderByDirection = dao.DomainDAOOrderByDirectionAscending

					} else if direction == "desc" {
						orderByDirection = dao.DomainDAOOrderByDirectionDescending

					} else {
						context.MessageResponse(http.StatusBadRequest, "invalid-query-order-by")
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
					context.MessageResponse(http.StatusBadRequest, "invalid-query-page-size")
					return
				}

			case "page":
				var err error
				pagination.Page, err = strconv.Atoi(value)
				if err != nil {
					context.MessageResponse(http.StatusBadRequest, "invalid-query-page")
					return
				}
			}
		}
	}
}

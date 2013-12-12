package rest

import (
	"net/http"
	"shelter/dao"
	"shelter/net/http/rest"
	"shelter/net/http/rest/language"
	"strconv"
	"strings"
	"time"
)

func init() {
	rest.HandleFunc("/domain/", handleDomain)
}

func handleDomain(r *http.Request, context *ShelterRESTContext) {
	if r.Method == "GET" {
		retrieveDomain(r, context)

	} else if r.Method == "PUT" {
		createUpdateDomain(r, context)

	} else if r.Method == "DELETE" {
		removeDomain(r, context)

	} else {
		context.Response(http.StatusMethodNotAllowed)
	}
}

func retrieveDomain(r *http.Request, context *ShelterRESTContext) {
	fqdn := getFQDNFromURI(r.URL.Path)
	if len(fqdn) == 0 {
		context.ResponseMessage(http.StatusBadRequest, "invalid-uri")
		return
	}

	domainDAO := dao.DomainDAO{
		Database: context.Database,
	}

	domain, err := domainDAO.FindByFQDN(fqdn)
	if err != nil {
		context.Response(http.StatusNotFound)
		return
	}

	if len(r.Header.Get("If-Modified-Since")) > 0 {
		ifModifiedSince, err := time.Parse(time.RFC1123, r.Header.Get("If-Modified-Since"))
		if err != nil {
			context.ResponseMessage(http.StatusBadRequest, "invalid-header-date")
			return
		}

		if domain.LastModifiedAt.Before(ifModifiedSince) || domain.LastModifiedAt.Equal(ifModifiedSince) {
			context.Response(http.StatusNotModified)
			return
		}
	}

	if len(r.Header.Get("If-Unmodified-Since")) > 0 {
		ifUnmodifiedSince, err := time.Parse(time.RFC1123, r.Header.Get("If-Unmodified-Since"))
		if err != nil {
			context.ResponseMessage(http.StatusBadRequest, "invalid-header-date")
			return
		}

		if domain.LastModifiedAt.After(ifModifiedSince) {
			context.Response(http.StatusPreconditionFailed)
			return
		}
	}

	context.AddHeader("ETag", fmt.Sprintf("%d", domain.Revision))
	context.AddHeader("Last-Modified", domain.LastModifiedAt.Format(time.RFC1123))
	context.ResponseJSON(http.StatusOK, domain)
}

func createUpdateDomain(r *http.Request, context *ShelterRESTContext) {

}

func removeDomain(r *http.Request, context *ShelterRESTContext) {

}

// Retrieve the FQDN from the URI. The FQDN is going to be the last part of the URI. For
// example "/domain/rafael.net.br" will return "rafael.net.br". If there's any error an
// empty string is returned
func getFQDNFromURI(uri string) string {
	idx := strings.LastIndex(uri, "/")
	if idx == -1 {
		return ""
	}

	return uri[idx:]
}

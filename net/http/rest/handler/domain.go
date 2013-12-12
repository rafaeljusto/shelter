package rest

import (
	"net/http"
	"shelter/dao"
	"shelter/net/http/rest"
	"shelter/net/http/rest/language"
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
	// TODO Check:
	//   If-Modified-Since
	//   If-Match
	//   If-None-Match

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

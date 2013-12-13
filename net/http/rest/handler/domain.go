package rest

import (
	"net/http"
	"shelter/dao"
	"shelter/net/http/rest"
	"shelter/net/http/rest/language"
	"shelter/net/http/rest/protocol"
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
	fqdn := getFQDNFromURI(r.URL.Path)
	if len(fqdn) == 0 {
		context.ResponseMessage(http.StatusBadRequest, "invalid-uri")
		return
	}

	var domainRequest protocol.DomainRequest
	if err := context.RequestJSON(&domainRequest); err != nil {
		context.ResponseMessage(http.StatusBadRequest, "invalid-json-content")
		return
	}

	// We need to set the FQDN in the domain request object because it is sent only in the
	// URI and not in the domain request body to avoid information redudancy
	domainRequest.FQDN = fqdn

	domainDAO := dao.DomainDAO{
		Database: context.Database,
	}

	// We need to load the domain from the database and merge it with the changes from the
	// user, if the domain does not exist yet thats alright because we will create it
	domain, _ := domainDAO.FindByFQDN(fqdn)
	domain, err = transcoder.Merge(domain, domainRequest)

	if err != nil {
		// TODO: Log!
		context.Response(http.StatusInternalServerError)
		return
	}

	if err := domainDAO.Save(&domain); err != nil {
		// TODO: Log!
		context.Response(http.StatusInternalServerError)
		return
	}

	context.AddHeader("ETag", fmt.Sprintf("%d", domain.Revision))
	context.AddHeader("Last-Modified", domain.LastModifiedAt.Format(time.RFC1123))
	context.AddHeader("Location", "/domain/"+domain.FQDN)

	if domain.Revision == 1 {
		context.Response(http.StatusCreated)
	} else {
		context.Response(http.StatusNoContent)
	}
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

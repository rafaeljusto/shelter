package rest

import (
	"fmt"
	"net/http"
	"shelter/dao"
	"shelter/net/http/rest"
	"shelter/net/http/rest/check"
	"shelter/net/http/rest/context"
	"shelter/net/http/rest/protocol"
	"strings"
	"time"
)

func init() {
	rest.HandleFunc("/domain/", handleDomain)
}

func handleDomain(r *http.Request, context *context.ShelterRESTContext) {
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

func retrieveDomain(r *http.Request, context *context.ShelterRESTContext) {
	fqdn := getFQDNFromURI(r.URL.Path)
	if len(fqdn) == 0 {
		context.MessageResponse(http.StatusBadRequest, "invalid-uri")
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

	modifiedSince, err := check.IfModifiedSince(r, domain.LastModifiedAt)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-header-date")
		return

	} else if !modifiedSince {
		// If the requested variant has not been modified since the time specified in this
		// field, an entity will not be returned from the server; instead, a 304 (not
		// modified) response will be returned without any message-body
		context.Response(http.StatusNotModified)
		return
	}

	unmodifiedSince, err := check.IfUnmodifiedSince(r, domain.LastModifiedAt)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-header-date")
		return

	} else if !unmodifiedSince {
		// If the requested variant has been modified since the specified time, the server
		// MUST NOT perform the requested operation, and MUST return a 412 (Precondition
		// Failed)
		context.Response(http.StatusPreconditionFailed)
		return
	}

	match, err := check.IfMatch(r, domain.Revision)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-if-match")
		return

	} else if !match {
		// If "*" is given and no current entity exists or if none of the entity tags match
		// the server MUST NOT perform the requested method, and MUST return a 412
		// (Precondition Failed) response
		context.MessageResponse(http.StatusPreconditionFailed, "if-match-failed")
	}

	context.AddHeader("ETag", fmt.Sprintf("%d", domain.Revision))
	context.AddHeader("Last-Modified", domain.LastModifiedAt.Format(time.RFC1123))
	context.JSONResponse(http.StatusOK, protocol.ToDomainResponse(domain))
}

func createUpdateDomain(r *http.Request, context *context.ShelterRESTContext) {
	fqdn := getFQDNFromURI(r.URL.Path)
	if len(fqdn) == 0 {
		context.MessageResponse(http.StatusBadRequest, "invalid-uri")
		return
	}

	var domainRequest protocol.DomainRequest
	if err := context.JSONRequest(&domainRequest); err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-json-content")
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

	modifiedSince, err := check.IfModifiedSince(r, domain.LastModifiedAt)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-header-date")
		return

	} else if !modifiedSince {
		// If the requested variant has not been modified since the time specified in this
		// field, an entity will not be returned from the server; instead, a 304 (not
		// modified) response will be returned without any message-body
		context.Response(http.StatusNotModified)
		return
	}

	unmodifiedSince, err := check.IfUnmodifiedSince(r, domain.LastModifiedAt)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-header-date")
		return

	} else if !unmodifiedSince {
		// If the requested variant has been modified since the specified time, the server
		// MUST NOT perform the requested operation, and MUST return a 412 (Precondition
		// Failed)
		context.Response(http.StatusPreconditionFailed)
		return
	}

	if domain, err = protocol.Merge(domain, domainRequest); err != nil {
		// TODO: Log!
		context.Response(http.StatusInternalServerError)
		return
	}

	match, err := check.IfMatch(r, domain.Revision)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-if-match")
		return

	} else if !match {
		// If "*" is given and no current entity exists or if none of the entity tags match
		// the server MUST NOT perform the requested method, and MUST return a 412
		// (Precondition Failed) response
		context.MessageResponse(http.StatusPreconditionFailed, "if-match-failed")
	}

	if err := domainDAO.Save(&domain); err != nil {
		if strings.Index(err.Error(), "duplicate key error index") != -1 {
			context.MessageResponse(http.StatusConflict, "conflict")

		} else {
			// TODO: Log!
			context.Response(http.StatusInternalServerError)
		}

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

func removeDomain(r *http.Request, context *context.ShelterRESTContext) {

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

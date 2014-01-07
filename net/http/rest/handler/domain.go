package handler

import (
	"fmt"
	"net/http"
	"shelter/dao"
	"shelter/net/http/rest/check"
	"shelter/net/http/rest/context"
	"shelter/net/http/rest/log"
	"shelter/net/http/rest/protocol"
	"strings"
	"time"
)

func init() {
	HandleFunc("/domain/", HandleDomain)
}

func HandleDomain(r *http.Request, context *context.Context) {
	fqdn := getFQDNFromURI(r.URL.Path)
	if len(fqdn) == 0 {
		context.MessageResponse(http.StatusBadRequest, "invalid-uri")
		return
	}

	if r.Method == "GET" {
		retrieveDomain(r, context, fqdn, true)

	} else if r.Method == "PUT" {
		createUpdateDomain(r, context, fqdn)

	} else if r.Method == "DELETE" {
		removeDomain(r, context, fqdn)

	} else if r.Method == "HEAD" {
		retrieveDomain(r, context, fqdn, false)

	} else {
		context.Response(http.StatusMethodNotAllowed)
	}
}

// The HEAD method is identical to GET except that the server MUST NOT return a message-
// body in the response. For that reason we have the domainInResponseParameter
func retrieveDomain(r *http.Request, context *context.Context, fqdn string, domainInResponse bool) {
	domainDAO := dao.DomainDAO{
		Database: context.Database,
	}

	domain, err := domainDAO.FindByFQDN(fqdn)
	if err != nil {
		context.Response(http.StatusNotFound)
		return
	}

	modifiedSince, err := check.HTTPIfModifiedSince(r, domain.LastModifiedAt)
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

	unmodifiedSince, err := check.HTTPIfUnmodifiedSince(r, domain.LastModifiedAt)
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

	match, err := check.HTTPIfMatch(r, domain.Revision)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-if-match")
		return

	} else if !match {
		// If "*" is given and no current entity exists or if none of the entity tags match
		// the server MUST NOT perform the requested method, and MUST return a 412
		// (Precondition Failed) response
		context.MessageResponse(http.StatusPreconditionFailed, "if-match-failed")
		return
	}

	noneMatch, err := check.HTTPIfNoneMatch(r, domain.Revision)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-if-none-match")
		return

	} else if !noneMatch {
		// Instead, if the request method was GET or HEAD, the server SHOULD respond with a
		// 304 (Not Modified) response, including the cache-related header fields
		// (particularly ETag) of one of the entities that matched. For all other request
		// methods, the server MUST respond with a status of 412 (Precondition Failed)
		context.AddHeader("ETag", fmt.Sprintf("%d", domain.Revision))
		context.MessageResponse(http.StatusNotModified, "if-match-none-failed")
		return
	}

	context.AddHeader("ETag", fmt.Sprintf("%d", domain.Revision))
	context.AddHeader("Last-Modified", domain.LastModifiedAt.Format(time.RFC1123))

	if domainInResponse {
		context.JSONResponse(http.StatusOK, protocol.ToDomainResponse(domain))
	} else {
		context.Response(http.StatusOK)
	}
}

func createUpdateDomain(r *http.Request, context *context.Context, fqdn string) {
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

	modifiedSince, err := check.HTTPIfModifiedSince(r, domain.LastModifiedAt)
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

	unmodifiedSince, err := check.HTTPIfUnmodifiedSince(r, domain.LastModifiedAt)
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

	match, err := check.HTTPIfMatch(r, domain.Revision)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-if-match")
		return

	} else if !match {
		// If "*" is given and no current entity exists or if none of the entity tags match
		// the server MUST NOT perform the requested method, and MUST return a 412
		// (Precondition Failed) response
		context.MessageResponse(http.StatusPreconditionFailed, "if-match-failed")
		return
	}

	noneMatch, err := check.HTTPIfNoneMatch(r, domain.Revision)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-if-none-match")
		return

	} else if !noneMatch {
		// Instead, if the request method was GET or HEAD, the server SHOULD respond with a
		// 304 (Not Modified) response, including the cache-related header fields
		// (particularly ETag) of one of the entities that matched. For all other request
		// methods, the server MUST respond with a status of 412 (Precondition Failed)
		context.MessageResponse(http.StatusPreconditionFailed, "if-match-none-failed")
		return
	}

	if domain, err = protocol.Merge(domain, domainRequest); err != nil {
		log.Println("Error while merging domain objects for create or "+
			"update operation. Details:", err)
		context.Response(http.StatusInternalServerError)
		return
	}

	if err := domainDAO.Save(&domain); err != nil {
		if strings.Index(err.Error(), "duplicate key error index") != -1 {
			context.MessageResponse(http.StatusConflict, "conflict")

		} else {
			log.Println("Error while saving domain object for create or "+
				"update operation. Details:", err)
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

func removeDomain(r *http.Request, context *context.Context, fqdn string) {
	domainDAO := dao.DomainDAO{
		Database: context.Database,
	}

	domain, err := domainDAO.FindByFQDN(fqdn)
	if err != nil {
		context.Response(http.StatusNotFound)
		return
	}

	modifiedSince, err := check.HTTPIfModifiedSince(r, domain.LastModifiedAt)
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

	unmodifiedSince, err := check.HTTPIfUnmodifiedSince(r, domain.LastModifiedAt)
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

	match, err := check.HTTPIfMatch(r, domain.Revision)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-if-match")
		return

	} else if !match {
		// If "*" is given and no current entity exists or if none of the entity tags match
		// the server MUST NOT perform the requested method, and MUST return a 412
		// (Precondition Failed) response
		context.MessageResponse(http.StatusPreconditionFailed, "if-match-failed")
		return
	}

	noneMatch, err := check.HTTPIfNoneMatch(r, domain.Revision)
	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-if-none-match")
		return

	} else if !noneMatch {
		// Instead, if the request method was GET or HEAD, the server SHOULD respond with a
		// 304 (Not Modified) response, including the cache-related header fields
		// (particularly ETag) of one of the entities that matched. For all other request
		// methods, the server MUST respond with a status of 412 (Precondition Failed)
		context.MessageResponse(http.StatusPreconditionFailed, "if-match-none-failed")
		return
	}

	if err := domainDAO.Remove(&domain); err != nil {
		log.Println("Error while removing domain object. Details:", err)
		context.Response(http.StatusInternalServerError)
		return
	}

	context.Response(http.StatusNoContent)
}

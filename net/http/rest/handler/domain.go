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
	"strings"
	"time"
)

func init() {
	HandleFunc(regexp.MustCompile(`^/domain/([[:alnum:]]|\-|\.)+$`), HandleDomain)
}

func HandleDomain(r *http.Request, context *context.Context) {
	fqdn := getFQDNFromURI(r.URL.Path)
	if len(fqdn) == 0 {
		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-uri", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "GET" || r.Method == "HEAD" {
		retrieveDomain(r, context, fqdn)

	} else if r.Method == "PUT" {
		createUpdateDomain(r, context, fqdn)

	} else if r.Method == "DELETE" {
		removeDomain(r, context, fqdn)

	} else {
		context.Response(http.StatusMethodNotAllowed)
	}
}

// The HEAD method is identical to GET except that the server MUST NOT return a message-
// body in the response. But now the responsability for don't adding the body is from the
// mux while writing the response
func retrieveDomain(r *http.Request, context *context.Context, fqdn string) {
	domainDAO := dao.DomainDAO{
		Database: context.Database,
	}

	domain, err := domainDAO.FindByFQDN(fqdn)
	if err != nil {
		context.Response(http.StatusNotFound)
		return
	}

	if !CheckHTTPCacheHeaders(r, context, domain.LastModifiedAt, domain.Revision) {
		return
	}

	context.AddHeader("ETag", strconv.Itoa(domain.Revision))
	context.AddHeader("Last-Modified", domain.LastModifiedAt.Format(time.RFC1123))

	if err := context.JSONResponse(http.StatusOK,
		protocol.ToDomainResponse(domain)); err != nil {

		log.Println("Error while writing response. Details:", err)
		context.Response(http.StatusInternalServerError)
	}
}

func createUpdateDomain(r *http.Request, context *context.Context, fqdn string) {
	var domainRequest protocol.DomainRequest
	if err := context.JSONRequest(&domainRequest); err != nil {
		log.Println("Received an invalid JSON. Details:", err)

		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-json-content", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
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

	if !CheckHTTPCacheHeaders(r, context, domain.LastModifiedAt, domain.Revision) {
		return
	}

	var err error
	if domain, err = protocol.Merge(domain, domainRequest); err != nil {
		messageId := ""

		switch err {
		case model.ErrInvalidFQDN:
			messageId = "invalid-fqdn"
		case protocol.ErrInvalidDNSKEY:
			messageId = "invalid-dnskey"
		case protocol.ErrInvalidDSAlgorithm:
			messageId = "invalid-ds-algorithm"
		case protocol.ErrInvalidDSDigestType:
			messageId = "invalid-ds-digest-type"
		case protocol.ErrInvalidIP:
			messageId = "invalid-ip"
		case protocol.ErrInvalidLanguage:
			messageId = "invalid-language"
		}

		if len(messageId) == 0 {
			log.Println("Error while merging domain objects for create or "+
				"update operation. Details:", err)
			context.Response(http.StatusInternalServerError)

		} else {
			if err := context.MessageResponse(http.StatusBadRequest,
				messageId, r.URL.RequestURI()); err != nil {

				log.Println("Error while writing response. Details:", err)
				context.Response(http.StatusInternalServerError)
			}
		}
		return
	}

	if err := domainDAO.Save(&domain); err != nil {
		if strings.Index(err.Error(), "duplicate key error index") != -1 {
			if err := context.MessageResponse(http.StatusConflict,
				"conflict", r.URL.RequestURI()); err != nil {

				log.Println("Error while writing response. Details:", err)
				context.Response(http.StatusInternalServerError)
			}

		} else {
			log.Println("Error while saving domain object for create or "+
				"update operation. Details:", err)
			context.Response(http.StatusInternalServerError)
		}

		return
	}

	context.AddHeader("ETag", strconv.Itoa(domain.Revision))
	context.AddHeader("Last-Modified", domain.LastModifiedAt.Format(time.RFC1123))

	if domain.Revision == 1 {
		context.AddHeader("Location", "/domain/"+domain.FQDN)
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

	if !CheckHTTPCacheHeaders(r, context, domain.LastModifiedAt, domain.Revision) {
		return
	}

	if err := domainDAO.Remove(&domain); err != nil {
		log.Println("Error while removing domain object. Details:", err)
		context.Response(http.StatusInternalServerError)
		return
	}

	context.Response(http.StatusNoContent)
}

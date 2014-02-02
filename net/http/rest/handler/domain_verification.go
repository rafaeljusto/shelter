package handler

import (
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
	"github.com/rafaeljusto/shelter/net/scan"
	"net/http"
	"regexp"
)

func init() {
	HandleFunc(regexp.MustCompile(`^/domain/([[:alnum:]]|\-|\.)+/verification$`), HandleDomainVerification)
}

func HandleDomainVerification(r *http.Request, context *context.Context) {
	fqdn := getFQDNFromURI(r.URL.Path)
	if len(fqdn) == 0 {
		if err := context.MessageResponse(http.StatusBadRequest,
			"invalid-uri", r.URL.RequestURI()); err != nil {

			log.Println("Error while writing response. Details:", err)
			context.Response(http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "PUT" {
		scanDomain(r, context, fqdn)
	}
}

// scanDomain is responsable for checking a domain object on-the-fly without persisting in database,
// useful for pre-registration validations in the registry
func scanDomain(r *http.Request, context *context.Context, fqdn string) {
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

	var domain model.Domain
	var err error

	if domain, err = protocol.Merge(domain, domainRequest); err != nil {
		log.Println("Error while merging domain objects for create or "+
			"update operation. Details:", err)
		context.Response(http.StatusInternalServerError)
		return
	}

	scan.ScanDomain(&domain)

	if err := context.JSONResponse(http.StatusOK,
		protocol.ToDomainResponse(domain)); err != nil {

		log.Println("Error while writing response. Details:", err)
		context.Response(http.StatusInternalServerError)
	}
}

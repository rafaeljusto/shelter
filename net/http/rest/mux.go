package rest

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"shelter/config"
	"shelter/database/mongodb"
	"shelter/net/http/rest/check"
	"shelter/net/http/rest/context"
	"shelter/net/http/rest/handler"
	"shelter/net/http/rest/log"
	"shelter/net/http/rest/messages"
	"strings"
	"time"
)

// List of possible errors that can occur when calling functions from this file. Other
// erros can also occurs from low level layers
var (
	ErrSecretNotFound = errors.New("Secret related to Authorization's secret id not found")
)

// Main router used by the Shelter REST system to manage the requests
var (
	mux Mux
)

// Mux is responsable for all initial HTTP checks before calling a specific handler, and
// for adding the system HTTP headers on each response
type Mux struct{}

// Main function of the REST server
func (mux Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := mux.findRoute(r.URL.Path)
	if handler == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	database, err := mongodb.Open(
		config.ShelterConfig.Database.URI,
		config.ShelterConfig.Database.Name,
	)

	context, err := context.NewContext(r, database)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating context. Details:", err)
		return
	}

	if mux.checkHTTPHeaders(w, r, &context) {
		handler(r, &context)
	}

	mux.writeResponse(w, context)
}

// Find the best handler for the given URI. The best handler is the most specific one
func (mux Mux) findRoute(uri string) handler.Handler {
	var selectedRoute string
	var selectedHandler handler.Handler

	for route, handler := range handler.Routes {
		if !strings.HasPrefix(uri, route) {
			continue
		}

		// Try to find the most specific route
		if len(selectedRoute) == 0 || strings.HasPrefix(route, selectedRoute) {
			selectedRoute = route
			selectedHandler = handler
		}
	}

	return selectedHandler
}

// Verify HTTP headers and fill context with user preferences
func (mux Mux) checkHTTPHeaders(w http.ResponseWriter,
	r *http.Request, context *context.Context) bool {

	// We first check the language header, because if it's acceptable the next messages are
	// going to be returned in the language choosen by the user
	if !check.HTTPAcceptLanguage(r, context) {
		context.MessageResponse(http.StatusNotAcceptable, "accept-language-error", r.RequestURI)
		return false
	}

	if !check.HTTPAccept(r) {
		context.MessageResponse(http.StatusNotAcceptable, "accept-error", r.RequestURI)
		return false
	}

	if !check.HTTPAcceptCharset(r) {
		context.MessageResponse(http.StatusNotAcceptable, "accept-charset-error", r.RequestURI)
		return false
	}

	if !check.HTTPContentType(r) {
		context.MessageResponse(http.StatusBadRequest, "invalid-content-type", r.RequestURI)
		return false
	}

	if !check.HTTPContentMD5(r, context) {
		context.MessageResponse(http.StatusBadRequest, "invalid-content-md5", r.RequestURI)
		return false
	}

	timeFrameOK, err := check.HTTPDate(r)

	if err != nil {
		context.MessageResponse(http.StatusBadRequest, "invalid-header-date", r.RequestURI)
		return false

	} else if !timeFrameOK {
		context.MessageResponse(http.StatusBadRequest, "invalid-date-time-frame", r.RequestURI)
		return false
	}

	authorized, err := check.HTTPAuthorization(r, func(secretId string) (string, error) {
		secret, ok := config.ShelterConfig.RESTServer.Secrets[secretId]

		if !ok {
			return "", ErrSecretNotFound
		}

		// In the near future the secret will be encrypted in the configuration file and the
		// decrypt process can generate problems
		return secret, nil
	})

	if err != nil {
		if err == check.ErrHTTPContentTypeNotFound {
			context.MessageResponse(http.StatusBadRequest, "content-type-missing", r.RequestURI)

		} else if err == check.ErrHTTPContentMD5NotFound {
			context.MessageResponse(http.StatusBadRequest, "content-md5-missing", r.RequestURI)

		} else if err == check.ErrHTTPDateNotFound {
			context.MessageResponse(http.StatusBadRequest, "date-missing", r.RequestURI)

		} else if err == check.ErrHTTPAuthorizationNotFound {
			context.MessageResponse(http.StatusBadRequest, "authorization-missing", r.RequestURI)

		} else if err == check.ErrInvalidHTTPAuthorization {
			context.MessageResponse(http.StatusBadRequest, "invalid-authorization", r.RequestURI)

		} else if err == ErrSecretNotFound {
			context.MessageResponse(http.StatusBadRequest, "secret-not-found", r.RequestURI)

		} else {
			context.Response(http.StatusInternalServerError)
			log.Println("Error checking authorization. Details:", err)
		}

		return false

	} else if !authorized {
		context.Response(http.StatusUnauthorized)
		return false
	}

	return true
}

// Write response with the defaults HTTP response headers
func (mux Mux) writeResponse(w http.ResponseWriter,
	context context.Context) {

	if len(context.ResponseContent) > 0 {
		w.Header().Add("Content-Type", fmt.Sprintf("application/vnd.shelter+json; charset=%s", check.SupportedCharset))
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(context.ResponseContent)))

		if context.Language != nil {
			w.Header().Add("Content-Language", context.Language.Name())
		}

		hash := md5.New()
		hash.Write(context.ResponseContent)
		hashBytes := hash.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)
		w.Header().Add("Content-MD5", hashBase64)
	}

	w.Header().Add("Accept", check.SupportedContentType)
	w.Header().Add("Accept-Language", messages.ShelterRESTLanguagePacks.Names())
	w.Header().Add("Accept-Charset", "utf-8")
	w.Header().Add("Date", time.Now().UTC().Format(time.RFC1123))

	for key, value := range context.HTTPHeader {
		w.Header().Add(key, value)
	}

	w.WriteHeader(context.ResponseHTTPStatus)

	if len(context.ResponseContent) > 0 {
		w.Write(context.ResponseContent)
	}
}

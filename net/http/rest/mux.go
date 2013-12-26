package rest

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"shelter/config"
	"shelter/net/http/rest/check"
	"shelter/net/http/rest/context"
	"shelter/net/http/rest/language"
	"shelter/net/http/rest/log"
	"strings"
	"time"
)

// Main router used by the Shelter REST system to manage the requests
var (
	mux shelterRESTMux
)

// List of possible errors that can occur when calling functions from this file. Other
// erros can also occurs from low level layers
var (
	ErrSecretNotFound = errors.New("Secret related to Authorization's secret id not found")
)

// Create this type to make it easy to reference a Shelter REST server handler
type shelterRESTHandler func(*http.Request, *context.ShelterRESTContext)

// shelterRESTMux is responsable for storing all routes. Beyond of searching the best
// route for each request, the mux will do all initial HTTP checks before calling the
// handler
type shelterRESTMux struct {
	routes map[string]shelterRESTHandler // Map of all available routes
}

// Function created only to register the handlers more easily in the mux
func HandleFunc(route string, handler shelterRESTHandler) {
	mux.routes[route] = handler
}

// Main function of the REST server
func (mux shelterRESTMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := mux.findRoute(r.URL.Path)
	if handler == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	context, err := context.NewShelterRESTContext(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating context. Details:", err)
		return
	}

	if !mux.checkHTTPHeaders(w, r, &context) {
		return
	}

	handler(r, &context)
	mux.writeResponse(w, context)
}

// Find the best handler for the given URI. The best handler is the most specific one
func (mux shelterRESTMux) findRoute(uri string) shelterRESTHandler {
	var selectedRoute string
	var selectedHandler shelterRESTHandler

	for route, handler := range mux.routes {
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
func (mux shelterRESTMux) checkHTTPHeaders(w http.ResponseWriter,
	r *http.Request, context *context.ShelterRESTContext) bool {

	// We first check the language header, because if it's acceptable the next messages are
	// going to be returned in the language choosen by the user
	if !check.HTTPAcceptLanguage(r, context) {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, context.Language.Messages["accept-language-error"])
		return false
	}

	if !check.HTTPAccept(r) {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, context.Language.Messages["accept-error"])
		return false
	}

	if !check.HTTPAcceptCharset(r) {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, context.Language.Messages["accept-charset-error"])
		return false
	}

	if !check.HTTPContentType(r) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, context.Language.Messages["invalid-content-type"])
		return false
	}

	if !check.HTTPContentMD5(r, context) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, context.Language.Messages["invalid-content-md5"])
		return false
	}

	timeFrameOK, err := check.HTTPDate(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, context.Language.Messages["invalid-header-date"])
		return false

	} else if !timeFrameOK {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, context.Language.Messages["invalid-date-time-frame"])
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, context.Language.Messages["content-type-missing"])

		} else if err == check.ErrHTTPContentMD5NotFound {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, context.Language.Messages["content-md5-missing"])

		} else if err == check.ErrHTTPDateNotFound {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, context.Language.Messages["date-missing"])

		} else if err == check.ErrHTTPAuthorizationNotFound {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, context.Language.Messages["authorization-missing"])

		} else if err == check.ErrInvalidHTTPAuthorization {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, context.Language.Messages["invalid-authorization"])

		} else if err == ErrSecretNotFound {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, context.Language.Messages["secret-not-found"])

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error checking authorization. Details:", err)
		}

		return false

	} else if !authorized {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}

// Write response with the defaults HTTP response headers
func (mux shelterRESTMux) writeResponse(w http.ResponseWriter,
	context context.ShelterRESTContext) {

	if len(context.ResponseContent) > 0 {
		w.Header().Add("Content-Type", "application/vnd.shelter+json")
		w.Header().Add("Content-Encoding", "utf-8")
		w.Header().Add("Content-Language", context.Language.Name())
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(context.ResponseContent)))

		hash := md5.New()
		hash.Write(context.ResponseContent)
		hashBytes := hash.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)
		w.Header().Add("Content-MD5", hashBase64)
	}

	w.Header().Add("Accept", check.SupportedContentType)
	w.Header().Add("Accept-Language", language.ShelterRESTLanguagePacks.Names())
	w.Header().Add("Accept-Charset", "utf-8")
	w.Header().Add("Date", time.Now().UTC().Format(time.RFC1123))

	for key, value := range context.HTTPHeader {
		w.Header().Add(key, value)
	}

	w.WriteHeader(context.ResponseHttpStatus)

	if len(context.ResponseContent) > 0 {
		w.Write(context.ResponseContent)
	}
}

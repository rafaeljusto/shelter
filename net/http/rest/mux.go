package rest

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"shelter/net/http/rest/language"
	"strings"
	"time"
)

// Main router used by the Shelter REST system to manage the requests
var (
	mux shelterRESTMux
)

// Create this type to make it easy to reference a Shelter REST server handler
type shelterRESTHandler func(*http.Request, *ShelterRESTContext)

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

// Supported HTTP headers in request:
//   Method
//   Date
//   Content-type
//   Content-Length
//   Content-MD5
//   Accept
//   Accept-Charset
//   Accept-Language
//   Authorization
//   If-Modified-Since
//   If-Match
//   If-None-Match

// Supported HTTP headers in response:
//   Content-Encoding
//   Content-Language
//   Content-Length
//   Content-MD5
//   Content-Type
//   Date
//   ETag
//   Last-Modified
//   Status
//   Accept-Language
//   Accept

func (mux shelterRESTMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var selectedRoute string
	var selectedHandler shelterRESTHandler

	for route, handler := range mux.routes {
		if !strings.HasPrefix(r.URL.Path, route) {
			continue
		}

		// Try to find the most specific route
		if len(selectedRoute) == 0 || strings.HasPrefix(route, selectedRoute) {
			selectedRoute = route
			selectedHandler = handler
		}
	}

	if len(selectedRoute) == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	context, err := newShelterRESTContext()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: Log!
		return
	}

	// We first check the language header, because if it's acceptable the next messages are
	// going to be returned in the language choosen by the user
	if !checkHTTPAcceptLanguage(r, &context) {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, context.Language.Messages["accept-language-error"])
		return
	}

	if !checkHTTPAccept(r) {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, context.Language.Messages["accept-error"])
		return
	}

	if !checkHTTPAcceptCharset(r) {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, context.Language.Messages["accept-charset-error"])
		return
	}

	// TODO Check:
	//   Content-MD5
	//   Content-Type
	//   Authorization

	selectedHandler(r, &context)

	if len(context.responseMessage) > 0 {
		w.Header().Add("Content-Type", "application/vnd.shelter+json")
		w.Header().Add("Content-Encoding", "utf-8")
		w.Header().Add("Content-Language", context.Language.Name())
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(context.responseMessage)))

		hash := md5.New()
		hash.Write(context.responseMessage)
		w.Header().Add("Content-MD5", fmt.Sprintf("%x", hash.Sum(nil)))
	}

	w.Header().Add("Accept", "application/vnd.shelter+json")
	w.Header().Add("Accept-Language", language.ShelterRESTLanguagePacks.Names())
	w.Header().Add("Accept-Charset", "utf-8")
	w.Header().Add("Date", time.Now().UTC().Format(time.RFC1123))

	for key, value := range context.httpHeader {
		w.Header().Add(key, value)
	}

	w.WriteHeader(context.responseHttpStatus)

	if len(context.responseMessage) > 0 {
		w.Write(context.responseMessage)
	}
}

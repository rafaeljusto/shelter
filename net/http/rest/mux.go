package rest

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"shelter/config"
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
	logger *log.Logger                   // REST server log
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

	context, err := newShelterRESTContext(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		mux.logger.Println("Error creating context. Details:", err)
		return
	}

	if !mux.checkHTTPHeaders(w, r, &context) {
		return
	}

	handler(r, &context)
	mux.writeResponse(w, context)
}

// Open REST server log file to print errors that can occur while managing the requests
func (mux *shelterRESTMux) initializeLogger() {
	restLogPath := fmt.Sprintf("%s/%s",
		config.ShelterConfig.Log.BasePath,
		config.ShelterConfig.Log.RESTFilename,
	)

	restLog, err := os.Create(restLogPath)
	if err != nil {
		log.Println(err)
		return
	}

	mux.logger = log.New(restLog, "", log.LstdFlags)
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
	r *http.Request, context *ShelterRESTContext) bool {

	// We first check the language header, because if it's acceptable the next messages are
	// going to be returned in the language choosen by the user
	if !checkHTTPAcceptLanguage(r, context) {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, context.Language.Messages["accept-language-error"])
		return false
	}

	if !checkHTTPAccept(r) {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, context.Language.Messages["accept-error"])
		return false
	}

	if !checkHTTPAcceptCharset(r) {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(w, context.Language.Messages["accept-charset-error"])
		return false
	}

	if !checkContentType(r) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, context.Language.Messages["invalid-content-type"])
		return false
	}

	if !checkHTTPContentMD5(r, context) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, context.Language.Messages["invalid-content-md5"])
		return false
	}

	// TODO Check authorization!

	return true
}

// Write response with the defaults HTTP response headers
func (mux shelterRESTMux) writeResponse(w http.ResponseWriter, context ShelterRESTContext) {
	if len(context.responseMessage) > 0 {
		w.Header().Add("Content-Type", "application/vnd.shelter+json")
		w.Header().Add("Content-Encoding", "utf-8")
		w.Header().Add("Content-Language", context.Language.Name())
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(context.responseMessage)))

		hash := md5.New()
		hash.Write(context.responseMessage)
		hashBytes := hash.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)
		w.Header().Add("Content-MD5", hashBase64)
	}

	w.Header().Add("Accept", supportedContentType)
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

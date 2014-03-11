package client

import (
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/net/http/client/handler"
	"net/http"
)

// Main router used by the Shelter web client system to manage the requests
var (
	mux Mux
)

// Mux is responsable for finding and calling a specific handler
type Mux struct {
}

// Main function of the client web server
func (mux Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// Something went really wrong while processing a request. Just send a HTTP status 500
		// to the client and log the error stacktrace
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Web client panic detected. Details:", r)
		}
	}()

	// Check if the URI exists in our system
	handler := mux.findRoute(r.URL.Path)
	if handler == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	handler(w, r)
}

// Find the best handler for the given URI. The first match is always returned
func (mux Mux) findRoute(uri string) handler.Handler {
	for route, handler := range handler.Routes {
		if route.MatchString(uri) {
			return handler
		}
	}

	return handler.StaticHandler
}

package client

import (
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/net/http/client/handler"
	"github.com/rafaeljusto/shelter/net/http/rest"
	"net"
	"net/http"
	"strings"
)

// Main router used by the Shelter web client system to manage the requests
var (
	mux Mux
)

// Mux is responsable for finding and calling a specific handler
type Mux struct {
	ACL []*net.IPNet // Network allowed ranges to send requests to the client web server
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

	// Verify if the user can send requests to this client web server
	if allowed, err := mux.checkACL(r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error checking CIDR whitelist. Details:", err)
		return

	} else if !allowed {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check if the URI exists in our system
	handler := mux.findRoute(r.URL.Path)
	if handler == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	handler(w, r)
}

// checkACL is responsable for checking if the user is allowed to send requests to the
// client web server
func (mux Mux) checkACL(r *http.Request) (bool, error) {
	// When there's nobody in the whitelist, everybody is allowed
	if len(mux.ACL) == 0 {
		return true, nil
	}

	remoteAddrParts := strings.Split(r.RemoteAddr, ":")
	if len(remoteAddrParts) != 2 {
		// Remote address without port
		return false, rest.ErrInvalidRemoteIP
	}

	ip := net.ParseIP(remoteAddrParts[0])
	if ip == nil {
		// Something wrong, because the client web server could not identify the remote
		// address properly. This is really awkward, because this is a responsability of the
		// server, maybe this error will never be throw
		return false, rest.ErrInvalidRemoteIP
	}

	for _, cidr := range mux.ACL {
		if cidr.Contains(ip) {
			return true, nil
		}
	}

	return false, nil
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

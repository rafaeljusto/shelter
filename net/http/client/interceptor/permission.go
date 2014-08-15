// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the Web request before calling the handler
package interceptor

import (
	"github.com/rafaeljusto/shelter/log"
	"github.com/trajber/handy/interceptor"
	"net"
	"net/http"
	"strings"
)

var (
	ACL []*net.IPNet
)

type Permission struct {
	interceptor.NoAfterInterceptor
}

// Permission is responsable for checking if the user is allowed to send requests to the
// Web client
func (i *Permission) Before(w http.ResponseWriter, r *http.Request) {
	// When there's nobody in the whitelist, everybody is allowed
	if len(ACL) == 0 {
		return
	}

	remoteAddrParts := strings.Split(r.RemoteAddr, ":")
	if len(remoteAddrParts) != 2 {
		// Remote address without port
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error checking CIDR whitelist. Details: Remote IP address could not be parsed")
		return
	}

	ip := net.ParseIP(remoteAddrParts[0])
	if ip == nil {
		// Something wrong, because the REST server could not identify the remote address
		// properly. This is really awkward, because this is a responsability of the server,
		// maybe this error will never be throw
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error checking CIDR whitelist. Details: Remote IP address could not be parsed")
		return
	}

	for _, cidr := range ACL {
		if cidr.Contains(ip) {
			return
		}
	}

	w.WriteHeader(http.StatusForbidden)
}

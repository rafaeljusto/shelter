// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the Web request before calling the handler
package interceptor

import (
	"github.com/rafaeljusto/handy/interceptor"
	"github.com/rafaeljusto/shelter/log"
	"net"
	"net/http"
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

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error checking CIDR whitelist. Details: Remote IP address '%s' could not be parsed. %s", r.RemoteAddr, err)
		return
	}

	ip := net.ParseIP(host)
	if ip == nil {
		// Something wrong, because the REST server could not identify the remote address
		// properly. This is really awkward, because this is a responsability of the server,
		// maybe this error will never be throw
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error checking CIDR whitelist. Details: Remote IP address '%s' could not be parsed", r.RemoteAddr)
		return
	}

	for _, cidr := range ACL {
		if cidr.Contains(ip) {
			return
		}
	}

	w.WriteHeader(http.StatusForbidden)
}

// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"github.com/rafaeljusto/handy/interceptor"
	"github.com/rafaeljusto/shelter/log"
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
// REST server
func (i *Permission) Before(w http.ResponseWriter, r *http.Request) {
	// When there's nobody in the whitelist, everybody is allowed
	if len(ACL) == 0 {
		return
	}

	var clientAddress string

	xff := r.Header.Get("X-Forwarded-For")
	xff = strings.TrimSpace(xff)

	if len(xff) > 0 {
		xffParts := strings.Split(xff, ",")
		if len(xffParts) == 1 {
			clientAddress = strings.TrimSpace(xffParts[0])
		} else if len(xffParts) > 1 {
			clientAddress = strings.TrimSpace(xffParts[len(xffParts)-2])
		}

	} else {
		clientAddress = strings.TrimSpace(r.Header.Get("X-Real-IP"))
	}

	if len(clientAddress) == 0 {
		var err error
		clientAddress, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error checking CIDR whitelist. Details: Remote IP address '%s' could not be parsed. %s", r.RemoteAddr, err)
			return
		}

	}

	ip := net.ParseIP(clientAddress)
	if ip == nil {
		// Something wrong, because the REST server could not identify the remote address properly. This
		// is really awkward, because this is a responsability of the server, maybe this error will
		// never be throw
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error checking CIDR whitelist. Details: IP address '%s' could not be parsed.", clientAddress)
		return
	}

	for _, cidr := range ACL {
		if cidr.Contains(ip) {
			return
		}
	}

	w.WriteHeader(http.StatusForbidden)
}

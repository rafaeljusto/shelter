// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"github.com/rafaeljusto/handy/interceptor"
	"net"
	"net/http"
	"strings"
)

type remoteAddresser interface {
	RemoteAddress() net.IP
	SetRemoteAddress(net.IP)
}

type RemoteAddress struct {
	interceptor.NopInterceptor
	handler remoteAddresser
}

func NewRemoteAddress(h remoteAddresser) *RemoteAddress {
	return &RemoteAddress{handler: h}
}

func (i *RemoteAddress) Before(w http.ResponseWriter, r *http.Request) {
	//log.Debug("Interceptor Before: Remote Address")

	var clientAddress string

	// http://en.wikipedia.org/wiki/X-Forwarded-For
	//
	// The general format of the field is:
	//
	//     X-Forwarded-For: client, proxy1, proxy2
	//
	// where the value is a comma+space separated list of IP addresses, the left-most being the
	// original client, and each successive proxy that passed the request adding the IP address where
	// it received the request from. In this example, the request passed through proxy1, proxy2, and
	// then proxy3 (not shown in the header). proxy3 appears as remote address of the request.
	//
	// Since it is easy to forge an X-Forwarded-For field the given information should be used with
	// care. The last IP address is always the IP address that connects to the last proxy, which means
	// it is the most reliable source of information. X-Forwarded-For data can be used in a forward or
	// reverse proxy scenario.
	xff := r.Header.Get("X-Forwarded-For")
	xff = strings.TrimSpace(xff)

	if len(xff) > 0 {
		xffParts := strings.Split(xff, ",")
		if len(xffParts) == 1 {
			// If we have only one address we will pick it as the end-user IP
			clientAddress = strings.TrimSpace(xffParts[0])
		} else if len(xffParts) > 1 {
			// If we have many addresses, we will pick second last as we assume that the client came
			// through a proxy
			clientAddress = strings.TrimSpace(xffParts[len(xffParts)-2])
		}

	} else {
		// Some proxies set the X-Real-IP to determinate the end-user address
		clientAddress = strings.TrimSpace(r.Header.Get("X-Real-IP"))
	}

	if len(clientAddress) > 0 {
		if address := net.ParseIP(clientAddress); address != nil {
			i.handler.SetRemoteAddress(net.ParseIP(clientAddress))
			return
		}
	}

	clientAddress, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		//log.Notice(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i.handler.SetRemoteAddress(net.ParseIP(clientAddress))
}

// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/normalize"
	"net"
	"time"
)

// Nameserver object used in the protocol to determinate what the user can update
type NameserverRequest struct {
	Host *string `json:"host,omitempty"` // Nameserver's name
	IPv4 *string `json:"ipv4,omitempty"` // Host's IPv4 (optional when don't need glue)
	IPv6 *string `json:"ipv6,omitempty"` // Host's IPv6 (optional)
}

func (n *NameserverRequest) Normalize() {
	if n.Host != nil {
		host, _ := normalize.NormalizeDomainName(*n.Host)
		n.Host = &host
	}
}

func (n *NameserverRequest) Validate() (Translator, error) {
	var messagesHolder MessagesHolder

	if n.Host == nil {
		messagesHolder.Add(NewMessageResponseWithField(ErrorCodeInvalidFQDN,
			"nameserver.host", "", nil))
	}

	if n.IPv4 != nil && len(*n.IPv4) > 0 {
		if ipv4 := net.ParseIP(*n.IPv4); ipv4 == nil {
			messagesHolder.Add(NewMessageResponseWithField(ErrorCodeInvalidIP,
				"nameserver.ipv4", *n.IPv4, nil))
		}
	}

	if n.IPv6 != nil && len(*n.IPv6) > 0 {
		if ipv6 := net.ParseIP(*n.IPv6); ipv6 == nil {
			messagesHolder.Add(NewMessageResponseWithField(ErrorCodeInvalidIP,
				"nameserver.ipv6", *n.IPv6, nil))
		}
	}

	return messagesHolder.Messages(), nil
}

// Namerserver object used in the protocol to determinate what the user can see. The
// status was converted to text format for easy interpretation
type NameserverResponse struct {
	Host        string    `json:"host,omitempty"`        // Nameserver's name
	IPv4        string    `json:"ipv4,omitempty"`        // Host's IPv4 (optional when don't need glue)
	IPv6        string    `json:"ipv6,omitempty"`        // Host's IPv6 (optional)
	LastStatus  string    `json:"lastStatus,omitempty"`  // Result of the last configuration check
	LastCheckAt time.Time `json:"lastCheckAt,omitempty"` // Time of the last configuration check
	LastOKAt    time.Time `json:"lastOKAt,omitempty"`    // Last time that the DNS configuration was OK
}

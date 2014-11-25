// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"github.com/rafaeljusto/shelter/protocol"
	"net"
	"strings"
	"time"
)

// List of possible nameserver status
const (
	NameserverStatusNotChecked        = iota // Nameserver not checked yet
	NameserverStatusOK                       // DNS configuration for this nameserver is OK
	NameserverStatusTimeout                  // Network timeout while trying to reach the nameserver
	NameserverStatusNoAuthority              // Nameserver does not have authority for this domain
	NameserverStatusUnknownDomainName        // Domain does not exists for this nameserver
	NameserverStatusUnknownHost              // Could not resolve nameserver (no glue)
	NameserverStatusServerFailure            // Nameserver configuration problem
	NameserverStatusQueryRefused             // DNS request rejected
	NameserverStatusConnectionRefused        // Connection refused by firewall or nameserver
	NameserverStatusCanonicalName            // Domain name is a link in the zone APEX
	NameserverStatusNotSynchronized          // Nameservers of this domain have a different version of the zone files
	NameserverStatusError                    // Generic error found in the nameserver
)

// NameserverStatus is a number that represents one of the possible nameserver status
// listed in the constant group above
type NameserverStatus int

// Convert the nameserver status enum to text for printing in reports or debugging
func NameserverStatusToString(status NameserverStatus) string {
	switch status {
	case NameserverStatusNotChecked:
		return "NOTCHECKED"
	case NameserverStatusOK:
		return "OK"
	case NameserverStatusTimeout:
		return "TIMEOUT"
	case NameserverStatusNoAuthority:
		return "NOAA"
	case NameserverStatusUnknownDomainName:
		return "UDN"
	case NameserverStatusUnknownHost:
		return "UH"
	case NameserverStatusServerFailure:
		return "SERVFAIL"
	case NameserverStatusQueryRefused:
		return "QREFUSED"
	case NameserverStatusConnectionRefused:
		return "CREFUSED"
	case NameserverStatusCanonicalName:
		return "CNAME"
	case NameserverStatusNotSynchronized:
		return "NOTSYNCH"
	case NameserverStatusError:
		return "ERROR"
	}

	return ""
}

// Nameserver store the information necessary to send the requests for a specific host and
// store the results of this requests
type Nameserver struct {
	Host        string           // Nameserver's name
	IPv4        net.IP           // Host's IPv4 (optional when don't need glue)
	IPv6        net.IP           // Host's IPv6 (optional)
	LastStatus  NameserverStatus // Result of the last configuration check
	LastCheckAt time.Time        // Time of the last configuration check
	LastOKAt    time.Time        // Last time that the DNS configuration was OK
}

// Method to check if the nameserver needs glue for a given domain name. A namerserver
// needs glue when the name of the domain is inside the nameserver (example: domain
// test.com.br and nameserver ns1.tes.com.br)
func (n Nameserver) NeedsGlue(fqdn string) bool {
	return strings.HasSuffix(n.Host, fqdn) ||
		strings.HasSuffix(n.Host, fqdn+".") ||
		strings.HasSuffix(n.Host+".", fqdn)
}

// ChangeStatus is a easy way to change the status of a nameserver because it also updates
// the last check date
func (n *Nameserver) ChangeStatus(status NameserverStatus) {
	n.LastStatus = status
	n.LastCheckAt = time.Now()

	if status == NameserverStatusOK {
		n.LastOKAt = n.LastCheckAt
	}
}

// Convert a nameserver request object into a nameserver model object. It can return
// errors related to the conversion of IP addresses and normalization of nameserver's
// hostname
func (n *Nameserver) Apply(nameserverRequest protocol.NameserverRequest) bool {
	if nameserverRequest.Host != nil {
		n.Host = *nameserverRequest.Host
	}

	if nameserverRequest.IPv4 != nil && len(*nameserverRequest.IPv4) > 0 {
		ipv4 := net.ParseIP(*nameserverRequest.IPv4)
		if ipv4 == nil {
			return false
		}
		n.IPv4 = ipv4
	}

	if nameserverRequest.IPv6 != nil && len(*nameserverRequest.IPv6) > 0 {
		ipv6 := net.ParseIP(*nameserverRequest.IPv6)
		if ipv6 == nil {
			return false
		}
		n.IPv6 = ipv6
	}

	return true
}

// Convert a nameserver of the system into a format with limited information to return it
// to the user
func (n *Nameserver) Protocol() protocol.NameserverResponse {
	ipv4 := ""
	if n.IPv4 != nil {
		ipv4 = n.IPv4.String()
	}

	ipv6 := ""
	if n.IPv6 != nil {
		ipv6 = n.IPv6.String()
	}

	return protocol.NameserverResponse{
		Host:        n.Host,
		IPv4:        ipv4,
		IPv6:        ipv6,
		LastStatus:  NameserverStatusToString(n.LastStatus),
		LastCheckAt: n.LastCheckAt,
		LastOKAt:    n.LastOKAt,
	}
}

type Nameservers []Nameserver

// Convert a list of nameserver requests objects into a list of nameserver model objects.
// Useful when merging domain object from the network with a domain object from the
// database. It can return errors related to the conversion of IP addresses and
// normalization of nameserver's hostname
func (n Nameservers) Apply(nameserversRequest []protocol.NameserverRequest) (Nameservers, bool) {
	for _, nameserverRequest := range nameserversRequest {
		if nameserverRequest.Host == nil {
			return n, false
		}

		found := false
		for i, nameserver := range n {
			if nameserver.Host == *nameserverRequest.Host {
				if !nameserver.Apply(nameserverRequest) {
					return n, false
				}

				n[i] = nameserver
				found = true
				break
			}
		}

		if !found {
			var nameserver Nameserver
			if !nameserver.Apply(nameserverRequest) {
				return n, false
			}

			n = append(n, nameserver)
		}
	}

	return n, true
}

// Convert a list of nameservers of the system into a format with limited information to
// return it to the user. This is only a easy way to call toNameserverResponse for each
// object in the list
func (n Nameservers) Protocol() []protocol.NameserverResponse {
	var nameserversResponse []protocol.NameserverResponse
	for _, nameserver := range n {
		nameserversResponse = append(nameserversResponse, nameserver.Protocol())
	}
	return nameserversResponse
}

// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package dnsutils contains useful features to manage DNS structures
package dnsutils

import (
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/miekg/dns"
)

// Useful function to retrieve all records of a specific type from the DNS response
// message. We assume that the resource record is an instance of the specific type based
// on the Rrtype attribute
func FilterRRs(rrs []dns.RR, rrType uint16) []dns.RR {
	var filtered []dns.RR
	for _, rr := range rrs {
		if rr.Header().Rrtype == rrType {
			filtered = append(filtered, rr)
		}
	}
	return filtered
}

// Useful function to return the first occurence of a resource record of a specific type.
// This method is faster than filterRRs when we are interested in only one record (like
// SOA)
func FilterFirstRR(rrs []dns.RR, rrType uint16) dns.RR {
	for _, rr := range rrs {
		if rr.Header().Rrtype == rrType {
			return rr
		}
	}
	return nil
}

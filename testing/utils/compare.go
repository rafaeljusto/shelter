// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// utils add features to make the test life easier
package utils

import (
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/protocol"
)

// Function to compare if two domains are equal, cannot use operator == because of the
// slices inside the domain object
func CompareDomain(d1, d2 model.Domain) bool {
	if d1.Id != d2.Id || d1.FQDN != d2.FQDN {
		return false
	}

	if len(d1.Nameservers) != len(d2.Nameservers) {
		return false
	}

	for i := 0; i < len(d1.Nameservers); i++ {
		// Cannot compare the nameservers directly with operator == because of the
		// pointers for IP addresses and dates
		if d1.Nameservers[i].Host != d2.Nameservers[i].Host ||
			d1.Nameservers[i].IPv4.String() != d2.Nameservers[i].IPv4.String() ||
			d1.Nameservers[i].IPv6.String() != d2.Nameservers[i].IPv6.String() ||
			d1.Nameservers[i].LastStatus != d2.Nameservers[i].LastStatus ||
			d1.Nameservers[i].LastCheckAt.Unix() != d2.Nameservers[i].LastCheckAt.Unix() ||
			d1.Nameservers[i].LastOKAt.Unix() != d2.Nameservers[i].LastOKAt.Unix() {
			return false
		}
	}

	if len(d1.DSSet) != len(d2.DSSet) {
		return false
	}

	for i := 0; i < len(d1.DSSet); i++ {
		// Cannot compare the nameservers directly with operator == because of the dates
		if d1.DSSet[i].Algorithm != d2.DSSet[i].Algorithm ||
			d1.DSSet[i].Digest != d2.DSSet[i].Digest ||
			d1.DSSet[i].DigestType != d2.DSSet[i].DigestType ||
			d1.DSSet[i].ExpiresAt.Unix() != d2.DSSet[i].ExpiresAt.Unix() ||
			d1.DSSet[i].Keytag != d2.DSSet[i].Keytag ||
			d1.DSSet[i].LastCheckAt.Unix() != d2.DSSet[i].LastCheckAt.Unix() ||
			d1.DSSet[i].LastOKAt.Unix() != d2.DSSet[i].LastOKAt.Unix() ||
			d1.DSSet[i].LastStatus != d2.DSSet[i].LastStatus {
			return false
		}
	}

	if len(d1.Owners) != len(d2.Owners) {
		return false
	}

	for i := 0; i < len(d1.Owners); i++ {
		if d1.Owners[i].Email.String() != d2.Owners[i].Email.String() ||
			d1.Owners[i].Language != d2.Owners[i].Language {
			return false
		}
	}

	return true
}

// Function to compare if two domains are equal, cannot use operator == because of the
// slices inside the domain object
func CompareProtocolDomain(d1, d2 protocol.DomainResponse) bool {
	if len(d1.Nameservers) != len(d2.Nameservers) {
		return false
	}

	for i := 0; i < len(d1.Nameservers); i++ {
		// Cannot compare the nameservers directly with operator == because of the
		// pointers for IP addresses and dates
		if d1.Nameservers[i].Host != d2.Nameservers[i].Host ||
			d1.Nameservers[i].IPv4 != d2.Nameservers[i].IPv4 ||
			d1.Nameservers[i].IPv6 != d2.Nameservers[i].IPv6 ||
			d1.Nameservers[i].LastStatus != d2.Nameservers[i].LastStatus ||
			d1.Nameservers[i].LastCheckAt.Unix() != d2.Nameservers[i].LastCheckAt.Unix() ||
			d1.Nameservers[i].LastOKAt.Unix() != d2.Nameservers[i].LastOKAt.Unix() {
			return false
		}
	}

	if len(d1.DSSet) != len(d2.DSSet) {
		return false
	}

	for i := 0; i < len(d1.DSSet); i++ {
		// Cannot compare the nameservers directly with operator == because of the dates
		if d1.DSSet[i].Algorithm != d2.DSSet[i].Algorithm ||
			d1.DSSet[i].Digest != d2.DSSet[i].Digest ||
			d1.DSSet[i].DigestType != d2.DSSet[i].DigestType ||
			d1.DSSet[i].ExpiresAt.Unix() != d2.DSSet[i].ExpiresAt.Unix() ||
			d1.DSSet[i].Keytag != d2.DSSet[i].Keytag ||
			d1.DSSet[i].LastCheckAt.Unix() != d2.DSSet[i].LastCheckAt.Unix() ||
			d1.DSSet[i].LastOKAt.Unix() != d2.DSSet[i].LastOKAt.Unix() ||
			d1.DSSet[i].LastStatus != d2.DSSet[i].LastStatus {
			return false
		}
	}

	if len(d1.Owners) != len(d2.Owners) {
		return false
	}

	for i := 0; i < len(d1.Owners); i++ {
		if d1.Owners[i].Email != d2.Owners[i].Email ||
			d1.Owners[i].Language != d2.Owners[i].Language {
			return false
		}
	}

	return true
}

// Function to compare if two domain lists are equal, cannot use operator == because of the
// slices inside the domain object
func CompareProtocolDomains(d1, d2 protocol.DomainsResponse) bool {
	if d1.NumberOfItems != d2.NumberOfItems ||
		d1.NumberOfPages != d2.NumberOfPages ||
		d1.Page != d2.Page ||
		d1.PageSize != d2.PageSize ||
		len(d1.Domains) != len(d2.Domains) ||
		len(d1.Links) != len(d2.Links) {

		return false
	}

	for index := range d1.Domains {
		if !CompareProtocolDomain(d1.Domains[index], d2.Domains[index]) {
			return false
		}
	}

	if !CompareProtocolLinks(d1.Links, d2.Links) {
		return false
	}

	return true
}

func CompareScan(s1, s2 model.Scan) bool {
	if s1.Id != s2.Id ||
		s1.Revision != s2.Revision ||
		s1.Status != s2.Status ||
		s1.StartedAt.Unix() != s2.StartedAt.Unix() ||
		s1.FinishedAt.Unix() != s2.FinishedAt.Unix() ||
		s1.DomainsScanned != s2.DomainsScanned ||
		s1.DomainsWithDNSSECScanned != s2.DomainsWithDNSSECScanned {
		return false
	}

	for key, value := range s1.NameserverStatistics {
		if otherValue, ok := s2.NameserverStatistics[key]; !ok || value != otherValue {
			return false
		}
	}

	for key, value := range s1.DSStatistics {
		if otherValue, ok := s2.DSStatistics[key]; !ok || value != otherValue {
			return false
		}
	}

	return true
}

func CompareProtocolScan(s1, s2 protocol.ScanResponse) bool {
	if s1.Status != s2.Status ||
		s1.StartedAt.Unix() != s2.StartedAt.Unix() ||
		s1.FinishedAt.Unix() != s2.FinishedAt.Unix() ||
		s1.DomainsScanned != s2.DomainsScanned ||
		s1.DomainsWithDNSSECScanned != s2.DomainsWithDNSSECScanned {

		return false
	}

	for key, value := range s1.NameserverStatistics {
		if otherValue, ok := s2.NameserverStatistics[key]; !ok || value != otherValue {
			return false
		}
	}

	for key, value := range s1.DSStatistics {
		if otherValue, ok := s2.DSStatistics[key]; !ok || value != otherValue {
			return false
		}
	}

	return true
}

func CompareProtocolScans(s1, s2 protocol.ScansResponse) bool {
	if s1.NumberOfItems != s2.NumberOfItems ||
		s1.NumberOfPages != s2.NumberOfPages ||
		s1.Page != s2.Page ||
		s1.PageSize != s2.PageSize ||
		len(s1.Scans) != len(s2.Scans) ||
		len(s1.Links) != len(s2.Links) {

		return false
	}

	for index := range s1.Scans {
		if !CompareProtocolScan(s1.Scans[index], s2.Scans[index]) {
			return false
		}
	}

	if !CompareProtocolLinks(s1.Links, s2.Links) {
		return false
	}

	return true
}

func CompareProtocolLinks(l1 []protocol.Link, l2 []protocol.Link) bool {
	if len(l1) != len(l2) {
		return false
	}

	for i := range l1 {
		if l1[i].HRef != l2[i].HRef || len(l1[i].Types) != len(l2[i].Types) {
			return false
		}

		for j := range l1[i].Types {
			if l1[i].Types[j] != l2[i].Types[j] {
				return false
			}
		}
	}

	return true
}

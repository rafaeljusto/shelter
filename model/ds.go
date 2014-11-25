// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"github.com/miekg/dns"
	"github.com/rafaeljusto/shelter/protocol"
	"time"
)

// List of possible DS status
const (
	DSStatusNotChecked       = iota // DS record not checked yet
	DSStatusOK                      // DNSSEC configuration for this DS is OK
	DSStatusTimeout                 // Network timeout while trying to retrieve the DNSKEY
	DSStatusNoSignature             // No RRSIG records found for the related DNSKEY
	DSStatusExpiredSignature        // At least one RRSIG record was expired
	DSStatusNoKey                   // No DNSKEY was found with the keytag of the DS
	DSStatusNoSEP                   // DNSKEY related to DS does not have the bit SEP on
	DSStatusSignatureError          // Error while checking DNSKEY signatures
	DSStatusDNSError                // DNS error (check nameserver status)
)

// DSStatus is a number that represents one of the possible DS status listed in the
// constant group above
type DSStatus int

// Convert the DS status enum to text for printing in reports or debugging
func DSStatusToString(status DSStatus) string {
	switch status {
	case DSStatusNotChecked:
		return "NOTCHECKED"
	case DSStatusOK:
		return "OK"
	case DSStatusTimeout:
		return "TIMEOUT"
	case DSStatusNoSignature:
		return "NOSIG"
	case DSStatusExpiredSignature:
		return "EXPSIG"
	case DSStatusNoKey:
		return "NOKEY"
	case DSStatusNoSEP:
		return "NOSEP"
	case DSStatusSignatureError:
		return "SIGERR"
	case DSStatusDNSError:
		return "DNSERR"
	}

	return ""
}

// When converting a DNSKEY into a DS we need to choose wich digest type are we going to
// use, as we don't want to bother the user asking this information we assume a default
// digest type
const (
	DefaultDigestType = protocol.DSDigestTypeSHA256
)

// DS store the information necessary to validate if a domain is configured correctly with
// DNSSEC, and it also stores the results of the validations. When the hosts have multiple
// DNSSEC problems, the worst problem (using a priority algorithm) will be stored in the
// DS
type DS struct {
	Keytag      uint16                // DNSKEY's identification number
	Algorithm   protocol.DSAlgorithm  // DNSKEY's algorithm
	Digest      string                // Hash of the DNSKEY content
	DigestType  protocol.DSDigestType // Hash type decided by user when generating the DS
	ExpiresAt   time.Time             // DNSKEY's signature expiration date
	LastStatus  DSStatus              // Result of the last configuration check
	LastCheckAt time.Time             // Time of the last configuration check
	LastOKAt    time.Time             // Last time that the DNSSEC configuration was OK
}

// ChangeStatus is a easy way to change the status of a DS because it also updates the
// last check date
func (d *DS) ChangeStatus(status DSStatus) {
	d.LastStatus = status
	d.LastCheckAt = time.Now()

	if status == DSStatusOK {
		d.LastOKAt = d.LastCheckAt
	}
}

// ApplyDNSKEY convert a DNSKEY request object into a DS model object. The DNSKEY request cannot
// contain any null fields, even if its a HTTP patch action, because we always need them to convert
// to a valid DS. On any problem false will be returned.
func (d *DS) ApplyDNSKEY(dnskeyRequest protocol.DNSKEYRequest) bool {
	if dnskeyRequest.Algorithm == nil ||
		dnskeyRequest.Flags == nil ||
		dnskeyRequest.PublicKey == nil {

		return false
	}

	rawDNSKEY := dns.DNSKEY{
		Hdr: dns.RR_Header{
			Name: dnskeyRequest.FQDN,
		},
		Flags:     *dnskeyRequest.Flags,
		Protocol:  uint8(3),
		Algorithm: *dnskeyRequest.Algorithm,
		PublicKey: *dnskeyRequest.PublicKey,
	}

	rawDS := rawDNSKEY.ToDS(int(DefaultDigestType))
	if rawDS == nil {
		return false
	}

	dsRequest := protocol.DSRequest{
		Keytag:     &rawDS.KeyTag,
		Algorithm:  &rawDS.Algorithm,
		Digest:     &rawDS.Digest,
		DigestType: &rawDS.DigestType,
	}

	return d.Apply(dsRequest)
}

// Convert a DS request object into a DS model object. It can return errors related to the
// conversion of the algorithm, when it is out of range
func (d *DS) Apply(dsRequest protocol.DSRequest) bool {
	if dsRequest.Keytag == nil {
		return false
	}

	d.Keytag = *dsRequest.Keytag

	if dsRequest.Digest != nil {
		d.Digest = *dsRequest.Digest
	}

	if dsRequest.Algorithm != nil {
		d.Algorithm = protocol.DSAlgorithm(*dsRequest.Algorithm)
	}

	if dsRequest.DigestType != nil {
		d.DigestType = protocol.DSDigestType(*dsRequest.DigestType)
	}

	return true
}

// Convert a DS of the system into a format with limited information to return it to the
// user
func (d *DS) Protocol() protocol.DSResponse {
	return protocol.DSResponse{
		Keytag:      d.Keytag,
		Algorithm:   uint8(d.Algorithm),
		Digest:      d.Digest,
		DigestType:  uint8(d.DigestType),
		ExpiresAt:   d.ExpiresAt,
		LastStatus:  DSStatusToString(d.LastStatus),
		LastCheckAt: d.LastCheckAt,
		LastOKAt:    d.LastOKAt,
	}
}

type DSSet []DS

// Convert a list of DS requests objects into a list of DS model objects. Useful when
// merging domain object from the network with a domain object from the database. It can
// return errors related to the conversion of the algorithm, when it is out of range
func (d DSSet) Apply(dssetRequest []protocol.DSRequest) (DSSet, bool) {
	for _, dsRequest := range dssetRequest {
		if dsRequest.Keytag == nil {
			return d, false
		}

		found := false
		for i, ds := range d {
			if ds.Keytag == *dsRequest.Keytag {
				if !ds.Apply(dsRequest) {
					return d, false
				}

				d[i] = ds
				found = true
				break
			}
		}

		if !found {
			var ds DS
			if !ds.Apply(dsRequest) {
				return d, false
			}

			d = append(d, ds)
		}
	}

	return d, true
}

// ApplyDNSKEYs convert a list of DNSKEY requests objects into a list of DS model objects. Useful
// when merging domain object from the network with a domain object from the database
func (d DSSet) ApplyDNSKEYs(dnskeysRequest []protocol.DNSKEYRequest) (DSSet, bool) {
	for _, dnskeyRequest := range dnskeysRequest {
		var newDS DS
		if !newDS.ApplyDNSKEY(dnskeyRequest) {
			return d, false
		}

		found := false
		for i, ds := range d {
			if ds.Keytag == newDS.Keytag {
				d[i] = newDS
				found = true
				break
			}
		}

		if !found {
			d = append(d, newDS)
		}
	}

	return d, true
}

// Convert a list of DS of the system into a format with limited information to return it
// to the user. This is only a easy way to call toDSResponse for each object in the list
func (d DSSet) Protocol() []protocol.DSResponse {
	var dsSetResponse []protocol.DSResponse
	for _, ds := range d {
		dsSetResponse = append(dsSetResponse, ds.Protocol())
	}
	return dsSetResponse
}

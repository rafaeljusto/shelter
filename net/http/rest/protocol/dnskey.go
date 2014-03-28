// protocol - REST protocol description
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package protocol

import (
	"errors"
	"github.com/miekg/dns"
	"github.com/rafaeljusto/shelter/model"
	"strings"
)

// When converting a DNSKEY into a DS we need to choose wich digest type are we going to
// use, as we don't want to bother the user asking this information we assume a default
// digest type
const (
	DefaultDigestType = model.DSDigestTypeSHA256
)

// List of possible errors that can occur when calling methods from this object
var (
	// Error returned when trying to convert a DNSKEY into a DS
	ErrInvalidDNSKEY = errors.New("DNSKEY data is invalid to generate DS")
)

// DNSKEY object from the protocol used to simplify the user lifes. The system will
// convert automatically DNSKEY objects into DS objects
type DNSKEYRequest struct {
	Flags     uint16 `json:"flags,omitempty"`     // RFC 4034 - 2.1.1 (SEP, ZONE, ...)
	Algorithm uint8  `json:"algorithm,omitempty"` // RFC 4034 - Appendix A.1
	PublicKey string `json:"publicKey,omitempty"` // Base64 enconded string of the public key
}

// Convert a DNSKEY request object into a DS model object. It can return errors related to
// the generation of the DS objects
func (d *DNSKEYRequest) toDSModel(fqdn string) (model.DS, error) {
	// The base64 decode method don't deal very well with spaces inside the public key raw
	// data. So we replace it before calculating the KeyTag
	publicKey := strings.Replace(d.PublicKey, " ", "", -1)

	rawDNSKEY := dns.DNSKEY{
		Hdr: dns.RR_Header{
			Name: fqdn,
		},
		Flags:     d.Flags,
		Protocol:  uint8(3),
		Algorithm: d.Algorithm,
		PublicKey: publicKey,
	}

	rawDS := rawDNSKEY.ToDS(int(DefaultDigestType))
	if rawDS == nil {
		return model.DS{}, ErrInvalidDNSKEY
	}

	dsRequest := DSRequest{
		Keytag:     rawDS.KeyTag,
		Algorithm:  rawDS.Algorithm,
		Digest:     rawDS.Digest,
		DigestType: rawDS.DigestType,
	}

	return dsRequest.toDSModel()
}

// Convert a list of DNSKEY requests objects into a list of DS model objects. Useful when
// merging domain object from the network with a domain object from the database. It can
// return errors related to the generation of the DS objects. The name of this method is
// not so good, because we already have the ToDSSetModel
func dnskeysRequestsToDSSetModel(fqdn string, dnskeysRequests []DNSKEYRequest) ([]model.DS, error) {
	var dsSet []model.DS
	for _, dnskeyRequest := range dnskeysRequests {
		ds, err := dnskeyRequest.toDSModel(fqdn)
		if err != nil {
			return nil, err
		}

		dsSet = append(dsSet, ds)
	}

	return dsSet, nil
}

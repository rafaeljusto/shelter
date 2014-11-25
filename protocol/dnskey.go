// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"strconv"
	"strings"
)

// DNSKEY object from the protocol used to simplify the user lifes. The system will
// convert automatically DNSKEY objects into DS objects
type DNSKEYRequest struct {
	FQDN      string  // Used to transform the DNSKEY into a DS (necessary for keytag calculation)
	Flags     *uint16 `json:"flags,omitempty"`     // RFC 4034 - 2.1.1 (SEP, ZONE, ...)
	Algorithm *uint8  `json:"algorithm,omitempty"` // RFC 4034 - Appendix A.1
	PublicKey *string `json:"publicKey,omitempty"` // Base64 enconded string of the public key
}

func (d *DNSKEYRequest) Normalize() {
	if d.PublicKey != nil {
		// The base64 decode method don't deal very well with spaces inside the public key raw
		// data. So we replace it before calculating the KeyTag
		k := strings.Replace(*d.PublicKey, " ", "", -1)
		d.PublicKey = &k
	}
}

func (d *DNSKEYRequest) Validate() (Translator, error) {
	var messagesHolder MessagesHolder

	if d.Algorithm == nil {
		messagesHolder.Add(NewMessageResponseWithField(ErrorCodeInvalidDNSKEY,
			"dnskey.algorithm", "", nil))

	} else if !IsValidDSAlgorithm(*d.Algorithm) {
		messagesHolder.Add(NewMessageResponseWithField(ErrorCodeInvalidDSAlgorithm,
			"dnskey.algorithm", strconv.Itoa(int(*d.Algorithm)), nil))
	}

	if d.Flags == nil {
		messagesHolder.Add(NewMessageResponseWithField(ErrorCodeInvalidDNSKEY,
			"dnskey.flags", "", nil))
	}

	if d.PublicKey == nil {
		messagesHolder.Add(NewMessageResponseWithField(ErrorCodeInvalidDNSKEY,
			"dnskey.publicKey", "", nil))
	}

	return messagesHolder.Messages(), nil
}

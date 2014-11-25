// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/normalize"
)

// Domain object from the protocol used to determinate what the user can update
type DomainRequest struct {
	FQDN        *string             `json:"-"`                     // Actual domain name
	Nameservers []NameserverRequest `json:"nameservers,omitempty"` // Nameservers that asnwer with authority for this domain
	DSSet       []DSRequest         `json:"dsset,omitempty"`       // Records for the DNS tree chain of trust
	DNSKEYS     []DNSKEYRequest     `json:"dnskeys,omitempty"`     // Records that can be converted into DS records
	Owners      []OwnerRequest      `json:"owners,omitempty"`      // E-mails that will be alerted on any problem
}

func (d *DomainRequest) Normalize() {
	if d.FQDN != nil {
		fqdn, _ := normalize.NormalizeDomainName(*d.FQDN)
		d.FQDN = &fqdn
	}
}

func (d *DomainRequest) Validate() (Translator, error) {
	var messagesHolder MessagesHolder

	if d.FQDN == nil {
		messagesHolder.Add(NewMessageResponseWithField(ErrorCodeInvalidFQDN,
			"domain.fqdn", "", nil))
	}

	return messagesHolder.Messages(), nil
}

// Domain object from the protocol used to determinate what the user can see. The last
// modified field is not here because it is sent in HTTP header field as it is with
// revision (ETag)
type DomainResponse struct {
	FQDN        string               `json:"fqdn"`                  // Actual domain name
	Nameservers []NameserverResponse `json:"nameservers,omitempty"` // Nameservers that asnwer with authority for this domain
	DSSet       []DSResponse         `json:"dsset,omitempty"`       // Records for the DNS tree chain of trust
	Owners      []OwnerResponse      `json:"owners,omitempty"`      // E-mails that will be alerted on any problem
	Links       []Link               `json:"links,omitempty"`       // Links to manipulate object
}

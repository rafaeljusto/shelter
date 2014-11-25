// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"github.com/rafaeljusto/shelter/protocol"
	"net/mail"
	"strings"
)

// Owner represents the responsable for the domain that can be alerted if any
// configuration problem is detected
type Owner struct {
	Email    *mail.Address // E-mail that will be alerted on any problem
	Language string        // Language used to send alerts
}

// Convert a owner request object into a owner model object. It can return errors related
// to the e-mail format
func (o *Owner) Apply(ownerRequest protocol.OwnerRequest) bool {
	if ownerRequest.Email != nil {
		email, err := mail.ParseAddress(*ownerRequest.Email)
		if err != nil {
			return false
		}

		o.Email = email
	}

	if ownerRequest.Language != nil {
		o.Language = *ownerRequest.Language
	}

	return true
}

// Convert a owner of the system into a format with limited information to return it to
// the user. For now we are not limiting any information
func (o *Owner) Protocol() protocol.OwnerResponse {
	// E-mail to string conversion formats the address as a valid RFC 5322 address. If the
	// address's name contains non-ASCII characters the name will be rendered according to
	// RFC 2047. We are going to remove the "<" and ">" from the e-mail address for better
	// look
	email := o.Email.String()
	email = strings.TrimLeft(email, "<")
	email = strings.TrimRight(email, ">")

	return protocol.OwnerResponse{
		Email:    email,
		Language: o.Language,
	}
}

type Owners []Owner

// Convert a list of owner requests objects into a list of owner model objects. Useful
// when merging domain object from the network with a domain object from the database. It
// can return errors related to the e-mail format in one of the converted owners
func (o Owners) Apply(ownersRequest []protocol.OwnerRequest) (Owners, bool) {
	for _, ownerRequest := range ownersRequest {
		var owner Owner
		if !owner.Apply(ownerRequest) {
			return o, false
		}

		o = append(o, owner)
	}

	return o, true
}

// Convert a list of owners of the system into a format with limited information to return
// it to the user. This is only a easy way to call toOwnerResponse for each object in the
// list
func (o Owners) Protocol() []protocol.OwnerResponse {
	var ownersResponse []protocol.OwnerResponse
	for _, owner := range o {
		ownersResponse = append(ownersResponse, owner.Protocol())
	}
	return ownersResponse
}

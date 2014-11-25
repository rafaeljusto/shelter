// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/language"
	"github.com/rafaeljusto/shelter/normalize"
	"net/mail"
)

// Owner object used in the protocol to determinate what the user can update, for this
// case, everything
type OwnerRequest struct {
	Email    *string `json:"email,omitempty"`    // E-mail that the owner wants to be alerted
	Language *string `json:"language,omitempty"` // Language that the owner wants to receive the messages
}

func (o *OwnerRequest) Normalize() {
	if o.Language != nil {
		l := normalize.NormalizeLanguage(*o.Language)
		o.Language = &l
	}
}

func (o *OwnerRequest) Validate() (Translator, error) {
	var messagesHolder MessagesHolder

	if o.Email != nil {
		if _, err := mail.ParseAddress(*o.Email); err != nil {
			messagesHolder.Add(NewMessageResponseWithField(ErrorCodeInvalidEmail,
				"owner.email", "", nil))
		}
	}

	if o.Language != nil && !language.IsValidLanguage(*o.Language) {
		messagesHolder.Add(NewMessageResponseWithField(ErrorCodeInvalidLanguage,
			"owner.language", "", nil))
	}

	return messagesHolder.Messages(), nil
}

// Owner object used in the protocol to determinate what the user can see
type OwnerResponse struct {
	Email    string `json:"email,omitempty"`    // E-mail that the owner wants to be alerted
	Language string `json:"language,omitempty"` // Language that the owner wants to receive the messages
}

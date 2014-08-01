// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/messages"
	"net/url"
)

// MessageResponse struct was created to return a message to the user with more
// information to easy integrate and solve problems
type MessageResponse struct {
	Id      string `json:"id,omitempty"`      // Code for integration systems to automatically solve the problem
	Message string `json:"message,omitempty"` // Message in the user's desired language
	Links   []Link `json:"links,omitempty"`   // Links associating this message with other resources
}

func NewMessageResponse(
	id string, roid string,
	language *messages.LanguagePack,
) (*MessageResponse, error) {

	messageResponse := &MessageResponse{
		Id: id,
	}

	if language != nil && language.Messages != nil {
		messageResponse.Message = language.Messages[id]
	} else {
		// When we don't find the message file, at least add the message id to make it easy to
		// identify the error. Useful on test scenarios
		messageResponse.Message = id
	}

	// Add the object related to the message according to RFC 4287
	if len(roid) != 0 {
		// roid must be an URI to be a valid link
		uri, err := url.Parse(roid)
		if err != nil {
			return nil, err
		}

		messageResponse.Links = []Link{
			{
				Types: []LinkType{LinkTypeRelated},
				HRef:  uri.RequestURI(),
			},
		}
	}

	return messageResponse, nil
}

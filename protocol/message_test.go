// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/messages"
	"testing"
)

func TestNewMessageResponse(t *testing.T) {
	messages.ShelterRESTLanguagePacks = messages.LanguagePacks{
		Default: "en-US",
		Packs: []messages.LanguagePack{
			{
				GenericName:  "en",
				SpecificName: "en-US",
				Messages: map[string]string{
					"message-id": "test message",
				},
			},
		},
	}
	messages.ShelterRESTLanguagePack = &messages.ShelterRESTLanguagePacks.Packs[0]

	messageResponse, err := NewMessageResponse("message-id", "/test/1", messages.ShelterRESTLanguagePack)
	if err != nil {
		t.Fatal("Should create the message response. Details:", err)
	}

	if messageResponse.Message != "test message" {
		t.Error("Did not retrieve the correct message")
	}

	messageResponse, err = NewMessageResponse("unknown-message-id", "/test/1", nil)
	if err != nil {
		t.Fatal("Should create the message response. Details:", err)
	}

	if messageResponse.Message != "unknown-message-id" {
		t.Errorf("Did not retrieve the correct message. Expected %s and got %s",
			"unknown-message-id", messageResponse.Message)
	}

	messageResponse, err = NewMessageResponse("message-id", ":", messages.ShelterRESTLanguagePack)
	if err == nil {
		t.Fatal("Accepting an invalid URI as ROID")
	}
}

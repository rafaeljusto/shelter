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

func TestNewMessageResponseWithFields(t *testing.T) {
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

	messageResponse, err := NewMessageResponseWithField("message-id", "field1", "value1", "/test/1",
		messages.ShelterRESTLanguagePack)
	if err != nil {
		t.Fatal("Should create the message response with fields. Details:", err)
	}

	if messageResponse.Message != "test message" {
		t.Error("Did not retrieve the correct message")
	}

	if messageResponse.Field != "field1" {
		t.Error("Did not set the correct field")
	}

	if messageResponse.Value != "value1" {
		t.Error("Did not set the correct field's value")
	}

	messageResponse, err = NewMessageResponseWithField("message-id", "field1", "value1", ":",
		messages.ShelterRESTLanguagePack)
	if err == nil {
		t.Fatal("Accepting an invalid URI as ROID")
	}
}

func TestMessageResponseTranslate(t *testing.T) {
	data := []struct {
		description         string
		message             MessageResponse
		language            messages.LanguagePack
		expectedResult      bool
		expectedTranslation string
	}{
		{
			description: "Testing translation message when language is not defined",
			message: MessageResponse{
				Id: "message-XX",
			},
			language: messages.LanguagePack{
				GenericName:  "en",
				SpecificName: "en-US",
				Messages:     nil,
			},
			expectedResult:      false,
			expectedTranslation: "message-XX",
		},
		{
			description: "Testing translation message not found",
			message: MessageResponse{
				Id: "message-XX",
			},
			language: messages.LanguagePack{
				GenericName:  "en",
				SpecificName: "en-US",
				Messages: map[string]string{
					"message-id": "test message",
				},
			},
			expectedResult:      false,
			expectedTranslation: "",
		},
		{
			description: "Testing successful message translation",
			message: MessageResponse{
				Id: "message-id",
			},
			language: messages.LanguagePack{
				GenericName:  "en",
				SpecificName: "en-US",
				Messages: map[string]string{
					"message-id": "test message",
				},
			},
			expectedResult:      true,
			expectedTranslation: "test message",
		},
	}

	for i, item := range data {
		ok := item.message.Translate(&item.language)

		if !ok && item.expectedResult {
			t.Errorf("Item %d, “%s”: mismatch results. Expecting true", i, item.description)

		} else if ok && !item.expectedResult {
			t.Errorf("Item %d, “%s”: mismatch results. Expecting false", i, item.description)
		}

		if item.message.Message != item.expectedTranslation {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%s'; found '%s'",
				i,
				item.description,
				item.expectedTranslation,
				item.message.Message,
			)
		}
	}
}

func TestMessagesTranslate(t *testing.T) {
	data := []struct {
		description          string
		messages             Messages
		language             messages.LanguagePack
		expectedResult       bool
		expectedTranslations []string
	}{
		{
			description: "Testing translations when language is not defined",
			messages: Messages{
				&MessageResponse{Id: "message-XX"},
				&MessageResponse{Id: "message-id1"},
			},
			language: messages.LanguagePack{
				GenericName:  "en",
				SpecificName: "en-US",
				Messages:     nil,
			},
			expectedResult: false,
			expectedTranslations: []string{
				"message-XX",
				"message-id1",
			},
		},
		{
			description: "Testing translations for unknown ids",
			messages: Messages{
				&MessageResponse{Id: "message-XX"},
				&MessageResponse{Id: "message-id1"},
			},
			language: messages.LanguagePack{
				GenericName:  "en",
				SpecificName: "en-US",
				Messages: map[string]string{
					"message-id1": "test message 1",
					"message-id2": "test message 2",
				},
			},
			expectedResult: false,
			expectedTranslations: []string{
				"",
				"test message 1",
			},
		},
		{
			description: "Testing successful messages translations",
			messages: Messages{
				&MessageResponse{Id: "message-id1"},
				&MessageResponse{Id: "message-id2"},
			},
			language: messages.LanguagePack{
				GenericName:  "en",
				SpecificName: "en-US",
				Messages: map[string]string{
					"message-id1": "test message 1",
					"message-id2": "test message 2",
				},
			},
			expectedResult: true,
			expectedTranslations: []string{
				"test message 1",
				"test message 2",
			},
		},
	}

	for i, item := range data {
		ok := item.messages.Translate(&item.language)
		if !ok && item.expectedResult {
			t.Errorf("Item %d, “%s”: mismatch results. Expecting true", i, item.description)

		} else if ok && !item.expectedResult {
			t.Errorf("Item %d, “%s”: mismatch results. Expecting false", i, item.description)
		}

		messages := []Translator(item.messages)
		if len(messages) != len(item.expectedTranslations) {
			t.Fatalf(
				"Item %d, “%s”: mismatch results. Expecting '%d'; found '%d'",
				i,
				item.description,
				len(item.expectedTranslations),
				len(messages),
			)
		}

		for i, message := range messages {
			m, ok := message.(*MessageResponse)
			if !ok {
				t.Fatalf("Item %d, “%s”: Not a message response!", i, item.description)
			}

			if m.Message != item.expectedTranslations[i] {
				t.Errorf(
					"Item %d, “%s”: mismatch results. Expecting '%s'; found '%s'",
					i,
					item.description,
					item.expectedTranslations[i],
					m.Message,
				)
			}
		}
	}
}

type strangeTranslator struct {
}

func (s strangeTranslator) Translate(language *messages.LanguagePack) bool {
	return true
}

func TestMessageHolderAdd(t *testing.T) {
	data := []struct {
		description    string
		translators    []Translator
		expectedResult bool
		expectedNumber int
	}{
		{
			description: "Strange translator object is being add",
			translators: []Translator{
				&MessageResponse{Id: "id1"},
				Messages{
					&MessageResponse{Id: "id2"},
					&MessageResponse{Id: "id3"},
				},
				strangeTranslator{},
			},
			expectedResult: false,
			expectedNumber: 3,
		},
		{
			description: "Nil translator object is being add",
			translators: []Translator{
				&MessageResponse{Id: "id1"},
				Messages{
					&MessageResponse{Id: "id2"},
					&MessageResponse{Id: "id3"},
				},
				nil,
			},
			expectedResult: false,
			expectedNumber: 3,
		},
		{
			description: "Normal case with mixed types",
			translators: []Translator{
				&MessageResponse{Id: "id1"},
				Messages{
					&MessageResponse{Id: "id2"},
					&MessageResponse{Id: "id3"},
				},
			},
			expectedResult: true,
			expectedNumber: 3,
		},
	}

	for i, item := range data {
		var messagesHolder MessagesHolder
		ok := messagesHolder.Add(item.translators...)

		if !ok && item.expectedResult {
			t.Errorf("Item %d, “%s”: mismatch results. Expecting true", i, item.description)

		} else if ok && !item.expectedResult {
			t.Errorf("Item %d, “%s”: mismatch results. Expecting false", i, item.description)
		}

		if len(messagesHolder.msgs) != item.expectedNumber {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%d'; found '%d'",
				i,
				item.description,
				item.expectedNumber,
				len(messagesHolder.msgs),
			)
		}
	}
}

func TestMessageHolderMessages(t *testing.T) {
	data := []struct {
		description    string
		translators    []Translator
		expectedResult bool
		expectedNil    bool
		expectedNumber int
	}{
		{
			description:    "No messages inside",
			translators:    nil,
			expectedResult: true,
			expectedNil:    true,
			expectedNumber: 0,
		},
		{
			description: "Normal case with mixed types",
			translators: []Translator{
				&MessageResponse{Id: "id1"},
				Messages{
					&MessageResponse{Id: "id2"},
					&MessageResponse{Id: "id3"},
				},
			},
			expectedResult: true,
			expectedNil:    false,
			expectedNumber: 3,
		},
	}

	for i, item := range data {
		var messagesHolder MessagesHolder

		if item.translators != nil {
			ok := messagesHolder.Add(item.translators...)

			if !ok && item.expectedResult {
				t.Errorf("Item %d, “%s”: mismatch results. Expecting true", i, item.description)

			} else if ok && !item.expectedResult {
				t.Errorf("Item %d, “%s”: mismatch results. Expecting false", i, item.description)
			}
		}

		messages := messagesHolder.Messages()
		if messages != nil && item.expectedNil {
			t.Fatalf("Item %d, “%s”: mismatch results. Expecting nil object", i, item.description)
		} else if messages == nil && !item.expectedNil {
			t.Fatalf("Item %d, “%s”: mismatch results. Unexpected nil object", i, item.description)
		}

		if item.expectedNil {
			continue
		}

		msgs, ok := messages.(Messages)
		if !ok {
			t.Fatal("Not a Messages object")
		}
		specificMsgs := []Translator(msgs)

		if len(specificMsgs) != item.expectedNumber {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%d'; found '%d'",
				i,
				item.description,
				item.expectedNumber,
				len(specificMsgs),
			)
		}
	}
}

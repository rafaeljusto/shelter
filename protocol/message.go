// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/messages"
	"net/url"
)

// List of possible input error types. This is useful for translating input errors to good messages
// for the end-user in the desired language
const (
	ErrorCodeAcceptCharset        ErrorCode = "accept-charset-error"
	ErrorCodeAccept               ErrorCode = "accept-error"
	ErrorCodeAcceptLanguage       ErrorCode = "accept-language-error"
	ErrorCodeAuthorizationMissing ErrorCode = "authorization-missing"
	ErrorCodeConflict             ErrorCode = "conflict"
	ErrorCodeContentMD5Missing    ErrorCode = "content-md5-missing"
	ErrorCodeContentTypeMissing   ErrorCode = "content-type-missing"
	ErrorCodeDateMissing          ErrorCode = "date-missing"
	ErrorCodeIfMatchFailed        ErrorCode = "if-match-failed"
	ErrorCodeIfNoneMatchFailed    ErrorCode = "if-none-match-failed"
	ErrorCodeInvalidAuthorization ErrorCode = "invalid-authorization"
	ErrorCodeInvalidContentMD5    ErrorCode = "invalid-content-md5"
	ErrorCodeInvalidContentType   ErrorCode = "invalid-content-type"
	ErrorCodeInvalidDateTimeFrame ErrorCode = "invalid-date-time-frame"
	ErrorCodeInvalidDNSKEY        ErrorCode = "invalid-dnskey"
	ErrorCodeInvalidDSAlgorithm   ErrorCode = "invalid-ds-algorithm"
	ErrorCodeInvalidDSDigestType  ErrorCode = "invalid-ds-digest-type"
	ErrorCodeInvalidFQDN          ErrorCode = "invalid-fqdn"
	ErrorCodeInvalidHeaderDate    ErrorCode = "invalid-header-date"
	ErrorCodeInvalidIfMatch       ErrorCode = "invalid-if-match"
	ErrorCodeInvalidIfNoneMatch   ErrorCode = "invalid-if-none-match"
	ErrorCodeInvalidIP            ErrorCode = "invalid-ip"
	ErrorCodeInvalidJSONContent   ErrorCode = "invalid-json-content"
	ErrorCodeInvalidLanguage      ErrorCode = "invalid-language"
	ErrorCodeInvalidQueryOrderBy  ErrorCode = "invalid-query-order-by"
	ErrorCodeInvalidQueryPage     ErrorCode = "invalid-query-page"
	ErrorCodeInvalidQueryPageSize ErrorCode = "invalid-query-page-size"
	ErrorCodeInvalidURI           ErrorCode = "invalid-uri"
	ErrorCodeSectetNotFound       ErrorCode = "secret-not-found"
)

// ErrorCode will garantee that only known errors are reported
type ErrorCode string

// Translator represents an object that can be translated to a specific language so we can return it
// to the client in a specific representation
type Translator interface {
	Translate(language *messages.LanguagePack) bool
}

// MessageResponse struct was created to return a message to the user with more
// information to easy integrate and solve problems
type MessageResponse struct {
	Id      string `json:"id,omitempty"`      // Code for integration systems to automatically solve the problem
	Field   string `json:"field,omitempty"`   // Field that the message is related to
	Value   string `json:"value,omitempty"`   // Current value of the field
	Message string `json:"message,omitempty"` // Message in the user's desired language
	Links   []Link `json:"links,omitempty"`   // Links associating this message with other resources
}

// NewMessageResponse builds a message adding necessary links for the response
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

// NewMessageResponseWithField builds a message with a specific field adding necessary links for the
// response. This type of message is better for the end-user determinates what was the problem
func NewMessageResponseWithField(id string, field, value, roid string,
	language *messages.LanguagePack) (*MessageResponse, error) {

	messageResponse, err := NewMessageResponse(id, roid, language)
	if err != nil {
		return nil, err
	}

	messageResponse.Field = field
	messageResponse.Value = value
	return messageResponse, nil
}

// Translate fills the convert the message content to the desired language before returning it to
// the end user. If the message is not found in the desired language it will return false
func (msg *MessageResponse) Translate(language *messages.LanguagePack) bool {
	if language != nil && language.Messages != nil {
		var ok bool
		msg.Message, ok = language.Messages[string(msg.Id)]
		return ok
	}

	// When we don't find the message file, at least add the message id to make it easy to
	// identify the error. Useful on test scenarios
	msg.Message = string(msg.Id)
	return false
}

// Messages is set of message objects to add methods to make our life easier and to be able to
// return many messages at once
type Messages []Translator

// Translate converts every message to the desired end-user language. If one of the messages is not
// found in the desired language it will return false
func (m Messages) Translate(language *messages.LanguagePack) bool {
	translationOK := true
	for i := range m {
		if ok := m[i].Translate(language); !ok {
			translationOK = false
		}
	}

	return translationOK
}

// MessagesHolder make it easy to interact with a single message or a group of messages
type MessagesHolder struct {
	msgs Messages
}

// Add is responsable for inserting new messages or list of messages into an internal buffer
func (m *MessagesHolder) Add(trs ...Translator) bool {
	for _, tr := range trs {
		if tr == nil {
			return false
		}

		switch tr.(type) {
		case *MessageResponse:
			m.msgs = append(m.msgs, tr.(*MessageResponse))
		case Messages:
			m.msgs = append(m.msgs, tr.(Messages)...)
		default:
			return false
		}
	}

	return true
}

// Messages returns all the messages that were stored in the MessagesHolder object
func (m *MessagesHolder) Messages() Translator {
	if len(m.msgs) == 0 {
		return nil
	}

	return m.msgs
}

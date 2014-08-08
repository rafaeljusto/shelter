// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package errors describe the low level errors of the Shelter system
package errors

import (
	"fmt"
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

// InputError was created to report all input user errors, containing all the necessary information
// to identify the field and the value that caused the problem
type InputError struct {
	Id    ErrorCode // Error code to be translated on the protocol
	Field string    // Field that caused the problem
	Value string    // Value used in the field that caused the problem
}

// NewInputError is an easy way to build an error
func NewInputError(id ErrorCode, field, value string) InputError {
	return InputError{
		Id:    id,
		Field: field,
		Value: value,
	}
}

// Error retrieve the description of the user input problen
func (e InputError) Error() string {
	return fmt.Sprintf("Error '%s' in field '%s' with the value '%s'", e.Id, e.Field, e.Value)
}

package check

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"github.com/rafaeljusto/shelter/net/http/rest/messages"
	"net/http"
	"strings"
	"time"
)

const (
	// Variable defining the supported content type for the system. For now we only support JSON, but
	// the idea is to support in a near future XML
	SupportedContentType = "application/vnd.shelter+json"

	// Define the supported charset of the system. For now we use for everything utf-8, from database
	// to data manipulation. There's no conversion for any special charset
	SupportedCharset = "utf-8"

	// Variable used to determinate the namespace in Authorization HTTP header. The format
	// is "<namespace> <secretId>:<secret>"
	SupportedNamespace = "shelter"

	// Date HTTP header field must be near the current time, otherwise we must assume that is a reply
	// attack. The variable bellow stores the tollerance for the date HTTP header field
	timeFrameDuration = time.Duration(10 * time.Minute)
)

// List of possible errors that can occur when calling functions from this file. Other
// erros can also occurs from low level layers
var (
	ErrHTTPAuthorizationNotFound = errors.New("Missing HTTP Authorization header")
	ErrInvalidHTTPAuthorization  = errors.New("Invalid HTTP Authorization header")
)

// This method check the content type that the user support for answers. If the user doesn't support
// the system content types we should return an HTTP error code
func HTTPAccept(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	accept = strings.TrimSpace(accept)
	accept = strings.ToLower(accept)

	if len(accept) == 0 {
		return true
	}

	for _, acceptPart := range strings.Split(accept, ",") {
		acceptPart = strings.TrimSpace(acceptPart)

		// For now we are ignoring versioning
		if idx := strings.Index(acceptPart, ";"); idx > 0 {
			acceptPart = acceptPart[0:idx]
		}

		if acceptPart == "*" || acceptPart == SupportedContentType {
			return true
		}
	}

	return false
}

// The accept language check beyond verifying if the language exists in out system, set
// the first language found in the context
func HTTPAcceptLanguage(r *http.Request, context *context.Context) bool {
	acceptLanguage := r.Header.Get("Accept-Language")
	acceptLanguage = strings.TrimSpace(acceptLanguage)
	acceptLanguage = strings.ToLower(acceptLanguage)

	if len(acceptLanguage) == 0 {
		return true
	}

	for _, acceptLanguagePart := range strings.Split(acceptLanguage, ",") {
		acceptLanguagePart = strings.TrimSpace(acceptLanguagePart)

		// For now we are ignoring language preference
		if idx := strings.Index(acceptLanguagePart, ";"); idx > 0 {
			acceptLanguagePart = acceptLanguagePart[0:idx]
		}

		if acceptLanguagePart == "*" {
			return true
		}

		languagePack := messages.ShelterRESTLanguagePacks.Select(acceptLanguagePart)
		if languagePack != nil {
			context.Language = languagePack
			return true
		}
	}

	return false
}

// Accept Charset HTTP header field is verified in this method. For now we only support UTF-8
// namesepace, there're no intentions to addopt ISO-8859-1
func HTTPAcceptCharset(r *http.Request) bool {
	acceptCharset := r.Header.Get("Accept-Charset")
	acceptCharset = strings.TrimSpace(acceptCharset)
	acceptCharset = strings.ToLower(acceptCharset)

	if len(acceptCharset) == 0 {
		return true
	}

	for _, acceptCharsetPart := range strings.Split(acceptCharset, ",") {
		acceptCharsetPart = strings.TrimSpace(acceptCharsetPart)

		// For now we are ignoring charset preference
		if idx := strings.Index(acceptCharsetPart, ";"); idx > 0 {
			acceptCharsetPart = acceptCharsetPart[0:idx]
		}

		if acceptCharsetPart == "*" || acceptCharsetPart == SupportedCharset {
			return true
		}
	}

	return false
}

// Check the user current content type format. For now we only accept JSON content respecting the
// Shelter protocol, but in a near future we plan to accept XML too
func HTTPContentType(r *http.Request) bool {
	contentType := getHTTPContentType(r)
	if len(contentType) == 0 {
		return true
	}

	// For now we are ignoring version and only checking the charset option
	if idx := strings.Index(contentType, ";"); idx > 0 {
		options := contentType[idx:]

		optionsParts := strings.Split(options, ";")
		for _, optionsPart := range optionsParts {
			option := strings.Split(optionsPart, "=")
			if len(option) != 2 {
				continue
			}

			key := option[0]
			key = strings.TrimSpace(key)
			key = strings.ToLower(key)

			value := option[1]
			value = strings.TrimSpace(value)
			value = strings.ToLower(value)

			if key == "charset" {
				if value != SupportedCharset {
					return false
				}
			}
		}

		// Removing options from content type to compare it with the supported types of the
		// system
		contentType = contentType[0:idx]
	}

	return contentType == SupportedContentType
}

// To garantee that the content was not modified during the network phase or is incomplete, we check
// the hash of the content and compare with the HTTP header field
func HTTPContentMD5(r *http.Request, context *context.Context) bool {
	contentMD5 := getHTTPContentMD5(r)
	if len(contentMD5) == 0 {
		return true
	}

	hash := md5.New()
	hash.Write(context.RequestContent)
	hashBytes := hash.Sum(nil)
	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

	return hashBase64 == contentMD5
}

// HTTPDate method is responsable for checking the time frame of the request, avoiding reply
// attacks, that's when an attacker use the same request again in a different moment to corrupt the
// data
func HTTPDate(r *http.Request) (bool, error) {
	dateStr := getHTTPDate(r)
	if len(dateStr) == 0 {
		return true, nil
	}

	date, err := time.Parse(time.RFC1123, dateStr)
	if err != nil {
		return false, err
	}

	// Check if the date is inside the time frame, avoiding reply attack
	now := time.Now().UTC()
	frameInception := now.Add(timeFrameDuration * -1)
	frameExpiration := now.Add(timeFrameDuration)
	return !date.UTC().Before(frameInception) && !date.UTC().After(frameExpiration), nil
}

// HTTPAuthorization garantees that the user was the only that really sent the information. Using a
// group of information of the request and a shared secret the server can recreate the authorization
// data and compare it with the header field. We are using the same approach that Amazon company
// used in their Cloud API. More information can be found in
// http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html#ConstructingTheAuthenticationHeader
func HTTPAuthorization(r *http.Request, secretFinder func(string) (string, error)) (bool, error) {
	authorization := r.Header.Get("Authorization")
	authorization = strings.TrimSpace(authorization)

	// Authorization header is mandatory in all requests
	if len(authorization) == 0 {
		return false, ErrHTTPAuthorizationNotFound
	}

	// The authorization should always follow the format: "<namespace> <secretId>:<secret>"
	authorizationParts := strings.Split(authorization, " ")
	if len(authorizationParts) != 2 {
		return false, ErrInvalidHTTPAuthorization
	}

	namespace := authorizationParts[0]
	namespace = strings.TrimSpace(namespace)
	namespace = strings.ToLower(namespace)

	if namespace != SupportedNamespace {
		return false, ErrInvalidHTTPAuthorization
	}

	secretParts := strings.Split(authorizationParts[1], ":")
	if len(secretParts) != 2 {
		return false, ErrInvalidHTTPAuthorization
	}

	secretId := secretParts[0]
	secretId = strings.TrimSpace(secretId)
	secretId = strings.ToLower(secretId)

	stringToSign, err := BuildStringToSign(r, secretId)
	if err != nil {
		return false, err
	}

	secret, err := secretFinder(secretId)
	if err != nil {
		return false, err
	}

	signature := GenerateSignature(stringToSign, secret)
	return signature == secretParts[1], nil
}

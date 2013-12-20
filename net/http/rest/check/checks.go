package check

import (
	"crypto/md5"
	"encoding/base64"
	"net/http"
	"shelter/net/http/rest/context"
	"shelter/net/http/rest/language"
	"strconv"
	"strings"
	"time"
)

const (
	SupportedContentType = "application/vnd.shelter+json"
	timeFrameDuration    = "10m"
)

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
func HTTPAcceptLanguage(r *http.Request, context *context.ShelterRESTContext) bool {
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

		languagePack := language.ShelterRESTLanguagePacks.Select(acceptLanguagePart)
		if languagePack != nil {
			context.Language = languagePack
			return true
		}
	}

	return false
}

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

		if acceptCharsetPart == "*" || acceptCharset == "utf-8" {
			return true
		}
	}

	return false
}

func ContentType(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	contentType = strings.TrimSpace(contentType)
	contentType = strings.ToLower(contentType)

	if len(contentType) == 0 {
		return true
	}

	// For now we are ignoring version
	if idx := strings.Index(contentType, ";"); idx > 0 {
		contentType = contentType[0:idx]
	}

	return contentType == SupportedContentType
}

func HTTPContentMD5(r *http.Request, context *context.ShelterRESTContext) bool {
	contentMD5 := r.Header.Get("Content-MD5")
	contentMD5 = strings.TrimSpace(contentMD5)

	if len(contentMD5) == 0 {
		return true
	}

	hash := md5.New()
	hash.Write(context.RequestContent)
	hashBytes := hash.Sum(nil)
	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

	return hashBase64 == contentMD5
}

func Date(r *http.Request) (bool, error) {
	dateStr := r.Header.Get("Date")
	dateStr = strings.TrimSpace(dateStr)

	if len(dateStr) == 0 {
		return true, nil
	}

	date, err := time.Parse(time.RFC1123, dateStr)
	if err != nil {
		return false, err
	}

	// Check if the date is inside the time frame, avoiding reply attack
	now := time.Now().UTC()
	duration, _ := time.ParseDuration(timeFrameDuration)
	frameInception := now.Add(duration * -1)
	frameExpiration := now.Add(duration)
	return !date.UTC().Before(frameInception) && !date.UTC().After(frameExpiration), nil
}

func IfModifiedSince(r *http.Request, lastModifiedAt time.Time) (bool, error) {
	ifModifiedSinceStr := r.Header.Get("If-Modified-Since")
	ifModifiedSinceStr = strings.TrimSpace(ifModifiedSinceStr)

	if len(ifModifiedSinceStr) == 0 {
		return true, nil
	}

	ifModifiedSince, err := time.Parse(time.RFC1123, ifModifiedSinceStr)
	if err != nil {
		return true, err
	}

	if lastModifiedAt.Before(ifModifiedSince) || lastModifiedAt.Equal(ifModifiedSince) {
		return false, nil
	}

	return true, nil
}

func IfUnmodifiedSince(r *http.Request, lastModifiedAt time.Time) (bool, error) {
	ifUnmodifiedSinceStr := r.Header.Get("If-Unmodified-Since")
	ifUnmodifiedSinceStr = strings.TrimSpace(ifUnmodifiedSinceStr)

	if len(ifUnmodifiedSinceStr) == 0 {
		return true, nil
	}

	ifUnmodifiedSince, err := time.Parse(time.RFC1123, ifUnmodifiedSinceStr)
	if err != nil {
		return true, err
	}

	if lastModifiedAt.After(ifUnmodifiedSince) {
		return false, nil
	}

	return true, nil
}

func IfMatch(r *http.Request, revision int) (bool, error) {
	ifMatch := r.Header.Get("If-Match")
	ifMatch = strings.TrimSpace(ifMatch)

	if len(ifMatch) == 0 {
		return true, nil
	}

	ifMatchParts := strings.Split(ifMatch, ",")

	for _, ifMatchPart := range ifMatchParts {
		ifMatchPart = strings.TrimSpace(ifMatchPart)

		// If "*" is given and no current entity exists, the server MUST NOT perform the
		// requested method, and MUST return a 412 (Precondition Failed) response
		if ifMatchPart == "*" {
			return (revision > 0), nil
		}

		etag, err := strconv.Atoi(ifMatchPart)
		if err != nil {
			return false, err
		}

		if etag == revision {
			return true, nil
		}
	}

	// RFC 2616 - 14.24 - If none of the entity tags match the server MUST NOT perform the
	// requested method, and MUST return a 412 (Precondition Failed) response
	return false, nil
}

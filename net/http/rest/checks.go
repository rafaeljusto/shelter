package rest

import (
	"crypto/md5"
	"encoding/base64"
	"net/http"
	"shelter/net/http/rest/language"
	"strings"
	"time"
)

const (
	supportedContentType = "application/vnd.shelter+json"
	timeFrameDuration    = "10m"
)

func checkHTTPAccept(r *http.Request) bool {
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

		if acceptPart == "*" || acceptPart == supportedContentType {
			return true
		}
	}

	return false
}

// The accept language check beyond verifying if the language exists in out system, set
// the first language found in the context
func checkHTTPAcceptLanguage(r *http.Request, context *ShelterRESTContext) bool {
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

func checkHTTPAcceptCharset(r *http.Request) bool {
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

func checkContentType(r *http.Request) bool {
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

	return contentType == supportedContentType
}

func checkHTTPContentMD5(r *http.Request, context *ShelterRESTContext) bool {
	contentMD5 := r.Header.Get("Content-MD5")
	contentMD5 = strings.TrimSpace(contentMD5)

	if len(contentMD5) == 0 {
		return true
	}

	hash := md5.New()
	hash.Write(context.requestContent)
	hashBytes := hash.Sum(nil)
	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

	return hashBase64 == contentMD5
}

func checkDate(r *http.Request) bool {
	dateStr := r.Header.Get("Date")
	dateStr = strings.TrimSpace(dateStr)

	if len(dateStr) == 0 {
		return true
	}

	date, err := time.Parse(time.RFC1123, dateStr)
	if err != nil {
		return false
	}

	// Check if the date is inside the time frame, avoiding reply attack
	now := time.Now().UTC()
	duration, _ := time.ParseDuration(timeFrameDuration)
	frameInception := now.Add(duration * -1)
	frameExpiration := now.Add(duration)
	return !date.UTC().Before(frameInception) && !date.UTC().After(frameExpiration)
}

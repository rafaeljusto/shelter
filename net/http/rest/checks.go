package rest

import (
	"net/http"
	"shelter/net/http/rest/language"
	"strings"
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

		if acceptPart == "*" || acceptPart == "application/vnd.shelter+json" {
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

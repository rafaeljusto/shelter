package check

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"shelter/net/http/rest/context"
	"shelter/net/http/rest/language"
	"sort"
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

func Authorization(r *http.Request) bool {
	// http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html#ConstructingTheAuthenticationHeader

	// StringToSign = HTTP-Verb + "\n" +
	// 	Content-MD5 + "\n" + // RFC1864
	// 	Content-Type + "\n" +
	// 	Date + "\n" +
	// 	AccessKeyID + "\n" +
	// 	Path + "\n" +
	// 	CanonicalizedQueryString;

	stringToSign := r.Method

	contentMD5 := ""
	contentType := ""

	if r.ContentLength > 0 {
		contentMD5 = r.Header.Get("Content-MD5")
		contentMD5 = strings.TrimSpace(contentMD5)

		if len(contentMD5) == 0 {
			// TODO: Error?
		}

		contentType = r.Header.Get("Content-Type")
		contentType = strings.TrimSpace(contentType)
		contentType = strings.ToLower(contentType)

		if len(contentType) == 0 {
			// TODO: Error?
		}

		// For now we are ignoring version
		if idx := strings.Index(contentType, ";"); idx > 0 {
			contentType = contentType[0:idx]
		}

		stringToSign = fmt.Sprintf("%s\n%s", stringToSign, contentMD5)
		stringToSign = fmt.Sprintf("%s\n%s", stringToSign, contentType)
	}

	dateStr := r.Header.Get("Date")
	dateStr = strings.TrimSpace(dateStr)

	if len(dateStr) == 0 {
		// TODO: Error?
	}

	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, dateStr)
	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, "secretId")
	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, r.URL.Path)

	var queryString []string
	for k, v := range r.URL.Query() {
		for _, vm := range v { // multiple values
			keyAndValue := fmt.Sprintf("%s=%s", k, vm)
			queryString = append(queryString, keyAndValue)
		}
	}

	sort.Strings(queryString)
	sortedQueryString := strings.Join(queryString, "&")
	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, sortedQueryString)

	// TODO

	return true
}

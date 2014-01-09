package check

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// List of possible errors that can occur when calling functions from this file. Other
// erros can also occurs from low level layers. Some HTTP headers are mandatory because of
// the system authentication mechanisms
var (
	ErrHTTPContentMD5NotFound  = errors.New("Content-MD5 HTTP header not found")
	ErrHTTPContentTypeNotFound = errors.New("Content-Type HTTP header not found")
	ErrHTTPDateNotFound        = errors.New("Date HTTP header not found")
)

// Build the string that will be used to generate the signature of the HTTP Authorization
// header field. The format is defined bellow.
//
// StringToSign = HTTP-Verb + "\n" +
//  Content-MD5 + "\n" + // RFC1864
//  Content-Type + "\n" +
//  Date + "\n" +
//  AccessKeyID + "\n" +
//  Path + "\n" +
//  CanonicalizedQueryString;
func BuildStringToSign(r *http.Request, secretId string) (string, error) {
	stringToSign := r.Method

	contentMD5 := ""
	contentType := ""

	if r.ContentLength > 0 {
		contentMD5 = getHTTPContentMD5(r)

		if len(contentMD5) == 0 {
			return "", ErrHTTPContentMD5NotFound
		}

		contentType = getHTTPContentType(r)

		if len(contentType) == 0 {
			return "", ErrHTTPContentTypeNotFound
		}

		// For now we are ignoring version
		if idx := strings.Index(contentType, ";"); idx > 0 {
			contentType = contentType[0:idx]
		}

		stringToSign = fmt.Sprintf("%s\n%s", stringToSign, contentMD5)
		stringToSign = fmt.Sprintf("%s\n%s", stringToSign, contentType)
	}

	dateStr := getHTTPDate(r)
	if len(dateStr) == 0 {
		return "", ErrHTTPDateNotFound
	}

	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, dateStr)
	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, secretId)
	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, r.URL.Path)

	var queryString []string
	for key, values := range r.URL.Query() {
		for _, value := range values {
			keyAndValue := fmt.Sprintf("%s=%s", key, value)
			queryString = append(queryString, keyAndValue)
		}
	}

	sort.Strings(queryString)
	sortedQueryString := strings.Join(queryString, "&")
	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, sortedQueryString)

	return stringToSign, nil
}

func GenerateSignature(stringToSign, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

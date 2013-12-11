package rest

import (
	"net/http"
)

func init() {
	http.Handle("/domain/", handleDomain)
}

func handleDomain(w http.ResponseWriter, r *http.Request) {
	// Supported HTTP headers in request:
	//   Method
	//   Date
	//   Content-type
	//   Content-Length
	//   Content-MD5
	//   Accept
	//   Accept-Charset
	//   Accept-Language
	//   Authorization
	//   If-Modified-Since
	//   If-Match
	//   If-None-Match

	// Supported HTTP headers in response:
	//   Content-Encoding
	//   Content-Language
	//   Content-Length
	//   Content-MD5
	//   Content-Type
	//   Date
	//   ETag
	//   Last-Modified
	//   Status
	//   Accept-Language
	//   Accept
}

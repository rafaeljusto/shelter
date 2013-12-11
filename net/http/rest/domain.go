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

	// TODO Check:
	//   Accept
	//   Accept-Charset
	//   Accept-Language
	//   Content-MD5
	//   Content-Type
	//   Authorization

	if r.Method == "GET" {
		retrieveDomain(w, r)

	} else if r.Method == "PUT" {
		createUpdateDomain(w, r)

	} else if r.Method == "DELETE" {
		removeDomain(w, r)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func retrieveDomain(w http.ResponseWriter, r *http.Request) {
	// TODO Check:
	//   If-Modified-Since
	//   If-Match
	//   If-None-Match
}

func createUpdateDomain(w http.ResponseWriter, r *http.Request) {

}

func removeDomain(w http.ResponseWriter, r *http.Request) {

}

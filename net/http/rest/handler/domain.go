package rest

import (
	"net/http"
	"shelter/net/http/rest"
	"shelter/net/http/rest/language"
)

func init() {
	rest.HandleFunc("/domain/", handleDomain)
}

func handleDomain(w http.ResponseWriter, r *http.Request, context *ShelterRESTContext) {
	if r.Method == "GET" {
		retrieveDomain(w, r, context)

	} else if r.Method == "PUT" {
		createUpdateDomain(w, r, context)

	} else if r.Method == "DELETE" {
		removeDomain(w, r, context)

	} else {
		context.Response(http.StatusMethodNotAllowed)
	}
}

func retrieveDomain(w http.ResponseWriter, r *http.Request, context *ShelterRESTContext) {
	// TODO Check:
	//   If-Modified-Since
	//   If-Match
	//   If-None-Match
}

func createUpdateDomain(w http.ResponseWriter, r *http.Request, context *ShelterRESTContext) {

}

func removeDomain(w http.ResponseWriter, r *http.Request, context *ShelterRESTContext) {

}

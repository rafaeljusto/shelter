package handler

import (
	"fmt"
	"github.com/rafaeljusto/shelter/log"
	"io"
	"net/http"
	"regexp"
)

func init() {
	HandleFunc(regexp.MustCompile(`^/domains(/.*)?$`), HandleDomains)
}

func HandleDomains(w http.ResponseWriter, r *http.Request) {
	restAddress, err := retrieveRESTAddress()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error while retrieving the REST address. Details:", err)
		return
	}

	request, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s%s", restAddress, r.RequestURI),
		nil,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating a request in web client. Details:", err)
		return
	}

	request.Header.Set("Accept-Language", r.Header.Get("Accept-Language"))

	response, err := signAndSend(request, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error signing and sending a request in web client. Details:", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		// TODO: Create a function to detect errors. Will be necessary when we receive a bad
		// request while creating a domain for example
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Sprintf("Exepected status code %d but received %d from "+
			"/domains result in web client", http.StatusOK, response.StatusCode))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	if _, err := io.Copy(w, response.Body); err != nil {
		// Here we already set the response code, so the client will receive a OK result
		// without body
		log.Println("Error copying REST response to web client response. Details:", err)
		return
	}
}

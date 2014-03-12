package handler

import (
	"fmt"
	"github.com/rafaeljusto/shelter/log"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func init() {
	HandleFunc(regexp.MustCompile(`^/domain/([[:alnum:]]|\-|\.)+$`), HandleDomain)
}

func HandleDomain(w http.ResponseWriter, r *http.Request) {
	restAddress, err := retrieveRESTAddress()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error while retrieving the REST address. Details:", err)
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error while reading request body in web client. Details:", err)
		return
	}

	request, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s%s", restAddress, r.RequestURI),
		strings.NewReader(string(content)),
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating a request in web client. Details:", err)
		return
	}

	request.Header.Set("Accept-Language", r.Header.Get("Accept-Language"))

	response, err := signAndSend(request, content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error signing and sending a request in web client. Details:", err)
		return
	}

	if response.StatusCode != http.StatusCreated &&
		response.StatusCode != http.StatusNoContent {

		// TODO: Create a function to detect errors. Will be necessary when we receive a bad
		// request while creating a domain for example
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Sprintf("Exepected status code %d or %d but received %d from "+
			"/domain result in web client", http.StatusCreated, http.StatusNoContent,
			response.StatusCode))
		return
	}

	w.WriteHeader(response.StatusCode)
}

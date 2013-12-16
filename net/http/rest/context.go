package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"labix.org/v2/mgo"
	"net/http"
	"shelter/config"
	"shelter/database/mongodb"
	"shelter/net/http/rest/language"
)

// ShelterRESTContext is responsable to store a state of a request during the request
// live. It's necessary to create the request enviroment, with the preferred settings of
// the user
type ShelterRESTContext struct {
	Database           *mgo.Database          // MongoDB Database
	Language           *language.LanguagePack // Language choosen by the user
	requestContent     []byte                 // Request body
	responseHttpStatus int                    // Response HTTP status
	responseMessage    []byte                 // Response HTTP message
	httpHeader         map[string]string      // Extra headers to be sent in the response
}

// Initialize a new context. By default use system choosen language pack. We also store
// the bytes from the request body because we can read only once from the buffer reader
// and we need to check it for some HTTP header verifications. We don't return a pointer
// of the context because we want to control the destruction of the object and don't leave
// it to the garbage collector
func newShelterRESTContext(r *http.Request) (ShelterRESTContext, error) {
	context := ShelterRESTContext{
		Language: language.ShelterRESTLanguagePack,
	}

	var err error
	if context.requestContent, err = ioutil.ReadAll(r.Body); err != nil {
		return context, err
	}

	database, err := mongodb.Open(
		config.ShelterConfig.Database.URI,
		config.ShelterConfig.Database.Name,
	)

	if err != nil {
		return context, err
	}

	context.Database = database
	return context, nil
}

// Transform the content body, that is in JSON format into an object
func (s *ShelterRESTContext) JSONRequest(object interface{}) error {
	decoder := json.NewDecoder(bytes.NewBuffer(s.requestContent))
	return decoder.Decode(object)
}

// Store only the HTTP status, for no content responses
func (s *ShelterRESTContext) Response(httpStatus int) {
	s.responseHttpStatus = httpStatus
}

// Store a message response, translating the message id to the proper language message
func (s *ShelterRESTContext) MessageResponse(httpStatus int, messageId string) {
	s.responseHttpStatus = httpStatus
	s.responseMessage = []byte(s.Language.Messages[messageId])
}

// Store a object in json format for the response
func (s *ShelterRESTContext) JSONResponse(httpStatus int, object interface{}) error {
	content, err := json.Marshal(object)

	if err != nil {
		return err
	}

	s.responseHttpStatus = httpStatus
	s.responseMessage = content
	return nil
}

// Add a custom HTTP header. Used for some types of response where you need to set ETag or
// LastModified fields
func (s *ShelterRESTContext) AddHeader(key, value string) {
	// Avoid adding headers that are automatically generated at the end of the request. We
	// don't allow header overwrite because in the low level MIMEHeader the HTTP header
	// value is appended instead of replaced
	if key == "Content-Type" ||
		key == "Content-Encoding" ||
		key == "Content-Charset" ||
		key == "Content-Language" ||
		key == "Content-Length" ||
		key == "Content-MD5" ||
		key == "Accept" ||
		key == "Accept-Language" ||
		key == "Accept-Charset" ||
		key == "Date" {
		return
	}

	s.httpHeader[key] = value
}

package rest

import (
	"encoding/json"
	"labix.org/v2/mgo"
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
	responseHttpStatus int                    // Response HTTP status
	responseMessage    []byte                 // Response HTTP message
	httpHeader         map[string]string      // Extra headers to be sent in the response
}

// Initialize a new context. By default use system choosen language pack
func newShelterRESTContext() (ShelterRESTContext, error) {
	context := ShelterRESTContext{
		Language: language.ShelterRESTLanguagePack,
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

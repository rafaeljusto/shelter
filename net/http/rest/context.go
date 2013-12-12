package rest

import (
	"encoding/json"
	"shelter/net/http/rest/language"
)

// ShelterRESTContext is responsable to store a state of a request during the request
// live. It's necessary to create the request enviroment, with the preferred settings of
// the user
type ShelterRESTContext struct {
	Language           *language.LanguagePack // Language choosen by the user
	responseHttpStatus int                    // Response HTTP status
	responseMessage    []byte                 // Response HTTP message
}

// Initialize a new context. By default use system choosen language pack
func newShelterRESTContext() ShelterRESTContext {
	return ShelterRESTContext{
		Language: language.ShelterRESTLanguagePack,
	}
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

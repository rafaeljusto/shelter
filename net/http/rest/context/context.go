package context

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"labix.org/v2/mgo"
	"net/http"
	"shelter/net/http/rest/messages"
)

// Context is responsable to store a state of a request during the request
// live. It's necessary to create the request enviroment, with the preferred settings of
// the user
type Context struct {
	Database           *mgo.Database          // MongoDB Database
	Language           *messages.LanguagePack // Language choosen by the user
	RequestContent     []byte                 // Request body
	ResponseHTTPStatus int                    // Response HTTP status
	ResponseContent    []byte                 // Response body
	HTTPHeader         map[string]string      // Extra headers to be sent in the response
}

// Initialize a new context. By default use system choosen messages pack. We also store
// the bytes from the request body because we can read only once from the buffer reader
// and we need to check it for some HTTP header verifications. We don't return a pointer
// of the context because we want to control the destruction of the object and don't leave
// it to the garbage collector. We are injecting the database connection into the context
// to allow unit tests in it. If the programmer forgot to use this constructor, the low
// level DAOs are going to detect a nil pointer
func NewContext(r *http.Request, database *mgo.Database) (Context, error) {
	context := Context{
		Database:           database,
		Language:           messages.ShelterRESTLanguagePack,
		ResponseHTTPStatus: http.StatusOK, // Default status code
		HTTPHeader:         make(map[string]string),
	}

	if r.ContentLength > 0 && r.Body != nil {
		var err error
		if context.RequestContent, err = ioutil.ReadAll(r.Body); err != nil {
			return context, err
		}
	}

	return context, nil
}

// Transform the content body, that is in JSON format into an object
func (s *Context) JSONRequest(object interface{}) error {
	decoder := json.NewDecoder(bytes.NewBuffer(s.RequestContent))
	return decoder.Decode(object)
}

// Store only the HTTP status, for no content responses
func (s *Context) Response(httpStatus int) {
	s.ResponseHTTPStatus = httpStatus
}

// Store a message response, translating the message id to the proper language message
func (s *Context) MessageResponse(httpStatus int, messageId string) {
	s.ResponseHTTPStatus = httpStatus
	s.ResponseContent = []byte(s.Language.Messages[messageId])
}

// Store a object in json format for the response
func (s *Context) JSONResponse(httpStatus int, object interface{}) error {
	content, err := json.Marshal(object)

	if err != nil {
		return err
	}

	s.ResponseHTTPStatus = httpStatus
	s.ResponseContent = content
	return nil
}

// Add a custom HTTP header. Used for some types of response where you need to set ETag or
// LastModified fields
func (s *Context) AddHeader(key, value string) {
	// Avoid adding headers that are automatically generated at the end of the request. We
	// don't allow header overwrite because in the low level MIMEHeader the HTTP header
	// value is appended instead of replaced
	if key == "Content-Type" ||
		key == "Content-Language" ||
		key == "Content-Length" ||
		key == "Content-MD5" ||
		key == "Accept" ||
		key == "Accept-Language" ||
		key == "Accept-Charset" ||
		key == "Date" {
		return
	}

	s.HTTPHeader[key] = value
}

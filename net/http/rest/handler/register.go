package handler

import (
	"net/http"
	"shelter/net/http/rest/context"
)

// Routes is responsable for storing the link beteween an URI and a handler
var (
	Routes map[string]Handler
)

// Create this type to make it easy to reference a handler
type Handler func(*http.Request, *context.ShelterRESTContext)

// Function created only to register the handlers more easily
func HandleFunc(route string, handler Handler) {
	Routes[route] = handler
}

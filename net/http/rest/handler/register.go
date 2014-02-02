package handler

import (
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"net/http"
	"regexp"
)

// Routes is responsable for storing the link beteween an URI and a handler
var (
	Routes map[*regexp.Regexp]Handler
)

// Create this type to make it easy to reference a handler
type Handler func(*http.Request, *context.Context)

// Function created only to register the handlers more easily
func HandleFunc(routeRegexp *regexp.Regexp, handler Handler) {
	// We are initializing the router here because if we do this in a init function there's no
	// garantee that the function will be called before the init functions of the other handlers
	if Routes == nil {
		Routes = make(map[*regexp.Regexp]Handler)
	}

	Routes[routeRegexp] = handler
}

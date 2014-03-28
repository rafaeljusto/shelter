// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package handler

import (
	"net/http"
	"regexp"
)

// Routes is responsable for storing the link beteween an URI and a handler. It uses
// regular expression to match the URI because it's faster and allows more complex URI
// matchs, like /domain/<something>/verification
var (
	Routes map[*regexp.Regexp]Handler
)

// Create this type to make it easy to reference a handler
type Handler func(w http.ResponseWriter, r *http.Request)

// Function created only to register the handlers more easily
func HandleFunc(routeRegexp *regexp.Regexp, handler Handler) {
	// We are initializing the router here because if we do this in a init function there's
	// no garantee that the function will be called before the init functions of the other
	// handlers
	if Routes == nil {
		Routes = make(map[*regexp.Regexp]Handler)
	}

	Routes[routeRegexp] = handler
}

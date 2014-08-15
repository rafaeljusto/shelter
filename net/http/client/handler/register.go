// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package handler store the web client handlers of specific URI
package handler

import (
	"github.com/rafaeljusto/handy"
)

// Routes is responsable for storing the link beteween an URI and a handler. It uses a
// library to match the URI because it's faster and allows more complex URI matchs, like
// /domain/<something>/verification
var (
	Routes map[string]func() handy.Handler
)

// Function created only to register the handlers more easily
func HandleFunc(pattern string, handler func() handy.Handler) {
	// We are initializing the router here because if we do this in a init function there's no
	// garantee that the function will be called before the init functions of the other handlers
	if Routes == nil {
		Routes = make(map[string]func() handy.Handler)
	}

	Routes[pattern] = handler
}

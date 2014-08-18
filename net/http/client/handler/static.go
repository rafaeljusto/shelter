// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package handler store the web client handlers of specific URI
package handler

import (
	"github.com/rafaeljusto/handy"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/net/http/client/interceptor"
	"net/http"
	"path/filepath"
)

var (
	staticHandler func(http.ResponseWriter, *http.Request)
)

func init() {
	HandleFunc("/", func() handy.Handler {
		return new(StaticHandler)
	})
}

func StartStaticHandler() {
	staticPath := filepath.Join(
		config.ShelterConfig.BasePath,
		config.ShelterConfig.WebClient.StaticPath,
	)

	staticHandler = http.FileServer(http.Dir(staticPath)).ServeHTTP
}

type StaticHandler struct {
	handy.DefaultHandler
}

func (h *StaticHandler) Get(w http.ResponseWriter, r *http.Request) {
	staticHandler(w, r)
}

func (h *StaticHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(new(interceptor.Permission))
}

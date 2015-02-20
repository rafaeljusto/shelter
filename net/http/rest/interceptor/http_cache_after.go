// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/rafaeljusto/handy/interceptor"
	"net/http"
)

type HTTPCacheAfter struct {
	interceptor.NoBeforeInterceptor
	httpCacheHandler HTTPCacheHandler
}

func NewHTTPCacheAfter(h HTTPCacheHandler) *HTTPCacheAfter {
	return &HTTPCacheAfter{httpCacheHandler: h}
}

func (i *HTTPCacheAfter) After(w http.ResponseWriter, r *http.Request) {
	CheckHTTPCache(w, r, i.httpCacheHandler)
}

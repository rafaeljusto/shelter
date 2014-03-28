// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package handler

import (
	"github.com/rafaeljusto/shelter/config"
	"net/http"
	"path/filepath"
)

var (
	StaticHandler Handler
)

func StartStaticHandler() {
	staticPath := filepath.Join(
		config.ShelterConfig.BasePath,
		config.ShelterConfig.WebClient.StaticPath,
	)

	StaticHandler = http.FileServer(http.Dir(staticPath)).ServeHTTP
}

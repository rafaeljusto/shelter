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

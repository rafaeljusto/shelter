// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package client is the web client service
package client

import (
	"crypto/tls"
	"github.com/rafaeljusto/handy"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/net/http/client/handler"
	"github.com/rafaeljusto/shelter/net/http/client/interceptor"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
)

// Function created to run the listeners before dropping privileges in the main binary, so
// that we can listen in low ports without keeping the program as a super user
func Listen() ([]net.Listener, error) {
	listeners := make([]net.Listener, 0, len(config.ShelterConfig.WebClient.Listeners))

	for _, v := range config.ShelterConfig.WebClient.Listeners {
		ipAndPort := net.JoinHostPort(v.IP, strconv.Itoa(v.Port))

		if v.TLS {
			certificatePath := filepath.Join(
				config.ShelterConfig.BasePath,
				config.ShelterConfig.WebClient.TLS.CertificatePath,
			)

			privateKeyPath := filepath.Join(
				config.ShelterConfig.BasePath,
				config.ShelterConfig.WebClient.TLS.PrivateKeyPath,
			)

			cert, err := tls.LoadX509KeyPair(certificatePath, privateKeyPath)
			if err != nil {
				return nil, err
			}

			tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}}

			ln, err := tls.Listen("tcp", ipAndPort, &tlsConfig)
			if err != nil {
				return nil, err
			}

			listeners = append(listeners, ln)

		} else {
			ln, err := net.Listen("tcp", ipAndPort)
			if err != nil {
				return nil, err
			}

			listeners = append(listeners, ln)
		}
	}

	return listeners, nil
}

func Start(listeners []net.Listener) error {
	// Initialize CIDR whitelist
	if err := loadACL(); err != nil {
		return err
	}

	// Static handler must be called directly because it uses configuration file values to
	// configure the root path. For that reason we can't build it in the init function
	handler.StartStaticHandler()

	// Handy logger should use the same logger of the Shelter system
	handy.Logger = log.Logger

	mux := handy.NewHandy()
	mux.Recover = func(r interface{}) {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		log.Printf("Web client panic detected. Details: %v\n%s", r, buf)
	}

	for pattern, handler := range handler.Routes {
		mux.Handle(pattern, handler)
	}

	server := http.Server{
		Handler: mux,
	}

	for _, v := range listeners {
		// We are not checking the error returned by Serve, because if we check for some
		// reason the HTTP server stop answering the requests
		go server.Serve(v)
	}

	return nil
}

func loadACL() error {
	for _, cidrStr := range config.ShelterConfig.RESTServer.ACL {
		if _, cidr, err := net.ParseCIDR(cidrStr); err == nil {
			interceptor.ACL = append(interceptor.ACL, cidr)
		} else {
			return err
		}
	}

	return nil
}

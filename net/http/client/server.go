package client

import (
	"crypto/tls"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/net/http/client/handler"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
)

// Function created to run the listeners before dropping privileges in the main binary, so
// that we can listen in low ports without keeping the program as a super user
func Listen() ([]net.Listener, error) {
	listeners := make([]net.Listener, 0, len(config.ShelterConfig.ClientServer.Listeners))

	for _, v := range config.ShelterConfig.ClientServer.Listeners {
		ipAndPort := net.JoinHostPort(v.IP, strconv.Itoa(v.Port))

		if v.TLS {
			certificatePath := filepath.Join(
				config.ShelterConfig.BasePath,
				config.ShelterConfig.ClientServer.TLS.CertificatePath,
			)

			privateKeyPath := filepath.Join(
				config.ShelterConfig.BasePath,
				config.ShelterConfig.ClientServer.TLS.PrivateKeyPath,
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
	// Static handler must be called directly because it uses configuration file values to
	// configure the root path. For that reason we can't build it in the init function
	handler.StartStaticHandler()

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

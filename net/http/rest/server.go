package rest

import (
	"crypto/tls"
	"net"
	"net/http"
	"shelter/config"
	"shelter/net/http/rest/language"
	"strconv"
)

func Start() error {
	if err := language.LoadConfig(config.ShelterConfig.RESTServer.LanguageConfigPath); err != nil {
		return err
	}

	listeners := make([]net.Listener, 0, len(config.ShelterConfig.RESTServer.Listeners))

	for _, v := range config.ShelterConfig.RESTServer.Listeners {
		ipAndPort := net.JoinHostPort(v.IP, strconv.Itoa(v.Port))

		if v.TLS {
			cert, err := tls.LoadX509KeyPair(config.ShelterConfig.RESTServer.TLS.CertificatePath,
				config.ShelterConfig.RESTServer.TLS.PrivateKeyPath)

			if err != nil {
				return err
			}

			tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}}

			ln, err := tls.Listen("tcp", ipAndPort, &tlsConfig)
			if err != nil {
				return err
			}

			listeners = append(listeners, ln)

		} else {
			ln, err := net.Listen("tcp", ipAndPort)
			if err != nil {
				return err
			}

			listeners = append(listeners, ln)
		}

	}

	for _, v := range listeners {
		go http.Serve(v, mux)
	}

	return nil
}

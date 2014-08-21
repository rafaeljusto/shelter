// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package handler store the web client handlers of specific URI
package handler

import (
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"testing"
)

func TestRetrieveRESTAddress(t *testing.T) {
	if _, err := retrieveRESTAddress(); err == nil {
		t.Error("Not returning error when there's no address configured")
	}

	listeners := []struct {
		IP   string
		Port int
		TLS  bool
	}{
		{
			IP:   "127.0.0.1",
			Port: 443,
			TLS:  true,
		},
	}

	config.ShelterConfig = config.Config{
		RESTServer: struct {
			Enabled            bool
			LanguageConfigPath string

			TLS struct {
				CertificatePath string
				PrivateKeyPath  string
			}

			Listeners []struct {
				IP   string
				Port int
				TLS  bool
			}

			Timeouts struct {
				ReadSeconds  int
				WriteSeconds int
			}

			ACL     []string
			Secrets map[string]string
		}{
			Listeners: listeners,
		},
	}

	if address, err := retrieveRESTAddress(); err != nil {
		t.Error("Returning error when there's an address. Details:", err)
	} else if address != "https://[127.0.0.1]:443" {
		t.Error("Returning wrong address for https")
	}

	listeners = []struct {
		IP   string
		Port int
		TLS  bool
	}{
		{
			IP:   "::1",
			Port: 80,
			TLS:  false,
		},
	}

	config.ShelterConfig = config.Config{
		RESTServer: struct {
			Enabled            bool
			LanguageConfigPath string

			TLS struct {
				CertificatePath string
				PrivateKeyPath  string
			}

			Listeners []struct {
				IP   string
				Port int
				TLS  bool
			}

			Timeouts struct {
				ReadSeconds  int
				WriteSeconds int
			}

			ACL     []string
			Secrets map[string]string
		}{
			Listeners: listeners,
		},
	}

	if address, err := retrieveRESTAddress(); err != nil {
		t.Error("Returning error when there's an address. Details:", err)
	} else if address != "http://[::1]:80" {
		t.Error("Returning wrong address for http")
	}
}

func TestSignAndSend(t *testing.T) {
	r, err := http.NewRequest("GET", "http://127.0.0.1/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := signAndSend(r, nil); err != ErrNoSecretFound {
		t.Error("Not throwing error when there's no secret")
	}

	config.ShelterConfig = config.Config{
		RESTServer: struct {
			Enabled            bool
			LanguageConfigPath string

			TLS struct {
				CertificatePath string
				PrivateKeyPath  string
			}

			Listeners []struct {
				IP   string
				Port int
				TLS  bool
			}

			Timeouts struct {
				ReadSeconds  int
				WriteSeconds int
			}

			ACL     []string
			Secrets map[string]string
		}{
			Secrets: map[string]string{
				"key01": "ohV43/9bKlVNaXeNTqEuHQp57LCPCQ==",
			},
		},
	}

	if _, err := signAndSend(r, nil); err != nil {
		// Avoid connection error
		switch specificErr := err.(type) {
		case *net.OpError:
			if specificErr.Op != "read" {
				t.Error(err)
			}

		case syscall.Errno:
			if specificErr != syscall.ECONNREFUSED {
				t.Error(err)
			}

		case *url.Error:
			if specificErr.Op != "Get" {
				t.Error(err)
			}

		default:
			t.Error(err)
		}
	}

	if len(r.Header.Get("Date")) == 0 || len(r.Header.Get("Authorization")) == 0 {
		t.Error("Not setting the correct HTTP headers")
	}

	content := []byte("This is a test")
	r, err = http.NewRequest("GET", "http://127.0.0.1/test",
		strings.NewReader(string(content)))

	if err != nil {
		t.Fatal(err)
	}

	if _, err := signAndSend(r, content); err != nil {
		// Avoid connection error
		switch specificErr := err.(type) {
		case *net.OpError:
			if specificErr.Op != "read" {
				t.Error(err)
			}

		case syscall.Errno:
			if specificErr != syscall.ECONNREFUSED {
				t.Error(err)
			}

		case *url.Error:
			if specificErr.Op != "Get" {
				t.Error(err)
			}

		default:
			t.Error(err)
		}
	}

	if len(r.Header.Get("Date")) == 0 ||
		len(r.Header.Get("Authorization")) == 0 ||
		r.Header.Get("Content-Type") != check.SupportedContentType ||
		len(r.Header.Get("Content-MD5")) == 0 {

		t.Error("Not setting the correct HTTP headers")
	}
}

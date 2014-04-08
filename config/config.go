// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package config defines the Shelter configuration parameters
package config

import (
	"encoding/json"
	"io/ioutil"
)

var (
	// ShelterConfig is the instance of the system configuration that all modules will use.
	// It doesn't use a locker strategy because is a read only structure. The main binary is
	// responsable for initializing this global variable
	ShelterConfig Config
)

// List of possible e-mail authentication types
const (
	AuthenticationTypeNone        AuthenticationType = ""            // No authentication
	AuthenticationTypePlain       AuthenticationType = "PLAIN"       // Basic user/password authentication
	AuthenticationTypeCRAMMD5Auth AuthenticationType = "CRAMMD5AUTH" // HMAC-MD5 challenge-response authentication
)

// AuthenticationType is the text that represents the authentication type
type AuthenticationType string

// Config structure describes all the configuration variables used in the Shelter system
type Config struct {
	// Base path of the system, all other paths will prepend this path. This is useful to
	// isolate the project in a specific directory
	BasePath string

	// Path and name of the file that will receive all the messages from the Shelter system
	LogFilename string

	// List of supported languages. When adding a new language, remember to respect the IANA
	// Language Subtag Registry file to define the name (e.g. en-US, pt-BR) and to add the
	// necessary information for the new language (REST messages, notification template and
	// web client languages)
	Languages []string

	// Stores the database information to start a connection. For now we are using MongoDB
	// for the persistence layer
	Database struct {
		// Name of the database
		Name string

		// URI used by MongoDB to connect to the database
		URI string
	}

	// Store all variables related to a scan job in the Shelter system
	Scan struct {
		Enabled           bool
		Time              string
		IntervalHours     int
		NumberOfQueriers  int
		DomainsBufferSize int
		ErrorsBufferSize  int
		UDPMaxSize        uint16
		SaveAtOnce        int
		ConnectionRetries int

		Resolver struct {
			Address string
			Port    int
		}

		Timeouts struct {
			DialSeconds  int
			ReadSeconds  int
			WriteSeconds int
		}

		VerificationIntervals struct {
			MaxOKDays              int
			MaxErrorDays           int
			MaxExpirationAlertDays int
		}
	}

	RESTServer struct {
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
	}

	WebClient struct {
		Enabled    bool
		StaticPath string

		TLS struct {
			CertificatePath string
			PrivateKeyPath  string
		}

		Listeners []struct {
			IP   string
			Port int
			TLS  bool
		}
	}

	Notification struct {
		Enabled                    bool
		Time                       string
		IntervalHours              int
		NameserverErrorAlertDays   int
		NameserverTimeoutAlertDays int
		DSErrorAlertDays           int
		DSTimeoutAlertDays         int
		From                       string

		// Define the path that has the template files. Each template file must have the
		// filename related to the language that it uses in lowercase (e.g. en-us.tmpl, pt-
		// br.tmpl). The basic structure of each template should be as described bellow. The
		// mail header and the parameters beteween "{{" and "}}" must not be removed, because
		// they are used to build the basic structure of the notification.
		//
		//     {{$domain := .}}
		//
		//     From: {{.From}}
		//     To: {{.To}}
		//     Subject: Misconfiguration on domain {{$domain.FQDN}}
		//
		//     Dear Sir/Madam,
		//
		//     During our periodically domain verification, a configuration problem was
		//     detected with the domain {{$domain.FQDN}}.
		//
		//     {{range $nameserver := $domain.Nameservers}}
		//       {{if nsStatusEq $nameserver.LastStatus "TIMEOUT"}}
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "NOAA"}}
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "UDN"}}
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "UH"}}
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "SERVFAIL"}}
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "QREFUSED"}}
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "CREFUSED"}}
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "CNAME"}}
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "NOTSYNCH"}}
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "ERROR"}}
		//
		//       {{end}}
		//     {{end}}
		//
		//     {{range $ds := $domain.DSSet}}
		//       {{if dsStatusEq $ds.LastStatus "TIMEOUT"}}
		//
		//       {{else if dsStatusEq $ds.LastStatus "NOSIG"}}
		//
		//       {{else if dsStatusEq $ds.LastStatus "EXPSIG"}}
		//
		//       {{else if dsStatusEq $ds.LastStatus "NOKEY"}}
		//
		//       {{else if dsStatusEq $ds.LastStatus "NOSEP"}}
		//
		//       {{else if dsStatusEq $ds.LastStatus "SIGERR"}}
		//
		//       {{else if dsStatusEq $ds.LastStatus "DNSERR"}}
		//
		//       {{else if isNearExpiration $ds}}
		//
		//       {{end}}
		//     {{end}}
		//
		//     Best regards,
		//     LACTLD
		//
		// You can also use other variables in the template file, like {{$nameserver.Host}} or
		// {{$ds.Keytag}} to create better user messages for the current scenario.
		TemplatesPath string

		SMTPServer struct {
			Server string
			Port   int

			Auth struct {
				Type     AuthenticationType
				Username string
				Password string
			}
		}
	}
}

// LoadConfig reads the system file and try to build the global configuration instance
// using the data
func LoadConfig(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &ShelterConfig); err != nil {
		return err
	}

	return nil
}

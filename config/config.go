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
		// Flag to enable or disable the scan module. Even if the module is disabled is a good
		// idea to fill some information, because other modules can use it (like the
		// VerificationIntervals for the notification module)
		Enabled bool

		// Time of the day that you want that the scan to run (e.g. 05:00:00 BRT or 08:00:00
		// UTC). The scan will not occur exactly in the given time because of the scheduler
		// architecture, it will probably have a delay of some minutes.
		Time string

		// Number of hours between each scan
		IntervalHours int

		// Number of parallel workers that will be sending queries to the registered
		// nameservers. Remember that ideal number of queriers is defined by the hardware that
		// you have
		NumberOfQueriers int

		// Size of the buffer that process the domains, retrieving from database, sending
		// requests and saving into database. This should be defined thinking on the number of
		// queriers, because if there're many queries they will empty the buffer faster
		DomainsBufferSize int

		// Size of the buffer that store errors of the scan
		ErrorsBufferSize int

		// The maximum size of a UDP package that your network can receive without being
		// droppped by a firewall. Please, check EDNS0 for more information
		UDPMaxSize uint16

		// Number of domains to accumulate before saving into database, this save us IO time
		// and should be defined according to the number of queriers
		SaveAtOnce int

		// Number of times that a querier will retry to send the DNS request to the host.
		// After that will consider a timeout problem
		ConnectionRetries int

		// Information about the recursive DNS server for specific services of the scan. Like
		// QueryDomain, that retrieves the nameservers and DS records from a domain name
		Resolver struct {
			// IP address from the resolver
			Address string

			// Port from the resolver
			Port int
		}

		// Timeouts define the number of seconds that the system will wait for network
		// operations
		Timeouts struct {
			// Number of seconds that the system will wait while trying to connect to a remote
			// host, for example, in a DNS request via TCP
			DialSeconds int

			// Number of seconds that the system will wait for a response from the remote host
			ReadSeconds int

			// Number of seconds that the system will wait to receive an ack for the message
			// written to a remote host via TCP
			WriteSeconds int
		}

		// Intervals days used for selecting registered domains to be scanned depending on the
		// current domain state
		VerificationIntervals struct {
			// Maximum number of days that a well configured domain will wait until it is
			// scanned again
			MaxOKDays int

			// Maximum number of days that a domain with DNS or DNSSEC configuration problems
			// will wait until it is scanned again
			MaxErrorDays int

			// Maximum number of days before the DNSSEC expiration date that a domain will be
			// verified for DNS or DNSSEC problems, and mainly to check if it was already
			// resigned
			MaxExpirationAlertDays int
		}
	}

	// Store all variables related to the REST server
	RESTServer struct {
		// Flag to enable or disable the REST server module. Even if the module is disabled is
		// a good idea to fill the ACL and Secrets (for the WebClient module)
		Enabled bool

		// Path and filename of the REST messages file (messages.conf) that store all the
		// messages that the REST server can return to a client on the desired language. The
		// format of this file must follow JSON syntax and the following structure:
		//
		//     {
		//       "default": "en-us",
		//       "packs": [
		//         {
		//           "GenericName": "en",
		//           "SpecificName": "en-us",
		//           "Messages": {
		//             "label1": "message1",
		//             "label2": "message2",
		//             ...
		//             "labelN": "messageN"
		//         }
		//     }
		LanguageConfigPath string

		// TLS store the necessary data to allow an encrypted connection
		TLS struct {
			// X509 certificate (.pem) file
			CertificatePath string

			// X509 private key (.pem) file
			PrivateKeyPath string
		}

		// Addresses that the REST server will listen to
		Listeners []struct {
			// IP address that can be IPv4 or IPv6
			IP string

			// Port number that will be used
			Port int

			// Flag that indicates if it will use HTTPS or not. All interfaces will use the same
			// TLS certificates
			TLS bool
		}

		// Number of seconds that the REST server will wait for network operations
		Timeouts struct {
			// Number of seconds that the system will wait for a response from the remote host
			ReadSeconds int

			// Number of seconds that the system will wait to receive an ack for the message
			// written to a remote host
			WriteSeconds int
		}

		// List of IP ranges that are allowed to connect to the REST server, all other will
		// receive a FORBIDDEN HTTP status code
		ACL []string

		// Store the shared secret keys used by the clients to sign the requests
		Secrets map[string]string
	}

	// Store all necessary information for the web client
	WebClient struct {
		// Flag to enable or disable the web client module
		Enabled bool

		// Base path where all html static content is found, like index.html, CSS and
		// javascripts
		StaticPath string

		// TLS store the necessary data to allow an encrypted connection
		TLS struct {
			// X509 certificate (.pem) file
			CertificatePath string

			// X509 private key (.pem) file
			PrivateKeyPath string
		}

		// Addresses that the web client will listen to
		Listeners []struct {
			// IP address that can be IPv4 or IPv6
			IP string

			// Port number that will be used
			Port int

			// Flag that indicates if it will use HTTPS or not. All interfaces will use the same
			// TLS certificates
			TLS bool
		}
	}

	// Store all notification related variables
	Notification struct {
		// Flag to enable or disable the notification module
		Enabled bool

		// Time of the day that you want that the notification to run (e.g. 06:00:00 BRT or
		// 09:00:00 UTC). The notification will not occur exactly in the given time because of
		// the scheduler architecture, it will probably have a delay of some minutes. It is
		// recommended to run the notification after the scan
		Time string

		// Number of hours between each notification
		IntervalHours int

		// How many days we will wait with a DNS misconfigured nameserver until we notify the
		// domain's owners
		NameserverErrorAlertDays int

		// How many days we will wait with a unresponsive nameserver until we notify the
		// domain's owners
		NameserverTimeoutAlertDays int

		// How many days we will wait with a DNSSEC misconfigured nameserver until we notify the
		// domain's owners
		DSErrorAlertDays int

		// How many days we will wait with a unresponsive nameserver for DNSSEC queries until
		// we notify the domain's owners
		DSTimeoutAlertDays int

		// All notification e-mails are sent with this From
		From string

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
		//     Subject: Title
		//
		//     Some introduction message.
		//
		//     {{range $nameserver := $domain.Nameservers}}
		//       {{if nsStatusEq $nameserver.LastStatus "TIMEOUT"}}
		//         Error description.
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "NOAA"}}
		//         Error description.
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "UDN"}}
		//         Error description.
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "UH"}}
		//         Error description.
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "SERVFAIL"}}
		//         Error description.
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "QREFUSED"}}
		//         Error description.
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "CREFUSED"}}
		//         Error description.
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "CNAME"}}
		//         Error description.
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "NOTSYNCH"}}
		//         Error description.
		//
		//       {{else if nsStatusEq $nameserver.LastStatus "ERROR"}}
		//         Error description.
		//
		//       {{end}}
		//     {{end}}
		//
		//     {{range $ds := $domain.DSSet}}
		//       {{if dsStatusEq $ds.LastStatus "TIMEOUT"}}
		//         Error description.
		//
		//       {{else if dsStatusEq $ds.LastStatus "NOSIG"}}
		//         Error description.
		//
		//       {{else if dsStatusEq $ds.LastStatus "EXPSIG"}}
		//         Error description.
		//
		//       {{else if dsStatusEq $ds.LastStatus "NOKEY"}}
		//         Error description.
		//
		//       {{else if dsStatusEq $ds.LastStatus "NOSEP"}}
		//         Error description.
		//
		//       {{else if dsStatusEq $ds.LastStatus "SIGERR"}}
		//         Error description.
		//
		//       {{else if dsStatusEq $ds.LastStatus "DNSERR"}}
		//         Error description.
		//
		//       {{else if isNearExpiration $ds}}
		//         Error description.
		//
		//       {{end}}
		//     {{end}}
		//
		//     Goodbye message.
		//
		// You can also use other variables in the template file, like {{$nameserver.Host}} or
		// {{$ds.Keytag}} to create better user messages for the current scenario.
		TemplatesPath string

		// Store all necessary information to send notification e-mails using an SMTP server
		SMTPServer struct {
			// Name or IP address of the SMTP server
			Server string

			// Port of the SMTP server
			Port int

			// Authentication information of the SMTP server
			Auth struct {
				// Type of authentication, that can be empty, "PLAIN" or "CRAMMD5AUTH"
				Type AuthenticationType

				// Username used for authentication
				Username string

				// Passowrd used for authentication
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

package config

import (
	"encoding/json"
	"io/ioutil"
)

var (
	ShelterConfig Config
)

const (
	AuthenticationTypeNone        AuthenticationType = ""
	AuthenticationTypePlain       AuthenticationType = "PLAIN"
	AuthenticationTypeCRAMMD5Auth AuthenticationType = "CRAMMD5AUTH"
)

type AuthenticationType string

type Config struct {
	BasePath    string
	LogFilename string
	Languages   []string

	Database struct {
		Name string
		URI  string
	}

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

	Notification struct {
		Enabled                    bool
		Time                       string
		IntervalHours              int
		NameserverErrorAlertDays   int
		NameserverTimeoutAlertDays int
		DSErrorAlertDays           int
		DSTimeoutAlertDays         int
		From                       string
		TemplatesPath              string

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

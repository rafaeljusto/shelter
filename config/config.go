package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

var (
	ShelterConfig Config
)

type Config struct {
	Database struct {
		Name string
		URI  string
	}

	Scan struct {
		NumberOfQueriers  int
		DomainsBufferSize int
		ErrorsBufferSize  int
		UDPMaxSize        uint16
		SaveAtOnce        int
		ConnectionRetries int

		Timeouts struct {
			DialSeconds  time.Duration
			ReadSeconds  time.Duration
			WriteSeconds time.Duration
		}

		VerificationIntervals struct {
			MaxOKDays              int
			MaxErrorDays           int
			MaxExpirationAlertDays int
		}
	}

	RESTServer struct {
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
	}

	Log struct {
		BasePath     string
		ScanFilename string
		RESTFilename string
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

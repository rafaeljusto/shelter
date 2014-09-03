// shelter - DNS/DNSSEC misconfiguration checker
//
// Copyright (C) 2014 Rafael Dantas Justo <adm@rafael.net.br>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

package main

import (
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/client"
	"github.com/rafaeljusto/shelter/net/http/rest"
	"github.com/rafaeljusto/shelter/net/mail/notification"
	"github.com/rafaeljusto/shelter/net/scan"
	"github.com/rafaeljusto/shelter/scheduler"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

const (
	copyright string = `Shelter version 0.3, Copyright (C) 2014 Rafael Dantas Justo
Shelter comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome
to redistribute it under certain conditions.
`
)

// List of arguments that can be filled in the program command line
var (
	configFilePath string // General configuration path
	showVersion    *bool  // Show system version
)

// We store all listeners to make it easier later to stop all in a system SIGTERM event
var (
	restListeners   []net.Listener
	clientListeners []net.Listener
)

// List of possible return codes of the program. This will be useful later to build a
// command line documentation
const (
	NoError = iota
	ErrInputParameters
	ErrLoadingConfig
	ErrListeningRESTInterfaces
	ErrStartingRESTServer
	ErrListeningClientInterfaces
	ErrStartingWebClient
	ErrScanTimeFormat
	ErrCurrentScanInitialize
	ErrNotificationTemplates
)

// We are going to use the initialization function to read command line arguments, load
// the configuration file and register system signals
func init() {
	flag.StringVar(&configFilePath, "config", "", "Configuration file")
	showVersion = flag.Bool("version", false, "System version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\nUsage of %s:\n", copyright, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion {
		fmt.Println(copyright)
		os.Exit(0)
	}

	if len(configFilePath) == 0 {
		fmt.Println("The configuration file was not informed")
		os.Exit(ErrInputParameters)
	}

	if err := loadSettings(); err != nil {
		log.Println("Error loading the configuration file. Details:", err)
		os.Exit(ErrLoadingConfig)
	}

	manageSystemSignals()
}

// The main function of the system is responsable for deploying all system components that
// are enabled. For now we have the REST server and the scan system
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	logPath := filepath.Join(
		config.ShelterConfig.BasePath,
		config.ShelterConfig.LogFilename,
	)

	if err := log.SetOutput(logPath); err != nil {
		log.Println(err)
		return
	}
	defer log.Close()

	if config.ShelterConfig.RESTServer.Enabled {
		var err error
		restListeners, err = rest.Listen()
		if err != nil {
			log.Println("Error while aquiring interfaces for REST server. Details:", err)
			os.Exit(ErrListeningRESTInterfaces)
		}

		if err := rest.Start(restListeners); err != nil {
			log.Println("Error starting the REST server. Details:", err)
			os.Exit(ErrStartingRESTServer)
		}
	}

	if config.ShelterConfig.WebClient.Enabled {
		var err error
		clientListeners, err = client.Listen()
		if err != nil {
			log.Println("Error while aquiring interfaces for Client server. Details:", err)
			os.Exit(ErrListeningClientInterfaces)
		}

		if err := client.Start(clientListeners); err != nil {
			log.Println("Error starting the Client server. Details:", err)
			os.Exit(ErrStartingWebClient)
		}
	}

	if config.ShelterConfig.Scan.Enabled {
		// Attention: Cannot use timezone abbreviations
		// http://stackoverflow.com/questions/25368415/golang-timezone-parsing
		scanTime, err := time.Parse("15:04:05 -0700", config.ShelterConfig.Scan.Time)
		if err != nil {
			log.Println("Scan time not in a valid format. Details:", err)
			os.Exit(ErrScanTimeFormat)
		}

		scanTime = scanTime.UTC()
		now := time.Now().UTC()

		nextExecution := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			scanTime.Hour(),
			scanTime.Minute(),
			scanTime.Second(),
			scanTime.Nanosecond(),
			scanTime.Location(),
		)

		scheduler.Register(scheduler.Job{
			Type:          scheduler.JobTypeScan,
			NextExecution: nextExecution,
			Interval:      time.Duration(config.ShelterConfig.Scan.IntervalHours) * time.Hour,
			Task:          scan.ScanDomains,
		})

		// Must be called after registering in scheduler, because we retrieve the next execution time
		// from it
		if err := model.InitializeCurrentScan(); err != nil {
			log.Println("Current scan information got an error while initializing. Details:", err)
			os.Exit(ErrCurrentScanInitialize)
		}
	}

	if config.ShelterConfig.Notification.Enabled {
		if err := notification.LoadTemplates(); err != nil {
			log.Println("Error loading notification templates. Details:", err)
			os.Exit(ErrNotificationTemplates)
		}

		// Attention: Cannot use timezone abbreviations
		// http://stackoverflow.com/questions/25368415/golang-timezone-parsing
		notificationTime, err := time.Parse("15:04:05 -0700", config.ShelterConfig.Scan.Time)
		if err != nil {
			log.Println("Scan time not in a valid format. Details:", err)
			os.Exit(ErrScanTimeFormat)
		}

		notificationTime = notificationTime.UTC()
		now := time.Now().UTC()

		nextExecution := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			notificationTime.Hour(),
			notificationTime.Minute(),
			notificationTime.Second(),
			notificationTime.Nanosecond(),
			notificationTime.Location(),
		)

		scheduler.Register(scheduler.Job{
			Type:          scheduler.JobTypeNotification,
			NextExecution: nextExecution,
			Interval:      time.Duration(config.ShelterConfig.Notification.IntervalHours) * time.Hour,
			Task:          notification.Notify,
		})
	}

	scheduler.Start()

	select {}
}

// Shelter could receive system signals for OS, so this method catch the signals to create
// smothly actions for each one. For example, when receives a KILL signal, we should wait
// to process all requests before finishing the server
func manageSystemSignals() {
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGHUP)

		for {
			sig := <-sigs

			if sig == syscall.SIGHUP {
				if err := loadSettings(); err != nil {
					log.Println("Error reloading confirguration file. Details:", err)
				}

			} else if sig == syscall.SIGTERM {
				for _, listener := range restListeners {
					if err := listener.Close(); err != nil {
						log.Println("Error closing listener. Details:", err)
					}
				}
				restListeners = []net.Listener{}

				for _, listener := range clientListeners {
					if err := listener.Close(); err != nil {
						log.Println("Error closing listener. Details:", err)
					}
				}
				clientListeners = []net.Listener{}

				// TODO: Wait the last requests to be processed? One possibly solution is to
				// create a request counter in MUX, we wait while this counter is non-zero. If
				// there's a scan running what are we going to do?

				os.Exit(NoError)
			}
		}
	}()
}

// loadSettings function is responsable for lading the configuration parameters from a
// file. It will be used when the system starts for the first time and when it receives a
// SIGHUP signal
func loadSettings() error {
	// TODO: Possible concurrent access problem while reloading the configuration file. And
	// we also should reload many structures that could change with the new configuration
	// files, like the network interfaces
	if err := config.LoadConfig(configFilePath); err != nil {
		return err
	}

	// Load languages to model. We don't do this in the configuration package, because we
	// don't want to create a dependency between the model and the config taht could become
	// a cross reference
	for _, language := range config.ShelterConfig.Languages {
		if !model.AddLanguage(language) {
			return fmt.Errorf("Language '%s' is not valid", language)
		}
	}

	return nil
}

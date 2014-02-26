package main

import (
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
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

// We store all listeners to make it easier later to stop all in a system SIGTERM event
var (
	restListeners []net.Listener
)

// List of possible return codes of the program. This will be useful later to build a
// command line documentation
const (
	NoError = iota
	ErrInputParameters
	ErrLoadingConfig
	ErrListeningRESTInterfaces
	ErrStartingRESTServer
	ErrScanTimeFormat
	ErrCurrentScanInitialize
	ErrNotificationTemplates
)

// We are going to use the initialization function to read command line arguments, load
// the configuration file and register system signals
func init() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("Usage: %s <configuration file>\n", os.Args[0])
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

	if config.ShelterConfig.Scan.Enabled {
		scanTime, err := time.Parse("15:04:05 MST", config.ShelterConfig.Scan.Time)
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

		notificationTime, err := time.Parse("15:04:05 MST", config.ShelterConfig.Scan.Time)
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

				// TODO: Wait the last requests to be processed? On epossibly solution is to
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
	return config.LoadConfig(flag.Arg(0))
}

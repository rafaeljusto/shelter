package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"shelter/config"
	"shelter/net/http/rest"
	"syscall"
)

// We store all listeners to make it easier later to stop all in a system signal event
var (
	restListeners []net.Listener
)

func init() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("Usage: %s <configuration file>\n", os.Args[0])
		os.Exit(1)
	}

	if err := loadSettings(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if config.ShelterConfig.RESTServer.Enabled {
		var err error
		restListeners, err = rest.Listen()
		if err != nil {
			fmt.Println("Error while aquiring interfaces for REST server. Details:", err)
			os.Exit(1)
		}
	}

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGHUP)

		for {
			sig := <-sigs

			if sig == syscall.SIGHUP {
				if err := loadSettings(); err != nil {
					// TODO!
				}

			} else if sig == syscall.SIGTERM {
				for _, listener := range restListeners {
					if err := listener.Close(); err != nil {
						// TODO!
					}
				}
				os.Exit(0)
			}
		}
	}()

	// Apparently we can't drop privileges because of the Go routines
	// if err := changePrivileges(); err != nil {
	// 	fmt.Println("Error changing process privileges. Details:", err)
	// }

	if config.ShelterConfig.RESTServer.Enabled {
		if err := rest.Start(restListeners); err != nil {
			fmt.Println("Error starting the REST server. Details:", err)
		}
	}
}

func changePrivileges() error {
	if err := syscall.Setuid(config.ShelterConfig.Shelter.UID); err != nil {
		return err
	}

	return nil
}

func loadSettings() error {
	return config.LoadConfig(flag.Arg(0))
}

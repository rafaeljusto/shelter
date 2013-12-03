package service

import (
	"fmt"
	"log"
	"os"
	"shelter/config"
	"shelter/database/mongodb"
	"shelter/net/scan"
	"sync"
)

func ScanDomains() {
	scanLogPath := fmt.Sprintf("%s/%s",
		config.ShelterConfig.Log.BasePath,
		config.ShelterConfig.Log.ScanFilename,
	)

	scanLog, err := os.Open(scanLogPath)
	if err != nil {
		log.Println(err)
		return
	}

	logger := log.New(scanLog, "", log.LstdFlags)

	database, err := mongodb.Open(
		config.ShelterConfig.Database.URI,
		config.ShelterConfig.Database.Name,
	)

	if err != nil {
		logger.Println(err)
		return
	}

	injector := scan.NewInjector(
		database,
		config.ShelterConfig.Scan.DomainsBufferSize,
		config.ShelterConfig.Scan.VerificationIntervals.MaxOKDays,
		config.ShelterConfig.Scan.VerificationIntervals.MaxErrorDays,
		config.ShelterConfig.Scan.VerificationIntervals.MaxExpirationAlertDays,
	)

	querierDispatcher := scan.NewQuerierDispatcher(
		config.ShelterConfig.Scan.NumberOfQueriers,
		config.ShelterConfig.Scan.DomainsBufferSize,
		config.ShelterConfig.Scan.UDPMaxSize,
		config.ShelterConfig.Scan.Timeouts.DialSeconds,
		config.ShelterConfig.Scan.Timeouts.ReadSeconds,
		config.ShelterConfig.Scan.Timeouts.WriteSeconds,
	)

	collector := scan.NewCollector(
		database,
		config.ShelterConfig.Scan.SaveAtOnce,
	)

	var scanGroup sync.WaitGroup
	errorsChannel := make(chan error, config.ShelterConfig.Scan.ErrorsBufferSize)
	domainsToQueryChannel := injector.Start(&scanGroup, errorsChannel)
	domainsToSaveChannel := querierDispatcher.Start(&scanGroup, domainsToQueryChannel)
	collector.Start(&scanGroup, domainsToSaveChannel, errorsChannel)

	go func() {
		for {
			select {
			case err := <-errorsChannel:
				logger.Println(err)
			}
		}
	}()

	scanGroup.Wait()
}

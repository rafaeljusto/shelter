package service

import (
	"fmt"
	"log"
	"os"
	"shelter/config"
	"shelter/database/mongodb"
	"shelter/net/scan"
	"sync"
	"time"
)

// Function responsable for running the domain scan system, checking the configuration of
// each domain in the database according to an algorithm. This method is synchronous and
// will return only after the scan proccess is done
func ScanDomains() {
	scanLogPath := fmt.Sprintf("%s/%s",
		config.ShelterConfig.Log.BasePath,
		config.ShelterConfig.Log.ScanFilename,
	)

	scanLog, err := os.Create(scanLogPath)
	if err != nil {
		log.Println(err)
		return
	}
	defer scanLog.Close()

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
		config.ShelterConfig.Scan.Timeouts.DialSeconds*time.Second,
		config.ShelterConfig.Scan.Timeouts.ReadSeconds*time.Second,
		config.ShelterConfig.Scan.Timeouts.WriteSeconds*time.Second,
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
				// Detect the poison pill to finish the error listener go routine. This poison
				// pill should be sent after all parts of the scan are done and we are sure that
				// we don't have any error to log anymore
				if err == nil {
					return
				} else {
					logger.Println(err)
				}
			}
		}
	}()

	// Wait for all parts of the scan to finish their job
	scanGroup.Wait()

	// Finish the error listener sending a poison pill
	errorsChannel <- nil
}

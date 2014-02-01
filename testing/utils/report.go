package utils

import (
	"io/ioutil"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// If we found a report file in the current path, rename it so we don't lose the old data.
// We are going to use the modification date from the file. We also don't check the errors
// because we really don't care
func WriteReport(reportFile string, report string) {
	if file, err := os.Open(reportFile); err == nil {
		newFilename := reportFile + ".old-"

		if fileStatus, err := file.Stat(); err == nil {
			newFilename += fileStatus.ModTime().Format("20060102150405")

		} else {
			// Did not find the modification date, so lets use now
			newFilename += time.Now().Format("20060102150405")
		}

		// We don't use defer because we want to rename it before the end of scope
		file.Close()

		os.Rename(reportFile, newFilename)
	}

	ioutil.WriteFile(reportFile, []byte(report), 0444)
}

// Initialize CPU profile. This function will return a function that MUST be defered to the end of
// the report test. If something goes wrong, the funciton will abort the program
func StartCPUProfile(name string) func() {
	profileFile, err := os.Create(name)
	if err != nil {
		Fatalln("Error creating CPU profile file", err)
	}

	if err := pprof.StartCPUProfile(profileFile); err != nil {
		Fatalln("Error starting CPU profile file", err)
	}

	return pprof.StopCPUProfile
}

// Initialize Memory profile. This function will return a function that MUST be defered to the end
// of the report test. If something goes wrong, the funciton will abort the program
func StartMemoryProfile(name string) func() {
	return func() {
		runtime.GC()

		profileFile, err := os.Create(name)
		if err != nil {
			Fatalln("Error creating memory profile file", err)
		}

		if err := pprof.Lookup("heap").WriteTo(profileFile, 1); err != nil {
			Fatalln("Error writing to memory profile file", err)
		}
		profileFile.Close()
	}
}

// Initialize Go Routines profile. This function will return a function that MUST be defered to the
// end of the report test. If something goes wrong, the funciton will abort the program
func StartGoRoutinesProfile(name string) func() {
	return func() {
		profileFile, err := os.Create(name)
		if err != nil {
			Fatalln("Error creating Go routines profile file", err)
		}

		if err := pprof.Lookup("goroutine").WriteTo(profileFile, 2); err != nil {
			Fatalln("Error writing to Go routines profile file", err)
		}
		profileFile.Close()
	}
}

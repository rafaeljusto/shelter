package utils

import (
	"io/ioutil"
	"os"
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

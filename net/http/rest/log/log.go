package log

import (
	"log"
	"os"
)

var (
	logger  *log.Logger // Main object used to control log messages
	logFile *os.File    // File used to log the messages
)

// Function responsable for creating the log file and initilizing the main object
// responsable for logging
func SetOutput(filename string) error {
	// For initilizing the log, we must first create the file, and if already exists, we
	// append the information on it. How are we going to rotate this file?
	var err error
	logFile, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	logger = log.New(logFile, "", log.Lshortfile|log.LstdFlags)
	return nil
}

// Basic log packages methods now writing to our local logger. If logger is not
// initialized we fallback to the standard log system
func Println(v ...interface{}) {
	if logger == nil {
		log.Println(v...)
	} else {
		logger.Println(v...)
	}
}

// Basic log packages methods now writing to our local logger. If logger is not
// initialized we fallback to the standard log system
func Printf(format string, v ...interface{}) {
	if logger == nil {
		log.Printf(format, v...)
	} else {
		logger.Printf(format, v...)
	}
}

// Close the log file and finish the logger
func Close() {
	logFile.Close()
}

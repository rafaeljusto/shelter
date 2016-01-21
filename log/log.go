// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package log is a centralized log of the Shelter system
package log

import (
	"log"
	"os"

	"github.com/rafaeljusto/shelter/config"
)

var (
	Logger  *log.Logger // Main object used to control log messages
	logFile *os.File    // File used to log the messages
)

// SetOutput function responsible for creating the log file and initializing the
// main object responsible for logging
func SetOutput(filename string) error {
	// For initializing the log, we must first create the file, and if already
	// exists, we append the information on it. How are we going to rotate this
	// file?
	var err error
	logFile, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	Logger = log.New(logFile, "", log.Lshortfile|log.LstdFlags)
	return nil
}

// Println basic log packages methods now writing to our local logger. If logger
// is not initialized we fallback to the standard log system
func Println(v ...interface{}) {
	label := "[normal]"
	v = append([]interface{}{label}, v...)

	if Logger == nil {
		log.Println(v...)
	} else {
		Logger.Println(v...)
	}
}

// Printf basic log packages methods now writing to our local logger. If logger
// is not initialized we fallback to the standard log system
func Printf(format string, v ...interface{}) {
	label := "[normal] "

	if Logger == nil {
		log.Printf(label+format, v...)
	} else {
		Logger.Printf(label+format, v...)
	}
}

// Info basic log packages methods now writing to our local logger. If logger
// is not initialized we fallback to the standard log system
func Info(v ...interface{}) {
	if config.ShelterConfig.LogLevel == config.LogLevelNormal {
		return
	}

	label := "[info]"
	v = append([]interface{}{label}, v...)

	if Logger == nil {
		log.Println(v...)
	} else {
		Logger.Println(v...)
	}
}

// Infof basic log packages methods now writing to our local logger. If logger
// is not initialized we fallback to the standard log system
func Infof(format string, v ...interface{}) {
	if config.ShelterConfig.LogLevel == config.LogLevelNormal {
		return
	}

	label := "[info] "
	if Logger == nil {
		log.Printf(label+format, v...)
	} else {
		Logger.Printf(label+format, v...)
	}
}

// Debug basic log packages methods now writing to our local logger. If logger
// is not initialized we fallback to the standard log system
func Debug(v ...interface{}) {
	if config.ShelterConfig.LogLevel != config.LogLevelDebug {
		return
	}

	label := "[debug]"
	v = append([]interface{}{label}, v...)

	if Logger == nil {
		log.Println(v...)
	} else {
		Logger.Println(v...)
	}
}

// Debugf basic log packages methods now writing to our local logger. If logger
// is not initialized we fallback to the standard log system
func Debugf(format string, v ...interface{}) {
	if config.ShelterConfig.LogLevel != config.LogLevelDebug {
		return
	}

	label := "[debug] "
	if Logger == nil {
		log.Printf(label+format, v...)
	} else {
		Logger.Printf(label+format, v...)
	}
}

// Close the log file and finish the logger
func Close() {
	logFile.Close()
}

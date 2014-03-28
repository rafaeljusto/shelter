// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"log"
)

var (
	TestName string
)

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints a simple message without ending the test
func Println(message string) {
	message = fmt.Sprintf("%s integration test: %s", TestName, message)
	log.Println(message)
}

// Function to allow printing a percentage inline, that is, replacing the last percentage
// with the new one. You should always use this method until we got 100%, or you will miss
// the breaking line and the text after the percentage will be mess up
func PrintProgress(label string, percentage int) {
	if percentage == 100 {
		fmt.Printf("\r... %s: %d%%\n", label, percentage)
	} else {
		fmt.Printf("\r... %s: %d%%", label, percentage)
	}
}

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints an error message without ending the test
func Errorln(message string, err error) {
	message = fmt.Sprintf("%s integration test: %s", TestName, message)
	if err != nil {
		message = fmt.Sprintf("%s. Details: %s", message, err.Error())
	}

	log.Println(message)
}

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints an error message and ends the test
func Fatalln(message string, err error) {
	message = fmt.Sprintf("%s integration test: %s", TestName, message)
	if err != nil {
		message = fmt.Sprintf("%s. Details: %s", message, err.Error())
	}

	log.Fatalln(message)
}

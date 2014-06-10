// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// utils add features to make the test life easier
package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// List of possible errors in this test. There can be also other errors from low level
// structures
var (
	// Config file path is a mandatory parameter
	ErrConfigFileUndefined = errors.New("Config file path undefined")
)

// Function to read the configuration file
func ReadConfigFile(configFilePath string, config interface{}) error {
	// Config file path is a mandatory program parameter
	if len(configFilePath) == 0 {
		return ErrConfigFileUndefined
	}

	confBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(confBytes, &config); err != nil {
		return err
	}

	return nil
}

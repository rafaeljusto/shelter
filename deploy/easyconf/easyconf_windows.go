// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/rafaeljusto/shelter/deploy/easyconf/conf"
	"log"
	"path/filepath"
)

var (
	sampleConfigFilPath string
)

func init() {
	flag.StringVar(&sampleConfigFilPath, "sample", "", "Initial file to use as a configuration reference")
}

func main() {
	flag.Parse()
	if len(sampleConfigFilPath) == 0 {
		log.Println("All arguments are mandatory!")
		flag.PrintDefaults()
		return
	}

	configFilePath := filepath.Join(filepath.Dir(sampleConfigFilPath), "shelter.conf")
	keysPath := filepath.Join(filepath.Dir(sampleConfigFilPath), "keys")
	conf.Run(configFilePath, sampleConfigFilPath, keysPath)
}

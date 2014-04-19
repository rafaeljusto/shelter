// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"github.com/rafaeljusto/shelter/deploy/easyconf/conf"
	"path/filepath"
)

// Unix like systems that uses debian packages cannot pass arguments in post intall script calls, so
// we are going to put it hard coded
func main() {
	sampleConfigFilPath := "/usr/shelter/etc/shelter.conf.unix.sample"
	configFilePath := filepath.Join(filepath.Dir(sampleConfigFilPath), "shelter.conf")
	keysPath := filepath.Join(filepath.Dir(sampleConfigFilPath), "keys")
	conf.Run(configFilePath, sampleConfigFilPath, keysPath)
}

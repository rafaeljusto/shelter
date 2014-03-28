// utils - Features for make the test life easier
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package utils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"net/http"
	"time"
)

func BuildHTTPHeader(r *http.Request, content []byte) {
	if r.ContentLength > 0 {
		r.Header.Set("Content-Type", check.SupportedContentType)

		hash := md5.New()
		hash.Write(content)
		hashBytes := hash.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

		r.Header.Set("Content-MD5", hashBase64)
	}

	r.Header.Set("Date", time.Now().Format(time.RFC1123))

	stringToSign, err := check.BuildStringToSign(r, "1")
	if err != nil {
		Fatalln("Error creating authorization", err)
	}

	signature := check.GenerateSignature(stringToSign, config.ShelterConfig.RESTServer.Secrets["1"])
	r.Header.Set("Authorization", fmt.Sprintf("%s %d:%s", check.SupportedNamespace, 1, signature))
}

// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/rafaeljusto/handy/interceptor"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	// isCurrentScan identify if the user desires to check only the scan in progress. As the current
	// scan doesn't have an id, we use the special world "current". We use the option i for case
	// insensitive
	isCurrentScan = regexp.MustCompile(`(?i:current)$`)
)

type ScanHandler interface {
	DatabaseHandler
	GetStartedAt() string
	SetScan(scan model.Scan)
	MessageResponse(string, string) error
}

type Scan struct {
	interceptor.NoAfterInterceptor
	scanHandler ScanHandler
}

func NewScan(h ScanHandler) *Scan {
	return &Scan{scanHandler: h}
}

func (i *Scan) Before(w http.ResponseWriter, r *http.Request) {
	date, err := time.Parse(time.RFC3339Nano, strings.ToUpper(i.scanHandler.GetStartedAt()))

	if err != nil {
		if err := i.scanHandler.MessageResponse("invalid-uri", r.URL.RequestURI()); err == nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Println("Error while writing response. Details:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	scanDAO := dao.ScanDAO{
		Database: i.scanHandler.GetDatabase(),
	}

	scan, err := scanDAO.FindByStartedAt(date)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	i.scanHandler.SetScan(scan)
}

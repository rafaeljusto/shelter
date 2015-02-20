// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// interceptor add steps to the REST request before calling the handler
package interceptor

import (
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/rafaeljusto/handy/interceptor"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/model"
	"net/http"
)

type DomainHandler interface {
	DatabaseHandler
	GetFQDN() string
	SetDomain(domain model.Domain)
}

type Domain struct {
	interceptor.NoAfterInterceptor
	domainHandler DomainHandler
}

func NewDomain(h DomainHandler) *Domain {
	return &Domain{domainHandler: h}
}

func (i *Domain) Before(w http.ResponseWriter, r *http.Request) {
	domainDAO := dao.DomainDAO{
		Database: i.domainHandler.GetDatabase(),
	}

	domain, err := domainDAO.FindByFQDN(i.domainHandler.GetFQDN())

	// For PUT method if the domain does not exist yet thats alright because we will create
	// it
	if r.Method != "PUT" && err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	i.domainHandler.SetDomain(domain)
}

// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/model"
	"testing"
)

func TestToDomainsResponse(t *testing.T) {
	domains := []model.Domain{
		{
			FQDN: "example1.com.br.",
		},
		{
			FQDN: "example2.com.br.",
		},
		{
			FQDN: "example3.com.br.",
		},
		{
			FQDN: "example4.com.br.",
		},
		{
			FQDN: "example5.com.br.",
		},
	}

	pagination := dao.DomainDAOPagination{
		PageSize: 10,
		Page:     1,
		OrderBy: []dao.DomainDAOSort{
			{
				Field:     dao.DomainDAOOrderByFieldFQDN,
				Direction: dao.DAOOrderByDirectionAscending,
			},
			{
				Field:     dao.DomainDAOOrderByFieldLastModifiedAt,
				Direction: dao.DAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: len(domains) / 10,
	}

	domainsResponse := ToDomainsResponse(domains, pagination, true, "example")

	if len(domainsResponse.Domains) != len(domains) {
		t.Error("Not converting domain model objects properly")
	}

	if domainsResponse.PageSize != 10 {
		t.Error("Pagination not storing the page size properly")
	}

	if domainsResponse.Page != 1 {
		t.Error("Pagination not storing the current page properly")
	}

	if domainsResponse.NumberOfItems != len(domains) {
		t.Error("Pagination not storing number of items properly")
	}

	if domainsResponse.NumberOfPages != len(domains)/10 {
		t.Error("Pagination not storing number of pages properly")
	}

	// When we are on the first page, and there's no other page, don't show any link
	if len(domainsResponse.Links) != 0 {
		t.Error("Response not adding the necessary links when there is only one page")
	}
}

func TestToDomainsResponseLinks(t *testing.T) {
	domains := []model.Domain{
		{
			FQDN: "example1.com.br.",
		},
		{
			FQDN: "example2.com.br.",
		},
		{
			FQDN: "example3.com.br.",
		},
		{
			FQDN: "example4.com.br.",
		},
		{
			FQDN: "example5.com.br.",
		},
	}

	pagination := dao.DomainDAOPagination{
		PageSize: 2,
		Page:     2,
		OrderBy: []dao.DomainDAOSort{
			{
				Field:     dao.DomainDAOOrderByFieldFQDN,
				Direction: dao.DAOOrderByDirectionAscending,
			},
			{
				Field:     dao.DomainDAOOrderByFieldLastModifiedAt,
				Direction: dao.DAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: 3,
	}

	domainsResponse := ToDomainsResponse(domains, pagination, true, "example")

	// Show all actions when navigating in the middle of the pagination
	if len(domainsResponse.Links) != 4 {
		t.Error("Response not adding the necessary links when we are navigating")
	}

	pagination = dao.DomainDAOPagination{
		PageSize: 2,
		Page:     1,
		OrderBy: []dao.DomainDAOSort{
			{
				Field:     dao.DomainDAOOrderByFieldFQDN,
				Direction: dao.DAOOrderByDirectionAscending,
			},
			{
				Field:     dao.DomainDAOOrderByFieldLastModifiedAt,
				Direction: dao.DAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: 3,
	}

	domainsResponse = ToDomainsResponse(domains, pagination, true, "example")

	// Don't show previous or fast backward when we are in the first page
	if len(domainsResponse.Links) != 2 {
		t.Error("Response not adding the necessary links when we are at the first page")
	}

	pagination = dao.DomainDAOPagination{
		PageSize: 2,
		Page:     3,
		OrderBy: []dao.DomainDAOSort{
			{
				Field:     dao.DomainDAOOrderByFieldFQDN,
				Direction: dao.DAOOrderByDirectionAscending,
			},
			{
				Field:     dao.DomainDAOOrderByFieldLastModifiedAt,
				Direction: dao.DAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: 3,
	}

	domainsResponse = ToDomainsResponse(domains, pagination, true, "example")

	// Don't show next or fast foward when we are in the last page
	if len(domainsResponse.Links) != 2 {
		t.Error("Response not adding the necessary links when we are in the last page")
	}
}

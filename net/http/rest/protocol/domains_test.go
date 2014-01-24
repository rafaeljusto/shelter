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
				Direction: dao.DomainDAOOrderByDirectionAscending,
			},
			{
				Field:     dao.DomainDAOOrderByFieldLastModifiedAt,
				Direction: dao.DomainDAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: len(domains) / 10,
	}

	domainsResponse := ToDomainsResponse(domains, pagination)

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

	// We can't show next and previous when there's only one page
	if len(domainsResponse.Links) != 2 {
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
				Direction: dao.DomainDAOOrderByDirectionAscending,
			},
			{
				Field:     dao.DomainDAOOrderByFieldLastModifiedAt,
				Direction: dao.DomainDAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: 3,
	}

	domainsResponse := ToDomainsResponse(domains, pagination)

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
				Direction: dao.DomainDAOOrderByDirectionAscending,
			},
			{
				Field:     dao.DomainDAOOrderByFieldLastModifiedAt,
				Direction: dao.DomainDAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: 3,
	}

	domainsResponse = ToDomainsResponse(domains, pagination)

	// Don't show previous when we are in the first page
	if len(domainsResponse.Links) != 3 {
		t.Error("Response not adding the necessary links when we are at the first page")
	}

	pagination = dao.DomainDAOPagination{
		PageSize: 2,
		Page:     3,
		OrderBy: []dao.DomainDAOSort{
			{
				Field:     dao.DomainDAOOrderByFieldFQDN,
				Direction: dao.DomainDAOOrderByDirectionAscending,
			},
			{
				Field:     dao.DomainDAOOrderByFieldLastModifiedAt,
				Direction: dao.DomainDAOOrderByDirectionDescending,
			},
		},
		NumberOfItems: len(domains),
		NumberOfPages: 3,
	}

	domainsResponse = ToDomainsResponse(domains, pagination)

	// Don't show next when we are in the last page
	if len(domainsResponse.Links) != 3 {
		t.Error("Response not adding the necessary links when we are in the last page")
	}
}

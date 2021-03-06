// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"fmt"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/model"
)

// DomainsResponse store multiple domains objects with pagination support
type DomainsResponse struct {
	Page          int              `json:"page"`              // Current page selected
	PageSize      int              `json:"pageSize"`          // Number of domains in a page
	NumberOfPages int              `json:"numberOfPages"`     // Total number of pages for the result set
	NumberOfItems int              `json:"numberOfItems"`     // Total number of domains in the result set
	Domains       []DomainResponse `json:"domains,omitempty"` // List of domain objects for the current page
	Links         []Link           `json:"links,omitempty"`   // Links for pagination managment
}

// Convert a list of domain objects into protocol format with pagination support
func ToDomainsResponse(
	domains []model.Domain,
	pagination dao.DomainDAOPagination,
	expand bool,
	filter string,
) DomainsResponse {

	var domainsResponses []DomainResponse
	for _, domain := range domains {
		domainsResponses = append(domainsResponses, ToDomainResponse(domain, true))
	}

	var orderBy string
	for _, sort := range pagination.OrderBy {
		if len(orderBy) > 0 {
			orderBy += "@"
		}

		orderBy += fmt.Sprintf("%s:%s",
			dao.DomainDAOOrderByFieldToString(sort.Field),
			dao.DAOOrderByDirectionToString(sort.Direction),
		)
	}

	expandParameter := ""
	if expand {
		expandParameter = "&expand"
	}

	// Add pagination managment links to the response. The URI is hard coded, I didn't have
	// any idea on how can we do this dynamically yet. We cannot get the URI from the
	// handler because we are going to have a cross-reference problem
	var links []Link

	// Only add fast backward if we aren't in the first page
	if pagination.Page > 1 {
		links = append(links, Link{
			Types: []LinkType{LinkTypeFirst},
			HRef: fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s&filter=%s%s",
				pagination.PageSize, 1, orderBy, filter, expandParameter),
		})
	}

	// Only add previous if theres a previous page
	if pagination.Page-1 >= 1 {
		links = append(links, Link{
			Types: []LinkType{LinkTypePrev},
			HRef: fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s&filter=%s%s",
				pagination.PageSize, pagination.Page-1, orderBy, filter, expandParameter),
		})
	}

	// Only add next if there's a next page
	if pagination.Page+1 <= pagination.NumberOfPages {
		links = append(links, Link{
			Types: []LinkType{LinkTypeNext},
			HRef: fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s&filter=%s%s",
				pagination.PageSize, pagination.Page+1, orderBy, filter, expandParameter),
		})
	}

	// Only add the fast forward if we aren't on the last page
	if pagination.Page < pagination.NumberOfPages {
		links = append(links, Link{
			Types: []LinkType{LinkTypeLast},
			HRef: fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s&filter=%s%s",
				pagination.PageSize, pagination.NumberOfPages, orderBy, filter, expandParameter),
		})
	}

	return DomainsResponse{
		Page:          pagination.Page,
		PageSize:      pagination.PageSize,
		NumberOfPages: pagination.NumberOfPages,
		NumberOfItems: pagination.NumberOfItems,
		Domains:       domainsResponses,
		Links:         links,
	}
}

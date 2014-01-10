package protocol

import (
	"fmt"
	"shelter/dao"
	"shelter/model"
)

// DomainsResponse store multiple domains objects with pagination support
type DomainsResponse struct {
	Page          int              `json:"page"`          // Current page selected
	PageSize      int              `json:"pageSize"`      // Number of domains in a page
	NumberOfPages int              `json:"numberOfPages"` // Total number of pages for the result set
	NumberOfItems int              `json:"numberOfItems"` // Total number of domains in the result set
	Domains       []DomainResponse `json:"domains"`       // List of domain objects for the current page
	Links         Links            `json:"links"`         // Links for pagination managment
}

// Convert a list of domain objects into protocol format with pagination support
func ToDomainsResponse(domains []model.Domain, pagination dao.DomainDAOPagination) DomainsResponse {
	var domainsResponses []DomainResponse
	for _, domain := range domains {
		domainsResponses = append(domainsResponses, ToDomainResponse(domain))
	}

	var orderBy string
	for _, sort := range pagination.OrderBy {
		if len(orderBy) > 0 {
			orderBy += "@"
		}

		orderBy += fmt.Sprintf("%s:%s",
			dao.DomainDAOOrderByFieldToString(sort.Field),
			dao.DomainDAOOrderByDirectionToString(sort.Direction),
		)
	}

	// Add pagination managment links to the response. The URI is hard coded, I didn't have
	// any idea on how can we do this dynamically yet. We cannot get the URI from the
	// handler because we are going to have a cross-reference problem
	links := Links{
		LinkTypeFirst:    fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s", pagination.PageSize, 1, orderBy),
		LinkTypeLast:     fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s", pagination.PageSize, pagination.NumberOfPages, orderBy),
		LinkTypeNext:     fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s", pagination.PageSize, pagination.Page+1, orderBy),
		LinkTypePrevious: fmt.Sprintf("/domains/?pagesize=%d&page=%d&orderby=%s", pagination.PageSize, pagination.Page-1, orderBy),
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

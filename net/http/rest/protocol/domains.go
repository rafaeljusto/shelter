package protocol

import (
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
}

// Convert a list of domain objects into protocol format with pagination support
func ToDomainsResponse(domains []model.Domain, pagination dao.DomainDAOPagination) DomainsResponse {
	var domainsResponses []DomainResponse
	for _, domain := range domains {
		domainsResponses = append(domainsResponses, ToDomainResponse(domain))
	}

	return DomainsResponse{
		Page:          pagination.Page,
		PageSize:      pagination.PageSize,
		NumberOfPages: pagination.NumberOfPages,
		NumberOfItems: pagination.NumberOfItems,
		Domains:       domainsResponses,
	}
}

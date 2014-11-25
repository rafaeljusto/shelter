// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

// DomainsResponse store multiple domains objects with pagination support
type DomainsResponse struct {
	Page          int              `json:"page"`              // Current page selected
	PageSize      int              `json:"pageSize"`          // Number of domains in a page
	NumberOfPages int              `json:"numberOfPages"`     // Total number of pages for the result set
	NumberOfItems int              `json:"numberOfItems"`     // Total number of domains in the result set
	Domains       []DomainResponse `json:"domains,omitempty"` // List of domain objects for the current page
	Links         []Link           `json:"links,omitempty"`   // Links for pagination managment
}

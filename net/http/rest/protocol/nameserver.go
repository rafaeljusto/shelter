package protocol

import (
	"time"
)

// Nameserver object used in the protocol to determinate what the user can update
type NameserverRequest struct {
	Host string `json:"host,omitempty"` // Nameserver's name
	IPv4 string `json:"ipv4,omitempty"` // Host's IPv4 (optional when don't need glue)
	IPv6 string `json:"ipv6,omitempty"` // Host's IPv6 (optional)
}

// Namerserver object used in the protocol to determinate what the user can see. The
// status was converted to text format for easy interpretation
type NameserverResponse struct {
	Host        string    `json:"host,omitempty"`        // Nameserver's name
	IPv4        string    `json:"ipv4,omitempty"`        // Host's IPv4 (optional when don't need glue)
	IPv6        string    `json:"ipv6,omitempty"`        // Host's IPv6 (optional)
	LastStatus  string    `json:"lastStatus,omitempty"`  // Result of the last configuration check
	LastCheckAt time.Time `json:"lastCheckAt,omitempty"` // Time of the last configuration check
	LastOKAt    time.Time `json:"lastOKAt,omitempty"`    // Last time that the DNS configuration was OK
}

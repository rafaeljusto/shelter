package model

import (
	"time"
)

// List of possible DS status
const (
	DSStatusOK = iota
	DSStatusTimeout
	DSStatusNoSignature
	DSStatusExpiredSignature
	DSStatusNoKey
	DSStatusNoSEP
	DSStatusSignatureError
	DSStatusNoDNSSEC
	DSStatusDNSError
)

type DSStatus int

// DS store the information necessary to validate if a domain is configured correctly with
// DNSSEC, and it also stores the results of the validations. When the hosts have multiple
// DNSSEC problems, the worst problem (using a priority algorithm) will be stored in the
// DS
type DS struct {
	Keytag    uint16
	Algorithm uint8
	Digest    string
	Status    DSStatus
	LastCheck time.Time
	LastOK    time.Time
}

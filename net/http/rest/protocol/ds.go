package protocol

import (
	"errors"
	"shelter/model"
	"time"
)

// List of possible errors that can occur when calling methods from this object
var (
	// Error returned when trying to convert a DS with an unknown DS algorithm number
	ErrInvalidDSAlgorithm = errors.New("DS algorithm invalid or not supported")
	// Error returned when trying to convert a DS with an unknown DS dogest type number
	ErrInvalidDSDigestType = errors.New("DS digest type invalid or not supported")
)

// DS object used in the protocol to determinate what the user can update
type DSRequest struct {
	Keytag     uint16 `json:"keytag,omitempty"`     // DNSKEY's identification number
	Algorithm  uint8  `json:"algorithm,omitempty"`  // DNSKEY's algorithm
	Digest     string `json:"digest,omitempty"`     // Hash of the DNSKEY content
	DigestType uint8  `json:"digestType,omitempty"` // Hash type decided by user when generating the DS
}

// DS object used in the protocol to determinate what the user can see. The status was
// converted to text format for easy interpretation
type DSResponse struct {
	Keytag      uint16    `json:"keytag,omitempty"`      // DNSKEY's identification number
	Algorithm   uint8     `json:"algorithm,omitempty"`   // DNSKEY's algorithm
	Digest      string    `json:"digest,omitempty"`      // Hash of the DNSKEY content
	DigestType  uint8     `json:"digestType,omitempty"`  // Hash type decided by user when generating the DS
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`   // DNSKEY's signature expiration date
	LastStatus  string    `json:"lastStatus,omitempty"`  // Result of the last configuration check
	LastCheckAt time.Time `json:"lastCheckAt,omitempty"` // Time of the last configuration check
	LastOKAt    time.Time `json:"lastOKAt,omitempty"`    // Last time that the DNSSEC configuration was OK
}

// Convert a list of DS requests objects into a list of DS model objects. Useful when
// merging domain object from the network with a domain object from the database. It can
// return errors related to the conversion of the algorithm, when it is out of range
func toDSSetModel(dsSetRequest []DSRequest) ([]model.DS, error) {
	var dsSet []model.DS
	for _, dsRequest := range dsSetRequest {
		ds, err := toDSModel(dsRequest)
		if err != nil {
			return nil, err
		}

		dsSet = append(dsSet, ds)
	}

	return dsSet, nil
}

// Convert a DS request object into a DS model object. It can return errors related to the
// conversion of the algorithm, when it is out of range
func toDSModel(dsRequest DSRequest) (model.DS, error) {
	ds := model.DS{
		Keytag: dsRequest.Keytag,
		Digest: dsRequest.Digest, // TODO: Normalize digest? (lowercase or uppercase)
	}

	if !model.IsValidDSAlgorithm(dsRequest.Algorithm) {
		return ds, ErrInvalidDSAlgorithm
	}

	ds.Algorithm = model.DSAlgorithm(dsRequest.Algorithm)

	if !model.IsValidDSDigestType(dsRequest.DigestType) {
		return ds, ErrInvalidDSDigestType
	}

	ds.DigestType = model.DSDigestType(dsRequest.DigestType)

	return ds, nil
}

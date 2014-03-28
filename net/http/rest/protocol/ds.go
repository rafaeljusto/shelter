// protocol - REST protocol description
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package protocol

import (
	"errors"
	"github.com/rafaeljusto/shelter/model"
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

// Convert a DS request object into a DS model object. It can return errors related to the
// conversion of the algorithm, when it is out of range
func (d *DSRequest) toDSModel() (model.DS, error) {
	ds := model.DS{
		Keytag: d.Keytag,
		Digest: model.NormalizeDSDigest(d.Digest),
	}

	if !model.IsValidDSAlgorithm(d.Algorithm) {
		return ds, ErrInvalidDSAlgorithm
	}

	ds.Algorithm = model.DSAlgorithm(d.Algorithm)

	if !model.IsValidDSDigestType(d.DigestType) {
		return ds, ErrInvalidDSDigestType
	}

	ds.DigestType = model.DSDigestType(d.DigestType)

	return ds, nil
}

// Convert a list of DS requests objects into a list of DS model objects. Useful when
// merging domain object from the network with a domain object from the database. It can
// return errors related to the conversion of the algorithm, when it is out of range
func toDSSetModel(dsSetRequest []DSRequest) ([]model.DS, error) {
	var dsSet []model.DS
	for _, dsRequest := range dsSetRequest {
		ds, err := dsRequest.toDSModel()
		if err != nil {
			return nil, err
		}

		dsSet = append(dsSet, ds)
	}

	return dsSet, nil
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

// Convert a DS of the system into a format with limited information to return it to the
// user
func toDSResponse(ds model.DS) DSResponse {
	return DSResponse{
		Keytag:      ds.Keytag,
		Algorithm:   uint8(ds.Algorithm),
		Digest:      ds.Digest,
		DigestType:  uint8(ds.DigestType),
		ExpiresAt:   ds.ExpiresAt,
		LastStatus:  model.DSStatusToString(ds.LastStatus),
		LastCheckAt: ds.LastCheckAt,
		LastOKAt:    ds.LastOKAt,
	}
}

// Convert a list of DS of the system into a format with limited information to return it
// to the user. This is only a easy way to call toDSResponse for each object in the list
func toDSSetResponse(dsSet []model.DS) []DSResponse {
	var dsSetResponse []DSResponse
	for _, ds := range dsSet {
		dsSetResponse = append(dsSetResponse, toDSResponse(ds))
	}
	return dsSetResponse
}

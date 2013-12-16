package transcoder

import (
	"shelter/model"
	"shelter/net/http/rest/protocol"
)

// Convert a list of DS requests objects into a list of DS model objects. Useful when
// merging domain object from the network with a domain object from the database. It can
// return errors related to the conversion of the algorithm, when it is out of range
func toDSSetModel(dsSetRequest []protocol.DSRequest) ([]model.DS, error) {
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
func toDSModel(dsRequest protocol.DSRequest) (model.DS, error) {
	ds := model.DS{
		Keytag:     dsRequest.Keytag,
		Algorithm:  dsRequest.Algorithm,
		Digest:     dsRequest.Digest,
		DigestType: dsRequest.DigestType,
	}

	return ds, nil
}

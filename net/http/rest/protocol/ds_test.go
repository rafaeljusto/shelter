package protocol

import (
	"shelter/model"
	"testing"
)

func TestToDSSetModel(t *testing.T) {
	dsSetRequest := []DSRequest{
		{
			Keytag:     41674,
			Algorithm:  5,
			Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
			DigestType: 1,
		},
		{
			Keytag:     45966,
			Algorithm:  7,
			Digest:     "B7C0BDE8F3C90E573B956B14A14CAF5001A3E841",
			DigestType: 1,
		},
	}

	_, err := toDSSetModel(dsSetRequest)
	if err != nil {
		t.Error(err)
	}

	dsSetRequest = []DSRequest{
		{
			Keytag:     41674,
			Algorithm:  5,
			Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
			DigestType: 1,
		},
		{
			Keytag:     45966,
			Algorithm:  0,
			Digest:     "B7C0BDE8F3C90E573B956B14A14CAF5001A3E841",
			DigestType: 1,
		},
	}

	_, err = toDSSetModel(dsSetRequest)
	if err == nil {
		t.Error("Not verifying errors in DS set conversion")
	}
}

func TestToDSModel(t *testing.T) {
	dsRequest := DSRequest{
		Keytag:     41674,
		Algorithm:  5,
		Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
		DigestType: 1,
	}

	ds, err := toDSModel(dsRequest)
	if err != nil {
		t.Fatal(err)
	}

	if ds.Keytag != 41674 {
		t.Error("Not keeping keytag in conversion")
	}

	if ds.Algorithm != model.DSAlgorithmRSASHA1 {
		t.Error("Not converting DS algorithm correctly")
	}

	if ds.Digest != "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C" {
		t.Error("Not keeping digest in conversion")
	}

	if ds.DigestType != model.DSDigestTypeSHA1 {
		t.Error("Not converting DS digest type correctly")
	}

	dsRequest = DSRequest{
		Keytag:     41674,
		Algorithm:  0,
		Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
		DigestType: 1,
	}

	ds, err = toDSModel(dsRequest)
	if err == nil {
		t.Error("Allowing an invalid DS algorithm")
	}

	dsRequest = DSRequest{
		Keytag:     41674,
		Algorithm:  5,
		Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
		DigestType: 6,
	}

	ds, err = toDSModel(dsRequest)
	if err == nil {
		t.Error("Allowing an invalid DS digest type")
	}
}

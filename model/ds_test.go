package model

import (
	"testing"
	"time"
)

func TestDSChangeStatus(t *testing.T) {
	timeMark := time.Now()

	ds := DS{
		LastStatus:  DSStatusDNSError,
		LastCheckAt: timeMark,
	}

	ds.ChangeStatus(DSStatusOK)

	if ds.LastStatus != DSStatusOK {
		t.Error("ChangeStatus method did not change DS attribute")
	}

	if ds.LastCheckAt.Before(timeMark) || ds.LastCheckAt.Equal(timeMark) {
		t.Error("ChangeStatus method did not update the last check date")
	}
}

func TestDSStatusToString(t *testing.T) {
	if DSStatusToString(DSStatusOK) != "OK" {
		t.Error("DS status OK not converting correctly to string")
	}

	if DSStatusToString(DSStatusTimeout) != "TIMEOUT" {
		t.Error("DS status TIMEOUT not converting correctly to string")
	}

	if DSStatusToString(DSStatusNoSignature) != "NOSIG" {
		t.Error("DS status NOSIG not converting correctly to string")
	}

	if DSStatusToString(DSStatusExpiredSignature) != "EXPSIG" {
		t.Error("DS status EXPSIG not converting correctly to string")
	}

	if DSStatusToString(DSStatusNoKey) != "NOKEY" {
		t.Error("DS status NOKEY not converting correctly to string")
	}

	if DSStatusToString(DSStatusNoSEP) != "NOSEP" {
		t.Error("DS status NOSEP not converting correctly to string")
	}

	if DSStatusToString(DSStatusSignatureError) != "SIGERR" {
		t.Error("DS status SIGERR not converting correctly to string")
	}

	if DSStatusToString(DSStatusDNSError) != "DNSERR" {
		t.Error("DS status DNSERR not converting correctly to string")
	}

	if DSStatusToString(999999) != "" {
		t.Error("Unknown DS status associated to some existing status")
	}
}

func TestValidDSAlgorithm(t *testing.T) {
	if !IsValidDSAlgorithm(1) || !IsValidDSAlgorithm(8) || !IsValidDSAlgorithm(254) {
		t.Error("Not accepting valid DS algorithms")
	}

	if IsValidDSAlgorithm(0) || IsValidDSAlgorithm(255) {
		t.Error("Accepting invalid DS algorithms")
	}
}

func TestValidDSDigestType(t *testing.T) {
	if !IsValidDSDigestType(0) || !IsValidDSDigestType(3) || !IsValidDSDigestType(5) {
		t.Error("Not accepting valid DS digest type")
	}

	if IsValidDSDigestType(6) {
		t.Error("Accepting invalid digest type")
	}
}

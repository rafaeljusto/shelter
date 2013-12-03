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

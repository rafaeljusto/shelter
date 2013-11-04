package model

import (
	"testing"
	"time"
)

func TestAllNameserversOK(t *testing.T) {
	d := Domain{
		Nameservers: []Nameserver{
			{
				LastStatus: NameserverStatusFail,
			},
			{
				LastStatus: NameserverStatusOK,
			},
		},
	}

	if d.AllNameserversOK() {
		t.Error("Fail to detect a nameserver with error")
	}

	d = Domain{
		Nameservers: []Nameserver{},
	}

	if !d.AllNameserversOK() {
		t.Error("Fail to inform all nameservers are OK when there are no nameservers")
	}

	d = Domain{
		Nameservers: []Nameserver{
			{
				LastStatus: NameserverStatusOK,
			},
			{
				LastStatus: NameserverStatusOK,
			},
		},
	}

	if !d.AllNameserversOK() {
		t.Error("Fail to inform that all nameservers are OK")
	}
}

func TestAllDSSetOK(t *testing.T) {
	d := Domain{
		DSSet: []DS{
			{
				LastStatus: DSStatusExpiredSignature,
			},
			{
				LastStatus: DSStatusOK,
			},
		},
	}

	if d.AllDSSetOK() {
		t.Error("Fail to detect a DS with error")
	}

	d = Domain{
		DSSet: []DS{},
	}

	if !d.AllDSSetOK() {
		t.Error("Fail to inform all DS set are OK when there are no DS records")
	}

	d = Domain{
		DSSet: []DS{
			{
				LastStatus: DSStatusOK,
			},
			{
				LastStatus: DSStatusOK,
			},
		},
	}

	if !d.AllDSSetOK() {
		t.Error("Fail to inform that all DS set are OK")
	}
}

func TestDaysSinceLastCheck(t *testing.T) {
	twoDays, _ := time.ParseDuration("48h")
	threeDays, _ := time.ParseDuration("72h")
	fourDays, _ := time.ParseDuration("96h")
	fiveDays, _ := time.ParseDuration("120h")

	d := Domain{
		Nameservers: []Nameserver{
			{
				LastCheckAt: time.Now().Add(-fiveDays),
			},
			{
				LastCheckAt: time.Now().Add(-twoDays),
			},
			{
				LastCheckAt: time.Now().Add(-twoDays),
			},
		},
		DSSet: []DS{
			{
				LastCheckAt: time.Now().Add(-fourDays),
			},
			{
				LastCheckAt: time.Now().Add(-threeDays),
			},
		},
	}

	if days := d.DaysSinceLastCheck(); days != 2 {
		t.Errorf("Not counting correctly the number of days since last "+
			"check for nameservers. Should be 2, but got %d", days)
	}

	d = Domain{
		Nameservers: []Nameserver{
			{
				LastCheckAt: time.Now().Add(-fiveDays),
			},
			{
				LastCheckAt: time.Now().Add(-fourDays),
			},
			{
				LastCheckAt: time.Now().Add(-fourDays),
			},
		},
		DSSet: []DS{
			{
				LastCheckAt: time.Now().Add(-fourDays),
			},
			{
				LastCheckAt: time.Now().Add(-threeDays),
			},
		},
	}

	if days := d.DaysSinceLastCheck(); days != 3 {
		t.Errorf("Not counting correctly the number of days since last "+
			"check for DS set. Should be 3, but got %d", days)
	}

	d = Domain{}

	if d.DaysSinceLastCheck() < 365 {
		t.Errorf("Not counting correctly the number of days since last check " +
			"when the object is empty")
	}
}

func TestNearDNSSECExpirationDate(t *testing.T) {
	tenDays, _ := time.ParseDuration("240h")
	elevenDays, _ := time.ParseDuration("264h")

	d := Domain{
		DSSet: []DS{
			{
				ExpiresAt: time.Now().Add(tenDays),
			},
		},
	}

	if !d.IsNearDNSSECExpirationDate(10) {
		t.Error("Could not detect when the DNSSEC expiration is near")
	}

	d = Domain{
		DSSet: []DS{
			{
				ExpiresAt: time.Now().Add(elevenDays),
			},
		},
	}

	if d.IsNearDNSSECExpirationDate(10) {
		t.Error("Could not detect when the DNSSEC expiration is far")
	}

	d = Domain{
		DSSet: []DS{},
	}

	if d.IsNearDNSSECExpirationDate(10) {
		t.Error("Could not detect when there's no expiration date")
	}
}

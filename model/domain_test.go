package model

import (
	"strconv"
	"testing"
	"time"
)

func TestShouldBeScanned(t *testing.T) {
	var (
		maxOKVerificationDays    = 7
		maxErrorVerificationDays = 3
		maxExpirationAlertDays   = 10
	)

	maxErrorVerificationDuration, _ :=
		time.ParseDuration("-" + strconv.Itoa(maxErrorVerificationDays*24) + "h")
	maxOKVerificationDuration, _ :=
		time.ParseDuration("-" + strconv.Itoa(maxOKVerificationDays*24) + "h")
	maxExpirationAlertDuration, _ :=
		time.ParseDuration("-" + strconv.Itoa(maxExpirationAlertDays*24) + "h")

	d := Domain{
		Nameservers: []Nameserver{
			{
				LastStatus:  NameserverStatusFail,
				LastCheckAt: time.Now().Add(maxErrorVerificationDuration),
			},
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now().Add(maxErrorVerificationDuration),
			},
		},
	}

	if !d.ShouldBeScanned(maxOKVerificationDays,
		maxErrorVerificationDays, maxExpirationAlertDays) {

		t.Error("Did not select a domain with DNS errors and in the limit " +
			"of errors max verification period")
	}

	d = Domain{
		DSSet: []DS{
			{
				LastStatus:  DSStatusTimeout,
				LastCheckAt: time.Now().Add(maxErrorVerificationDuration),
			},
			{
				LastStatus:  DSStatusTimeout,
				LastCheckAt: time.Now().Add(maxErrorVerificationDuration),
			},
		},
	}

	if !d.ShouldBeScanned(maxOKVerificationDays,
		maxErrorVerificationDays, maxExpirationAlertDays) {

		t.Error("Did not select a domain with DNSSEC errors and in the limit " +
			"of errors max verification period")
	}

	d = Domain{
		Nameservers: []Nameserver{
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now().Add(maxOKVerificationDuration),
			},
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now().Add(maxOKVerificationDuration),
			},
		},
	}

	if !d.ShouldBeScanned(maxOKVerificationDays,
		maxErrorVerificationDays, maxExpirationAlertDays) {

		t.Error("Did not select a domain configured correctly but in the limit " +
			"of ok max verification period")
	}

	d = Domain{
		Nameservers: []Nameserver{
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now(),
			},
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now(),
			},
		},

		DSSet: []DS{
			{
				LastStatus:  DSStatusOK,
				LastCheckAt: time.Now(),
				ExpiresAt:   time.Now().Add(maxExpirationAlertDuration),
			},
			{
				LastStatus:  DSStatusOK,
				LastCheckAt: time.Now(),
			},
		},
	}

	if !d.ShouldBeScanned(maxOKVerificationDays,
		maxErrorVerificationDays, maxExpirationAlertDays) {

		t.Error("Did not select a domain with DNSSEC signatures near expiration")
	}

	d = Domain{
		Nameservers: []Nameserver{
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now(),
			},
			{
				LastStatus:  NameserverStatusOK,
				LastCheckAt: time.Now(),
			},
		},
	}

	if d.ShouldBeScanned(maxOKVerificationDays,
		maxErrorVerificationDays, maxExpirationAlertDays) {

		t.Error("Selected a domain configured correctly and checked now, " +
			"this shouldn't happen, or this domain is very lucky, we should " +
			"take a look in the algorithm")
	}
}

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

	if d.allNameserversOK() {
		t.Error("Fail to detect a nameserver with error")
	}

	d = Domain{
		Nameservers: []Nameserver{},
	}

	if !d.allNameserversOK() {
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

	if !d.allNameserversOK() {
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

	if d.allDSSetOK() {
		t.Error("Fail to detect a DS with error")
	}

	d = Domain{
		DSSet: []DS{},
	}

	if !d.allDSSetOK() {
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

	if !d.allDSSetOK() {
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

	if days := d.daysSinceLastCheck(); days != 2 {
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

	if days := d.daysSinceLastCheck(); days != 3 {
		t.Errorf("Not counting correctly the number of days since last "+
			"check for DS set. Should be 3, but got %d", days)
	}

	d = Domain{}

	if d.daysSinceLastCheck() < 365 {
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

	if !d.isNearDNSSECExpirationDate(10) {
		t.Error("Could not detect when the DNSSEC expiration is near")
	}

	d = Domain{
		DSSet: []DS{
			{
				ExpiresAt: time.Now().Add(elevenDays),
			},
		},
	}

	if d.isNearDNSSECExpirationDate(10) {
		t.Error("Could not detect when the DNSSEC expiration is far")
	}

	d = Domain{
		DSSet: []DS{},
	}

	if d.isNearDNSSECExpirationDate(10) {
		t.Error("Could not detect when there's no expiration date")
	}
}

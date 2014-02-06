package protocol

import (
	"github.com/rafaeljusto/shelter/model"
	"labix.org/v2/mgo/bson"
	"testing"
	"time"
)

func TestScanToScanResponse(t *testing.T) {
	scanResponse := ScanToScanResponse(model.Scan{
		Id:                      bson.NewObjectId(),
		Revision:                1,
		Status:                  model.ScanStatusExecuted,
		StartedAt:               time.Now().Add(-1 * time.Hour),
		FinishedAt:              time.Now().Add(-30 * time.Minute),
		LastModifiedAt:          time.Now().Add(-30 * time.Minute),
		DomainsScanned:          10,
		DomainsWihDNSSECScanned: 4,
		NameserverStatistics: map[string]uint64{
			model.NameserverStatusToString(model.NameserverStatusOK):      16,
			model.NameserverStatusToString(model.NameserverStatusTimeout): 4,
		},
		DSStatistics: map[string]uint64{
			model.DSStatusToString(model.DSStatusOK):               3,
			model.DSStatusToString(model.DSStatusExpiredSignature): 1,
		},
	})

	if scanResponse.Status != "EXECUTED" {
		t.Error("Status is not being translated correctly for a scan")
	}

	// TODO
}

func TestCurrentScanToScanResponse(t *testing.T) {
	scanResponse := CurrentScanToScanResponse(model.CurrentScan{
		DomainsToBeScanned: 4,
		Scan: model.Scan{
			Status:                  model.ScanStatusRunning,
			StartedAt:               time.Now().Add(-1 * time.Hour),
			DomainsScanned:          2,
			DomainsWihDNSSECScanned: 0,
		},
	})

	if scanResponse.Status != "RUNNING" {
		t.Error("Status is not being translated correctly for a current scan")
	}

	// TODO
}

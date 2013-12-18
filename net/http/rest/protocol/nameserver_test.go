package protocol

import (
	"net"
	"shelter/model"
	"strings"
	"testing"
	"time"
)

func TestToNameserversModel(t *testing.T) {
	nameserversRequest := []NameserverRequest{
		{
			Host: "ns1.example.com.br",
			IPv4: "127.0.0.1",
			IPv6: "::1",
		},
		{
			Host: "ns2.example.com.br",
			IPv4: "127.0.0.2",
			IPv6: "::2",
		},
	}

	if _, err := toNameserversModel(nameserversRequest); err != nil {
		t.Error(err)
	}

	nameserversRequest = []NameserverRequest{
		{
			Host: "ns1.example.com.br",
			IPv4: "127.ABC.0.1",
			IPv6: "::1",
		},
		{
			Host: "ns2.example.com.br",
			IPv4: "127.0.0.2",
			IPv6: "::2",
		},
	}

	if _, err := toNameserversModel(nameserversRequest); err == nil {
		t.Error("Not checking errors from nameservers conversion")
	}
}

func TestToNameserverModel(t *testing.T) {
	nameserverRequest := NameserverRequest{
		Host: "ns1.example.com.br",
		IPv4: "127.0.0.1",
		IPv6: "::1",
	}

	nameserver, err := toNameserverModel(nameserverRequest)
	if err != nil {
		t.Fatal(err)
	}

	if nameserver.Host != "ns1.example.com.br." {
		t.Error("Not normalizing nameserver's host")
	}

	if nameserver.IPv4.String() != "127.0.0.1" {
		t.Error("Not parsing correctly IPv4")
	}

	if nameserver.IPv6.String() != "::1" {
		t.Error("Not parsing correctly IPv6")
	}

	nameserverRequest = NameserverRequest{
		Host: "ns1.example.com.br",
		IPv4: "127..0.0.1",
		IPv6: "::1",
	}

	nameserver, err = toNameserverModel(nameserverRequest)
	if err == nil {
		t.Error("Accepting an invalid IPv4")
	}

	nameserverRequest = NameserverRequest{
		Host: "ns1.example.com.br",
		IPv4: "127.0.0.1",
		IPv6: ":::1",
	}

	nameserver, err = toNameserverModel(nameserverRequest)
	if err == nil {
		t.Error("Accepting an invalid IPv6")
	}

	nameserverRequest = NameserverRequest{
		Host: strings.Repeat("x", 65536) + "\uff00", // int32 overflow
		IPv4: "127.0.0.1",
		IPv6: "::1",
	}

	nameserver, err = toNameserverModel(nameserverRequest)
	if err == nil {
		t.Error("Accepting an invalid FQDN for IDNA")
	}
}

func TestToNameserversResponse(t *testing.T) {
	now := time.Now()

	nameservers := []model.Nameserver{
		{
			Host:        "ns1.example.com.br.",
			IPv4:        net.ParseIP("127.0.0.1"),
			IPv6:        net.ParseIP("::1"),
			LastStatus:  model.NameserverStatusOK,
			LastCheckAt: now,
			LastOKAt:    now,
		},
		{
			Host:        "ns2.example.com.br.",
			IPv4:        net.ParseIP("127.0.0.2"),
			IPv6:        net.ParseIP("::2"),
			LastStatus:  model.NameserverStatusError,
			LastCheckAt: now,
		},
	}

	nameserversResponse := toNameserversResponse(nameservers)

	if len(nameserversResponse) != 2 {
		t.Error("Fail to convert multiple nameservers")
	}
}

func TestToNameserverResponse(t *testing.T) {
	now := time.Now()

	nameserver := model.Nameserver{
		Host:        "ns1.example.com.br.",
		IPv4:        net.ParseIP("127.0.0.1"),
		IPv6:        net.ParseIP("::1"),
		LastStatus:  model.NameserverStatusOK,
		LastCheckAt: now,
		LastOKAt:    now,
	}

	nameserverResponse := toNameserverResponse(nameserver)

	if nameserverResponse.Host != "ns1.example.com.br." {
		t.Error("Fail to convert host")
	}

	if nameserverResponse.IPv4 != "127.0.0.1" {
		t.Error("Fail to convert IPv4")
	}

	if nameserverResponse.IPv6 != "::1" {
		t.Error("Fail to convert IPv6")
	}

	if nameserverResponse.LastStatus !=
		model.NameserverStatusToString(model.NameserverStatusOK) {

		t.Error("Fail to convert last status")
	}

	if nameserverResponse.LastCheckAt.Unix() != now.Unix() ||
		nameserverResponse.LastOKAt.Unix() != now.Unix() {

		t.Error("Fail to convert dates")
	}
}

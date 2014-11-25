// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"testing"
)

func TestOwnerApply(t *testing.T) {
	ownerRequest := OwnerRequest{
		Email:    "example01@example.com.br",
		Language: "pt-br",
	}

	var owner Owner
	if !owner.Apply(ownerRequest) {
		t.Fatal("Not applying a valid request")
	}

	if owner.Email.Address != "example01@example.com.br" {
		t.Error("Not converting e-mail properly")
	}

	if owner.Language != "pt-BR" {
		t.Error("Not converting language properly")
	}

	ownerRequest = OwnerRequest{
		Email:    "notavalidemail",
		Language: "pt-br",
	}

	owner, err = ownerRequest.toOwnerModel()
	if owner.Apply(ownerRequest) {
		t.Error("Not checking e-mail format on conversion")
	}

	ownerRequest = OwnerRequest{
		Email:    "example01@example.com.br",
		Language: "zzz",
	}

	if owner.Apply(ownerRequest) {
		t.Error("Not checking language on conversion")
	}
}

func TestOwnerProtocol(t *testing.T) {
	email, err := mail.ParseAddress("example@example.com.br")
	if err != nil {
		t.Fatal(err)
	}

	owner := model.Owner{
		Email:    email,
		Language: "en-US",
	}

	ownerResponse := owner.Protocol()

	if ownerResponse.Email != "example@example.com.br" {
		t.Error("Not converting e-mail properly")
	}

	if ownerResponse.Language != "en-US" {
		t.Error("Not converting language properly")
	}
}

func TestOwnersApply(t *testing.T) {
	ownersRequest := []OwnerRequest{
		{
			Email:    "example01@example.com.br.",
			Language: "pt-br",
		},
		{
			Email:    "example02@example.com.br.",
			Language: "en-US",
		},
	}

	var owners Owners
	if !owners.Apply(ownersRequest) {
		t.Error("Should apply a valid owners request")
	}

	ownersRequest = []OwnerRequest{
		{
			Email:    "notavalidemail",
			Language: "pt-br",
		},
		{
			Email:    "example02@example.com.br.",
			Language: "en-US",
		},
	}

	if owners.Apply(ownersRequest) {
		t.Error("Not checking invalid e-mail on conversion")
	}

	ownersRequest = []OwnerRequest{
		{
			Email:    "example01@example.com.br.",
			Language: "zzzzz",
		},
		{
			Email:    "example02@example.com.br.",
			Language: "en-US",
		},
	}

	if owners.Apply(ownersRequest) {
		t.Error("Not checking invalid language on conversion")
	}
}

func TestOwnersProtocol(t *testing.T) {
	email1, err := mail.ParseAddress("example1@example.com.br")
	if err != nil {
		t.Fatal(err)
	}

	email2, err := mail.ParseAddress("example2@example.com.br")
	if err != nil {
		t.Fatal(err)
	}

	owners := Owners{
		{
			Email:    email1,
			Language: "en-US",
		},
		{
			Email:    email2,
			Language: "pt-BR",
		},
	}

	ownersResponse := owners.Protocol()

	if len(ownersResponse) != 2 {
		t.Fatal("Not adding all owners to response")
	}

	if ownersResponse[0].Email != "example1@example.com.br" ||
		ownersResponse[0].Language != "en-US" {
		t.Error("Owner 1 was converted with problems")
	}

	if ownersResponse[1].Email != "example2@example.com.br" ||
		ownersResponse[1].Language != "pt-BR" {
		t.Error("Owner 2 was converted with problems")
	}
}

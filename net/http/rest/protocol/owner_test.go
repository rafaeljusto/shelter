// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package protocol

import (
	"github.com/rafaeljusto/shelter/model"
	"net/mail"
	"testing"
)

func TestToOwnerModel(t *testing.T) {
	ownerRequest := OwnerRequest{
		Email:    "example01@example.com.br",
		Language: "pt-br",
	}

	owner, err := ownerRequest.toOwnerModel()
	if err != nil {
		t.Fatal(err)
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
	if err == nil {
		t.Error("Not checking e-mail format on conversion")
	}

	ownerRequest = OwnerRequest{
		Email:    "example01@example.com.br",
		Language: "zzz",
	}

	owner, err = ownerRequest.toOwnerModel()
	if err == nil {
		t.Error("Not checking language on conversion")
	}
}

func TestToOwnersModel(t *testing.T) {
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

	if _, err := toOwnersModel(ownersRequest); err != nil {
		t.Error(err)
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

	if _, err := toOwnersModel(ownersRequest); err == nil {
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

	if _, err := toOwnersModel(ownersRequest); err == nil {
		t.Error("Not checking invalid language on conversion")
	}
}

func TestToOwnerResponse(t *testing.T) {
	email, err := mail.ParseAddress("example@example.com.br")
	if err != nil {
		t.Fatal(err)
	}

	owner := model.Owner{
		Email:    email,
		Language: "en-US",
	}

	ownerResponse := toOwnerResponse(owner)

	if ownerResponse.Email != "example@example.com.br" {
		t.Error("Not converting e-mail properly")
	}

	if ownerResponse.Language != "en-US" {
		t.Error("Not converting language properly")
	}
}

func TestToOwnersResponse(t *testing.T) {
	email1, err := mail.ParseAddress("example1@example.com.br")
	if err != nil {
		t.Fatal(err)
	}

	email2, err := mail.ParseAddress("example2@example.com.br")
	if err != nil {
		t.Fatal(err)
	}

	owners := []model.Owner{
		{
			Email:    email1,
			Language: "en-US",
		},
		{
			Email:    email2,
			Language: "pt-BR",
		},
	}

	ownersResponse := toOwnersResponse(owners)

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

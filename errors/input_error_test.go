package errors

import (
	"testing"
)

func TestNewInputError(t *testing.T) {
	err := NewInputError(ErrorCodeInvalidAuthorization, "Authorization", "abc123")

	if err.Id != ErrorCodeInvalidAuthorization {
		t.Error("Not storing id properly")
	}

	if err.Field != "Authorization" {
		t.Error("Not storing field properly")
	}

	if err.Value != "abc123" {
		t.Error("Not storing value properly")
	}
}

func TestInputErrorError(t *testing.T) {
	err := NewInputError(ErrorCodeInvalidAuthorization, "Authorization", "abc123")
	if err.Error() != "Error 'invalid-authorization' in field 'Authorization' with the value 'abc123'" {
		t.Error("Not reporting the correct message")
	}
}

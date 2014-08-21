package secret

import (
	"testing"
)

func TestKey(t *testing.T) {
	key := key()
	if len(key) != 16 {
		t.Error("Wrong key size!")
	}
}

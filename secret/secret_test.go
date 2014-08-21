package secret

import (
	"math/rand"
	"testing"
	"time"
)

type randomDataMaker struct {
	rand.Source
}

// Reference:
// https://github.com/dustin/randbo
func (r *randomDataMaker) Read(p []byte) (n int, err error) {
	todo := len(p)
	offset := 0
	for {
		val := int64(r.Int63())
		for i := 0; i < 8; i++ {
			p[offset] = byte(val)
			todo--
			if todo == 0 {
				return len(p), nil
			}
			offset++
			val >>= 8
		}
	}
}

func TestEncryptDecrypt(t *testing.T) {
	randomDataMaker := randomDataMaker{
		Source: rand.NewSource(time.Now().UnixNano()),
	}

	for i := 0; i < 1000; i++ {
		input := make([]byte, 128)
		_, err := randomDataMaker.Read(input)
		if err != nil {
			t.Error(err)
		}

		encrypted, err := Encrypt(string(input))
		if err != nil {
			t.Error(err)
		}

		decrypted, err := Decrypt(encrypted)
		if err != nil {
			t.Error(err)
		}

		if decrypted != string(input) {
			t.Errorf("Expected '%s' and got '%s'", string(input), decrypted)
		}
	}
}

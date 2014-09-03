package secret

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/rafaeljusto/shelter/errors"
)

// Encrypt uses the secret to encode the password in the configuration file
func Encrypt(input string) (string, error) {
	block, err := aes.NewCipher(key())
	if err != nil {
		return "", errors.NewSystemError(err)
	}

	iv := make([]byte, block.BlockSize())
	_, err = rand.Read(iv)
	if err != nil {
		return "", errors.NewSystemError(err)
	}

	output := make([]byte, len(input))
	ofbStream := cipher.NewOFB(block, iv)
	ofbStream.XORKeyStream(output, []byte(input))

	buffer := bytes.NewBuffer(iv)
	buffer.Write(output)

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

// Decrypt decodes a encrypted password in configuration file
func Decrypt(input string) (string, error) {
	block, err := aes.NewCipher(key())
	if err != nil {
		return "", errors.NewSystemError(err)
	}

	inputBytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", errors.NewSystemError(err)
	}

	iv := inputBytes[:block.BlockSize()]
	inputBytes = inputBytes[block.BlockSize():]

	output := make([]byte, len(inputBytes))
	ofbStream := cipher.NewOFB(block, iv)
	ofbStream.XORKeyStream(output, inputBytes)

	return string(output), nil
}

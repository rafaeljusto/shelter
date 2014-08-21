package main

import (
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/secret"
	"strings"
)

const (
	operationEncrypt operation = "encrypt"
	operationDecrypt operation = "decrypt"
)

type operation string

var (
	op       string
	password string
)

func init() {
	flag.StringVar(&op, "operation", "", fmt.Sprintf("%s|%s", operationEncrypt, operationDecrypt))
	flag.StringVar(&password, "password", "", "plain or encrypted password")
}

func main() {
	flag.Parse()

	if len(op) == 0 {
		fmt.Println("Operation not informed")
		flag.PrintDefaults()
		return
	}
	op = strings.ToLower(op)

	if len(password) == 0 {
		fmt.Println("Password not informed")
		flag.PrintDefaults()
		return
	}

	if operation(op) == operationEncrypt {
		var err error
		if password, err = secret.Encrypt(password); err == nil {
			fmt.Println("Encrypted password:", password)

		} else {
			fmt.Println(err)
		}

	} else if operation(op) == operationDecrypt {
		var err error
		if password, err = secret.Decrypt(password); err == nil {
			fmt.Println("Decrypted password:", password)

		} else {
			fmt.Println(err)
		}

	} else {
		fmt.Println("Invalid operation")
		flag.PrintDefaults()
	}
}

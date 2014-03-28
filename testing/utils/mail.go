// utils - Features for make the test life easier
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package utils

import (
	"bufio"
	"fmt"
	"net"
	"net/mail"
	"strings"
	"time"
)

// E-mail server created to test if the messages are being contructed and delivered correctly in the
// notification module. The server was created based on the Go-Guerrilla SMTPd from
// https://github.com/flashmob/go-guerrilla

const (
	server = "localhost"
)

func StartMailServer(port int) (chan *mail.Message, chan error, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		return nil, nil, err
	}

	messageChannel := make(chan *mail.Message)
	errorChannel := make(chan error)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				errorChannel <- err
				continue
			}

			mailMessage, err := handleMailClient(conn)
			if err != nil {
				errorChannel <- err
			} else {
				messageChannel <- mailMessage
			}
		}
	}()

	// Give some time for the mail server goes up
	time.Sleep(1 * time.Second)

	return messageChannel, errorChannel, nil
}

func handleMailClient(conn net.Conn) (*mail.Message, error) {
	defer conn.Close()

	var mailMessage string
	processing := true
	phase := 0

	for processing {
		switch phase {
		case 0:
			response := fmt.Sprintf("220 %s SMTP ShelterMail\r\n", server)
			if _, err := conn.Write([]byte(response)); err != nil {
				return nil, err
			}

			phase = 1

		case 1:
			data, err := readMailClient(conn, false)
			if err != nil {
				return nil, err
			}

			data = strings.Trim(data, " \n\r")
			cmd := strings.ToUpper(data)
			response := ""

			switch {
			case strings.Index(cmd, "HELO") == 0:
				response = fmt.Sprintf("250 %s Hello\r\n", server)

			case strings.Index(cmd, "EHLO") == 0:
				response = fmt.Sprintf("250-%s Hello %s[%s]\r\n250-SIZE 131072\r\n250 HELP\r\n",
					server, data[5:], conn.RemoteAddr().String(),
				)

			case strings.Index(cmd, "RCPT TO:") == 0:
				response = "250 Accepted\r\n"

			case strings.Index(cmd, "DATA") == 0:
				phase = 2
				response = "354 Enter message, ending with \".\" on a line by itself\r\n"

			case strings.Index(cmd, "QUIT") == 0:
				processing = false
				response = "221 Bye\r\n"

			default:
				response = "250 Ok\r\n"
			}

			if _, err := conn.Write([]byte(response)); err != nil {
				return nil, err
			}

		case 2:
			data, err := readMailClient(conn, true)
			if err != nil {
				return nil, err
			}

			mailMessage += data

			response := "250 Ok: queued as 12345\r\n"
			if _, err := conn.Write([]byte(response)); err != nil {
				return nil, err
			}

			phase = 1
		}
	}

	return mail.ReadMessage(strings.NewReader(mailMessage))
}

func readMailClient(conn net.Conn, isBody bool) (string, error) {
	endMark := "\r\n"
	if isBody {
		endMark = "\r\n.\r\n"
	}

	data := ""

	if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		return data, err
	}

	reader := bufio.NewReader(conn)

	for {
		partialData, err := reader.ReadString('\n')
		if err != nil {
			return data, err
		}

		data += partialData

		if strings.HasSuffix(data, endMark) {
			break
		}
	}

	return data, nil
}

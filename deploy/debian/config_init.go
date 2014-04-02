// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/rafaeljusto/shelter/config"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	configFilePath       = "/usr/shelter/etc/shelter.conf"
	sampleConfigFilePath = "/usr/shelter/etc/shelter.conf.sample"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Initialize configuration file

		if err := config.LoadConfig(sampleConfigFilePath); err != nil {
			log.Fatalln(err)
			return
		}

		readRESTListeners()

		jsonConfig, err := json.MarshalIndent(config.ShelterConfig, " ", " ")
		if err != nil {
			log.Fatalln(err)
			return
		}

		if err := ioutil.WriteFile(configFilePath, jsonConfig, 0664); err != nil {
			log.Fatalln(err)
			return
		}

	} else {
		// Update current configuration file
	}
}

func readRESTListeners() bool {
	addresses, continueProcessing := readRESTAddresses()
	if !continueProcessing {
		return false
	}

	port, continueProcessing := readRESTPort()
	if !continueProcessing {
		return false
	}

	useTLS, generateCerts, continueProcessing := readRESTTLS()
	if !continueProcessing {
		return false
	}

	if generateCerts {
		hostname, continueProcessing := readRESTCertsParams()
		if !continueProcessing {
			return false
		}

		cmd := exec.Command("/usr/shelter/bin/generate_cert", "--host", hostname)
		if err := cmd.Run(); err != nil {
			log.Println("Error generating certificates. Details:", err)

		} else {
			err := os.MkdirAll("/usr/shelter/etc/keys", os.ModeDir|0600)
			if err != nil {
				log.Println("Error creating certificates directory. Details:", err)

			} else {
				if err := moveFile("/usr/shelter/etc/keys/cert.pem", "cert.pem"); err != nil {
					log.Println("Error copying file cert.pem. Details:", err)
				}

				if err := moveFile("/usr/shelter/etc/keys/key.pem", "key.pem"); err != nil {
					log.Println("Error copying file key.pem. Details:", err)
				}
			}
		}
	}

	config.ShelterConfig.RESTServer.Listeners = []struct {
		IP   string
		Port int
		TLS  bool
	}{}

	for _, address := range addresses {
		config.ShelterConfig.RESTServer.Listeners = append(config.ShelterConfig.RESTServer.Listeners, struct {
			IP   string
			Port int
			TLS  bool
		}{
			IP:   address,
			Port: port,
			TLS:  useTLS,
		})
	}

	return true
}

func readRESTAddresses() ([]string, bool) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return nil, true
	}

	var options []string
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println(err)
			return nil, true
		}

		for _, a := range addrs {
			ip := a.String()
			ip = ip[:strings.Index(ip, "/")]
			options = append(options, fmt.Sprintf("%s", ip))
		}
	}

	overOption := -1
	var selectedOptions []int

	// By default all interfaces are going to be used
	for index, _ := range options {
		selectedOptions = append(selectedOptions, index)
	}

	restInputsDraw := func() {
		writeTitle("REST Configurations", 2, 4)
		writeText("Please select the IP addresses that you want to listen:", 2, 7)
		writeOptions(options, 2, 9)

		_, windowsHeight := termbox.Size()
		writeText("[TAB] Move over options", 2, windowsHeight-4)
		writeText("[SPACE] Select an option", 2, windowsHeight-3)
		writeText("[ENTER] Continue", 2, windowsHeight-2)

		for _, selectedOption := range selectedOptions {
			termbox.SetCell(3, 9+selectedOption, 0x221a, termbox.ColorYellow, termbox.ColorBlue)
		}

		if overOption > -1 {
			termbox.SetCell(2, 9+overOption, rune('['), termbox.ColorYellow, termbox.ColorBlue)
			termbox.SetCell(4, 9+overOption, rune(']'), termbox.ColorYellow, termbox.ColorBlue)
		}
	}

	restInputsAction := func(ev termbox.Event) bool {
		switch ev.Key {
		case termbox.KeyTab:
			// Move to the next option

			overOption += 1
			if overOption >= len(options) {
				overOption = 0
			}

		case termbox.KeySpace:
			// Select the option

			found := false
			for index, selectedOption := range selectedOptions {
				if selectedOption == overOption {
					if len(selectedOptions) == 1 {
						selectedOptions = []int{}
					} else if index == 0 {
						selectedOptions = selectedOptions[index+1:]
					} else if index == len(selectedOptions)-1 {
						selectedOptions = selectedOptions[:index]
					} else {
						selectedOptions = append(selectedOptions[:index], selectedOptions[index+1:]...)
					}
					found = true
				}
			}

			if !found {
				selectedOptions = append(selectedOptions, overOption)
			}

		case termbox.KeyEnter:
			// Finish reading inputs
			return false
		}

		draw(restInputsDraw)
		return true
	}

	if !readInput(restInputsDraw, restInputsAction) {
		return nil, false
	}

	var selectedAddresses []string
	for _, option := range selectedOptions {
		selectedAddresses = append(selectedAddresses, options[option])
	}
	return selectedAddresses, true
}

func readRESTPort() (int, bool) {
	port := "4443_"
	portPosition := 0

	restInputsDraw := func() {
		writeTitle("REST Configurations", 2, 4)
		writeText("Please inform the port that you want to listen:", 2, 7)
		writeText(port, 2, 9)

		if portPosition < len(port) {
			termbox.SetCell(2+portPosition, 9, rune(port[portPosition]), termbox.ColorWhite, termbox.ColorYellow)
		}

		_, windowsHeight := termbox.Size()
		writeText("[ENTER] Continue", 2, windowsHeight-2)
	}

	restInputsAction := func(ev termbox.Event) bool {
		switch ev.Key {
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			portPosition -= 1
			if portPosition < 0 {
				portPosition = 0
			}

			port = port[:portPosition] + "_" + port[portPosition+1:]

		case termbox.KeyDelete:
			if portPosition < len(port) {
				port = port[:portPosition] + port[portPosition+1:] + "_"
			}

		case termbox.KeyEnter:
			// Finish reading inputs
			return false

		default:
			if ev.Ch >= 48 && ev.Ch < 58 && portPosition < len(port) {
				port = port[:portPosition] + string(ev.Ch) + port[portPosition+1:]

				portPosition += 1
				if portPosition > len(port) {
					portPosition = len(port)
				}
			}
		}

		draw(restInputsDraw)
		return true
	}

	if !readInput(restInputsDraw, restInputsAction) {
		return 0, false
	}

	port = strings.Replace(port, "_", "", -1)
	portNumber, _ := strconv.Atoi(port)
	return portNumber, true
}

func readRESTTLS() (useTLS, generateCerts, continueProcessing bool) {
	options := []string{
		"Use TLS on interfaces (HTTPS)",
		"Generate self-signed certificates automatically (valid for 1 year)",
	}

	overOption := -1
	selectedOptions := []int{0, 1}

	restInputsDraw := func() {
		writeTitle("REST Configurations", 2, 4)
		writeText("Please select the following TLS options:", 2, 7)
		writeOptions(options, 2, 9)

		_, windowsHeight := termbox.Size()
		writeText("[TAB] Move over options", 2, windowsHeight-4)
		writeText("[SPACE] Select an option", 2, windowsHeight-3)
		writeText("[ENTER] Continue", 2, windowsHeight-2)

		for _, selectedOption := range selectedOptions {
			termbox.SetCell(3, 9+selectedOption, 0x221a, termbox.ColorYellow, termbox.ColorBlue)
		}

		if overOption > -1 {
			termbox.SetCell(2, 9+overOption, rune('['), termbox.ColorYellow, termbox.ColorBlue)
			termbox.SetCell(4, 9+overOption, rune(']'), termbox.ColorYellow, termbox.ColorBlue)
		}
	}

	restInputsAction := func(ev termbox.Event) bool {
		switch ev.Key {
		case termbox.KeyTab:
			// Move to the next option

			overOption += 1
			if overOption >= len(options) {
				overOption = 0
			}

		case termbox.KeySpace:
			// Select the option

			found := false
			for index, selectedOption := range selectedOptions {
				if selectedOption == overOption {
					if len(selectedOptions) == 1 {
						selectedOptions = []int{}
					} else if index == 0 {
						selectedOptions = selectedOptions[index+1:]
					} else if index == len(selectedOptions)-1 {
						selectedOptions = selectedOptions[:index]
					} else {
						selectedOptions = append(selectedOptions[:index], selectedOptions[index+1:]...)
					}
					found = true
				}
			}

			if !found {
				selectedOptions = append(selectedOptions, overOption)

				// Automatically certificates generation cannot exist withou TLS
				if overOption == 1 && len(selectedOptions) == 1 {
					selectedOptions = append(selectedOptions, 0)
				}
			}

		case termbox.KeyEnter:
			// Finish reading inputs
			return false
		}

		draw(restInputsDraw)
		return true
	}

	if !readInput(restInputsDraw, restInputsAction) {
		continueProcessing = false
		return
	}

	continueProcessing = true
	for _, option := range selectedOptions {
		if option == 0 {
			useTLS = true
		} else if option == 1 {
			generateCerts = true
		}
	}

	return
}

func readRESTCertsParams() (string, bool) {
	host := "localhost_________________________________________"
	hostPosition := 0

	restInputsDraw := func() {
		writeTitle("REST Configurations", 2, 4)
		writeText("Please inform the hostname of the certificate:", 2, 7)
		writeText(host, 2, 9)

		if hostPosition < len(host) {
			termbox.SetCell(2+hostPosition, 9, rune(host[hostPosition]), termbox.ColorWhite, termbox.ColorYellow)
		}

		_, windowsHeight := termbox.Size()
		writeText("[ENTER] Continue", 2, windowsHeight-2)
	}

	restInputsAction := func(ev termbox.Event) bool {
		switch ev.Key {
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			hostPosition -= 1
			if hostPosition < 0 {
				hostPosition = 0
			}

			host = host[:hostPosition] + "_" + host[hostPosition+1:]

		case termbox.KeyDelete:
			if hostPosition < len(host) {
				host = host[:hostPosition] + host[hostPosition+1:] + "_"
			}

		case termbox.KeyEnter:
			if len(strings.Replace(host, "_", "", -1)) == 0 {
				draw(func() {
					restInputsDraw()
					writeText("ERROR: Hostname cannot be empty!", 2, 11)
				})

				return true
			}

			// Finish reading inputs
			return false

		default:
			if ((ev.Ch >= 48 && ev.Ch < 58) || // 0-9
				(ev.Ch >= 65 && ev.Ch < 91) || // A-Z
				(ev.Ch >= 97 && ev.Ch < 123) || // a-z
				ev.Ch == 45 || ev.Ch == 46) && // - .
				hostPosition < len(host) {

				host = host[:hostPosition] + string(ev.Ch) + host[hostPosition+1:]

				hostPosition += 1
				if hostPosition > len(host) {
					hostPosition = len(host)
				}
			}
		}

		draw(restInputsDraw)
		return true
	}

	if !readInput(restInputsDraw, restInputsAction) {
		return "", false
	}

	return strings.Replace(host, "_", "", -1), true
}

func readInput(inputsDraw func(), inputsAction func(termbox.Event) bool) bool {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	draw(inputsDraw)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return false
			}

			if !inputsAction(ev) {
				return true
			}

		case termbox.EventResize:
			draw(inputsDraw)
		}
	}
}

func draw(inputsDraw func()) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	windowWidth, windowsHeight := termbox.Size()

	// Background
	for i := 0; i < windowWidth; i++ {
		for j := 0; j < windowsHeight; j++ {
			termbox.SetCell(i, j, 0x0, termbox.ColorWhite, termbox.ColorBlue)
		}
	}

	// Edges
	termbox.SetCell(0, 0, 0x250C, termbox.ColorWhite, termbox.ColorBlue)
	termbox.SetCell(windowWidth-1, 0, 0x2510, termbox.ColorWhite, termbox.ColorBlue)
	termbox.SetCell(0, windowsHeight-1, 0x2514, termbox.ColorWhite, termbox.ColorBlue)
	termbox.SetCell(windowWidth-1, windowsHeight-1, 0x2518, termbox.ColorWhite, termbox.ColorBlue)

	// Logo
	termbox.SetCell((windowWidth/2)-7, 2, 0xff33, termbox.ColorWhite, termbox.ColorBlue) // S
	termbox.SetCell((windowWidth/2)-5, 2, 0xff28, termbox.ColorWhite, termbox.ColorBlue) // H
	termbox.SetCell((windowWidth/2)-3, 2, 0xff25, termbox.ColorWhite, termbox.ColorBlue) // E
	termbox.SetCell((windowWidth/2)-1, 2, 0xff2c, termbox.ColorWhite, termbox.ColorBlue) // L
	termbox.SetCell((windowWidth/2)+1, 2, 0xff34, termbox.ColorWhite, termbox.ColorBlue) // T
	termbox.SetCell((windowWidth/2)+3, 2, 0xff25, termbox.ColorWhite, termbox.ColorBlue) // E
	termbox.SetCell((windowWidth/2)+5, 2, 0xff32, termbox.ColorWhite, termbox.ColorBlue) // R

	// Footer tip
	writeText("Press ESC to quit", windowWidth-20, windowsHeight-2)

	inputsDraw()

	termbox.Flush()
}

func writeTitle(text string, x, y int) {
	xTmp := x
	for _, character := range text {
		termbox.SetCell(xTmp, y, rune(character), termbox.ColorWhite, termbox.ColorBlue)
		xTmp += 1
	}

	xTmp = x
	for xTmp <= len(text)+1 {
		termbox.SetCell(xTmp, y+1, 0x2015, termbox.ColorWhite, termbox.ColorBlue)
		xTmp += 1
	}
}

func writeText(text string, x, y int) {
	for _, character := range text {
		termbox.SetCell(x, y, rune(character), termbox.ColorWhite, termbox.ColorBlue)
		x += 1
	}
}

func writeOptions(options []string, x, y int) {
	for index, option := range options {
		writeText("[ ] "+option, x, y+index)
	}
}

func moveFile(dst, orig string) error {
	file, err := os.Open(orig)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(dstFile, file); err != nil {
		return err
	}

	if err := os.Remove(orig); err != nil {
		return err
	}

	return nil
}

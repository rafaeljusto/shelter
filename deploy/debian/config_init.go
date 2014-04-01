// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/rafaeljusto/shelter/config"
	"log"
	"net"
	"os"
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

		configFileOutput, err := os.Create(configFilePath)
		if err != nil {
			log.Fatalln(err)
			return
		}
		defer configFileOutput.Close()

		configFileEncoder := json.NewEncoder(configFileOutput)
		if err := configFileEncoder.Encode(config.ShelterConfig); err != nil {
			log.Fatalln(err)
			return
		}

	} else {
		// Update current configuration file
	}
}

func readRESTListeners() {
	readRESTInterfaces()
	readRESTPort()
	readRESTTLS()
}

func readRESTInterfaces() []net.Interface {
	interfaces, err := net.Interfaces()
	if err != nil {
		// TODO: Manual input?!
		return nil
	}

	var options []string
	for _, i := range interfaces {
		var addresses []string

		addrs, err := i.Addrs()
		if err != nil {
			// TODO: Manual input?!
			return nil
		}

		for _, a := range addrs {
			addresses = append(addresses, a.String())
		}
		options = append(options, fmt.Sprintf("%s (%s)", i.Name, strings.Join(addresses, ", ")))
	}

	overOption := -1
	var selectedOptions []int

	restInputsDraw := func() {
		writeTitle("REST Configurations", 2, 4)
		writeText("Please select the interfaces that you want to listen:", 2, 7)
		writeOptions(options, 2, 9)

		_, windowsHeight := termbox.Size()
		writeText("[TAB] Move over options", 2, windowsHeight-4)
		writeText("[SPACE] Select an option", 2, windowsHeight-3)
		writeText("[ENTER] Finish", 2, windowsHeight-2)

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

	readInput(restInputsDraw, restInputsAction)

	var selectedInterfaces []net.Interface
	for _, option := range selectedOptions {
		selectedInterfaces = append(selectedInterfaces, interfaces[option])
	}
	return selectedInterfaces
}

func readRESTPort() int {
	port := "_____"
	portPosition := 0

	restInputsDraw := func() {
		writeTitle("REST Configurations", 2, 4)
		writeText("Please inform the port that you want to listen:", 2, 7)
		writeText(port, 2, 9)

		_, windowsHeight := termbox.Size()
		writeText("[ENTER] Finish", 2, windowsHeight-2)
	}

	restInputsAction := func(ev termbox.Event) bool {
		switch ev.Key {
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			portPosition -= 1
			if portPosition < 0 {
				portPosition = 0
			}

			port = port[:portPosition] + "_" + port[portPosition+1:]

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

	readInput(restInputsDraw, restInputsAction)

	strings.Replace(port, "_", "", -1)
	portNumber, _ := strconv.Atoi(port)
	return portNumber
}

func readRESTTLS() (useTLS, generateCerts bool) {
	options := []string{
		"Use TLS on interfaces (HTTPS)",
		"Generate self-signed certificates automatically (valid for 1 year)",
	}

	overOption := -1
	var selectedOptions []int

	restInputsDraw := func() {
		writeTitle("REST Configurations", 2, 4)
		writeText("Please select the following TLS options:", 2, 7)
		writeOptions(options, 2, 9)

		_, windowsHeight := termbox.Size()
		writeText("[TAB] Move over options", 2, windowsHeight-4)
		writeText("[SPACE] Select an option", 2, windowsHeight-3)
		writeText("[ENTER] Finish", 2, windowsHeight-2)

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

	readInput(restInputsDraw, restInputsAction)

	for _, option := range selectedOptions {
		if option == 0 {
			useTLS = true
		} else if option == 1 {
			generateCerts = true
		}
	}
	return
}

func readInput(inputsDraw func(), inputsAction func(termbox.Event) bool) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	draw(inputsDraw)

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				break loop

			}

			if !inputsAction(ev) {
				break loop
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

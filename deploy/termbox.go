// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package deploy has the necessary structures to build terminal screens
package deploy

import (
	"github.com/nsf/termbox-go"
	"regexp"
	"strconv"
	"strings"
)

var (
	HostnameOrIPInput = regexp.MustCompile("[0-9A-Za-z\\-\\.\\:]")
	HostnameInput     = regexp.MustCompile("[0-9A-Za-z\\-\\.]")
	AlphaNumericInput = regexp.MustCompile("[0-9A-Za-z]")
	NumericInput      = regexp.MustCompile("[0-9]")
	IPRangeInput      = regexp.MustCompile("[0-9a-fA-F\\:\\./]")
)

type Option struct {
	Value    string
	Selected bool
}

func ManageInputOptionsScreen(
	title, description string,
	options []Option,
	checkConsistency func([]Option, int) []Option,
) ([]Option, bool) {

	overOption := -1

	inputsDraw := func() {
		writeTitle(title, 2, 4)
		writeText(description, 2, 7)
		writeOptions(options, 2, 9)

		_, windowsHeight := termbox.Size()
		writeText("[TAB] Move over options", 2, windowsHeight-4)
		writeText("[SPACE] Select an option", 2, windowsHeight-3)
		writeText("[ENTER] Continue", 2, windowsHeight-2)

		for index, option := range options {
			if option.Selected {
				termbox.SetCell(3, 9+index, 0x221a, termbox.ColorYellow, termbox.ColorBlue)
			}
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

			var optionIndex int
			for index, option := range options {
				if index == overOption {
					optionIndex = index
					options[index].Selected = !option.Selected
					break
				}
			}

			if checkConsistency != nil {
				options = checkConsistency(options, optionIndex)
			}

		case termbox.KeyEnter:
			// Finish reading inputs
			return false
		}

		draw(inputsDraw)
		return true
	}

	if !readInput(inputsDraw, restInputsAction) {
		return nil, false
	}

	return options, true
}

func ManageInputTextScreen(
	title, description, input string,
	allowedInput *regexp.Regexp,
	validate func(string) (bool, string),
) (string, bool) {

	inputPosition := 0

	inputsDraw := func() {
		writeTitle(title, 2, 4)
		writeText(description, 2, 7)
		writeText(input, 2, 9)

		if inputPosition < len(input) {
			termbox.SetCell(2+inputPosition, 9, rune(input[inputPosition]), termbox.ColorWhite, termbox.ColorYellow)
		}

		_, windowsHeight := termbox.Size()
		writeText("[ENTER] Continue", 2, windowsHeight-2)
	}

	restInputsAction := func(ev termbox.Event) bool {
		switch ev.Key {
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			inputPosition -= 1
			if inputPosition < 0 {
				inputPosition = 0
			}

			input = input[:inputPosition] + "_" + input[inputPosition+1:]

		case termbox.KeyDelete:
			if inputPosition < len(input) {
				input = input[:inputPosition] + input[inputPosition+1:] + "_"
			}

		case termbox.KeyEnter:
			if validate != nil {
				inputTmp := strings.Replace(input, "_", "", -1)
				valid, msg := validate(inputTmp)

				if !valid {
					draw(func() {
						inputsDraw()
						writeText("ERROR: "+msg, 2, 11)
					})

					return true
				}
			}

			// Finish reading inputs
			return false

		default:
			if allowedInput.MatchString(string(ev.Ch)) &&
				inputPosition < len(input) {

				input = input[:inputPosition] + string(ev.Ch) + input[inputPosition+1:]

				inputPosition += 1
				if inputPosition > len(input) {
					inputPosition = len(input)
				}
			}
		}

		draw(inputsDraw)
		return true
	}

	if !readInput(inputsDraw, restInputsAction) {
		return "", false
	}

	return strings.Replace(input, "_", "", -1), true
}

func ManageInputNumberScreen(title, description string, number string) (int, bool) {
	number, continueProcessing :=
		ManageInputTextScreen(title, description, number, NumericInput,
			func(input string) (bool, string) {
				if len(input) == 0 {
					return false, "You must inform a number"
				}

				return true, ""
			})

	if !continueProcessing {
		return 0, false
	}

	numberConverted, _ := strconv.Atoi(number)
	return numberConverted, true
}

func ManageInputTextOptionsScreen(
	title, description string,
	input string,
	allowedInput *regexp.Regexp,
	validate func(string) (bool, string),
	options []Option,
	checkConsistency func([]Option, int) []Option,
) (string, []Option, bool) {

	inputPosition := 0
	overOption := 0

	inputsDraw := func() {
		writeTitle(title, 2, 4)
		writeText(description, 2, 7)
		writeText(input, 2, 9)
		writeOptions(options, 2, 11)

		_, windowsHeight := termbox.Size()
		writeText("[TAB] Move over options", 2, windowsHeight-4)
		writeText("[SPACE] Select an option", 2, windowsHeight-3)
		writeText("[ENTER] Continue", 2, windowsHeight-2)

		for index, option := range options {
			if option.Selected {
				termbox.SetCell(3, 11+index, 0x221a, termbox.ColorYellow, termbox.ColorBlue)
			}
		}

		if inputPosition < len(input) && overOption == 0 {
			termbox.SetCell(2+inputPosition, 9, rune(input[inputPosition]), termbox.ColorWhite, termbox.ColorYellow)
		}

		if overOption > 0 {
			termbox.SetCell(2, 10+overOption, rune('['), termbox.ColorYellow, termbox.ColorBlue)
			termbox.SetCell(4, 10+overOption, rune(']'), termbox.ColorYellow, termbox.ColorBlue)
		}
	}

	restInputsAction := func(ev termbox.Event) bool {
		switch ev.Key {
		case termbox.KeyTab:
			// Move to the next option

			overOption += 1
			if overOption >= len(options)+1 {
				overOption = 0
			}

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			if overOption != 0 {
				break
			}

			inputPosition -= 1
			if inputPosition < 0 {
				inputPosition = 0
			}

			input = input[:inputPosition] + "_" + input[inputPosition+1:]

		case termbox.KeyDelete:
			if overOption != 0 {
				break
			}

			if inputPosition < len(input) {
				input = input[:inputPosition] + input[inputPosition+1:] + "_"
			}

		case termbox.KeySpace:
			if overOption == 0 {
				break
			}

			// Select the option

			var optionIndex int
			for index, option := range options {
				if index == overOption-1 {
					optionIndex = index
					options[index].Selected = !option.Selected
					break
				}
			}

			if checkConsistency != nil {
				options = checkConsistency(options, optionIndex)
			}

		case termbox.KeyEnter:
			if validate != nil {
				inputTmp := strings.Replace(input, "_", "", -1)
				valid, msg := validate(inputTmp)

				if !valid {
					draw(func() {
						inputsDraw()
						writeText("ERROR: "+msg, 2, 12+len(options))
					})

					return true
				}
			}

			// Finish reading inputs
			return false

		default:
			if overOption != 0 {
				break
			}

			if allowedInput.MatchString(string(ev.Ch)) &&
				inputPosition < len(input) {

				input = input[:inputPosition] + string(ev.Ch) + input[inputPosition+1:]

				inputPosition += 1
				if inputPosition > len(input) {
					inputPosition = len(input)
				}
			}
		}

		draw(inputsDraw)
		return true
	}

	if !readInput(inputsDraw, restInputsAction) {
		return "", nil, false
	}

	return strings.Replace(input, "_", "", -1), options, true
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

func writeOptions(options []Option, x, y int) {
	for index, option := range options {
		writeText("[ ] "+option.Value, x, y+index)
	}
}

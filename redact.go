package redact

import (
	"fmt"
	"strings"
)

type Redaction struct {
	used    []bool
	matches []match
}

func New() *Redaction {
	return &Redaction{
		used:    make([]bool, 256),
		matches: make([]match, 0, 16),
	}
}

func (this *Redaction) All(input string) string {
	//this.matchCreditCard(input)
	this.matchCC(input)
	this.matchEmail(input)
	this.matchPhone(input)
	this.matchSSN(input)
	this.matchDOB(input)
	result := this.redactMatches(input)
	this.clear()
	return result
}
func (this *Redaction) clear() {
	this.matches = this.matches[0:0]
	for i := range this.used {
		this.used[i] = false
	}
}
func (this *Redaction) redactMatches(input string) string {
	if len(this.matches) == 0 {
		return input // no changes to redact
	}

	buffer := []byte(input)
	bufferLength := len(buffer)
	var lowIndex, highIndex int

	for _, match := range this.matches {
		lowIndex = match.InputIndex
		highIndex = lowIndex + match.Length
		if lowIndex < 0 {
			continue
		}
		if highIndex >= bufferLength {
			continue
		}
		for ; lowIndex < highIndex; lowIndex++ {
			buffer[lowIndex] = '*'
		}
	}

	return string(buffer)
}

func (this *Redaction) matchCC(input string) {
	var lastDigit int
	var length int
	var numbers int
	var isOdd bool
	var isCandidate bool
	var total int
	for i := len(input) - 1; i > 0; i-- {
		character := input[i]
		strChar := fmt.Sprintf("%c", character)
		_ = strChar

		if !isNumeric(input[i]) {
			if numbers > 12 && total%10 == 0 {
				this.appendMatch(lastDigit-length+1, length)
				length = 0
				total = 0
				isOdd = false
				lastDigit = i - 1
				numbers = 0
				continue
			}

			if ccBreakNotFound(character) && !isNumeric(input[i-1]) {
				lastDigit = i - 1
				length = 0
				isCandidate = false
				continue
			}
			if i < len(input)-1 && !ccBreakNotFound(input[i+1]) {
				continue
			}
			length++
		} else {
			isOdd = !isOdd
			numbers++
			number := int(character - '0')
			if !isOdd {
				number += number
				if number > 9 {
					number -= 9
				}
			}
			total += number

			if isCandidate {
				length++
			} else {
				isCandidate = true
				lastDigit = i
				numbers = 1
			}
		}
	}
}

func (this *Redaction) matchEmail(input string) {
	var start int
	var length int
	for i := 0; i < len(input); i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !emailBreakNotFound(character) {
			start = i + 1
			length = 0
			continue
		} else {
			if character == '@' {
				this.appendMatch(start, length)
				start = i + 1
				length = 0
			}
			length++
		}
	}
}

func (this *Redaction) matchPhone(input string) {
	var start int
	var length int
	var isCandidate bool
	for i := 0; i < len(input)-1; i++ {
		if this.used[i] {
			continue
		}
		character := input[i]
		if !isNumeric(character) {
			if character == '+' {
				start = i + 2
				length--
				continue
			}
			if phoneBreakNotFound(character) {
				start = i + 1
				length = 0
				isCandidate = false
				continue
			}
			if isPhoneNumber(length) {
				this.appendMatch(start, length)
				length = 0
				start = i + 1
				continue
			}
		}
		if isCandidate {
			length++
		} else {
			isCandidate = true
			start = i + 1
			length = 0
		}
	}
	if isNumeric(input[len(input)-1]) {
		length++
	}
	if isPhoneNumber(length) {
		this.appendMatch(start, length)
	}
}
func isPhoneNumber(length int) bool {
	return length >= 10 && length <= 14
}

func (this *Redaction) matchSSN(input string) {
	var start int
	var length int
	var isCandidate bool
	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !isNumeric(character) {
			if isSSN(length) {
				this.appendMatch(start, length)
				length = 0
				start = i + 1
				continue
			}
			if ssnBreakNotFound(character) {
				start = i + 1
				length = 0
				isCandidate = false
				continue
			}
		}
		if isCandidate {
			length++
		} else {
			isCandidate = true
			start = i + 1
		}
	}
	if isNumeric(input[len(input)-1]) {
		// length++ // TODO: test coverage
	}
	if isSSN(length) {
		// matches = appendMatches(matches, start, length) // TODO: test coverage
	}
}
func isSSN(length int) bool {
	return length >= 9 && length <= 11
}

func (this *Redaction) matchDOB(input string) {
	var start int
	var length int
	var isCandidate bool
	var monthStart int
	var monthLength int
	var monthCandidate bool
	//var startChar byte

	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !isNumeric(character) {
			if dobBreakNotFound(character) {
				monthLength++
				start = i + 1
				length = 0
				isCandidate = false
				continue
			}
			// TODO: this is allocating another string. Don't do dat.
			// Instead, pass in the starting index and length and find out if it's a month
			if isMonth(input[monthStart : monthStart+monthLength]) {
				monthCandidate = true
				continue
			}
			//if monthCheck(startChar, character) {
			//	monthCandidate = true
			//	continue
			//}
			monthStart = i + 1
			//startChar = input[i + 1]
			monthLength = 0
			if isDOB(length) {
				this.appendMatch(start, length)
				length = 0
				start = i + 1
				continue
			}
		}
		if isCandidate || monthCandidate {
			length++
		} else {
			isCandidate = true
			start = i + 1
		}
		if length == 2 && monthCandidate {
			this.appendMatch(monthStart, monthLength+length+1)
			monthCandidate = false
			length = 0
			start = 0
			monthStart = 0
			monthLength = 0
		}
	}
	if isNumeric(input[len(input)-1]) {
		// length++ // TODO: test coverage
	}
	if isDOB(length) {
		// matches = appendMatches(matches, start, length) // TODO: test coverage
	}
}
func isDOB(length int) bool {
	return length >= 6 && length <= 10
}
func isMonth(month string) bool {
	_, found := months[strings.ToLower(month)]
	return found
}
func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

func ccBreakNotFound(character byte) bool {
	return character != '-' && character != ' '
}
func emailBreakNotFound(character byte) bool {
	return character != '.' && character != ' '
}
func phoneBreakNotFound(character byte) bool {
	return character != '-' && character != ' ' && character != '(' && character != ')'
}
func ssnBreakNotFound(character byte) bool {
	return character != '-' && character != ' '
}
func dobBreakNotFound(character byte) bool {
	return character != '-' && character != ' ' && character != '/'
}
func (this *Redaction) appendMatch(start, length int) {
	for i := start; i <= start+length; i++ {
		this.used[i] = true
	}

	this.matches = append(this.matches, match{InputIndex: start, Length: length})
}

type match struct {
	InputIndex int
	Length     int
}

var months = map[string]struct{}{
	"january":   {},
	"jan":       {},
	"february":  {},
	"feb":       {},
	"march":     {},
	"mar":       {},
	"april":     {},
	"apr":       {},
	"may":       {},
	"june":      {},
	"jun":       {},
	"july":      {},
	"jul":       {},
	"august":    {},
	"aug":       {},
	"september": {},
	"sep":       {},
	"sept":      {},
	"october":   {},
	"oct":       {},
	"november":  {},
	"nov":       {},
	"december":  {},
	"dec":       {},
}

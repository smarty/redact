package redact

import "strings"

type Redaction struct {
	used    map[int]struct{}
	matches []match
}

func New() *Redaction {
	return &Redaction{
		used:    make(map[int]struct{}, 128),
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
	for key, _ := range this.used {
		delete(this.used, key)
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

func (this *Redaction) matchCC(input string){
	var lastDigit int
	var length int
	var numbers int
	var isOdd bool
	var isCandidate bool
	var total int
	for i := len(input) - 1; i > 0; i-- {
		character := input[i]

		if !isNumeric(input[i]) {
			if numbers > 12 && total % 10 == 0 {
				this.appendMatch(lastDigit-length + 1, length)
				length = 0
				total = 0
				isOdd = false // TODO
				lastDigit = i - 1
				numbers = 0
				isCandidate = false
				continue
			}

			if breakNotFound(character) && !isNumeric(input[i-1]) {
				lastDigit = i - 1
				length = 0
				isCandidate = false
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

func (this *Redaction) matchCreditCard(input string){
	var start int
	var length int
	var isCandidate bool
	var card []int
	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if !isNumeric(character) {
			// TODO: inline isCreditCard and checkLuhn into method--this avoids having to create a string to ask if it's a credit card
			// instead we track each numeric digit here and run a tally as we go along

			// TODO: Iterate backwards to be able to inline luhnfun
			if checkCard(length, card) {
				this.appendMatch(start, length)
				length = 0
				start = i + 1
				card = card[0:0]
				continue
			}
			if breakNotFound(character) && !isNumeric(input[i+1]) {
				start = i + 1
				length = 0
				isCandidate = false
				card = card[0:0]
				continue
			}
			length++
		} else {
			number := character - '0'
			card = append(card, int(number))
			if isCandidate {
				length++
			} else {
				isCandidate = true
				start = i
			}
		}
	}
	if isNumeric(input[len(input)-1]) {
		//length++ // TODO: test coverage
	}
	if checkCard(length, card) {
		// matches = appendMatches(matches, start, length) // TODO: test coverage
	}

}
func checkCard(length int, card []int) bool {
	return length >= 13 && length <= 24 && luhnCheck(card)
}

func luhnCheck(card []int) bool {
	var total int
	var isOdd bool
	for i := len(card) - 1; i >= 0; i-- {
		isOdd = !isOdd
		num := card[i]
		if !isOdd {
			num = card[i] * 2
			if num > 9 {
				num -= 9
			}
		}
		total += num
	}
	return total%10 == 0
}

func (this *Redaction) matchEmail(input string) {
	var start int
	var length int
	var isCandidate bool
	for i := 0; i < len(input); i++ {
		character := input[i]
		if _, found := this.used[i]; found {
			continue
		}
		if !breakNotFound(character) {
			if isCandidate {
				this.appendMatch(start, length)
			}
			start = i + 1
			length = 0
			isCandidate = false
			continue
		} else {
			length++
			if character == '@' {
				isCandidate = true
			}
		}
	}
}

func (this *Redaction) matchPhone(input string) {
	var start int
	var length int
	var isCandidate bool
	for i := 0; i < len(input)-1; i++ {
		if _, found := this.used[i]; found {
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
		if _, found := this.used[i]; found {
			continue
		}
		if !isNumeric(character) {
			if isSSN(length) {
				this.appendMatch(start, length)
				length = 0
				start = i + 1
				continue
			}
			if breakNotFound(character) {
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

	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if _, found := this.used[i]; found {
			continue
		}
		if !isNumeric(character) {
			if breakNotFound(character) {
				monthLength++
				start = i + 1
				length = 0
				isCandidate = false
				continue
			}
			// TODO: this is allocating another string. Don't do dat.
			// Instead, pass in the starting index and length and found out if it's a month
			if isMonth(input[monthStart : monthStart+monthLength]) {
				monthCandidate = true
				continue
			}
			monthStart = i + 1
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
func isMonth2(input string, start, length int) bool {
	_, found := months[input[start:start+length]]
	return found
}
func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

func breakNotFound(character byte) bool {
	return character != '-' && character != ' ' && character != '.' && character != '/'
}
func phoneBreakNotFound(character byte) bool {
	return breakNotFound(character) && character != '(' && character != ')'
}
func (this *Redaction) appendMatch(start, length int) {
	for i := start; i <= start+length; i++ {
		this.used[i] = struct{}{}
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

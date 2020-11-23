package redact

import (
	"strings"
)

func All(input string) string {
	var matches []match
	matches = append(matches, matchCreditCard(input)...)
	matches = append(matches, matchEmail(input)...)
	matches = append(matches, matchPhoneNum(input)...)
	matches = append(matches, matchSSN(input)...)
	matches = append(matches, matchDOB(input)...)
	return redactMatches(input, matches)
}

func redactMatches(input string, matches []match) string {
	for _, match := range matches {
		replace := input[match.InputIndex : match.InputIndex+match.Length]
		input = strings.ReplaceAll(input, replace, strings.Repeat("*", match.Length))
	}
	return input
}

func matchCreditCard(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool
	var total int
	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if !isNumeric(character) {
			if isCreditCard(length, input[start:start+length]) {
				matches = appendMatches(matches, start, length)
				length = 0
				start = i + 1
				total = 0
				continue
			}
			if breakNotFound(character) && !isNumeric(input[i+1]) {
				start = i + 1
				length = 0
				isCandidate = false
				continue
			}
			length++
		} else {
			number := int(character - '0')
			if isCandidate {
				length++
				total += number
			} else {
				isCandidate = true
				start = i
				total = number
			}
		}
	}
	if isNumeric(input[len(input)-1]) {
		//length++ // TODO: test coverage
	}
	if isCreditCard(length, input[start:start+length]) {
		// matches = appendMatches(matches, start, length) // TODO: test coverage
	}

	return matches
}
func isCreditCard(length int, input string) bool {
	return length >= 13 && length <= 24 && checkLuhn(input)
}
func checkLuhn(input string) bool {
	nDigits := len(input)
	var nSum int
	var isSecond bool
	for i := nDigits - 2; i >= 0; i-- {
		if !isNumeric(input[i]) {
			continue
		}
		d := input[i] - '0'
		if isSecond == false {
			d = d * 2
		}
		if d > 9 {
			d -= 9
		}
		digit := int(d)
		nSum += digit
		isSecond = !isSecond
	}
	temp := int(input[nDigits-1] - '0')
	mod := nSum % 10
	return mod == temp
}

func matchEmail(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool
	for i := 0; i < len(input); i++ {
		character := input[i]
		if _, found := used[i]; found {
			continue
		}
		if !breakNotFound(character) {
			if isCandidate {
				matches = appendMatches(matches, start, length)
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
	return matches
}

func matchPhoneNum(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool
	for i := 0; i < len(input)-1; i++ {
		if _, found := used[i]; found {
			continue
		}
		character := input[i]
		if !isNumeric(character) {
			if character == '+' {
				start = i + 2
				length--
				continue
			}
			if breakNotFound(character) {
				start = i + 1
				length = 0
				isCandidate = false
				continue
			}
			if isPhoneNumber(length) {
				matches = appendMatches(matches, start, length)
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
		matches = appendMatches(matches, start, length)
	}
	return matches
}
func isPhoneNumber(length int) bool {
	return length >= 10 && length <= 14
}
func matchSSN(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool
	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if _, found := used[i]; found {
			continue
		}
		if !isNumeric(character) {
			if isSSN(length) {
				matches = appendMatches(matches, start, length)
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
	return matches
}

func isSSN(length int) bool {
	return length >= 9 && length <= 11
}
func matchDOB(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool
	var monthStart int
	var monthLength int
	var monthCandidate bool

	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if _, found := used[i]; found {
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
			if isMonth(input[monthStart : monthStart+monthLength]) {
				monthCandidate = true
				continue
			}
			monthStart = i + 1
			monthLength = 0
			if isDOB(length) {
				matches = appendMatches(matches, start, length)
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
			matches = appendMatches(matches, monthStart, monthLength+length+1)
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
	return matches
}

func isDOB(length int) bool {
	return length >= 6 && length <= 10
}
func isMonth(month string) bool {
	_, found := months[month]
	return found
}
func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

func breakNotFound(character byte) bool {
	return character != '-' && character != ' ' && character != '.' && character != '(' && character != ')' && character != '/'
}
func appendMatches(matches []match, start, length int) []match {
	for i := start; i <= start+length; i++ {
		used[i] = struct{}{}
	}
	return append(matches, match{InputIndex: start, Length: length})
}

type match struct {
	InputIndex int
	Length     int
}

var (
	used   = make(map[int]struct{})
	months = map[string]struct{}{
		"January":   {},
		"Jan":       {},
		"February":  {},
		"Feb":       {},
		"March":     {},
		"Mar":       {},
		"April":     {},
		"Apr":       {},
		"May":       {},
		"June":      {},
		"Jun":       {},
		"July":      {},
		"Jul":       {},
		"August":    {},
		"Aug":       {},
		"September": {},
		"Sep":       {},
		"Sept":      {},
		"October":   {},
		"Oct":       {},
		"November":  {},
		"Nov":       {},
		"December":  {},
		"Dec":       {},
	}
)

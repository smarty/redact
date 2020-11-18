package redact

import (
	"regexp"
	"strconv"
	"strings"
)

func All(input string) string {
	var matches []match
	matches = append(matches, matchCreditCard(input)...)
	matches = append(matches, matchEmail(input)...)
	matches = append(matches, matchPhoneNum(input)...)
	return redactMatches(input, matches)
}

func redactMatches(input string, matches []match) string {
	for _, match := range matches {
		input = strings.ReplaceAll(input, input[match.InputIndex:match.InputIndex+match.Length], strings.Repeat("*", match.Length))
	}
	return input
}

//
//func DateOfBirth(input string) string {
//	found := findDOB(input)
//	for _, item := range found {
//		input = strings.ReplaceAll(input, item, "[DOB REDACTED]")
//	}
//	return input
//}
//func findDOB(input string) (dates []string) {
//	dates = append(dates, datesWithDashes.FindAllString(input, 1000000)...)
//	dates = append(dates, datesWithSlashes.FindAllString(input, 1000000)...)
//	dates = append(dates, datesInUSFormat.FindAllString(input, 1000000)...)
//	dates = append(dates, datesInEUFormat.FindAllString(input, 1000000)...)
//
//	return dates
//}
//

func matchEmail(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool // @
	for i := range input {
		character := input[i]
		if !breakNotFound(character) {
			if isCandidate {
				matches = append(matches, match{InputIndex: start, Length: length})
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

func matchCreditCard(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool
	var total int
	for i := 0; i < len(input); i++ {
		character := input[i]
		if !isNumeric(character) {
			if isCreditCard(length, input[start: start + length]) {
				matches = append(matches, match{InputIndex: start, Length: length})
				length = 0
				start = i + 1
				total = 0
				continue
			}
			if breakNotFound(character) {
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

	if isCreditCard(length, input[start:len(input) - 1]) {
		matches = append(matches, match{InputIndex: start, Length: length})
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
		if d > 9{
			d -= 9
		}
		digit := int(d)
		nSum += digit
		isSecond = !isSecond
	}
	temp := int(input[nDigits - 1] - '0')
	mod := nSum%10
	return  mod == temp
}

func breakNotFound(character byte) bool {
	return character != '-' && character != ' ' && character != '.' && character != '(' && character != ')'
}

func matchPhoneNum(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool
	for i := 0; i < len(input); i++ {
		character := input[i]
		if !isNumeric(character) {
			if isPhoneNumber(length) {
				matches = append(matches, match{InputIndex: start, Length: length})
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
	if isPhoneNumber(length) {
		matches = append(matches, match{InputIndex: start, Length: length})
	}
	return matches
}
func isPhoneNumber(length int) bool {
	return length >= 10 && length <= 14
}

//func SSN(input string) string {
//	found := findSSN(input)
//	for _, item := range found {
//		input = strings.ReplaceAll(input, item, "[SSN REDACTED]")
//	}
//	return input
//}
//func findSSN(input string) (SSNs []string) {
//	var temp string
//	for i, character := range input {
//		if breakNotFound(character) {
//			temp = fmt.Sprintf("%s%c", temp, character)
//		} else {
//			if spaceDelimitedCandidate(input, i) {
//				temp += " "
//			} else {
//				// Find index and store location
//				appendCandidate(temp, &SSNs, 9, 11)
//				temp = ""
//			}
//		}
//	}
//	return SSNs
//}
//
//func Phone(input string) string {
//	found := findPhone(input)
//	for _, item := range found {
//		input = strings.ReplaceAll(input, item, "[PHONE REDACTED]")
//	}
//	return input
//}
//func findPhone(input string) (telNums []string) {
//	var temp string
//	for i, character := range input {
//		if breakNotFound(character) {
//			temp = fmt.Sprintf("%s%c", temp, character)
//		} else {
//			if spaceDelimitedCandidate(input, i) {
//				temp += " "
//			} else {
//				appendTelCandidate(temp, &telNums, 10, 16)
//				temp = ""
//			}
//		}
//	}
//	return telNums
//}

func appendCandidate(temp string, items *[]string, min, max int) {
	lengthTemp := len(temp)
	if lengthTemp >= min && lengthTemp <= max {
		new := strings.ReplaceAll(temp, "-", "")
		new = strings.ReplaceAll(new, " ", "")
		new = strings.ReplaceAll(new, "(", "")
		new = strings.ReplaceAll(new, ")", "")
		if _, err := strconv.Atoi(new); err == nil {
			*items = append(*items, temp)
		}
	}
}
func appendTelCandidate(temp string, items *[]string, min, max int) {
	lengthTemp := len(temp)
	if lengthTemp == 11 && strings.Contains(temp, "-") {
		return
	}
	if resemblesCC(temp, lengthTemp) {
		return
	}
	appendCandidate(temp, items, min, max)
}
func resemblesCC(temp string, lengthTemp int) bool {
	return lengthTemp > 12 && !((strings.Contains(temp, "(") || strings.Contains(temp, ")")) || strings.Contains(temp, " ") || strings.Contains(temp, "-"))
}
func spaceDelimitedCandidate(input string, i int) bool {
	return i+1 < len(input) && ('0' <= input[i+1] && input[i+1] <= '9') && ('0' <= input[i-1] && input[i-1] <= '9' || input[i-1] == ')')
}
func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

type match struct {
	InputIndex int
	Length     int
}

var (
	datesWithDashes  = regexp.MustCompile(`\d{1,2}-\d{1,2}-\d{4}`)
	datesWithSlashes = regexp.MustCompile(`\d{1,2}/\d{1,2}/\d{4}`)
	datesInUSFormat  = regexp.MustCompile(`[a-zA-Z]{3,9} \d{1,2}, \d{4}`)
	datesInEUFormat  = regexp.MustCompile(`\d{1,2} [a-zA-Z]{3,9} \d{4}`)
	Separator        = map[int32]struct{}{'.': {}, '/': {}, ' ': {}, '-': {}}
)

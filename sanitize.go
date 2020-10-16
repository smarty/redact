package sanitize

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func RedactDateOfBirth(input string) string {
	foundDOB := FindDOB(input)
	for _, DOB := range foundDOB {
		input = strings.ReplaceAll(input, DOB, "[DOB REDACTED]")
	}
	return input
}

func FindDOB(input string) (dates []string) {
	dash := regexp.MustCompile(`\d{1,2}-\d{1,2}-\d{4}`)
	slash := regexp.MustCompile(`\d{1,2}/\d{1,2}/\d{4}`)
	standardUS := regexp.MustCompile(`[a-zA-Z]{3,9} \d{1,2}, \d{4}`)
	standardEU := regexp.MustCompile(`\d{1,2} [a-zA-Z]{3,9} \d{4}`)

	dates = append(dates, dash.FindAllString(input, 1000000)...)
	dates = append(dates, slash.FindAllString(input, 1000000)...)
	dates = append(dates, standardUS.FindAllString(input, 1000000)...)
	dates = append(dates, standardEU.FindAllString(input, 1000000)...)

	return dates
}

//////////////////////////////////////////////////////////////////

func RedactEmail(input string) string {
	foundEmails := FindEmails(input)
	for _, email := range foundEmails {
		input = strings.ReplaceAll(input, email, "[EMAIL REDACTED]")
	}
	return input
}

func FindEmails(input string) (emails []string) {
	var temp string
	var email bool
	for i, character := range input {
		if breakNotFound(character) {
			temp = fmt.Sprintf("%s%c", temp, character)
		} else {
			if email {
				emails = append(emails, buildEmail(input, temp, i))
				email = false
			}
			temp = ""
			continue
		}
		if character == '@' {
			email = true
		}
	}
	return emails
}
func buildEmail(input, temp string, i int) string {
	return fmt.Sprintf("%s%c%c%c%c", temp, input[i], input[i+1], input[i+2], input[i+3])
}

//////////////////////////////////////////////////////////////////

func RedactCreditCard(input string) string {
	cards := FindCreditCards(input)
	return SanitizeCreditCard(cards, input)
}

func FindCreditCards(input string) (cards []string) {
	var temp string
	for i, character := range input {
		if breakNotFound(character) {
			temp = fmt.Sprintf("%s%c", temp, character)
		} else {
			if spaceDelimitedCandidate(input, i) {
				temp += " "
			} else {
				appendCandidate(temp, &cards, 13, 19)
				temp = ""
			}
		}
	}
	return cards
}

func SanitizeCreditCard(cards []string, input string) string {
	for _, card := range cards {
		new := strings.ReplaceAll(card, "-", "")
		new = strings.ReplaceAll(new, " ", "")
		new = fmt.Sprintf("[MASKED %s****%s]", new[:4], new[12:])
		input = strings.ReplaceAll(input, card, new)
	}
	return input
}

//////////////////////////////////////////////////////////////////

func RedactSSN(input string) string {
	cards := FindSSN(input)
	return SanitizeSSN(cards, input)
}

func FindSSN(input string) (SSNs []string) {
	var temp string
	for i, character := range input {
		if breakNotFound(character) {
			temp = fmt.Sprintf("%s%c", temp, character)
		} else {
			if spaceDelimitedCandidate(input, i) {
				temp += " "
			} else {
				appendCandidate(temp, &SSNs, 9, 11)
				temp = ""
			}
		}
	}
	return SSNs
}

func SanitizeSSN(SSNs []string, input string) string {
	for _, SSN := range SSNs {
		input = strings.ReplaceAll(input, SSN,"[SSN REDACTED]")
	}
	return input
}

//////////////////////////////////////////////////////////////////

func RedactTel(input string) string {
	telNums := FindTel(input)
	return SanitizeTel(telNums, input)
}

func FindTel(input string) (telNums []string) {
	var temp string
	for i, character := range input {
		if breakNotFound(character) {
			temp = fmt.Sprintf("%s%c", temp, character)
		} else {
			if spaceDelimitedCandidate(input, i) {
				temp += " "
			} else {
				appendCandidate(temp, &telNums, 10, 16)
				temp = ""
			}
		}
	}
	return telNums
}

func SanitizeTel(telNums []string, input string) string {
	for _, num := range telNums {
		input = strings.ReplaceAll(input, num,"[TEL REDACTED]")
	}
	return input
}

//////////////////////////////////////////////////////////////////

func RedactAll(input string) string{
	input = RedactEmail(input)
	input = RedactDateOfBirth(input)
	input = RedactTel(input)
	input = RedactSSN(input)
	input = RedactCreditCard(input)
	return input
}

func breakNotFound(character int32) bool {
	return character != ' ' && character != '.' && character != ','
}

func appendCandidate(temp string, items *[]string, min, max int) {
	lengthTemp := len(temp)
	if lengthTemp == 11 && strings.Contains(temp, "-"){
		return
	}
	if lengthTemp > 12 && (strings.Contains(temp, "(") || strings.Contains(temp, ")")){
		return
	}
	if lengthTemp >= min && lengthTemp <= max { // SSN 9-11 //NUM 10-16 //CC 13-19
		new := strings.ReplaceAll(temp, "-", "")
		new = strings.ReplaceAll(new, " ", "")
		new = strings.ReplaceAll(new, "(", "")
		new = strings.ReplaceAll(new, ")", "")
		if _, err := strconv.Atoi(new); err == nil {
			*items = append(*items, temp)
		}
	}
}

func spaceDelimitedCandidate(input string, i int) bool {
	return i+1 < len(input) && ('0' <= input[i+1] && input[i+1] <= '9') && ('0' <= input[i-1] && input[i-1] <= '9' || input[i-1] == ')')
}
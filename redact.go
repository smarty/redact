package redact

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func All(input string) string {
	input = Email(input)
	input = DateOfBirth(input)
	input = Phone(input)
	input = SSN(input)
	input = CreditCard(input)
	return input
}

func DateOfBirth(input string) string {
	found := findDOB(input)
	for _, item := range found {
		input = strings.ReplaceAll(input, item, "[DOB REDACTED]")
	}
	return input
}
func findDOB(input string) (dates []string) {
	dates = append(dates, datesWithDashes.FindAllString(input, 1000000)...)
	dates = append(dates, datesWithSlashes.FindAllString(input, 1000000)...)
	dates = append(dates, datesInUSFormat.FindAllString(input, 1000000)...)
	dates = append(dates, datesInEUFormat.FindAllString(input, 1000000)...)

	return dates
}

func Email(input string) string {
	found := findEmails(input)
	for _, item := range found {
		// TODO: Display domain name
		input = strings.ReplaceAll(input, item, "[EMAIL REDACTED]")
	}
	return input
}
func findEmails(input string) (emails []string) {
	var temp string
	var email bool
	for i, character := range input {
		if breakNotFound(character) {
			// TODO: Avoid allocations
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
	// TODO: .co.uk (TLDs)
	return fmt.Sprintf("%s%c%c%c%c", temp, input[i], input[i+1], input[i+2], input[i+3])
}

func CreditCard(input string) string {
	cards := findCreditCards(input)
	return sanitizeCreditCard(cards, input)
}
func findCreditCards(input string) (cards []string) { // TODO: LUHN check
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
func sanitizeCreditCard(cards []string, input string) string {
	for _, card := range cards {
		new := strings.ReplaceAll(card, "-", "")
		new = strings.ReplaceAll(new, " ", "")
		// TODO: add an if statement for length check & change to only output the last 4 not 12-on
		// TODO: must pass LUHN/MOD10 algorithm
		new = fmt.Sprintf("[CARD %s****%s]", new[:4], new[len(new)-4:])
		input = strings.ReplaceAll(input, card, new)
	}
	return input
}

func SSN(input string) string {
	found := findSSN(input)
	for _, item := range found {
		input = strings.ReplaceAll(input, item, "[SSN REDACTED]")
	}
	return input
}
func findSSN(input string) (SSNs []string) {
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

func Phone(input string) string {
	found := findPhone(input)
	for _, item := range found {
		input = strings.ReplaceAll(input, item, "[PHONE REDACTED]")
	}
	return input
}
func findPhone(input string) (telNums []string) {
	var temp string
	for i, character := range input {
		if breakNotFound(character) {
			temp = fmt.Sprintf("%s%c", temp, character)
		} else {
			if spaceDelimitedCandidate(input, i) {
				temp += " "
			} else {
				appendTelCandidate(temp, &telNums, 10, 16)
				temp = ""
			}
		}
	}
	return telNums
}

func breakNotFound(character int32) bool {
	return character != ' ' && character != '.' && character != ',' && character != '<'
}
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

var (
	datesWithDashes  = regexp.MustCompile(`\d{1,2}-\d{1,2}-\d{4}`)
	datesWithSlashes = regexp.MustCompile(`\d{1,2}/\d{1,2}/\d{4}`)
	datesInUSFormat  = regexp.MustCompile(`[a-zA-Z]{3,9} \d{1,2}, \d{4}`)
	datesInEUFormat  = regexp.MustCompile(`\d{1,2} [a-zA-Z]{3,9} \d{4}`)
)

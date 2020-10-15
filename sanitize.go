package sanitize

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func FindDOB(input string) (dates []string){
	// TODO: define regular expressions as package-level variables (compilation only needs to happen once instead of every time FindDOB is called)
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

func FindEmails(input string)  (emails []string){
	var temp string
	var email bool
	for i, character := range input {
		if character ==  ' ' || character == '.'{
			if email{
				emails = append(emails, fmt.Sprintf("%s%c%c%c%c", temp, input[i], input[i+1], input[i+2], input[i+3]))
				email = false
			}
			temp = ""
			continue
		}
		if character == '@'{
			email = true
		}
		temp = fmt.Sprintf("%s%c", temp, character)
	}
	return emails
}

func RedactDateOfBirth(input string) string {
	foundDOB := FindDOB(input)
	for _, DOB := range foundDOB{
		input = strings.ReplaceAll(input, DOB, "[DOB REDACTED]")
	}
	return input
}

func RedactEmail(input string) string {
	foundEmails := FindEmails(input)
	for _, email := range foundEmails{
		input = strings.ReplaceAll(input, email, "[EMAIL REDACTED]")
	}
	return input
}

func RedactPhone(phone string) string {
	return  redact("TEL")
}

func RedactSSN(SSN string) string {
	return  redact("SSN")
}

func RedactCreditCard(input string) string {
	//Find CC and then sanitize
	cards := FindCreditCards(input)
	return SanitizeCreditCard(cards, input)
}

func FindCreditCards(input string) (cards []string){
	var temp string
	for _, character := range input {
		if character == ' ' || character == '.' {
			lengthTemp := len(temp)
			if lengthTemp >= 13 && lengthTemp <= 19 {
				new := strings.ReplaceAll(temp, "-", "")
				if _, err := strconv.Atoi(new) ; err == nil {
				cards = append(cards, temp)
				}
			}
			temp = ""
			continue
		}
		temp = fmt.Sprintf("%s%c", temp, character)
	}
	return cards
}

func SanitizeCreditCard(cards []string, input string) string {
	for _, card := range cards{
		new := strings.ReplaceAll(card, "-", "")
		new = fmt.Sprintf("[MASKED %s****%s]", new[:4], new[12:])
		input = strings.ReplaceAll(input, card, new)
	}
	return input
}

func redact(statement string) string {
	return fmt.Sprintf("[%s REDACTED]", statement)
}
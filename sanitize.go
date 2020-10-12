package sanitize

import (
	"fmt"
	"regexp"
	"strings"
)

var months = []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"}

func FindDOB(input string) (dates []string){
	dash := regexp.MustCompile(`\d{1,2}-\d{1,2}-\d{4}`)
	slash := regexp.MustCompile(`\d{1,2}/\d{1,2}/\d{4}`)
	standardUS := regexp.MustCompile(`[a-zA-Z]{3,9} \d{1,2}, \d{4}`)
	standardEU := regexp.MustCompile(`\d{1,2} [a-zA-Z]{3,9} \d{4}`)

	dates = append(dates, dash.FindAllString(input, 1000000)...)
	dates = append(dates, slash.FindAllString(input, 1000000)...)
	dates = append(dates, standardUS.FindAllString(input, 10000)...)
	dates = append(dates, standardEU.FindAllString(input, 10000)...)

	return dates
}
//s{3,9} \d{1,2},\d{4}
func standardDOB(input string) (dates []string){

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

func RedactCreditCard(card string) string {
	card = strings.ReplaceAll(card, "-", "")
	card = strings.ReplaceAll(card, " ", "")
	return fmt.Sprintf("[MASKED %s****%s]" , card[:4], card[12:])
}

func redact(statement string) string {
	return fmt.Sprintf("[%s REDACTED]", statement)
}


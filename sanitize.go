package sanitize

import (
	"fmt"
	"strings"
)


func FindEmails(input string) (emails []string) {
	//given a string we want to find the @ and go from the last space and find the next space
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



func RedactDateOfBirth(dateOfBirth string) string {
	return redact("DOB")
}

func RedactEmail(email string) string {
	return  redact("EMAIL")
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


package redact

import (
	"fmt"
	"strings"
)

func (this *creditCardRedact) match(input []byte) {
	if len(input) <= 0 {
		return
	}

	//input := 'Hi, my name is 24601 Caroline. I like to sing 4024007103939509 in choirs.'
	//First, take a character and check to see if you care about it.
	//Second, if you do add it to the overall value
	// --> Go look at Savannah's reset
	// third check to see if the value passes the luhn algorithm?
	for i := 0; i < len(input); i++ {
		character := input[i]

		if this.isInterestingValue(character) {
			if len(this.value) == 1 {
				this.lastDigitIndex = i
			}
			continue
		}
		fmt.Println(this.value)
		if luhn(this.value){
			this.appendMatch(this.lastDigitIndex, len(this.value))
			//reset the values
		}
	}
}

func luhn(CardNumber []byte) bool {
	odd := len(CardNumber) & 1
	var sum int
	for i, c := range CardNumber {
		if c < '0' || c > '9' {
			return false
		}
		if i&1 == odd {
			sum += t[c-'0']
		} else {
			sum += int(c - '0')
		}
	}
	return sum%10 == 0
}

var t = [...]int{0, 2, 4, 6, 8, 1, 3, 5, 7, 9}

func(this *creditCardRedact) isInterestingValue(value byte) bool {
	interesting := false
	switch {
	case value == '-':
		interesting = true
	case value == ' ':
		interesting = true
	case isNumeric(value):
		this.value = append(this.value, value)
	}
	return interesting
}

func stripBreaks(input []byte) string {
	returnString := string(input[:])
	strings.ReplaceAll(returnString, " ", "")
	strings.ReplaceAll(returnString, "-", "")
	return returnString
}

package redact
//TODO: Create testing... I have no idea if this would work
type creditCardRedact struct {
	*matched
	lastDigitIndex int //when did the value begin? --> Do we need to check for multiple credit cards? I guess that would be in appeneded matches..
	value          []byte
}

func (this *creditCardRedact) match(input []byte) {
	if len(input) <= 0 {
		return
	}
	//input := 'Hi, my name is 24601 Caroline. I like to sing 4024007103939509 in choirs.'
	for i := 0; i < len(input); i++ {
		character := input[i]

		if isNumeric(character) {
			this.value = append(this.value, character)

			if len(this.value) == 1{
				this.lastDigitIndex = i - 1
			}

			if !isNumeric(input[i+1]) && !luhn(string(this.value[:])) {
				this.value = nil
			} else if luhn(string(this.value[:])){
				this.appendMatch(this.lastDigitIndex, len(this.value))
			}
		}
	}
}


func luhn(CardNumber string) bool {
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
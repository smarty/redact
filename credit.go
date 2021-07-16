package redact
func (this *creditCardRedact) match(input []byte) {
	if len(input) <= 0 {
		return
	}
	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if this.length == 0 && !isNumeric(character) {
			this.reset(i)
			continue
		}
		if !isInterestingCharacter(character) && !isNumeric(character) {
			this.reset(i)
			continue
			//else if is very redundant
		} else if isNumeric(character) {
			this.length++
			this.numericLength++

			if this.numericLength > 19 {
				this.reset(i)
				continue
			}

			if this.length == 1 {
				this.lastDigitIndex = i
			}
			//this.value will only hold the bytes in a []byte
			this.value = append(this.value, character)
			//else is here is redundant
		} else if isInterestingCharacter(character) && this.breakType == byte(0) {
			this.length++
			this.breakType = character
		} else if character != this.breakType {
			this.reset(i)
		} else if character == this.breakType {
			if input[i+1] == this.breakType{
				this.reset(i)
				continue
			}
			this.length++
		}
		if i == (len(input)-2) && isNumeric(input[i+1]){
			this.value = append(this.value, input[i+1])
			if luhn(this.value){
				this.appendMatch(this.lastDigitIndex, len(input[this.lastDigitIndex:(this.lastDigitIndex + (this.length + 1))]))
				this.reset(i)
				continue
			}
		}
		if this.numericLength >= 12 && this.numericLength <= 19 && luhn(this.value) && !isNumeric(input[i+1]) {
			this.appendMatch(this.lastDigitIndex, len(input[this.lastDigitIndex:(this.lastDigitIndex + this.length)]))
			this.reset(i)
		}
	}
	this.reset(0)
}
func isInterestingCharacter(character byte) bool {
	interesting := false
	switch character {
	case ' ':
		interesting = true
	case '-':
		interesting = true
	default:
		return interesting
	}
	return interesting
}

func (this *creditCardRedact) reset(i int) {
	this.length = 0
	this.lastDigitIndex = i + 1
	this.value = nil
	this.numericLength = 0
	this.breakType = byte(0)
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

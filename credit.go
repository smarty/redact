package redact

func (this *creditCardRedaction) clear() {
	this.start = 0
	this.length = 0
	this.totalSum = 0
	this.isSecond = false
}

//Can be valid if there are letters at the very beginning or very end but NOT in the middle
//The only valid breaks for a card are: ' ' and '-'
//The breaks must be in the correct positions (i.e. 1010-1010-1010 1010, 1010-1010-0101-0, 0101-0101-0101-0101-010)
func (this *creditCardRedaction) match(input []byte) {
	var numericValue byte
	for i := len(input) - 1; i > 0; i-- {
		if i < len(this.used)-1 && this.used[i] {
			continue
		}
		if isNumeric(input[i]) {
			numericValue = input[i] - '0'
			this.luhnCheck(numericValue)
			this.length++
			continue
		}
		//We want to check for valid breaks ' ' or '-'
		//If it is a valid break and the next character is a number then we continue. Otherwise we reset.
		//We only validate if we find a letter + the length is sufficient
		if input[i] == ' ' || input[i] == '-' {
			if i > 0 && isNumeric(input[i-1]) && this.length < 20 {
				this.length++
				continue
			}
			if i < len(input)-1 && this.validateCreditCard(input[i+1]) {
				this.appendMatch(i+1, this.length)
			}
		}
		this.resetCount(i)
	}
	//Need to validate the card one more time and append just in case we missed one
}

func (this *creditCardRedaction) luhnCheck(value byte) {
	if this.isSecond {
		value *= 2
	}
	this.totalSum += value / 10
	this.totalSum += value % 10
	this.isSecond = !this.isSecond
}

func (this *creditCardRedaction) validateCreditCard(input byte) bool {
	if this.length < 12 || this.length > 20 {
		return false
	}
	if !validateNetwork(input) {
		return false
	}
	return this.totalSum%10 == 0
}

func (this *creditCardRedaction) resetCount(i int) {
	this.start = i - 1
	this.length = 0
	this.totalSum = 0
}

func validateNetwork(input byte) bool {
	return input >= '3' && input <= '6'
}

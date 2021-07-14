package redact

//TODO: Do I need this?
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
	//var numericValue byte
	// i = 18 - stay in as long as i doesnt become 0
	for i := len(input) - 1; i > 0; i-- {

		//some error checking
		if i < len(this.used)-1 && this.used[i] {
			continue
		}

		//if the ith position is numeric, increase the length
		if isNumeric(input[i]) {
			//I don't think this is the right way to be using the luhn algorithm
			//numericValue = input[i] - '0'
			//this.luhnCheck(numericValue)
			this.creditCardNumber = append(this.creditCardNumber, input[i])
			this.length++

			continue
		}

		//Im wondering if we just keep appending numbers and sending them to the luhn to be checked...? But then you could have a number like: 1928362 912646 9128641 183712456lksjdf 92386432487 187236498273 4673 271643 2763

		//We want to check for valid breaks ' ' or '-'
		//If it is a valid break and the next character is a number then we continue. Otherwise we reset.
		//We only validate if we find a letter + the length is sufficient
		//if the break changes then we reset
		//TODO: create a validBreak function - that makes sure that the breaks are in the correct position?
		if input[i] == ' ' || input[i] == '-' {
			//if i is not at the end of the input and the next character is a number and the length is less than twenty
			//increase the length and continue to the next step
			//TODO: add the luhn algorithm into here?
			//so if the luhn algorithm doesn't check out then move to the next character?
			if i > 0 && isNumeric(input[i-1]) && this.length < 20 {
				this.length++
				continue
			}
			//if we are in range and the valid card comes out true
			//valid credit card checks to see if the length is valid (this could be redundant)
			//and if the network is gucci
			//then append the match
			if i < len(input)-1{// && this.validateCreditCard(input[i-1])
				this.appendMatch(i+1, this.length)
				this.resetCount(i)
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

var t = [...]int{0, 2, 4, 6, 8, 1, 3, 5, 7, 9}

func luhn(s []byte) bool {
	odd := len(s) & 1
	var sum int
	for i, c := range s {
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

func validBreakPositions([]byte) bool {
	return false
}
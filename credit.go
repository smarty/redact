package redact

func (this *creditCardRedaction) clear() {
	this.start = 0
	this.length = 0
	this.totalSum = 0
	this.isSecond = false
	this.numBreaks = 0
}
func (this *creditCardRedaction) match(input []byte) {
	inputLength := len(input) - 1
	for i := inputLength; i > 0; i-- {
		if i < len(this.used)-1 && this.used[i] {
			continue
		}
		if isNumeric(input[i]) {
			this.numericLength++
			this.length++
			this.start = i
			this.luhnCheck(input[i])
			continue
		}
		if this.validateCard(input[this.start]) {
			this.appendMatch(this.start, this.length)
			this.resetCount(i, inputLength)
			continue
		}
		if this.length > 0 && this.validateBreaks(input[i]) && !this.validateBreaks(input[i+1]) {
			if this.breakType == 0 {
				this.breakType = input[i]
			}
			if !validateBreakPositions(this.previousNumericLength, this.numericLength) && isNumeric(input[i+1]) {
				this.resetCount(i, inputLength)
				continue
			}
			this.previousNumericLength = this.numericLength
			this.numericLength = 0
			this.numBreaks++
			this.length++
			continue
		}
		this.resetCount(i, inputLength)
	}
	if this.length != 0 && isNumeric(input[0]) {
		this.length++
		this.start = 0
		this.luhnCheck(input[this.start])
	}
	if this.validateCard(input[this.start]) {
		this.appendMatch(this.start, this.length)
		this.resetCount(0, inputLength)
	}
}

func (this *creditCardRedaction) resetCount(i, length int) {
	if this.length < length {
		this.start = i - 1
	}
	this.length = 0
	this.totalSum = 0
	this.isSecond = false
	this.numBreaks = 0
	this.breakType = 0
	this.numericLength = 0
	this.previousNumericLength = 0
}

func (this *creditCardRedaction) luhnCheck(input byte) {
	var value uint64
	value = uint64(input - '0')

	if this.isSecond {
		value *= 2
	}
	this.totalSum += value / 10
	this.totalSum += value % 10
	this.isSecond = !this.isSecond
}
func (this *creditCardRedaction) validateBreaks(input byte) bool {
	return (input == ' ' || input == '-') && (this.breakType == input || this.breakType == 0)
}
func (this *creditCardRedaction) validateCard(input byte) bool {
	if this.numBreaks != 0 && (this.numBreaks < MinCreditBreakLength || this.numBreaks > MaxCreditBreakLength) {
		return false
	}
	numericLength := this.length - this.numBreaks
	if numericLength < MinCreditLength_NoBreaks || numericLength > MaxCreditLength_NoBreaks {
		return false
	}
	if !validateNetwork(input) {
		return false
	}
	return this.totalSum%10 == 0
}
func validateNetwork(input byte) bool {
	return input >= '3' && input <= '6'
}

func validateBreakPositions(prev, curr int) bool {
	switch prev {
	case 0:
		return true
	case 4:
		if curr == 4 {
			return true
		}
	case 3:
		if curr == 4 {
			return true
		}
	case 5:
		if curr == 6 || curr == 4 {
			return true
		}
	case 6:
		if curr == 4 {
			return true
		}
	default:
		return false
	}
	return false
}

// Valid credit card lengths are:
/*
#### #### ##### (4-4-5)
#### ###### ##### (4-6-5)
#### #### #### #### (4-4-4-4)
#### #### #### #### ### (4-4-4-4-3)
*/

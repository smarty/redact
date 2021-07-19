package redact

func (this *creditCardRedaction) clear() {
	this.start = 0
	this.length = 0
	this.totalSum = 0
	this.isSecond = false
}

func (this *creditCardRedaction) match(input []byte) {
	for i := len(input) - 1; i > 0; i-- {
		if i < len(this.used)-1 && this.used[i] {
			continue
		}
		if isNumeric(input[i]) {
			this.length++
			this.start = i
			this.luhnCheck(input[i])
			continue
		}
		if this.length > 0 && !this.validateBreaks(input[i]){
			this.resetCount(i)
			continue
		}
		if i < len(input)-1 && this.validateCard(input[i+1]) {
			this.appendMatch(this.start, this.length)
			this.resetCount(i)
			continue
		}
	}
	if this.length != 0 && isNumeric(input[0]) {
		this.length++
		this.start--
		this.luhnCheck(input[this.start])
	}
	if this.start >= 0 && this.validateCard(input[this.start]) {
		this.appendMatch(this.start, this.length)
	}
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

func (this *creditCardRedaction) validateCard(input byte) bool {
	if this.breakLength != 0 && (this.breakLength < 2 || this.breakLength > 4) {
		return false
	}
	numericLength := this.length - this.breakLength
	if numericLength > 19 || numericLength < 12 {
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

func (this *creditCardRedaction) resetCount(i int) {
	this.start = i - 1
	this.length = 0
	this.totalSum = 0
	this.isSecond = false
	this.breakLength = 0
}

//TODO: Should be boolean?
func (this *creditCardRedaction) validateBreaks(input byte) bool {
	switch {
	case input == ' ':
		this.breakLength++
		this.length++
		return true
	case input == '-':
		this.breakLength++
		this.length++
		return true
	default:
		return false
	}
}

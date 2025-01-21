package redact

func (this *phoneRedaction) clear() {
	this.length = 0
	this.start = 0
	this.breakLength = 0
	this.numericLength = 0
}
func (this *phoneRedaction) match(input []byte) {
	var previousBreak byte
	for i, val := range input {
		if i < len(this.used)-1 && this.used[i] {
			this.resetCount(i)
			continue
		}
		if isNumeric(val) {
			this.numericLength++
			this.length++
			switch {
			case this.length < MaxPhoneLength_WithBreaks && this.numericLength != MinPhoneLength_WithNoBreaks:
				continue
			case this.length >= MinPhoneLength_WithNoBreaks && this.length <= MaxPhoneLength_WithBreaks && this.breakLength <= MaxPhoneBreakLength:
				this.validateMatch(input[this.start:i])
			}
			if i < len(input)-1 {
				this.resetCount(i)
			}
			continue
		}
		if i < len(input)-1 {
			this.validateBreaks(val, previousBreak, i)
			previousBreak = input[i]
		}
	}
}

func (this *phoneRedaction) resetCount(i int) {
	this.start = i + 1
	this.length = 0
	this.breakLength = 0
	this.numericLength = 0
}
func (this *phoneRedaction) validateBreaks(input, previousBreak byte, i int) {
	if !validBreak(input) {
		this.resetCount(i)
		return
	}
	if input == ' ' && previousBreak != ')' {
		this.resetCount(i)
		return
	}
	if input == '+' {
		if this.start != i {
			this.resetCount(i)
		}
		this.start = i + 1
		this.length = 1
		return
	}

	this.length++
	this.breakLength++
}
func validBreak(input byte) bool {
	return input == '-' || input == '(' || input == ')' || input == ' ' || input == '+'
}

func (this *phoneRedaction) validateMatch(testMatch []byte) {
	switch {
	case this.length == MinPhoneLength_WithBreaks && this.breakLength == MinPhoneBreakLength:
		if testMatch[3] == '-' && testMatch[7] == '-' {
			this.appendMatch(this.start, this.length)
		}

	case this.length == MaxPhoneLength_WithBreaks-1 && this.breakLength == MaxPhoneBreakLength-1: // No space
		if testMatch[1] == '(' && testMatch[5] == ')' && testMatch[9] == '-' {
			this.appendMatch(this.start, this.length)
		}
	case this.length == MaxPhoneLength_WithBreaks && this.breakLength == MaxPhoneBreakLength: // With space
		if testMatch[1] == '(' && testMatch[5] == ')' && testMatch[6] == ' ' && testMatch[10] == '-' {
			this.appendMatch(this.start, this.length)
		}
	}
}

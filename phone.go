package redact

func (this *phoneRedaction) clear() {
	this.length = 0
	this.start = 0
	this.breakLength = 0
	this.numericLength = 0
}
func (this *phoneRedaction) match(input []byte) {
	for i := 0; i < len(input)-1; i++ {
		if i < len(this.used)-1 && this.used[i] {
			continue
		}
		if isNumeric(input[i]) {
			this.numericLength++
			this.length++
			switch {
			case this.length < MaxPhoneLength_WithBreaks && this.numericLength != MinPhoneLength_WithNoBreaks:
				continue
			case this.length >= MinPhoneLength_WithNoBreaks && this.length <= MaxPhoneLength_WithBreaks && this.breakLength <= MaxPhoneBreakLength:
				this.validateMatch(input[this.start : i])
			}
			if i < len(input)-1 {
				this.resetCount(i)
			}
			continue
		}
		if i < len(input)-1 {
			this.validateBreaks(input[i], i)
		}
	}
}

func (this *phoneRedaction) resetCount(i int) {
	this.start = i + 1
	this.length = 0
	this.breakLength = 0
	this.numericLength = 0
}
func (this *phoneRedaction) validateBreaks(input byte, i int) {
	switch input {
	case '-':
		this.length++
		this.breakLength++
	case '(':
		this.length++
		this.breakLength++
	case ')':
		this.length++
		this.breakLength++
	case '+':
		this.start = i + 1
		this.length = 1
	default:
		this.resetCount(i)
	}
}
func (this *phoneRedaction) validateMatch(testMatch []byte) {
	switch {
	case this.length == MinPhoneLength_WithBreaks && this.breakLength == MinPhoneBreakLength:
		if testMatch[3] == '-' && testMatch[7] == '-' {
			this.appendMatch(this.start, this.length)
		}
	case this.length == MaxPhoneLength_WithBreaks && this.breakLength == MaxPhoneBreakLength:
		if testMatch[1] == '(' && testMatch[5] == ')' && testMatch[9] == '-' {
			this.appendMatch(this.start, this.length)
		}
	}
}

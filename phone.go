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
			case this.length < MaxLength_WithBreaks && this.numericLength != MinLength_WithNoBreaks:
				continue
			case this.length >= MinLength_WithNoBreaks && this.length <= MaxLength_WithBreaks && this.breakLength <= MaxBreakLength:
				this.validateMatch(input[this.start : this.start+this.length])
			}
			if i < len(input)-1 {
				this.resetCount(i)
			}
			continue
		}

		this.validateBreaks(input, i)
	}
}

func (this *phoneRedaction) validateBreaks(input []byte, i int) {
	switch input[i] {
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
		if i < len(input)-1 {
			this.start = i + 1
		}
		this.length++
	default:
		if i < len(input)-1 {
			this.resetCount(i)
		}
	}
}
func (this *phoneRedaction) validateMatch(testMatch []byte) {
	switch {
	case this.length == MinLength_WithBreaks && this.breakLength == MinBreakLength:
		if testMatch[3] == '-' && testMatch[7] == '-' {
			this.appendMatch(this.start, this.length)
		}
	case this.length == MaxLength_WithBreaks && this.breakLength == MaxBreakLength:
		if testMatch[1] == '(' && testMatch[5] == ')' && testMatch[9] == '-' {
			this.appendMatch(this.start, this.length)
		}
	}
}
func (this *phoneRedaction) resetCount(i int) {
	this.start = i + 1
	this.length = 0
	this.breakLength = 0
	this.numericLength = 0
}

const (
	MaxLength_WithBreaks   = 14
	MinLength_WithBreaks   = 12
	MinLength_WithNoBreaks = 10
	MaxBreakLength         = 3
	MinBreakLength         = 2
)

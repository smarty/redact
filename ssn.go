package redact

func (this *ssnRedaction) clear() {
	this.start = 0
	this.length = 0
	this.breakLength = 0
}
func (this *ssnRedaction) match(input []byte) {
	for i := 0; i < len(input)-1; i++ {
		if i < len(this.used)-1 && this.used[i] {
			continue
		}
		if isNumeric(input[i]) {
			this.length++
			switch {
			case this.length < MaxSSNLength_WithBreaks && this.breakLength <= MaxSSNBreakLength:
				continue
			case this.length == MaxSSNLength_WithBreaks && this.breakLength == MaxSSNBreakLength:
				this.validateMatch(input[this.start : this.start+this.length])
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
	this.resetCount(0)
}

func (this *ssnRedaction) resetCount(i int) {
	this.start = i + 1
	this.length = 0
	this.breakLength = 0
}

func (this *ssnRedaction) validateMatch(testMatch []byte) {
	switch {
	case testMatch[MinSSNBreakPosition] == '-' && testMatch[MaxSSNBreakPosition] == '-':
		this.appendMatch(this.start, this.length)
	case testMatch[MinSSNBreakPosition] == ' ' && testMatch[MaxSSNBreakPosition] == ' ':
		this.appendMatch(this.start, this.length)
	}
}
func (this *ssnRedaction) validateBreaks(input byte, i int) {
	switch {
	case input == ' ' && this.length >= 3:
		this.length++
		this.breakLength++
	case input == '-' && this.length >= 3:
		this.length++
		this.breakLength++
	default:
		this.resetCount(i)
	}
}

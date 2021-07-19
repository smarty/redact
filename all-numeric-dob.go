package redact

type AllNumericDOB struct {
	redact            *dobRedaction
	validNumericMonth bool
}

func (this *AllNumericDOB) findMatch(input []byte) {
	for i := 0; i < len(input); i++ {
		if i < len(this.redact.used)-1 && this.redact.used[i] {
			continue
		}
		if isNumeric(input[i]) {
			this.redact.numericLength++
			this.redact.length++
			continue
		}
		if this.redact.length > MinDOBLength_WithBreaks && this.redact.length < MaxDOBLength_WithBreaks {
			if this.redact.breakLength == MaxDOBBreakLength && this.validNumericMonth {
				this.redact.appendMatch(this.redact.start, this.redact.length)
				this.resetCount(i)
				continue
			}
		}
		this.validateBreaks(input[i], i)
		if !this.validateDOB(input[i-this.redact.numericLength : i]) {
			this.resetCount(i)
		}
	}
	if this.redact.length >= MinDOBLength_WithBreaks && this.redact.length <= 10 && this.redact.breakLength == MaxDOBBreakLength && this.validNumericMonth {
		if this.validateDOB(input[this.redact.length-this.redact.numericLength : this.redact.length]) {
			this.redact.appendMatch(this.redact.start, this.redact.length)
			this.resetCount(0)
		}
	}
}

func (this *AllNumericDOB) validateDOB(input []byte) bool {
	this.redact.numericLength = 0
	switch len(input) {
	case 0:
		return true
	case 4:
		return this.redact.validateYear(input)
	case 2:
		return this.validateDate(input)
	case 1:
		if input[0] != '0' {
			return true
		}
	}
	return false
}
func (this *AllNumericDOB) validateDate(input []byte) bool {
	switch {
	case input[0] == '0' && input[1] <= '9' && input[1] > 0:
		this.validNumericMonth = true
		return true
	case input[0] == '1' && input[1] <= '2':
		this.validNumericMonth = true
		return true
	case input[0] == '3' && input[1] > 1:
		return false
	default:
		return false
	}
}

func (this *AllNumericDOB) resetCount(i int) {
	this.redact.start = i + 1
	this.redact.length = 0
	this.redact.breakLength = 0
	this.redact.numericLength = 0
	this.validNumericMonth = false
}

func (this *AllNumericDOB) validateBreaks(input byte, i int) {
	switch {
	case input == '/':
		this.redact.length++
		this.redact.breakLength++
	case input == '-':
		this.redact.length++
		this.redact.breakLength++
	case this.redact.breakLength > 2:
		this.resetCount(i)
	default:
		this.resetCount(i)
	}
}

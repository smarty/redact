package redact

type AllNumericDOB struct {
	dobRedaction *dobRedaction
}

func (this *AllNumericDOB) findMatch(input []byte) {
	for i := 0; i < len(input); i++ {
		if i < len(this.dobRedaction.used)-1 && this.dobRedaction.used[i] {
			continue
		}
		if isNumeric(input[i]) {
			this.dobRedaction.numericLength++
			this.dobRedaction.length++
			continue
		}
		if this.dobRedaction.length > 8 && this.dobRedaction.length < 11 {
			if this.dobRedaction.breakLength == 2 && this.dobRedaction.validNumericMonth {
				this.dobRedaction.appendMatch(this.dobRedaction.start, this.dobRedaction.length)
				this.dobRedaction.resetCount(i)
				continue
			}
		}
		this.dobRedaction.validateBreaks(input[i], i)
		if !this.validateDOB(input[i-this.dobRedaction.numericLength : i]) {
			this.dobRedaction.resetCount(i)
		}
	}
	if this.dobRedaction.length >= 8 && this.dobRedaction.length <= 10 && this.dobRedaction.breakLength == 2 && this.dobRedaction.validNumericMonth {
		if this.validateDOB(input[this.dobRedaction.length-this.dobRedaction.numericLength : this.dobRedaction.length]) {
			this.dobRedaction.appendMatch(this.dobRedaction.start, this.dobRedaction.length)
		}
	}
}

func (this *AllNumericDOB) validateDOB(input []byte) bool {
	this.dobRedaction.numericLength = 0
	switch len(input) {
	case 0:
		return true
	case 4:
		return this.validateYear(input)
	case 2:
		return this.validateDate(input)
	case 1:
		if input[0] != '0' {
			return true
		}
	default:
		return false
	}
	return false
}
func (this *AllNumericDOB) validateDate(input []byte) bool {
	switch {
	case input[0] == '0' && input[1] <= '9' && input[1] > 0:
		this.dobRedaction.validNumericMonth = true
		return true
	case input[0] == '1' && input[1] <= '2':
		this.dobRedaction.validNumericMonth = true
		return true
	case input[0] == '3' && input[1] > 1:
		return false
	default:
		return false
	}
}
func (this *AllNumericDOB) validateYear(input []byte) bool {
	switch {
	case input[0] == '1' && input[1] == '9':
		return true
	case input[0] == '2' && input[1] == '0' && input[2] < '3':
		if input[2] == '2' && input[3] > '1' {
			return false
		}
		return true
	default:
		return false
	}
	return false
}

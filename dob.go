package redact

func (this *dobRedaction) clear() {
	this.start = 0
	this.length = 0
	this.numericLength = 0
	this.breakLength = 0
}
func (this *dobRedaction) match(input []byte) {
	numericDOB := allNumericDOB{redact: this}
	fullDOB := fullDOB{redact: this}

	numericDOB.findMatch(input)
	this.clear()
	fullDOB.findMatch(input)
}

func validateYear(input []byte) bool {
	switch {
	case input[0] == '1' && input[1] == '9':
		return true
	case input[0] == '2' && input[1] == '0' && input[2] < '3':
		if input[2] == '2' && input[3] > '1' {
			return false
		}
		return true
	}
	return false
}

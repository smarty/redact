package redact

func (this *dobRedaction) clear() {
	this.start = 0
	this.length = 0
	this.numericLength = 0
	this.breakLength = 0
	this.validNumericMonth = false
}

func (this *dobRedaction) match(input []byte) {
	numericDOB := AllNumericDOB{dobRedaction: this}
	numericDOB.findMatch(input)
	//fullDOB := FullDOB{dobRedaction: this}
	//fullDOB.findMatch(input)
}

func (this *dobRedaction) resetCount(i int) {
	this.start = i + 1
	this.length = 0
	this.breakLength = 0
	this.numericLength = 0
	this.validNumericMonth = false
}
func (this *dobRedaction) validateBreaks(input byte, i int) {
	switch {
	case input == '/':
		this.length++
		this.breakLength++
	case input == '-':
		this.length++
		this.breakLength++
	case this.breakLength > 2:
		this.resetCount(i)
	default:
		this.resetCount(i)
	}
}

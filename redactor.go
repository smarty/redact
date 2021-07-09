package redact

func (this *Redactor) RedactAll(input []byte) []byte {
	this.clear(this.phone, this.email, this.dob, this.credit, this.ssn)
	if len(input) <= 0 {
		return input
	}
	this.match(input, this.phone, this.email,this.ssn, this.dob, this.credit)
	result := this.redactMatches(input)
	return result
}

func (this *Redactor) redactMatches(input []byte) []byte {
	count := len(this.matches)
	if count == 0 {
		return input
	}
	this.monitor.Redacted(count)

	buffer := input
	bufferLength := len(buffer)
	var lowIndex, highIndex int

	for _, match := range this.matches {
		lowIndex = match.InputIndex
		highIndex = lowIndex + match.Length
		if lowIndex < 0 {
			continue
		}
		if highIndex > bufferLength {
			continue
		}
		for ; lowIndex < highIndex; lowIndex++ {
			buffer[lowIndex] = '*'
		}
	}

	output := buffer
	return output
}

func (this *Redactor) match(input []byte, matchMethod ...Redaction) {
	for _, method := range matchMethod {
		method.match(input)
	}
}
func (this *Redactor) clear(matchMethod ...Redaction) {
	for _, method := range matchMethod {
		method.clear()
	}
	this.matched.clear()
}

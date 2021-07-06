package redact

func (this *ssnRedaction) clear() {
	this.resetMatchValues()
	this.start = 0
}

func (this *ssnRedaction) match(input []byte) {
	if len(input) <= 0 {
		return
	}
	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if i > len(this.used)-1 {
			return
		}
		if this.used[i] {
			continue
		}
		if !isNumeric(character) {
			if isSSN(this.numbers) && this.breaks && this.numBreaks == 2 {
				this.appendMatch(this.start, this.length)
				this.resetMatchValues()
				this.start = i + 1
				continue
			}
			if ssnBreakNotFound(character) {
				this.start = i + 1
				this.resetMatchValues()
				continue
			}
			if i < len(input)-1 && isNumeric(input[i+1]) && this.isCandidate {
				this.length++
				if this.breakType == character && this.numbers == 5 {
					this.breaks = true
				}
				if this.numbers == 3 {
					this.breakType = character
				}
			}
			this.numBreaks++
			continue
		}
		if this.isCandidate {
			this.incrementLength()
		} else {
			this.isCandidate = true
			this.breaks = false
			this.start = i
			this.incrementLength()
		}
	}
	if isNumeric(input[len(input)-1]) {
		this.incrementLength()
	}
	if isSSN(this.numbers) && this.breaks {
		this.appendMatch(this.start, this.length)
	}
}

func (this *ssnRedaction) incrementLength() {
	this.numbers++
	this.length++
}

func (this *ssnRedaction) resetMatchValues() {
	this.numbers = 0
	this.breaks = false
	this.breakType = 'x'
	this.numBreaks = 0
	this.length = 0
	this.isCandidate = false
}
func ssnBreakNotFound(character byte) bool {
	return character != '-' && character != ' '
}
func isSSN(length int) bool {
	return length == 9
}

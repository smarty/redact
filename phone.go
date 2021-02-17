package redact

type phoneRedaction struct {
	*matched
	start       int
	length      int
	numbers     int
	breaks      int
	parenBreak  bool
	matchBreaks bool
	breakType   byte
	isCandidate bool
}

func (this *phoneRedaction) clear() {
	this.resetMatchValues()
	this.start = 0
	this.breakType = 0
	this.isCandidate = false
}

func (this *phoneRedaction) match(input string) {
	for i := 0; i < len(input)-1; i++ {
		if this.used[i] {
			continue
		}
		character := input[i]
		if !isNumeric(character) {
			if character == '+' {
				this.start = i + 2
				this.length--
				this.numbers--
				continue
			}
			if isPhoneNumber(this.numbers) {
				if correctBreaks(this.breaks, this.parenBreak, this.matchBreaks) {
					this.appendMatch(this.start, this.length)
				}
				this.resetMatchValues()
				this.start = i + 1
				continue
			}
			if phoneBreakNotFound(character) {
				this.start = i + 1
				this.length = 0
				this.numbers = 0
				this.isCandidate = false
				continue
			}
			if character == '(' {
				this.start = i
				this.isCandidate = true
				this.length++
				this.breaks++
				this.parenBreak = true
				continue
			}
			if character == ')' {
				this.length++
				this.breaks++
				this.parenBreak = true
				continue
			}
			if i < len(input)-1 && !isNumeric(input[i+1]) {
				this.length = 0
				this.numbers = 0
				this.start = i + 1
				this.breaks = 0
				this.parenBreak = false
				continue
			}
			if this.isCandidate {
				this.length++
				if this.breakType == character && this.numbers == 6 {
					this.matchBreaks = true
				} else {
					this.matchBreaks = false
				}
				if this.numbers == 3 {
					this.breakType = character
				}
			}
			this.breaks++
			continue
		}
		if this.isCandidate {
			this.incrementLength()

		} else {
			this.isCandidate = true
			this.start = i
			this.incrementLength()
			this.breaks = 0
			this.matchBreaks = false
			this.parenBreak = false
		}
	}
	if isNumeric(input[len(input)-1]) {
		this.incrementLength()
	}
	if isPhoneNumber(this.numbers) {
		if correctBreaks(this.breaks, this.parenBreak, this.matchBreaks) {
			this.appendMatch(this.start, this.length)
		}
		this.resetMatchValues()
	}
}

func (this *phoneRedaction) incrementLength() {
	this.length++
	this.numbers++
}

func (this *phoneRedaction) resetMatchValues() {
	this.length = 0
	this.numbers = 0
	this.breaks = 0
	this.parenBreak = false
	this.matchBreaks = false
}
func phoneBreakNotFound(character byte) bool {
	return character != '-' && character != '(' && character != ')'
}
func isPhoneNumber(length int) bool {
	return length == 10
}

func correctBreaks(breaks int, parenBreak, matchBreak bool) bool {
	if breaks == 3 && parenBreak {
		return true
	}
	if breaks == 4 && parenBreak {
		return true
	}
	if breaks == 2 && matchBreak {
		return true
	}
	return false
}

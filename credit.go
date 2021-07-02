package redact

type creditCardRedaction struct {
	*matched
	lastDigitIndex int
	length         int
	totalNumbers   int
	isOdd          bool
	isCandidate    bool
	isLetter       bool
	totalSum       int
	numBreaks      int
	groupLength    int
	numGroups      int
	breakType      byte
}

func (this *creditCardRedaction) clear() {
	this.resetMatchValues()
	this.lastDigitIndex = 0
	this.groupLength = 0
}

func (this *creditCardRedaction) match(input []byte) {
	if len(input) <= 0 {
		return
	}
	for i := len(input) - 1; i > 0; i-- {
		character := input[i]
		if !isNumeric(character) {
			if this.validCardCheck() && isValidNetwork(input[i+1]) {
				if this.numBreaks == 0 || this.numBreaks > 1 {
					this.appendMatch(this.lastDigitIndex-this.length+1, this.length)
					this.isLetter = false
				}
				this.lastDigitIndex = i - 1
				temp := this.isCandidate
				this.resetMatchValues()
				this.isCandidate = temp
				continue
			}
			if this.groupLength > 6 || this.groupLength < 4 {
				this.groupLength = 0
			} else {
				this.numGroups++
				this.groupLength = 0
			}
			if creditCardBreakNotFound(character) && i != len(input)-1 && !isNumeric(input[i-1]) {
				this.lastDigitIndex = i - 1
				this.length = 0
				this.totalSum = 0
				this.totalNumbers = 0
				this.isCandidate = false
				this.breakType = 'x'
				this.numBreaks = 0
				this.numGroups = 0
				continue
			}
			if this.isCandidate || isLetter(character) {
				if this.breakType == character && !creditCardBreakNotFound(character) {
					if i != len(input)-1 && isNumeric(input[i+1]) {
						this.numBreaks++
					} else {
						this.resetMatchValues()
						this.lastDigitIndex = i - 1
						this.length = 0
						this.totalSum = 0
						this.totalNumbers = 0
						this.isCandidate = false
						this.breakType = 'x'
						this.numBreaks = 0
						this.numGroups = 0
					}
				}
				if this.breakType == 'x' && !creditCardBreakNotFound(character) {
					this.breakType = character
					this.numBreaks++
				}
				if character != this.breakType {
					if i < len(input)-1 && isNumeric(input[i+1]) {
						temp := this.numBreaks
						this.resetMatchValues()
						this.numBreaks = temp
						this.numBreaks++
						continue
					}
					this.lastDigitIndex = i - 1
					this.length = 0
					this.totalSum = 0
					this.totalNumbers = 0
					this.isCandidate = false
					this.breakType = 'x'
					this.numBreaks = 0
					this.numGroups = 0
					continue
				}
			}
			if i < len(input)-1 && !creditCardBreakNotFound(input[i+1]) {
				continue
			}
			this.length++
		} else {
			this.isOdd = !this.isOdd
			this.totalNumbers++
			this.groupLength++
			number := int(character - '0')
			if !this.isOdd {
				number += number
				if number > 9 {
					number -= 9
				}
			}
			this.totalSum += number

			if this.isCandidate {
				this.length++
			} else {
				this.isCandidate = true
				this.breakType = 'x'
				this.lastDigitIndex = i
				this.totalNumbers = 1
				if this.length == 0 {
					this.length++
				}
			}
		}
	}
	if len(input) > 0 && isNumeric(input[0]) {
		this.isOdd = !this.isOdd
		this.totalNumbers++
		number := int(input[0] - '0')
		if !this.isOdd {
			number += number
			if number > 9 {
				number -= 9
			}
		}
		if this.numBreaks > 0 {
			this.numGroups++
		}
		this.totalSum += number
		this.length++
	}
	if this.validCardCheck() {
		if this.numBreaks > 0 && this.numGroups == 0 {
			this.resetMatchValues()
			this.isLetter = false
		}

		if this.length > 0 && (this.numBreaks == 0 || this.numBreaks > 1) {
			start := this.lastDigitIndex + 1 - this.length
			index := start
			if start < 0 {
				index = start * -1
				start = 0
			}
			if isValidNetwork(input[index]) {
				this.appendMatch(start, this.length)
				this.isLetter = false
			}
			this.resetMatchValues()
			this.isLetter = false
		}
	}
}

func (this *creditCardRedaction) validCardCheck() bool {
	if this.totalNumbers <= 12 {
		return false
	}
	if this.totalNumbers >= 20 {
		return false
	}
	if this.totalSum%10 != 0 {
		return false
	}
	if this.numGroups < 2 && this.numGroups != 0 {
		return false
	}
	return true
}

func (this *creditCardRedaction) resetMatchValues() {
	this.breakType = 'x'
	this.length = 0
	this.totalSum = 0
	this.isOdd = false
	this.totalNumbers = 0
	this.numBreaks = 0
	this.numGroups = 0
	this.isCandidate = false
}

func creditCardBreakNotFound(character byte) bool {
	return character != '-' && character != ' '
}

func isValidNetwork(character byte) bool {
	return character >= '0' && character <= '9'
}

func isLetter(character byte) bool {
	if character >= 65 || character <= 90 {
		return true
	}
	if character >= 97 || character <= 122 {
		return true
	}
	return false
}

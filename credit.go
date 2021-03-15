package redact

type creditCardRedaction struct {
	*matched
	lastDigitIndex int
	length         int
	totalNumbers   int
	isOdd          bool
	isCandidate    bool
	totalSum       int
	numBreaks      int
	lengthGroup    int
	numGroups      int
	breakType      byte
}

func (this *creditCardRedaction) clear() {
	this.resetMatchValues()
	this.lastDigitIndex = 0
	this.lengthGroup = 0
}

func (this *creditCardRedaction) match(input []byte) {
	if len(input) <= 0 {
		return
	}
	for i := len(input) - 1; i > 0; i-- {
		character := input[i]
		if !isNumeric(input[i]) {
			if this.validCardCheck(input) && isValidNetwork(input[i+1]) {
				if this.numBreaks == 0 || this.numBreaks > 1 {
					this.appendMatch(this.lastDigitIndex-this.length+1, this.length)
				}
				this.lastDigitIndex = i - 1
				temp := this.isCandidate
				this.resetMatchValues()
				this.isCandidate = temp
				continue
			}
			if this.lengthGroup > 6 || this.lengthGroup < 4 {
				this.lengthGroup = 0
			} else {
				this.numGroups++
				this.lengthGroup = 0
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
			if this.isCandidate {
				if this.breakType == character && !creditCardBreakNotFound(character) {
					this.numBreaks++
				}
				if this.breakType == 'x' && !creditCardBreakNotFound(character) {
					this.breakType = character
					this.numBreaks++
				}
				if this.breakType != character{
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
			this.lengthGroup++
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
	if this.validCardCheck(input){
		if this.numBreaks == 0 || this.numBreaks > 1 {
			start := (this.lastDigitIndex + 1) - this.length
			if isValidNetwork(input[start]) {
				this.appendMatch(start, this.length)
			}
			this.resetMatchValues()
		}
	}
}

func (this *creditCardRedaction) validCardCheck(input []byte) bool {
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
	////return this.totalNumbers > 12 && this.totalNumbers < 20 && this.totalSum%10 == 0 && isValidNetwork(input[0]) && (this.numGroups < 7 && this.numGroups > 2 || this.numGroups == 0) && this.breaks
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
	return character >= '3' && character <= '6'
}

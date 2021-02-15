package redact

type creditCardRedaction struct {
	lastDigit    int
	length       int
	totalNumbers int
	isOdd        bool
	isCandidate  bool
	totalSum     int
	breaks       bool
	numBreaks    int
	lengthGroup  int
	numGroups    int
	breakType    byte
	used         []bool
	matches      []match
}

func (this *creditCardRedaction) appendMatch(start int, length int) {
	for i := start; i <= start+length; i++ {
		this.used[i] = true
	}

	this.matches = append(this.matches, match{InputIndex: start, Length: length})
}

func (this *creditCardRedaction) clear() {
	this.resetMatchValues()
	this.lastDigit = 0
	this.lengthGroup = 0
}

func (this *creditCardRedaction) match(input string) {
	for i := len(input) - 1; i > 0; i-- {
		character := input[i]
		if !isNumeric(input[i]) {
			if this.validCardCheck(input) && isValidNetwork(input[i+1]) {
				if this.validNumBreaks() {
					this.appendMatch(this.lastDigit-this.length+1, this.length)
				}
				this.lastDigit = i - 1
				this.resetMatchValues()
				continue
			}
			if this.lengthGroup > 6 || this.lengthGroup < 4 {
				this.lengthGroup = 0
			} else {
				this.numGroups++
				this.lengthGroup = 0
			}
			if creditCardBreakNotFound(character) && i != len(input)-1 && !isNumeric(input[i-1]) {
				this.lastDigit = i - 1
				this.resetMatchValues()
				continue
			}
			if this.isCandidate {
				if this.breakType == character && !creditCardBreakNotFound(character) {
					this.breaks = true
					this.numBreaks++
				}
				if this.breakType == 'x' && !creditCardBreakNotFound(character) {
					this.breakType = character
					this.numBreaks++
				}
				if this.breakType != character {
					if i < len(input)-1 && isNumeric(input[i+1]) {
						temp := this.numBreaks
						this.resetMatchValues()
						this.numBreaks = temp
						this.numBreaks++
						continue
					}
					this.resetMatchValues()
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
				this.breaks = false
				this.lastDigit = i
				this.totalNumbers = 1
				if this.length == 0 {
					this.length++
				}
			}
		}
	}
	if isNumeric(input[0]) {
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
	if this.numBreaks == 0 {
		this.breaks = true
	}
	if this.validCardCheck(input) && this.breaks {
		if this.validNumBreaks(){
			this.appendMatch(this.lastDigit-this.length+1, this.length)
		}
		this.breaks = false
	}

}

func (this *creditCardRedaction) validNumBreaks() bool {
	return this.numBreaks == 0 || this.numBreaks > 1 && this.numBreaks < 5
}

func (this *creditCardRedaction) validCardCheck(input string) bool {
	return this.totalNumbers > 12 && this.totalNumbers < 20 && this.totalSum%10 == 0 && isValidNetwork(input[0]) && (this.numGroups < 7 && this.numGroups > 2 || this.numGroups == 0) && this.breaks
}

func (this *creditCardRedaction) resetMatchValues() {
	this.breaks = false
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

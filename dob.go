package redact

type dobRedaction struct {
	*matched
	start            int
	length           int
	isCandidate      bool
	monthStart       int
	monthLength      int
	monthCandidate   bool
	startChar        byte
	totalGroupLength int
	breaks           bool
	numBreaks        int
	breakType        byte
	groupLength      int
	firstDigit       byte
	secondDigit      byte
	thirdDigit       byte
	fourthDigit      byte
	validMonth       bool
	validYear        bool
}

func(this* dobRedaction) clear(){
	this.resetMatchValues()
	this.resetYearValues()
	this.startChar = 'x'
	this.numBreaks = 0
	this.breakType = 0
}

func (this *dobRedaction) match(input []byte) {
	this.resetYearValues()
	this.startChar = 'x'

	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !isNumeric(character) {
			if validMonthFirstLetter(character) && this.startChar == 'x' {
				this.startChar = character
				this.monthStart = i
			}
			if this.isDOB() {
				if this.groupLength == 2 && validDayDigit(this.firstDigit, this.secondDigit) && this.validMonth || this.groupLength == 4 {
					if i != len(input)-1 && input[i+1] != this.breakType {
						this.appendMatch(this.start, this.length)
					}
				}
				this.start = i + 1
				this.resetMatchValues()
				continue
			}
			if this.numBreaks == 2 {
				this.breaks = false
			}
			if character == ' ' {
				if this.monthLength > 2 && isMonth(this.startChar, input[i-1], this.monthLength) {
					this.monthCandidate = true
					this.monthLength++
					continue
				} else {
					this.resetMatchValues()
				}
			}
			if dobBreakNotFound(character) || (i < len(input)-1 && doubleBreak(character, input[i+1])) {
				if character == ',' && this.monthCandidate && this.groupLength <= 2 && this.groupLength != 0 {
					this.appendMatch(this.monthStart, this.monthLength+this.length+1)
					this.resetMatchValues()
					continue
				}
				if this.startChar != 'x' && character != ' ' {
					this.monthLength++
				}
				this.start = i + 1
				this.length = 0
				this.totalGroupLength = 0
				this.breaks = false
				this.isCandidate = false
				this.groupLength = 0
				this.numBreaks = 0
				this.validMonth = false
				this.monthCandidate = false
				continue
			}
			this.monthStart = i + 1
			this.monthLength = 0

			if this.isCandidate {
				this.length++
			}
			this.validDateCheck(character)
			this.groupLength = 0
			continue
		}
		this.totalGroupLength++
		this.groupLength++

		this.dateCalculator(character, i)

		if this.length == 2 && this.monthCandidate && this.groupLength <= 2 {
			if i < len(input)-1 {
				if input[i+1] == ',' {
					this.appendMatch(this.monthStart, this.monthLength+this.length+1)
				}
			}
			this.resetMatchValues()
		}
	}

	if isNumeric(input[len(input)-1]) {
		this.length++
		this.totalGroupLength++
		this.groupLength++
		if this.groupLength == 4 {
			this.fourthDigit = input[len(input)-1]
			this.validYear = validYearDigit(this.firstDigit, this.secondDigit, this.thirdDigit, this.fourthDigit)
		}
	}
	if this.isDOB() {
		this.appendMatch(this.start, this.length)
		this.resetMatchValues()
	}
}

func (this *dobRedaction) dateCalculator(character uint8, i int) {
	if this.firstDigit == 100 && this.groupLength < 3 {
		this.firstDigit = character
	} else {
		if this.groupLength == 2 {
			this.secondDigit = character
		}
	}
	if this.groupLength == 3 {
		this.thirdDigit = character
	}
	if this.groupLength == 4 {
		this.fourthDigit = character
		this.validYear = validYearDigit(this.firstDigit, this.secondDigit, this.thirdDigit, this.fourthDigit)
		this.resetYearValues()
	}
	if this.isCandidate || this.monthCandidate {
		this.length++
	} else {
		this.isCandidate = true
		this.breakType = 'x'
		this.start = i
		this.breaks = false
		this.length++
	}
}

func dobBreakNotFound(character byte) bool {
	return character != '/' && character != '-'
}

func (this *dobRedaction) isDOB() bool {
	return this.totalGroupLength >= 6 && this.totalGroupLength <= 8 && this.breaks && this.numBreaks == 2 && this.validYear
}

func (this *dobRedaction) validDateCheck(character uint8) {
	if this.firstDigit == '1' && this.secondDigit <= '2' && this.groupLength != 4 {
		this.validMonth = true
	}
	if validDayDigit(this.firstDigit, this.secondDigit) || (this.totalGroupLength == 4 && this.validYear) {
		if character == this.breakType && validGroupLength(this.groupLength) {
			this.breaks = true
			this.numBreaks++
		}
		if validGroupLength(this.groupLength) && this.totalGroupLength < 3 || this.validYear && !this.breaks {
			this.breakType = character
			this.numBreaks++
		}
		if this.secondDigit == 100 && this.groupLength != 4 {
			this.validMonth = true
		}
		this.resetYearValues()
	}
}
func validDayDigit(first, last byte) bool {
	if last == 100 {
		return true
	}
	if first == '3' && last > '1' {
		return false
	}
	if first > '3' && last != 100 {
		return false
	}
	return true
}
func validYearDigit(first, second, third, fourth byte) bool {
	if first > '2' {
		return false
	}
	if first == '1' && second != '9' {
		return false
	}
	if first == '2' && second > '0' {
		return false
	}
	if first == '2' && (second > '0' || third > '2') {
		return false
	}
	if first == '2' && second == '0' && third == '2' && fourth > '1' {
		return false
	}
	return true
}
func validGroupLength(length int) bool {
	return length == 1 || length == 4 || length == 2
}
func validMonthFirstLetter(first byte) bool {
	_, found := validFirst[first]
	return found
}
func isMonth(first, last byte, length int) bool {
	candidates, found := months[first]
	if !found {
		return false
	}
	candidate, found := candidates[last]
	if !found {
		return false
	}
	for _, number := range candidate {
		if number == length {
			return true
		}
	}
	return false
}
func doubleBreak(character, next byte) bool {
	return !dobBreakNotFound(character) && !dobBreakNotFound(next)
}

func (this *dobRedaction) resetYearValues() {
	this.firstDigit = 100
	this.secondDigit = 100
	this.thirdDigit = 100
	this.fourthDigit = 100
}
func (this *dobRedaction) resetMatchValues() {
	this.start = 0
	this.length = 0
	this.isCandidate = false
	this.monthStart = 0
	this.monthLength = 0
	this.monthCandidate = false
	this.totalGroupLength = 0
	this.breaks = false
	this.groupLength = 0
	this.firstDigit = 100
	this.secondDigit = 100
	this.validMonth = false
	this.validYear = false
}

var (
	months = map[byte]map[byte][]int{
		'J': {'n': []int{3}, 'y': []int{7, 4}, 'e': []int{4}, 'l': []int{3}, 'N': []int{3}, 'Y': []int{7, 4}, 'E': []int{4}, 'L': []int{3}},
		'F': {'b': []int{3}, 'y': []int{8}, 'B': []int{3}, 'Y': []int{8}},
		'M': {'h': []int{5}, 'r': []int{3}, 'y': []int{3}, 'H': []int{5}, 'R': []int{3}, 'Y': []int{3}},
		'A': {'g': []int{3}, 't': []int{6}, 'l': []int{5}, 'r': []int{3}, 'G': []int{3}, 'T': []int{6}, 'L': []int{5}, 'R': []int{3}},
		'S': {'r': []int{9}, 'p': []int{3}, 'R': []int{9}, 'P': []int{3}},
		'O': {'t': []int{3}, 'r': []int{7}, 'T': []int{3}, 'R': []int{7}},
		'N': {'v': []int{3}, 'r': []int{9}, 'V': []int{3}, 'R': []int{9}},
		'D': {'r': []int{8}, 'c': []int{3}, 'R': []int{8}, 'C': []int{3}},
	}
	validFirst = map[byte][]int{
		'J': {0},
		'F': {0},
		'M': {0},
		'A': {0},
		'S': {0},
		'O': {0},
		'N': {0},
		'D': {0},
	}
)
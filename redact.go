package redact

type Redaction struct {
	used    []bool
	matches []match
}

func New() *Redaction {
	return &Redaction{
		used:    make([]bool, 256),
		matches: make([]match, 0, 16),
	}
}

func (this *Redaction) All(input string) string {
	this.matchCreditCard(input)
	this.matchEmail(input)
	this.matchSSN(input)
	this.matchPhone(input)
	this.matchDOB(input)
	result := this.redactMatches(input)
	this.clear()
	return result
}
func (this *Redaction) clear() {
	this.matches = this.matches[0:0]
	for i := range this.used {
		this.used[i] = false
	}
}
func (this *Redaction) appendMatch(start, length int) {
	for i := start; i <= start+length; i++ {
		this.used[i] = true
	}

	this.matches = append(this.matches, match{InputIndex: start, Length: length})
}
func (this *Redaction) redactMatches(input string) string {
	if len(this.matches) == 0 {
		return input // no changes to redact
	}

	buffer := []byte(input)
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

	output := string(buffer)
	return output
}

func (this *Redaction) matchCreditCard(input string) {
	var lastDigit int
	var length int
	var numbers int
	var isOdd bool
	var isCandidate bool
	var total int
	for i := len(input) - 1; i > 0; i-- {
		character := input[i]

		if !isNumeric(input[i]) {
			if numbers > 12 && total%10 == 0 {
				this.appendMatch(lastDigit-length+1, length)
				length = 0
				total = 0
				isOdd = false
				lastDigit = i - 1
				numbers = 0
				continue
			}

			if creditCardBreakNotFound(character) && !isNumeric(input[i-1]) {
				lastDigit = i - 1
				length = 0
				isCandidate = false
				continue
			}
			if i < len(input)-1 && !creditCardBreakNotFound(input[i+1]) {
				continue
			}
			length++
		} else {
			isOdd = !isOdd
			numbers++
			number := int(character - '0')
			if !isOdd {
				number += number
				if number > 9 {
					number -= 9
				}
			}
			total += number

			if isCandidate {
				length++
			} else {
				isCandidate = true
				lastDigit = i
				numbers = 1
				if length == 0 {
					length++
				}
			}
		}
	}
	if isNumeric(input[0]) {
		isOdd = !isOdd
		numbers++
		number := int(input[0] - '0')
		if !isOdd {
			number += number
			if number > 9 {
				number -= 9
			}
		}
		total += number
		length++
	}
	if numbers > 12 && total%10 == 0 {
		this.appendMatch(lastDigit-length+1, length)
	}
}
func creditCardBreakNotFound(character byte) bool {
	return character != '-' && character != ' '
}

func (this *Redaction) matchEmail(input string) {
	var start int
	var length int
	for i := 0; i < len(input); i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !emailBreakNotFound(character) {
			start = i + 1
			length = 0
			continue
		} else {
			if character == '@' {
				this.appendMatch(start, length)
				start = i + 1
				length = 0
			}
			length++
		}
	}
}
func emailBreakNotFound(character byte) bool {
	return character != '.' && character != ' '
}

func (this *Redaction) matchPhone(input string) {
	var start int
	var length int
	var numbers int
	var isCandidate bool
	for i := 0; i < len(input)-1; i++ {
		if this.used[i] {
			continue
		}
		character := input[i]
		if !isNumeric(character) {
			if character == '+' {
				start = i + 2
				length--
				continue
			}
			if phoneBreakNotFound(character) {
				start = i + 1
				length = 0
				numbers = 0
				isCandidate = false
				continue
			}
			if isPhoneNumber(numbers) {
				this.appendMatch(start, length)
				length = 0
				numbers = 0
				start = i + 1
				continue
			}
		}
		numbers++
		if isCandidate {
			length++
		} else {
			isCandidate = true
			start = i + 1
			length = 0
			numbers = 0
		}
	}
	if isNumeric(input[len(input)-1]) {
		length++
	}
	if isPhoneNumber(length) {
		this.appendMatch(start, length)
	}
}
func phoneBreakNotFound(character byte) bool {
	return character != '-' && character != ' ' && character != '(' && character != ')'
}
func isPhoneNumber(length int) bool {
	return length >= 10 && length <= 14
}

func (this *Redaction) matchSSN(input string) {
	var start int
	var length int
	var numbers int
	var isCandidate bool
	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !isNumeric(character) {
			if isSSN(numbers) {
				this.appendMatch(start, length)
				numbers = 0
				length = 0
				start = i + 1
				isCandidate = false
				continue
			}
			if ssnBreakNotFound(character) {
				start = i + 1
				numbers = 0
				length = 0
				isCandidate = false
				continue
			}
			if isCandidate {
				length++
			}
			continue
		}
		if isCandidate {
			numbers++
			length++
		} else {
			isCandidate = true
			start = i
			numbers++
			length++
		}
	}
	if isNumeric(input[len(input)-1]) {
		numbers++
		length++
	}
	if isSSN(numbers) {
		this.appendMatch(start, length)
	}
}
func ssnBreakNotFound(character byte) bool {
	return character != '-' && character != ' '
}
func isSSN(length int) bool {
	return length == 9
}

func (this *Redaction) matchDOB(input string) {
	var start int
	var length int
	var isCandidate bool
	var monthStart int
	var monthLength int
	var monthCandidate bool
	var startChar byte
	var numbers int
	var numBreaks int

	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !isNumeric(character) {
			if dobBreakNotFound(character) || (i < len(input)-1 && doubleBreak(character, input[i+1])){
				monthLength++
				start = i + 1
				length = 0
				numbers = 0
				numBreaks = 0
				isCandidate = false
				continue
			}

			if monthLength > 2 && isMonth(startChar, input[i-1]) {
				monthCandidate = true
				continue
			}
			monthStart = i + 1
			startChar = input[i+1]
			monthLength = 0
			if isDOB(numbers, numBreaks) {
				this.appendMatch(start, length)
				length = 0
				numbers = 0
				numBreaks = 0
				start = i + 1
				isCandidate = false
				continue
			}
			if isCandidate {
				length++
			}
			numBreaks++
			continue
		}
		numbers++
		if isCandidate || monthCandidate {
			length++
		} else {
			isCandidate = true
			start = i
			numBreaks = 0
			length++
		}
		if length == 2 && monthCandidate {
			this.appendMatch(monthStart, monthLength+length+1)
			monthCandidate = false
			length = 0
			numBreaks = 0
			start = 0
			numbers = 0
			monthStart = 0
			monthLength = 0
			isCandidate = false
		}
	}
	if isNumeric(input[len(input)-1]) {
		length++
		numbers++
	}
	if isDOB(numbers, numBreaks) {
		this.appendMatch(start, length)
		numbers = 0
		numBreaks = 0
		isCandidate = false
	}
	if numbers > 8 {
		numbers = 0
		numBreaks = 0
	}
}

func doubleBreak(character, next byte) bool{
	return !dobBreakNotFound(character) && !dobBreakNotFound(next)
}

func dobBreakNotFound(character byte) bool {
	return character != '-' && character != ' ' && character != '/'
}
func isDOB(numbers, numBreaks int) bool {
	return numBreaks == 2 && numbers >= 4 && numbers <= 8
}
func isMonth(first, last byte) bool {
	_, foundFirst := firstChars[first]
	_, foundLast := lastChars[last]
	return foundFirst && foundLast
}

func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

type match struct {
	InputIndex int
	Length     int
}

var (
	firstChars = map[byte]struct{}{
		'j': {},
		'f': {},
		'm': {},
		'a': {},
		's': {},
		'o': {},
		'n': {},
		'd': {},
		'J': {},
		'F': {},
		'M': {},
		'A': {},
		'S': {},
		'O': {},
		'N': {},
		'D': {},
	}

	lastChars = map[byte]struct{}{
		'b': {},
		'h': {},
		'e': {},
		'n': {},
		'y': {},
		'l': {},
		'g': {},
		'p': {},
		't': {},
		'v': {},
		'r': {},
		'c': {},
	}
)

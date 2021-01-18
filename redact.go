package redact

type Redaction struct {
	used    []bool
	matches []match
}

func New() *Redaction {
	return &Redaction{
		used:    make([]bool, 512),
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

func isValidNetwork(character byte) bool {
	return character >= '3' && character <= '6'
}

func (this *Redaction) matchCreditCard(input string) {
	var lastDigit int
	var length int
	var numbers int
	var isOdd bool
	var isCandidate bool
	var total int
	var breaks bool
	var numBreaks int
	var breakType byte = 'x'

	for i := len(input) - 1; i > 0; i-- {
		character := input[i]
		if !isNumeric(input[i]) {
			if numbers > 12 && total%10 == 0 && breaks && isValidNetwork(input[i+1]) {
				if numBreaks == 0 || numBreaks == 3 || numBreaks == 4 {
					this.appendMatch(lastDigit-length+1, length)
				}
				breaks = false
				breakType = 'x'
				length = 0
				total = 0
				isOdd = false
				lastDigit = i - 1
				numbers = 0
				numBreaks = 0
				continue
			}
			if creditCardBreakNotFound(character) || !isNumeric(input[i-1]) {
				lastDigit = i - 1
				length = 0
				total = 0
				numbers = 0
				isCandidate = false
				breaks = false
				breakType = 'x'
				numBreaks = 0
				continue
			}
			if isCandidate {
				if breakType == character && i > 0 && isNumeric(input[i-1]) {
					breaks = true
				}
				breakType = character
				numBreaks++
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
				breakType = 'x'
				breaks = false
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
	if numbers > 12 && total%10 == 0 && isValidNetwork(input[0]) {
		this.appendMatch(lastDigit-length+1, length)
		breaks = false
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

//Check for paren
//Check types, if they are different then not a match. <- will skip if it is a paren
// breaks must be greater than 1 and less than or equal to 3
func (this *Redaction) matchPhone(input string) {
	var start int
	var length int
	var numbers int
	var parenBreak bool
	var breakType byte
	var matchBreaks bool
	var breaks int
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
				numbers--
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
				if correctBreaks(breaks, parenBreak, matchBreaks) {
					this.appendMatch(start, length)
				}
				length = 0
				numbers = 0
				breaks = 0
				start = i + 1
				matchBreaks = false
				parenBreak = false
				continue
			}
			if character == '(' {
				start = i
				isCandidate = true
				length++
				breaks++
				parenBreak = true
				continue
			}
			if character == ')' {
				length++
				breaks++
				parenBreak = true
				continue
			}
			if i < len(input)-1 && !isNumeric(input[i+1]) {
				length = 0
				numbers = 0
				start = i + 1
				breaks = 0
				parenBreak = false
				continue
			}
			if isCandidate {
				length++
				if breakType == character && numbers == 6 {
					matchBreaks = true
				} else {
					matchBreaks = false
				}
				if numbers == 3 {
					breakType = character
				}
			}
			breaks++
			continue
		}
		if isCandidate {
			length++
			numbers++
		} else {
			isCandidate = true
			start = i
			length++
			numbers++
			breaks = 0
			matchBreaks = false
			parenBreak = false
		}
	}
	if isNumeric(input[len(input)-1]) {
		length++
		numbers++
	}
	if isPhoneNumber(numbers) {
		if correctBreaks(breaks, parenBreak, matchBreaks) {
			this.appendMatch(start, length)
		}
		length = 0
		numbers = 0
		breaks = 0
		parenBreak = false
		matchBreaks = false
	}
}
func phoneBreakNotFound(character byte) bool {
	return character != '-' && character != ' ' && character != '(' && character != ')'
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

func (this *Redaction) matchSSN(input string) {
	var start int
	var length int
	var numbers int
	var breaks bool
	var numBreaks int
	var breakType byte = 'x'
	var isCandidate bool
	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !isNumeric(character) {
			if isSSN(numbers) && breaks && numBreaks == 2 {
				this.appendMatch(start, length)
				numbers = 0
				breaks = false
				breakType = 'x'
				numBreaks = 0
				length = 0
				start = i + 1
				isCandidate = false
				continue
			}
			if ssnBreakNotFound(character) {
				start = i + 1
				numbers = 0
				breaks = false
				breakType = 'x'
				numBreaks = 0
				length = 0
				isCandidate = false
				continue
			}
			if isCandidate && i < len(input)-1 && isNumeric(input[i+1]) {
				length++
				if breakType == character && numbers == 5 {
					breaks = true
				}
				if numbers == 3 {
					breakType = character
				}
			}
			numBreaks++
			continue
		}
		if isCandidate {
			numbers++
			length++
		} else {
			isCandidate = true
			breaks = false
			start = i
			numbers++
			length++
		}
	}
	if isNumeric(input[len(input)-1]) {
		numbers++
		length++
	}
	if isSSN(numbers) && breaks {
		this.appendMatch(start, length)
	}
}
func ssnBreakNotFound(character byte) bool {
	return character != '-' && character != ' '
}
func isSSN(length int) bool {
	return length == 9
}

//Breaks must be in correct place
func (this *Redaction) matchDOB(input string) {
	var start int
	var length int
	var isValidFirstChar bool
	var firstByte byte
	var isValidMonth bool
	var numberFound bool

	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !isNumeric(character) {
			if !isValidFirstChar && isValidFirstMonthCharacter(character) {
				firstByte = character
				isValidFirstChar = true
				start = i
				length++
				continue
			}
			if i > 1 && isMonth(firstByte, input[i-1], length){
				isValidMonth = true
			}
			if !dobBreakNotFound(character) {
				if !isValidMonth {
					isValidFirstChar = false
					firstByte = 'x'
					start = i + 1
					length = 0
					isValidMonth = false
					numberFound = false
					continue
				}
			}
			if !dobBreakNotFound(character) && isValidMonth && numberFound{
				this.appendMatch(start, length)
				isValidFirstChar = false
				firstByte = 'x'
				start = i + 1
				length = 0
				isValidMonth = false
				numberFound = false
				continue
			}
			length++
		} else {
			numberFound = true
			length++
			continue
		}
	}
}
func dobBreakNotFound(character byte) bool {
	return character != ' ' && character != ','
}

func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
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
func isValidFirstMonthCharacter(first byte) bool {
	_, found := firstLetterOfMonth[first]
	if found {
		return true
	}
	return false
}

type match struct {
	InputIndex int
	Length     int
}

var (
	months = map[byte]map[byte][]int{
		'j': {'n': []int{3}, 'y': []int{7, 4}, 'e': []int{4}, 'l': []int{3}},
		'f': {'b': []int{3}, 'y': []int{8}},
		'm': {'h': []int{5}, 'r': []int{3}, 'y': []int{3}},
		'a': {'g': []int{3}, 't': []int{6}, 'l': []int{5}, 'r': []int{3}},
		's': {'r': []int{9}, 'p': []int{3}, 't': []int{4}},
		'o': {'t': []int{3}, 'r': []int{7}},
		'n': {'v': []int{3}, 'r': []int{9}},
		'd': {'r': []int{8}, 'c': []int{3}},
		'J': {'n': []int{3}, 'y': []int{7, 4}, 'e': []int{4}, 'l': []int{3}},
		'F': {'b': []int{3}, 'y': []int{8}},
		'M': {'h': []int{5}, 'r': []int{3}, 'y': []int{3}},
		'A': {'g': []int{3}, 't': []int{6}, 'l': []int{5}, 'r': []int{3}},
		'S': {'r': []int{9}, 'p': []int{3}, 't': []int{4}},
		'O': {'t': []int{3}, 'r': []int{7}},
		'N': {'v': []int{3}, 'r': []int{9}},
		'D': {'r': []int{8}, 'c': []int{3}},
	}
	firstLetterOfMonth = map[byte][]int{
		'j': {1},
		'f': {1},
		'm': {1},
		'a': {1},
		's': {1},
		'o': {1},
		'n': {1},
		'd': {1},
		'J': {1},
		'F': {1},
		'M': {1},
		'A': {1},
		'S': {1},
		'O': {1},
		'N': {1},
		'D': {1},
	}
)

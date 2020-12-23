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
				if breakType == character{
					matchBreaks = true
				}else{
					matchBreaks = false
				}
				breakType = character
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
	if breaks == 4 && parenBreak{
		return true
	}
	if breaks == 2 && matchBreak{
		return true
	}
	return false
}

func (this *Redaction) matchSSN(input string) {
	var start int
	var length int
	var numbers int
	var breaks bool
	var breakType byte = 'x'
	var isCandidate bool
	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !isNumeric(character) {
			if isSSN(numbers) && (breakType == 'x' || breaks) {
				this.appendMatch(start, length)
				numbers = 0
				breaks = false
				breakType = 'x'
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
				length = 0
				isCandidate = false
				continue
			}
			if isCandidate && i < len(input)-1 && isNumeric(input[i+1]) {
				length++
				if breakType == character {
					breaks = true
				}
				breakType = character
			}
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
	if isSSN(numbers) && (breakType == 'x' || breaks) {
		this.appendMatch(start, length)
		breaks = false
		breakType = 'x'
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
	var breaks bool
	var breakType byte

	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !isNumeric(character) {
			if dobBreakNotFound(character) || (i < len(input)-1 && doubleBreak(character, input[i+1])) {
				monthLength++
				start = i + 1
				length = 0
				numbers = 0
				breaks = false
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
			if isDOB(numbers) && breaks {
				this.appendMatch(start, length)
				length = 0
				numbers = 0
				breaks = false
				start = i + 1
				isCandidate = false
				continue
			}
			if isCandidate {
				length++
			}
			if character == breakType {
				breaks = true
			}else{
				breaks = false
			}
			breakType = character
			continue
		}
		numbers++
		if isCandidate || monthCandidate {
			length++
		} else {
			isCandidate = true
			breakType = 'x'
			start = i
			breaks = false
			length++
		}
		if length == 2 && monthCandidate {
			this.appendMatch(monthStart, monthLength+length+1)
			monthCandidate = false
			length = 0
			breaks = false
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
	if isDOB(numbers){
		this.appendMatch(start, length)
		numbers = 0
		breaks = false
		isCandidate = false
	}
	if numbers > 8 {
		numbers = 0
		breaks = false
	}
}

func doubleBreak(character, next byte) bool {
	return !dobBreakNotFound(character) && !dobBreakNotFound(next)
}

func dobBreakNotFound(character byte) bool {
	return character != '-' && character != ' ' && character != '/'
}
func isDOB(numbers int) bool {
	return numbers >= 4 && numbers <= 8
}

func isMonth(first, last byte) bool {
	candidates, found := months[first]
	if !found {
		return false
	}
	for _, candidate := range candidates {
		if candidate == last {
			return true
		}
	}
	return false
}

func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

type match struct {
	InputIndex int
	Length     int
}

var (
	months = map[byte][]byte{
		'j': {'n', 'y', 'e', 'l'},
		'f': {'b', 'y'},
		'm': {'h', 'r', 'y'},
		'a': {'g', 't', 'l', 'r'},
		's': {'r', 'p', 't'},
		'o': {'t', 'r'},
		'n': {'v', 'r'},
		'd': {'r', 'c'},
		'J': {'n', 'y', 'e', 'l'},
		'F': {'b', 'y'},
		'M': {'h', 'r', 'y'},
		'A': {'g', 't', 'l', 'r'},
		'S': {'r', 'p', 't'},
		'O': {'t', 'r'},
		'N': {'v', 'r'},
		'D': {'r', 'c'},
	}
)

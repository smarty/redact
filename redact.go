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

func (this *Redaction) matchDOB(input string) {
	var start int
	var length int
	var isCandidate bool
	var monthStart int
	var monthLength int
	var monthCandidate bool
	var startChar byte
	var totalGroupLength int
	var breaks bool
	var numBreaks int
	var breakType byte
	var groupLength int
	var firstDigit byte
	var secondDigit byte
	var thirdDigit byte
	var fourthDigit byte
	var validMonth bool
	var validYear bool

	firstDigit = 100
	secondDigit = 100
	thirdDigit = 100
	fourthDigit = 100
	startChar = 'x'

	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		variable := string(character)
		_ = variable
		if this.used[i] {
			continue
		}
		if !isNumeric(character) {
			if isValidFirstLetter(character) && startChar == 'x' {
				startChar = character
				monthStart = i
			}
			if isDOB(totalGroupLength) && breaks && numBreaks == 2 && validYear {
				if groupLength == 2 && validDateDigit(firstDigit, secondDigit) && validMonth || groupLength == 4 {
					this.appendMatch(start, length)
					startChar = 'x'
				}
				if firstDigit > '1' {
					validMonth = true
				}
				breakType = 'x'
				numBreaks = 0
				length = 0
				totalGroupLength = 0
				breaks = false
				start = i + 1
				isCandidate = false
				groupLength = 0
				firstDigit = 100
				secondDigit = 100
				validYear = false
				continue
			}
			if numBreaks == 2 {
				breaks = false
			}
			if character == ' ' {
				if monthLength > 2 && isMonth(startChar, input[i-1], monthLength) {
					monthCandidate = true
					monthLength++
					continue
				} else{
					monthCandidate = false
					monthStart = 0
					monthLength = 0
					startChar = 'x'
				}
			}
			if dobBreakNotFound(character) || (i < len(input)-1 && doubleBreak(character, input[i+1])) {
				if character == ',' && monthCandidate {
					this.appendMatch(monthStart, monthLength+length+1)
					monthCandidate = false
					length = 0
					breaks = false
					start = 0
					totalGroupLength = 0
					monthStart = 0
					monthLength = 0
					isCandidate = false
					groupLength = 0
					firstDigit = 100
					secondDigit = 100
					validMonth = false
					validYear = false
					continue
				}
				if startChar != 'x' && character != ' ' {
					monthLength++
				}
				start = i + 1
				length = 0
				totalGroupLength = 0
				breaks = false
				isCandidate = false
				groupLength = 0
				numBreaks = 0
				validMonth = false
				monthCandidate = false
				continue
			}
			monthStart = i + 1
			monthLength = 0

			if isCandidate {
				length++
			}
			if firstDigit == '1' && secondDigit <= '2' && groupLength != 4 {
				validMonth = true
			}
			if validDateDigit(firstDigit, secondDigit) || (totalGroupLength == 4 && validYear) {
				if character == breakType && validGroupLength(groupLength) {
					breaks = true
					numBreaks++
				}
				if validGroupLength(groupLength) && totalGroupLength < 3 || validYear && !breaks {
					breakType = character
					numBreaks++
				}
				if secondDigit == 100 && groupLength != 4 {
					validMonth = true
				}
				firstDigit = 100
				secondDigit = 100
			}
			groupLength = 0
			continue
		}

		totalGroupLength++
		groupLength++
		if firstDigit == 100 && groupLength < 3 {
			firstDigit = character
		} else {
			if groupLength == 2 {
				secondDigit = character
			}
		}
		if groupLength == 3 {
			thirdDigit = character
		}
		if groupLength == 4 {
			fourthDigit = character
			validYear = validYearDigit(firstDigit, secondDigit, thirdDigit, fourthDigit)
			firstDigit = 100
			secondDigit = 100
			thirdDigit = 100
			fourthDigit = 100
		}
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
			breakType = 'x'
			startChar = 'x'
			numBreaks = 0
			monthCandidate = false
			length = 0
			breaks = false
			start = 0
			totalGroupLength = 0
			monthStart = 0
			monthLength = 0
			isCandidate = false
			groupLength = 0
			firstDigit = 100
			secondDigit = 100
			validMonth = false
		}
	}
	if isNumeric(input[len(input)-1]) {
		length++
		totalGroupLength++
		groupLength++
		if groupLength == 4 {
			fourthDigit = input[len(input) - 1]
			validYear = validYearDigit(firstDigit, secondDigit, thirdDigit, fourthDigit)
		}
	}
	if isDOB(totalGroupLength) && breaks && numBreaks == 2 && validYear && validMonth {
		this.appendMatch(start, length)
		totalGroupLength = 0
		breaks = false
		isCandidate = false
		validMonth = false
		firstDigit = 100
		secondDigit = 100
		groupLength = 0
		startChar = 'x'
		validYear = false
	}
	if totalGroupLength > 8 {
		totalGroupLength = 0
		breaks = false
	}
}
func dobBreakNotFound(character byte) bool {
	return character != '/' && character != '-'
}
func isDOB(numbers int) bool {
	return numbers >= 6 && numbers <= 8
}

func validDateDigit(first, last byte) bool {
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
	if first == '2' && (second > '0' || third > '2'){
		return false
	}
	if first == '2' && second == '0' && third == '2' && fourth > '1'{
		return false
	}
	return true
}
func validGroupLength(length int) bool {
	return length == 1 || length == 4 || length == 2
}
func isValidFirstLetter(first byte) bool {
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
func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

type match struct {
	InputIndex int
	Length     int
}

var (
	months = map[byte]map[byte][]int{
		'J': {'n': []int{3}, 'y': []int{7, 4}, 'e': []int{4}, 'l': []int{3}},
		'F': {'b': []int{3}, 'y': []int{8}},
		'M': {'h': []int{5}, 'r': []int{3}, 'y': []int{3}},
		'A': {'g': []int{3}, 't': []int{6}, 'l': []int{5}, 'r': []int{3}},
		'S': {'r': []int{9}, 'p': []int{3}},
		'O': {'t': []int{3}, 'r': []int{7}},
		'N': {'v': []int{3}, 'r': []int{9}},
		'D': {'r': []int{8}, 'c': []int{3}},
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

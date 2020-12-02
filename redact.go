package redact

type Redaction struct {
	used    map[int]struct{}
	matches []match
}

func New() *Redaction {
	return &Redaction{used: make(map[int]struct{}, 128)}
}

func (this *Redaction) All(input string) string {
	this.matchCreditCard(input)
	this.matchEmail(input)
	this.matchPhoneNum(input)
	this.matchSSN(input)
	this.matchDOB(input)
	result := this.redactMatches(input)
	this.clear()
	return result
}

func (this *Redaction) clear() {
	this.matches = this.matches[0:0]
	for key, _ := range this.used {
		delete(this.used, key)
	}
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
		if highIndex >= bufferLength {
			continue
		}
		for ; lowIndex < highIndex; lowIndex++ {
			buffer[lowIndex] = '*'
		}
	}

	return string(buffer)
}

func (this *Redaction) matchCreditCard(input string) (matches []match) {
	// TODO: inline checkLuhn into this algorithm--this avoids having to create a string to ask if it's a credit card
	// instead we track each numeric digit here and run a tally as we go along
	var start int
	var length int
	var isCandidate bool
	var total int
	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if !isNumeric(character) {
			if isCreditCard(length, input[start:start+length]) {
				this.appendMatch(start, length)
				length = 0
				start = i + 1
				total = 0
				continue
			}
			if breakNotFound(character) && !isNumeric(input[i+1]) {
				start = i + 1
				length = 0
				isCandidate = false
				continue
			}
			length++
		} else {
			number := int(character - '0')
			if isCandidate {
				length++
				total += number
			} else {
				isCandidate = true
				start = i
				total = number
			}
		}
	}
	if isNumeric(input[len(input)-1]) {
		//length++ // TODO: test coverage
	}
	if isCreditCard(length, input[start:start+length]) {
		// matches = appendMatches(matches, start, length) // TODO: test coverage
	}

	return matches
}
func isCreditCard(length int, input string) bool {
	return length >= 13 && length <= 24 && checkLuhn(input)
}
func checkLuhn(input string) bool {
	nDigits := len(input)
	var nSum int
	var isSecond bool
	for i := nDigits - 2; i >= 0; i-- {
		if !isNumeric(input[i]) {
			continue
		}
		d := input[i] - '0'
		if isSecond == false {
			d = d * 2
		}
		if d > 9 {
			d -= 9
		}
		digit := int(d)
		nSum += digit
		isSecond = !isSecond
	}
	temp := int(input[nDigits-1] - '0')
	mod := nSum % 10
	return mod == temp
}

func (this *Redaction) matchEmail(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool
	for i := 0; i < len(input); i++ {
		character := input[i]
		if _, found := this.used[i]; found {
			continue
		}
		if !breakNotFound(character) {
			if isCandidate {
				this.appendMatch(start, length)
			}
			start = i + 1
			length = 0
			isCandidate = false
			continue
		} else {
			length++
			if character == '@' {
				isCandidate = true
			}
		}
	}
	return matches
}

func (this *Redaction) matchPhoneNum(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool
	for i := 0; i < len(input)-1; i++ {
		if _, found := this.used[i]; found {
			continue
		}
		character := input[i]
		if !isNumeric(character) {
			if character == '+' {
				start = i + 2
				length--
				continue
			}
			if breakNotFound(character) {
				start = i + 1
				length = 0
				isCandidate = false
				continue
			}
			if isPhoneNumber(length) {
				this.appendMatch(start, length)
				length = 0
				start = i + 1
				continue
			}
		}
		if isCandidate {
			length++
		} else {
			isCandidate = true
			start = i + 1
			length = 0
		}
	}
	if isNumeric(input[len(input)-1]) {
		length++
	}
	if isPhoneNumber(length) {
		this.appendMatch(start, length)
	}
	return matches
}
func isPhoneNumber(length int) bool {
	return length >= 10 && length <= 14
}

func (this *Redaction) matchSSN(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool
	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if _, found := this.used[i]; found {
			continue
		}
		if !isNumeric(character) {
			if isSSN(length) {
				this.appendMatch(start, length)
				length = 0
				start = i + 1
				continue
			}
			if breakNotFound(character) {
				start = i + 1
				length = 0
				isCandidate = false
				continue
			}
		}
		if isCandidate {
			length++
		} else {
			isCandidate = true
			start = i + 1
		}
	}
	if isNumeric(input[len(input)-1]) {
		// length++ // TODO: test coverage
	}
	if isSSN(length) {
		// matches = appendMatches(matches, start, length) // TODO: test coverage
	}
	return matches
}
func isSSN(length int) bool {
	return length >= 9 && length <= 11
}

func (this *Redaction) matchDOB(input string) (matches []match) {
	var start int
	var length int
	var isCandidate bool
	var monthStart int
	var monthLength int
	var monthCandidate bool

	for i := 0; i < len(input)-1; i++ {
		character := input[i]
		if _, found := this.used[i]; found {
			continue
		}
		if !isNumeric(character) {
			if breakNotFound(character) {
				monthLength++
				start = i + 1
				length = 0
				isCandidate = false
				continue
			}
			if isMonth(input[monthStart : monthStart+monthLength]) {
				monthCandidate = true
				continue
			}
			monthStart = i + 1
			monthLength = 0
			if isDOB(length) {
				this.appendMatch(start, length)
				length = 0
				start = i + 1
				continue
			}
		}
		if isCandidate || monthCandidate {
			length++
		} else {
			isCandidate = true
			start = i + 1
		}
		if length == 2 && monthCandidate {
			this.appendMatch(monthStart, monthLength+length+1)
			monthCandidate = false
			length = 0
			start = 0
			monthStart = 0
			monthLength = 0
		}
	}
	if isNumeric(input[len(input)-1]) {
		// length++ // TODO: test coverage
	}
	if isDOB(length) {
		// matches = appendMatches(matches, start, length) // TODO: test coverage
	}
	return matches
}
func isDOB(length int) bool {
	return length >= 6 && length <= 10
}
func isMonth(month string) bool {
	_, found := months[month]
	return found
}
func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

func breakNotFound(character byte) bool {
	return character != '-' && character != ' ' && character != '.' && character != '(' && character != ')' && character != '/'
}
func (this *Redaction) appendMatch(start, length int) {
	for i := start; i <= start+length; i++ {
		this.used[i] = struct{}{}
	}

	this.matches = append(this.matches, match{InputIndex: start, Length: length})
}

type match struct {
	InputIndex int
	Length     int
}

var (
	months = map[string]struct{}{
		"January":   {},
		"Jan":       {},
		"February":  {},
		"Feb":       {},
		"March":     {},
		"Mar":       {},
		"April":     {},
		"Apr":       {},
		"May":       {},
		"June":      {},
		"Jun":       {},
		"July":      {},
		"Jul":       {},
		"August":    {},
		"Aug":       {},
		"September": {},
		"Sep":       {},
		"Sept":      {},
		"October":   {},
		"Oct":       {},
		"November":  {},
		"Nov":       {},
		"December":  {},
		"Dec":       {},
	}
)

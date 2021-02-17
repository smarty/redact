package redact

type Redaction struct {
	used    []bool
	matches []match
	phone   *phoneRedaction
	ssn     *ssnRedaction
	credit  *creditCardRedaction
	dob     *dobRedaction
}

func New() *Redaction {
	return &Redaction{
		used:    make([]bool, 512),
		matches: make([]match, 0, 16),
		phone:   &phoneRedaction{},
		ssn:     &ssnRedaction{},
		credit:  &creditCardRedaction{},
		dob:     &dobRedaction{},
	}
}
func (this *Redaction) All(input string) string {
	this.clear()
	this.matchCreditCard(input)
	this.matchEmail(input)
	this.matchSSN(input)
	this.matchPhone(input)
	this.matchDOB(input)
	result := this.redactMatches(input)
	return result
}
func (this *Redaction) clear() {
	this.matches = this.matches[0:0]
	for i := range this.used {
		this.used[i] = false
	}
	this.phone.clear()
	this.ssn.clear()
	this.credit.clear()
	this.dob.clear()
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
	this.credit.used = this.used
	this.credit.matches = this.matches

	this.credit.match(input)

	this.used = this.credit.used
	this.matches = this.credit.matches
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
	this.phone.used = this.used
	this.phone.matches = this.matches

	this.phone.match(input)

	this.used = this.phone.used
	this.matches = this.phone.matches
}

func (this *Redaction) matchSSN(input string) {
	this.ssn.matches = this.matches
	this.ssn.used = this.used

	this.ssn.match(input)

	this.matches = this.ssn.matches
	this.used = this.ssn.used
}

func (this *Redaction) matchDOB(input string) {
	this.dob.used = this.used
	this.dob.matches = this.matches

	this.dob.match(input)
	this.used = this.dob.used
	this.matches = this.dob.matches
}


func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

type match struct {
	InputIndex int
	Length     int
}


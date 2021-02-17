package redact

type Redaction struct {
	*matched
	phone  *phoneRedaction
	ssn    *ssnRedaction
	credit *creditCardRedaction
	dob    *dobRedaction
	email  *emailRedaction
}

func New() *Redaction {
	matched := &matched{
		used:    make([]bool, 512),
		matches: make([]match, 0, 16),
	}
	return &Redaction{
		matched: matched,
		phone:   &phoneRedaction{matched: matched},
		ssn:     &ssnRedaction{matched: matched},
		credit:  &creditCardRedaction{matched: matched},
		dob:     &dobRedaction{matched: matched},
		email:   &emailRedaction{matched: matched},
	}
}

func (this *Redaction) All(input string) string {
	this.clear()
	this.credit.match(input)
	this.email.match(input)
	this.ssn.match(input)
	this.phone.match(input)
	this.dob.match(input)
	result := this.redactMatches(input)
	return result
}
func (this *Redaction) clear() {
	this.matched.clear()
	this.phone.clear()
	this.ssn.clear()
	this.credit.clear()
	this.dob.clear()
	this.email.clear()
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

func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

type match struct {
	InputIndex int
	Length     int
}
type matched struct {
	used    []bool
	matches []match
}

func (this *matched) appendMatch(start, length int) {
	for i := start; i <= start+length; i++ {
		this.used[i] = true
	}

	this.matches = append(this.matches, match{InputIndex: start, Length: length})
}
func (this *matched) clear() {
	this.matches = this.matches[0:0]
	for i := range this.used {
		this.used[i] = false
	}
}

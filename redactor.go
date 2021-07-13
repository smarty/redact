package redact

type Redactor struct {
	*matched
	phone   *phoneRedaction
	ssn     *ssnRedaction
	credit  *creditCardRedact
	dob     *dobRedaction
	email   *emailRedaction
	monitor monitor
	//credit *creditCardRedact
}

func (this *Redactor) All(input []byte) []byte {
	this.clear()
	//TODO: ADD
	//this is your added method
	this.credit.match(input)
	this.email.match(input)
	this.ssn.match(input)
	this.phone.match(input)
	this.dob.match(input)
	result := this.redactMatches(input)
	return result
}
func (this *Redactor) clear() {
	this.matched.clear()
	this.credit.clear()
	this.email.clear()
	this.ssn.clear()
	this.phone.clear()
	this.dob.clear()
}

func (this *Redactor) redactMatches(input []byte) []byte {
	count := len(this.matches)
	if count == 0 {
		return input // no changes to redact
	}
	this.monitor.Redacted(count)

	buffer := input
	var lowIndex, highIndex int

	for _, match := range this.matches {
		lowIndex = match.InputIndex
		highIndex = lowIndex + match.Length
		if lowIndex < 0 {
			continue
		}
		if highIndex > len(buffer) {
			continue
		}
		for ; lowIndex < highIndex; lowIndex++ {
			buffer[lowIndex] = '*'
		}
	}

	output := buffer
	return output
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

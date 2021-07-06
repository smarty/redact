package redact

type phoneRedaction struct {
	*matched
	start       int
	length      int
	numbers     int
	breaks      int
	parenBreak  bool
	matchBreaks bool
	breakType   byte
	isCandidate bool
}
type ssnRedaction struct {
	*matched
	start       int
	length      int
	numbers     int
	breaks      bool
	numBreaks   int
	breakType   byte
	isCandidate bool
}
type creditCardRedaction struct {
	*matched
	lastDigitIndex int
	length         int
	totalNumbers   int
	isOdd          bool
	isCandidate    bool
	isLetter       bool
	totalSum       int
	numBreaks      int
	groupLength    int
	numGroups      int
	breakType      byte
}
type emailRedaction struct {
	*matched
	start  int
	length int
}
func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

type Redactor struct {
	*matched
	phone   *phoneRedaction
	ssn     *ssnRedaction
	credit  *creditCardRedaction
	dob     *dobRedaction
	email   *emailRedaction
	monitor monitor
}
type match struct {
	InputIndex int
	Length     int
}
type matched struct {
	used    []bool
	matches []match
}
func (this *matched) clear() {
	this.matches = this.matches[0:0]
	for i := range this.used {
		this.used[i] = false
	}
}

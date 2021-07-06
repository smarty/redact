package redact

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
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Redactor struct {
	*matched
	phone   *phoneRedaction
	ssn     *ssnRedaction
	credit  *creditCardRedaction
	dob     *dobRedaction
	email   *emailRedaction
	monitor monitor
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type emailRedaction struct {
	*matched
	start  int
	length int
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (this *matched) appendMatch(start, length int) {
	for i := start; i <= start+length; i++ {
		this.used[i] = true
	}
	this.matches = append(this.matches, match{InputIndex: start, Length: length})
}
func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

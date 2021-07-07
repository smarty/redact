package redact

type Redactor struct {
	*matched
	phone   *phoneRedaction // FIXME: Could we seperate these into it's own struct?
	ssn     *ssnRedaction
	credit  *creditCardRedaction
	dob     *dobRedaction
	email   *emailRedaction
	monitor monitor
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type phoneRedaction struct {
	*matched
	start         int
	length        int
	breakLength   int
	numericLength int
}
type emailRedaction struct {
	*matched
	start  int
	length int
}
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
type ssnRedaction struct {
	*matched
	start       int
	length      int
	breakLength int
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func isNumeric(value byte) bool {
	return value >= '0' && value <= '9'
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	MaxEmailLength = 254

	MaxPhoneLength_WithBreaks   = 14
	MinPhoneLength_WithBreaks   = 12
	MinPhoneLength_WithNoBreaks = 10
	MaxPhoneBreakLength         = 3
	MinPhoneBreakLength         = 2

	MaxSSNLength_WithBreaks = 11
	MaxSSNBreakLength       = 2
	MaxSSNBreakPosition     = 6
	MinSSNBreakPosition     = 3
)

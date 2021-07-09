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
type ssnRedaction struct {
	*matched
	start       int
	length      int
	breakLength int
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
type dobRedaction struct {
	*matched
	start           int
	length          int
	numBreaks       int
	numValues       int
	validMonthFound bool
}
type emailRedaction struct {
	*matched
	start  int
	length int
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

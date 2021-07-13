package redact

type Redactor2 struct {
	*matched
	phone   *phoneRedaction // FIXME: Could we seperate these into it's own struct?
	ssn     *ssnRedaction
	credit  *creditCardRedact
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
type ssnRedaction2 struct {
	*matched
	start       int
	length      int
	breakLength int
}

type creditCardRedact struct {
	*matched
	lastDigitIndex int
	value          []byte
	length         int
	numericLength  int
	breakType      byte
}

type dobRedaction2 struct {
	*matched
	start           int
	length          int
	numBreaks       int
	numValues       int
	validMonthFound bool
}
type emailRedaction2 struct {
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

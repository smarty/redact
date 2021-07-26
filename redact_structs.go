package redact

type Redactor struct {
	*matched
	phone   *phoneRedaction
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
	start     int
	length    int
	isSecond  bool
	totalSum  uint64
	numBreaks int
	breakType byte
}
type dobRedaction struct {
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

	MinDOBLength_WithBreaks = 8
	MaxDOBLength_WithBreaks = 11
	MaxDOBBreakLength       = 2

	MinCreditLength_NoBreaks = 12
	MaxCreditLength_NoBreaks = 19
	MaxCreditBreakLength     = 4
	MinCreditBreakLength     = 2
)

package redact

func (this *dobRedaction) clear() {
	this.start = 0
	this.length = 0
	this.numValues = 0
	this.numBreaks = 0
	this.validMonthFound = false
}

/* Valid Dates:
- Use a specific pattern (1/1/1111, 1111/1/1, Mar 3, 2222, etc.)
	-11/11/1111, 1/1/1111, 1/11/1111, 11/1/1111
	- Be careful with leap years, february will have more or less available days
- Must be a valid date (months + days match up) year is not out of range/ within 100 year span
- Must have at least 2 breaks, potentially 3
	- A break may be: / , ' '
- Months may come in all CAPS or all lowercase
- May also have abbreviated months
*/

/* Months with 31 days
- Jan, march, may, july, august, october, december

  Months with 30 days
- april, june, september, november

  Weird months
- February
*/
func (this *dobRedaction) match(input []byte) {
	for i := 0; i < len(input); i++ {
		if i < len(this.used)-1 && this.used[i] {
			continue
		}
		if isNumeric(input[i]) {
			this.numValues++
			this.length++
			continue
		}
		if this.length > 8 && this.length < 11{
			if this.numBreaks == 2 && this.validMonthFound {
				this.appendMatch(this.start, this.length)
				this.resetCount(i)
				continue
			}
		}
		this.validateBreaks(input[i], i)
		if !this.validateDOB(input[i-this.numValues : i]) {
			this.resetCount(i)
		}
		if this.length > 11 {
			this.resetCount(i)
		}
	}
	if this.length >= 8 && this.length <= 10 && this.numBreaks == 2 && this.validMonthFound {
		if this.validateDOB(input[this.length - this.numValues: this.length]){
			this.appendMatch(this.start, this.length)
		}
	}
}

func (this *dobRedaction) validateDOB(input []byte) bool {
	this.numValues = 0
	switch len(input) {
	case 0:
		return true
	case 4:
		return this.validateYear(input)
 	case 2:
		return this.validateDate(input)
	case 1:
		if input[0] != '0' {
			return true
		}
	default:
		return false
	}
	return false
}
func (this *dobRedaction) validateDate(input []byte) bool { //Valid day:
	switch {
	case input[0] == '0' && input[1] <= '9' && input[1] > 0:
		this.validMonthFound = true
		return true
	case input[0] == '1' && input[1] <= '2':
		this.validMonthFound = true
		return true
	case input[0] == '3' && input[1] > 1:
		return false
	default:
		return false
	}
}
func (this *dobRedaction) validateYear(input []byte) bool {
	switch {
	case input[0] == '1' && input[1] == '9':
		return true
	case input[0] == '2' && input[1] == '0' && input[2] < '3':
		if input[2] == '2' && input[3] > '1' {
			return false
		}
		return true
	default:
		return false
	}
	return false
}

func (this *dobRedaction) validateBreaks(input byte, i int) {
	switch {
	case input == '/':
		this.length++
		this.numBreaks++
	case input == '-':
		this.length++
		this.numBreaks++
	case this.numBreaks > 2:
		this.resetCount(i)
	default:
		this.resetCount(i)
	}
}

func (this *dobRedaction) resetCount(i int) {
	this.start = i + 1
	this.length = 0
	this.numBreaks = 0
	this.numValues = 0
	this.validMonthFound = false
}

/*
// MINIMUM: 2
// MAXIMUM: 43
 1 - 12 ** months
 1 - 31 ** days
*/

func validDayDigit(first, last byte) bool {
	if last == 100 {
		return true
	}
	if first == '3' && last > '1' {
		return false
	}
	if first > '3' && last != 100 {
		return false
	}
	return true
}
func validYearDigit(first, second, third, fourth byte) bool {
	if first > '2' {
		return false
	}
	if first == '1' && second != '9' {
		return false
	}
	if first == '2' && second > '0' {
		return false
	}
	if first == '2' && (second > '0' || third > '2') {
		return false
	}
	if first == '2' && second == '0' && third == '2' && fourth > '1' {
		return false
	}
	return true
}
func validGroupLength(length int) bool {
	return length == 1 || length == 4 || length == 2
}
func validMonthFirstLetter(first byte) bool {
	_, found := validFirst[first]
	return found
}
func isMonth(first, last byte, length int) bool {
	candidates, found := months[first]
	if !found {
		return false
	}
	candidate, found := candidates[last]
	if !found {
		return false
	}
	for _, number := range candidate {
		if number == length {
			return true
		}
	}
	return false
}

var (
	months = map[byte]map[byte][]int{
		'J': {'n': []int{3}, 'y': []int{7, 4}, 'e': []int{4}, 'l': []int{3}, 'N': []int{3}, 'Y': []int{7, 4}, 'E': []int{4}, 'L': []int{3}},
		'F': {'b': []int{3}, 'y': []int{8}, 'B': []int{3}, 'Y': []int{8}},
		'M': {'h': []int{5}, 'r': []int{3}, 'y': []int{3}, 'H': []int{5}, 'R': []int{3}, 'Y': []int{3}},
		'A': {'g': []int{3}, 't': []int{6}, 'l': []int{5}, 'r': []int{3}, 'G': []int{3}, 'T': []int{6}, 'L': []int{5}, 'R': []int{3}},
		'S': {'r': []int{9}, 'p': []int{3}, 'R': []int{9}, 'P': []int{3}},
		'O': {'t': []int{3}, 'r': []int{7}, 'T': []int{3}, 'R': []int{7}},
		'N': {'v': []int{3}, 'r': []int{9}, 'V': []int{3}, 'R': []int{9}},
		'D': {'r': []int{8}, 'c': []int{3}, 'R': []int{8}, 'C': []int{3}},
	}
	validFirst = map[byte][]int{
		'J': {0},
		'F': {0},
		'M': {0},
		'A': {0},
		'S': {0},
		'O': {0},
		'N': {0},
		'D': {0},
	}
)

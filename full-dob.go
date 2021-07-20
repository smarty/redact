package redact

type FullDOB struct {
	redact          *dobRedaction
	validMonthFound bool
	validDayFound   bool
}

func (this *FullDOB) findMatch(input []byte) {
	for i := 0; i < len(input) - 1; i++ {
		if i < len(this.redact.used)-1 && this.redact.used[i] {
			continue
		}
		if i < len(input)-1 && i > 0 && !this.validMonthFound && input[i] == ' ' {
			if !this.isMonth(input[this.redact.start], input[i-1], this.redact.length) {
				this.resetCount(i)
				continue
			}
			this.redact.length++
			this.validMonthFound = true
			continue
		}
		if i < len(input)-1 && this.validMonthFound && isNumeric(input[i]) {
			this.redact.numericLength++
			this.redact.length++
			switch {
			case this.redact.numericLength == 1 && !isNumeric(input[i+1]):
				this.validDayFound = true
				this.redact.numericLength = 0
				continue
			case this.redact.numericLength == 2 && !isNumeric(input[i+1]):
				if !this.validateDay(input[i-2 : i]) {
					this.resetCount(i)
					continue
				}
				this.validDayFound = true
				this.redact.numericLength = 0
				continue
			case this.redact.numericLength == 4 && this.validDayFound:
				this.redact.validateYear(input[i-4 : i])
				this.redact.appendMatch(this.redact.start, this.redact.length)
				this.resetCount(i)
				continue
			}
		}
		this.resetCount(i)
	}
	if this.validMonthFound && this.validDayFound && this.redact.validateYear(input[this.redact.length-4:this.redact.length]) {
		this.redact.appendMatch(this.redact.start, this.redact.length)
	}
}

func (this *FullDOB) validateDay(input []byte) bool {
	if len(input) < 2 || input[0] >= '3' && input[1] > '1' {
		return false
	}
	return true
}

func (this *FullDOB) isMonth(first, last byte, length int) bool {
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
func (this *FullDOB) resetCount(i int) {
	this.redact.start = i + 1
	this.redact.length = 0
	this.redact.breakLength = 0
	this.redact.numericLength = 0
	this.validMonthFound = false
	this.validDayFound = false
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
)

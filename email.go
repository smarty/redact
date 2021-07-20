package redact

func (this *emailRedaction) clear() {
	this.start = 0
	this.length = 0
}
func (this *emailRedaction) match(input []byte) {
	inputLength := len(input) - 1
	for i := 0; i < inputLength; i++ {
		if i < len(this.used) && this.used[i] {
			continue
		}
		if i < inputLength {
			this.checkMatch(input[i], i)
		}
		if this.length > MaxEmailLength {
			this.resetCount(i)
		}
	}
}

func (this *emailRedaction) resetCount(i int) {
	this.start = i + 1
	this.length = 0
}
func (this *emailRedaction) checkMatch(input byte, i int) {
	switch input {
	case '@':
		if this.start == 0 {
			this.length++
		}
		this.appendMatch(this.start, this.length-1)
		this.resetCount(i)
	case ' ':
		this.resetCount(i)
	}
	this.length++
}

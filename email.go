package redact

func (this *emailRedaction) clear() {
	this.start = 0
	this.length = 0
}
func (this *emailRedaction) match(input []byte) {
	var maxEmailLength = 254

	for i := 0; i < len(input); i++ {
		if i < len(this.used) && this.used[i] {
			continue
		}

		this.checkMatch(input, i)

		if this.length > maxEmailLength {
			this.resetCount(i)
		}
	}
}

func (this *emailRedaction) resetCount(i int) {
	this.start = i + 1
	this.length = 0
}
func (this *emailRedaction) checkMatch(input []byte, i int) {
	switch input[i] {
	case '@':
		this.appendMatch(this.start, this.length-1)
		this.resetCount(i)
	case ' ':
		this.resetCount(i)
	}
	this.length++
}

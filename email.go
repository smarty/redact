package redact

type emailRedaction struct {
	*matched
	start  int
	length int
}

func (this *emailRedaction) match(input []byte) {
	if len(input) <= 0 {
		return
	}
	for i := 0; i < len(input); i++ {
		character := input[i]
		if this.used[i] {
			continue
		}
		if !emailBreakNotFound(character) {
			this.start = i + 1
			this.length = 0
			continue
		} else {
			if character == '@' {
				this.appendMatch(this.start, this.length)
				this.start = i + 1
				this.length = 0
			}
			this.length++
		}
	}
}

func (this *emailRedaction) clear() {
	this.start = 0
	this.length = 0
}

func emailBreakNotFound(character byte) bool {
	return character != '.' && character != ' '
}

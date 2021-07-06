package redact

func (this *emailRedaction) match(input []byte) {
	var maxEmailLength = 254
	if len(input) <= 0 {
		return
	}
	for i := 0; i < len(input); i++ {
		character := input[i]
		if i > len(this.used)-1 {
			return
		}
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
		if this.length > maxEmailLength {
			this.length = 0
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

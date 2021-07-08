package redact

type phoneRedaction struct {
	*matched
	start       int
	length      int
}

func (this *phoneRedaction) clear() {
	this.start = 0
	this.length = 0
}

func (this *phoneRedaction) match(input []byte) {
	if len(input) <= 0 {
		return
	}
	for i := 0; i < len(input)-1; i++ {
		if i < len(this.used) && this.used[i] {
			continue
		}
		character := input[i]

		switch character{
		case '+':
			switch input[i+ 1] {
			case

			}
		}
	}}
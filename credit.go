package redact

func (this *creditCardRedaction) clear() {
	this.start = 0
	this.length = 0
	this.totalSum = 0
	this.isSecond = false
}

//Can be valid if there are letters at the very beginning or very end but NOT in the middle
//The only valid breaks for a card are: ' ' and '-'
//The breaks must be in the correct positions (i.e. 1010-1010-1010 1010, 1010-1010-0101-0, 0101-0101-0101-0101-010)
func (this *creditCardRedaction) match(input []byte) {
	for i := len(input) - 1; i > 0; i-- {
		if i < len(this.used)-1 && this.used[i] {
			continue
		}

		if isNumeric(input[i]) {
			this.length++
			this.start = i
			this.luhnCheck(input[i])
			continue
		}
		if i < len(input) - 1 && this.validateCard(input[i+1]){
			this.appendMatch(this.start, this.length)
			this.resetCount(i)
		}
	}
if this.validateCard(input[this.start]){
	this.appendMatch(this.start, this.length)
}
}

func (this *creditCardRedaction) luhnCheck(input byte) {
	var value uint64
	value = uint64(input - '0')

	if this.isSecond{
		value *= 2
	}
	this.totalSum += value / 10
	this.totalSum += value % 10
	this.isSecond = !this.isSecond
}

func (this *creditCardRedaction) validateCard(input byte) bool{
	if this.length > 19 || this.length < 12{
		return false
	}

	if !validateNetwork(input){
		return false
	}
	return this.totalSum%10 == 0
}
func validateNetwork(input byte) bool {
	return input >= '3' && input <= '6'
}

func (this *creditCardRedaction) resetCount(i int) {
	this.start = i - 1
	this.length = 0
	this.totalSum = 0
}
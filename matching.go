package redact

type match struct {
	InputIndex int
	Length     int
}
type matched struct {
	used    []bool
	matches []match
}

func (this *matched) clear() {
	this.matches = this.matches[0:0]
	for i := range this.used {
		this.used[i] = false
	}
}
func (this *matched) appendMatch(start, length int) {
	for i := start; i <= start+length; i++ {
		this.used[i] = true
	}
	this.matches = append(this.matches, match{InputIndex: start, Length: length})
}

package redact

type Redaction interface {
	match([]byte)
	clear()
}

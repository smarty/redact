package redact

type monitor interface {
	Redacted(count int)
}

type configuration struct {
	MaxLength  int
	BufferSize int
	Monitor    monitor
}

func New(options ...option) *Redactor {
	var config configuration
	Options.apply(options...)(&config)

	matched := &matched{
		used:    make([]bool, config.MaxLength),
		matches: make([]match, 0, config.BufferSize),
	}

	return &Redactor{
		matched: matched,
		phone:   &phoneRedaction{matched: matched},
		ssn:     &ssnRedaction{matched: matched},
		//credit:  &creditCardRedaction{matched: matched},
		credit:  &creditCardRedaction{matched: matched},
		dob:     &dobRedaction{matched: matched},
		email:   &emailRedaction{matched: matched},
		monitor: config.Monitor,
	}
}

var Options singleton

type singleton struct{}
type option func(*configuration)

func (singleton) MaxLength(value int) option {
	return func(this *configuration) { this.MaxLength = value }
}
func (singleton) BufferSize(value int) option {
	return func(this *configuration) { this.BufferSize = value }
}
func (singleton) Monitor(value monitor) option {
	return func(this *configuration) { this.Monitor = value }
}

func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		for _, option := range Options.defaults(options...) {
			option(this)
		}
	}
}
func (singleton) defaults(options ...option) []option {
	return append([]option{
		Options.MaxLength(512),
		Options.BufferSize(16),
		Options.Monitor(nop{}),
	}, options...)
}

type nop struct{}

func (nop) Redacted(int) {}

package kitty

import (
	"errors"
	"github.com/Kamva/gutil"
)

var (
	errNilConfig     = errors.New("config is Nil in the pack")
	errNilLogger     = errors.New("logger is Nil in the pack")
	errNilTranslator = errors.New("translator is Nil in the pack")
	errNilJobs       = errors.New("jobs is nil in the pack")
)

// Pack contains all of services in one place to manage our services.
type Pack struct {
	// must specify that should panic if one service is
	//nil and user request to get that service or just
	//returns nil.
	must bool

	config     Config
	log        Logger
	translator Translator
	jobs       Jobs
}

// SetConfig sets the config service.
func (p *Pack) SetConfig(config Config) {
	p.config = config
}

// SetLogger sets the logger service.
func (p *Pack) SetLogger(logger Logger) {
	p.log = logger
}

// SetTranslator sets the translator service.
func (p *Pack) SetTranslator(translator Translator) {
	p.translator = translator
}

// SetJobs sets the Jobs service.
func (p *Pack) SetJobs(jobs Jobs) {
	p.jobs = jobs
}

// Config returns the config service.
func (p *Pack) Config() Config {
	gutil.PanicNil(p.config, errNilConfig)

	return p.config
}

// Logger returns the logger service.
func (p *Pack) Log() Logger {
	gutil.PanicNil(p.log, errNilLogger)

	return p.log
}

// Translator returns the translator service.
func (p *Pack) Translator() Translator {
	gutil.PanicNil(p.translator, errNilTranslator)

	return p.translator
}

// Translator returns the translator service.
func (p *Pack) Jobs() Jobs {
	gutil.PanicNil(p.jobs, errNilJobs)

	return p.jobs
}

// NewPack returns new instance of the pack.
func NewPack(must bool) *Pack {
	return &Pack{
		must: must,
	}
}

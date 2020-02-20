package kitty

import (
	"errors"
	"github.com/Kamva/gutil"
)

var (
	errNilConfig     = errors.New("config value is Nil")
	errNilLogger     = errors.New("config value is Nil")
	errNilTranslator = errors.New("config value is Nil")
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

// Config returns the config service.
func (p *Pack) Config() Config {
	gutil.PanicNil(p.config, errNilConfig)

	return p.config
}

// Log returns the logger service.
func (p *Pack) Log() Logger {
	gutil.PanicNil(p.log, errNilLogger)

	return p.log
}

// Translator returns the translator service.
func (p *Pack) Translator() Translator {
	gutil.PanicNil(p.translator, errNilTranslator)

	return p.translator
}

// NewPack returns new instance of the pack.
func NewPack(must bool) *Pack {
	return &Pack{
		must: must,
	}
}

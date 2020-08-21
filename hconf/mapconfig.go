package hconf

import (
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
)

type mapConfig struct {
	conf hexa.Map
}

func (m *mapConfig) Unmarshal(instance interface{}) error {
	return gutil.MapToStruct(m.conf, instance)
}

// NewViperDriver returns new instance of viper driver.
func NewMapDriver(conf hexa.Map) hexa.Config {
	return &mapConfig{conf: conf}
}

// Assert viperConfig is type of hexa Config
var _ hexa.Config = &mapConfig{}

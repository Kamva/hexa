package hconf

import (
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"github.com/spf13/viper"
)

type viperConfig struct {
	v *viper.Viper
}

func (v *viperConfig) Unmarshal(instance interface{}) error {
	return tracer.Trace(v.v.Unmarshal(instance))
}

// NewViperDriver returns new instance of viper driver.
func NewViperDriver(viper *viper.Viper) hexa.Config {
	return &viperConfig{v: viper}
}

// Assert viperConfig is type of hexa Config
var _ hexa.Config = &viperConfig{}

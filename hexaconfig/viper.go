package hexaconfig

import (
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"github.com/spf13/viper"
)

type viperConfig struct {
	*viper.Viper
}

func (v *viperConfig) Unmarshal(instance interface{}) error {
	return tracer.Trace(v.Viper.Unmarshal(instance))
}

func (v *viperConfig) GetList(key string) []string {
	return v.GetStringSlice(key)
}

// NewViperDriver returns new instance of viper driver.
func NewViperDriver(viper *viper.Viper) hexa.Config {
	return &viperConfig{Viper: viper}
}

// Assert viperConfig is type of hexa Config
var _ hexa.Config = &viperConfig{}

package hexaconfig

import (
	"errors"
	"github.com/Kamva/hexa"
	"github.com/Kamva/tracer"
)

type mapConfig struct {
	conf hexa.Map
}

func (m *mapConfig) Get(key string) interface{} {
	return m.conf[key]
}

func (m *mapConfig) GetString(key string) string {
	return m.conf[key].(string)
}

func (m *mapConfig) GetInt64(key string) int64 {
	return m.conf[key].(int64)
}

func (m *mapConfig) GetFloat64(key string) float64 {
	return m.conf[key].(float64)
}

func (m *mapConfig) GetBool(key string) bool {
	return m.conf[key].(bool)
}

func (m *mapConfig) Unmarshal(instance interface{}) error {
	conf, ok := instance.(hexa.Map)
	if !ok {
		return tracer.Trace(errors.New("provided config for map config must be a map[string]interface"))
	}
	m.conf = conf
	return nil
}

func (m *mapConfig) GetList(key string) []string {
	return m.conf[key].([]string)
}

// NewViperDriver returns new instance of viper driver.
func NewMapDriver() hexa.Config {
	return &mapConfig{}
}

// Assert viperConfig is type of hexa Config
var _ hexa.Config = &mapConfig{}

package hexa

// Config is viperDriver interface that should embed in each Config object.
type Config interface {
	// Unmarshal function load viperDriver into viperDriver instance. pass instance by reference.
	Unmarshal(configInstance interface{}) error

	// Get method return config value in any type.
	Get(key string) interface{}

	// GetString method return config value as string
	GetString(key string) string

	// GetInt64 return int64 config value.
	GetInt64(key string) int64

	// GetFloat64 return float64 config value.
	GetFloat64(key string) float64

	// GetBool return bool config value.
	GetBool(key string) bool

	// GetList return config value as string list.
	GetList(key string) []string
}

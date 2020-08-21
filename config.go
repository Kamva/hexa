package hexa

type Config interface {
	// Unmarshal unmarshal config values to the provided struct.
	Unmarshal(instance interface{}) error
}

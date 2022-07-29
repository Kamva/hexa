package hexa

type Config interface {
	// Unmarshal unmarshal config values to the provided struct.
	Unmarshal(instance any) error
}

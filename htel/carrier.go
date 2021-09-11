package htel


type HexaCarrier map[string][]byte

func (hc HexaCarrier) Get(key string) string {
	return string(hc[key])
}

func (hc HexaCarrier) Set(key string, value string) {
	hc[key] = []byte(value)
}

func (hc HexaCarrier) Keys() []string {
	keys := make([]string, len(hc))

	i := 0
	for k, _ := range hc {
		keys[i] = k
		i++
	}

	return keys
}


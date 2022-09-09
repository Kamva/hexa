package hexa

func extendBytesMap(dest, src map[string][]byte, overwrite bool) {
	for key, val := range src {
		// If key exists in dest and we can not overwrite it, so continue.
		if _, ok := dest[key]; ok && !overwrite {
			continue
		}
		dest[key] = val
	}
}

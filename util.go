package hexa

import (
	"context"
)

// getMissedKeyInContext returns missed key in the context.
func getMissedKeyInContext(c context.Context, keys ...string) string {
	for _, k := range keys {
		if val := c.Value(k); val == nil {
			return k
		}
	}
	return ""
}

func extendBytesMap(dest, src map[string][]byte, overwrite bool) {
	for key, val := range src {
		// If key exists in dest and we can not overwrite it, so continue.
		if _, ok := dest[key]; ok && !overwrite {
			continue
		}
		dest[key] = val
	}
}

package common

// GetKeyForValue returns key for the given value
func GetKeyForValue(data map[string]string, val string) string {
	for k, v := range data {
		if v == val {
			return k
		}
	}
	return ""
}

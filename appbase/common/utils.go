package common

import (
	"encoding/json"
	"github.com/appbaseio/abc/log"
	"strings"
)

// GetKeyForValue returns key for the given value
func GetKeyForValue(data map[string]string, val string) string {
	for k, v := range data {
		if v == val {
			return k
		}
	}
	return ""
}

// JSONNumberToString converts a json.Number to string, properly
// i.e. no decimal points for a integer
// json.Number is required instead of normal types in map[..].. based decoding
func JSONNumberToString(number json.Number) string {
	str := number.String()
	if strings.HasSuffix(str, ".0") {
		return str[0 : len(str)-2]
	}
	return str
}

// JSONNumberToInt ...
func JSONNumberToInt(number json.Number) int64 {
	f, err := number.Float64()
	if err != nil {
		log.Errorln(err)
		return 0
	}
	return int64(f)
}

// StringInSlice checks if string is in list or not
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ColonPad pads spaces after colon
func ColonPad(text string, length int) string {
	textLen := len(text)
	text += ":"
	for i := 0; i < (length - textLen - 1); i++ {
		text += " "
	}
	return text
}

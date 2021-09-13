package mysql

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
)

func castType(colType string, value sql.RawBytes) interface{} {
	// https://dev.mysql.com/doc/refman/5.7/en/data-types.html
	colType = strings.ToLower(colType)
	val := string(value)

	if strings.Contains(colType, "date") || strings.Contains(colType, "time") {
		// TODO: parse to datetime object
		return val
	} else if strings.Contains(colType, "int") {
		i, _ := strconv.Atoi(val)
		return i
	} else if stringInSlice(colType, []string{"decimal", "numeric", "float", "double"}) {
		f, _ := strconv.ParseFloat(val, 64)
		return f
	} else if strings.Contains(colType, "json") {
		var m map[string]interface{}
		json.Unmarshal([]byte(value), &m)
		return m
	}
	return val
}

// stringInSlice checks if string is in list or not
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

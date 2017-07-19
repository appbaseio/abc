package app

import (
	"strings"
)

// SortOptions holds the list of valid sort options
var SortOptions = []string{"id", "name", "api-calls", "records", "storage"}

type fullApp struct {
	id   string
	name string
	appBody
}

// appsSorter helps to sort apps
// https://golang.org/pkg/sort/
type appsSorter struct {
	key  string
	apps []fullApp
}

func (a appsSorter) Len() int {
	return len(a.apps)
}

func (a appsSorter) Swap(i, j int) {
	a.apps[i], a.apps[j] = a.apps[j], a.apps[i]
}

func (a appsSorter) Less(i, j int) bool {
	switch a.key {
	case "id":
		return a.apps[i].id < a.apps[j].id
	case "name":
		return strings.ToLower(a.apps[i].name) < strings.ToLower(a.apps[j].name)
	case "api-calls":
		return a.apps[i].APICalls > a.apps[j].APICalls
	case "records":
		return a.apps[i].Records > a.apps[j].Records
	case "storage":
		return a.apps[i].Storage > a.apps[j].Storage
	default:
		return a.apps[i].id < a.apps[j].id
	}
}

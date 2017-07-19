package app

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
	return a.apps[i].id < a.apps[j].id
}

package spinner

import (
	"github.com/briandowns/spinner"
	"time"
)

var s *spinner.Spinner
var active bool

// Start starts the spinner
func Start() {
	startSpinner(9)
}

// StartText starts a text spinner
func StartText(text string) {
	startSpinner(11) // ⣾⣽⣻⢿⡿⣟⣯⣷
	s.Suffix = "  " + text
}

// Stop stops the spinner
func Stop() {
	s.Stop()
	active = false
}

// startSpinner ...
func startSpinner(sType int) {
	if active {
		Stop()
	}
	s = spinner.New(spinner.CharSets[sType], 100*time.Millisecond)
	s.Start()
	active = true
}

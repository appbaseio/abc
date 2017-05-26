package function

import "github.com/aviaryan/abc/message"

// Function has a single defined function to serve the purpose of apply logic to a message in order to return
// a message.
type Function interface {
	Apply(message.Msg) (message.Msg, error)
}

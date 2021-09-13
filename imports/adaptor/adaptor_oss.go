// +build oss

package adaptor

import (
	// ...
	"github.com/appbaseio/abc/importer/adaptor"
)

// Adaptors export
var Adaptors = adaptor.Adaptors

// GetAdaptor export
var GetAdaptor = adaptor.GetAdaptor

// Adaptor export
type Adaptor interface {
	adaptor.Adaptor
}

// Describable ...
type Describable adaptor.Describable

// RegisteredAdaptors ...
func RegisteredAdaptors() []string {
	return adaptor.RegisteredAdaptors()
}

// Error ...
// type Error adaptor.Error

// // ERROR ...
// const ERROR = adaptor.ERROR

// // CRITICAL ...
// const CRITICAL = adaptor.CRITICAL

// // Error returns the error as a string
// func (t Error) Error() string {
// 	return fmt.Sprintf("%s: %s", t.Lvl, t.Err)
// }

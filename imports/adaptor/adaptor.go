// +build !oss

package adaptor

import (
	// import
	adaptorx "github.com/appbaseio/abc/importer/adaptor"
)

// Adaptors export
var Adaptors = adaptorx.Adaptors

// GetAdaptor export
var GetAdaptor = adaptorx.GetAdaptor

// Adaptor export
type Adaptor interface {
	adaptorx.Adaptor
}

// Describable ...
type Describable adaptorx.Describable

// RegisteredAdaptors ...
func RegisteredAdaptors() []string {
	return adaptorx.RegisteredAdaptors()
}

// Error ...
// type Error adaptorx.Error

// // ERROR ...
// const ERROR = adaptorx.ERROR

// // CRITICAL ...
// const CRITICAL = adaptorx.CRITICAL

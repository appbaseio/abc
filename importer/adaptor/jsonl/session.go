package jsonl

import (
	"os"

	"github.com/appbaseio/abc/importer/client"
)

// Session serves as a wrapper for the underlying file
type Session struct {
	file *os.File
}

var _ client.Session = &Session{}

package json

import (
	"encoding/json"
	"os"

	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/importer/message"
)

var _ client.Writer = &Writer{}

// Writer implements client.Writer for use with Files
type Writer struct{}

func newWriter() *Writer {
	w := &Writer{}
	return w
}

func (w *Writer) Write(msg message.Msg) func(client.Session) (message.Msg, error) {
	return func(s client.Session) (message.Msg, error) {
		if err := dumpMessage(msg, s.(*Session).file); err != nil {
			return nil, err
		}
		if msg.Confirms() != nil {
			close(msg.Confirms())
		}
		return msg, nil
	}
}

func dumpMessage(msg message.Msg, f *os.File) error {
	return json.NewEncoder(f).Encode(msg.Data())
}

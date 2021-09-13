package json

import (
	"bytes"
	"encoding/json"

	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/importer/message"
	"github.com/appbaseio/abc/log"

	// "github.com/appbaseio/abc/importer/message/data"
	"github.com/appbaseio/abc/importer/message/ops"
)

var _ client.Reader = &Reader{}

// Reader implements the behavior defined by client.Reader for interfacing with the file.
type Reader struct {
	fileName string
	typeName string
}

func newReader(typeName string) client.Reader {
	return &Reader{"", typeName}
}

func (r *Reader) Read(_ map[string]client.MessageSet, filterFn client.NsFilterFunc) client.MessageChanFunc {
	return func(s client.Session, done chan struct{}) (chan client.MessageSet, error) {
		out := make(chan client.MessageSet)
		session := s.(*Session)
		r.fileName = session.file.Name()

		go func() {
			defer close(out)
			log.With("file", r.fileName).Infoln("starting Read func")
			// iterate rows
			iterationComplete := r.decodeFile(session, out, done)
			func() {
				for {
					select {
					case _, ok := <-iterationComplete:
						if !ok {
							return
						}
					case <-done:
						return
					}
				}
			}()
			// end
			log.With("file", r.fileName).Infoln("Read completed")
			return
		}()
		return out, nil
	}
}

func (r *Reader) decodeFile(s *Session, out chan<- client.MessageSet, done chan struct{}) <-chan string {
	fileDone := make(chan string)
	go func() {
		defer close(fileDone)
		buf := new(bytes.Buffer)

		buf.ReadFrom(s.file)
		streamBytes := buf.Bytes()
		jsonData := bytes.TrimLeft(streamBytes, " \t\r\n")
		if len(jsonData) <= 0 {
			return
		}
		isArray := jsonData[0] == '['
		isObject := jsonData[0] == '{'

		if isArray {
			file := make([]interface{}, 0)
			// read file
			err := json.Unmarshal(jsonData, &file)
			if err != nil {
				return
			}

			for _, row := range file {
				out <- client.MessageSet{
					Msg: message.From(ops.Insert, r.typeName, row.(map[string]interface{})),
				}
				// return early?
				select {
				default:
				case <-done:
					log.With("file", r.fileName).Infoln("Reading file stopped midway")
					return
				}
			}
		}

		if isObject {
			var file map[string]interface{}
			// read file
			err := json.Unmarshal(jsonData, &file)
			if err != nil {
				return
			}

			out <- client.MessageSet{
				Msg: message.From(ops.Insert, r.typeName, file),
			}
		}
	}()
	return fileDone
}

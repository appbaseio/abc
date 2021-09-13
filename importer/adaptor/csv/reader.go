package csv

import (
	"encoding/csv"
	"io"

	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/importer/message"
	"github.com/appbaseio/abc/importer/message/ops"
	"github.com/appbaseio/abc/log"
)

var _ client.Reader = &Reader{}

// Reader fulfills the client.Reader interface
type Reader struct {
	fName    string
	typeName string
}

func newReader(typeName string) client.Reader {
	return &Reader{"", typeName}
}

func (r *Reader) Read(_ map[string]client.MessageSet, filterFn client.NsFilterFunc) client.MessageChanFunc {
	return func(s client.Session, done chan struct{}) (chan client.MessageSet, error) {
		out := make(chan client.MessageSet)
		r.fName = s.(*Session).fName // set file name
		// ^^ important for good logging
		reader := s.(*Session).reader
		log.Infof("connection = %v", reader)

		go func() {
			defer close(out)
			log.With("file", r.fName).Infoln("starting Read func")
			// iterate rows
			iterationComplete := r.iterateRows(reader, out, done)
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
			log.With("file", r.fName).With("type", r.typeName).Infoln("Read completed")
			return
		}()
		return out, nil
	}
}

// iterateRows takes care of a table
func (r *Reader) iterateRows(reader *csv.Reader, out chan<- client.MessageSet, done chan struct{}) <-chan string {
	fileDone := make(chan string)
	go func() {
		defer close(fileDone)

		log.With("file", r.fName).Infoln("iterating...")

		ct := 0
		var columns []string

		for {
			record, err := reader.Read()
			// Stop at EOF.
			if err == io.EOF {
				break
			}
			ct = ct + 1

			if ct == 1 {
				columns = make([]string, len(record))
				columns = record
				continue
			}

			data := make(map[string]interface{})
			colCount := len(record)
			for i := 0; i < colCount; i++ {
				data[columns[i]] = record[i]
			}

			out <- client.MessageSet{
				Msg: message.From(ops.Insert, r.typeName, data),
			}

			// return early?
			select {
			default:
			case <-done:
				log.With("file", r.fName).Infoln("Reading file stopped midway")
				return
			}
		}

		fileDone <- r.typeName
	}()
	return fileDone
}

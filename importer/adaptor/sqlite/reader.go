package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/importer/message"
	"github.com/appbaseio/abc/importer/message/ops"
	"github.com/appbaseio/abc/log"
)

var _ client.Reader = &Reader{}

// Reader fulfills the client.Reader interface for use with both copying and tailing a MSSQL database.
type Reader struct {
	tail   bool
	dbName string
}

func newReader(tail bool) client.Reader {
	return &Reader{tail, ""}
}

func (r *Reader) Read(_ map[string]client.MessageSet, filterFn client.NsFilterFunc) client.MessageChanFunc {
	return func(s client.Session, done chan struct{}) (chan client.MessageSet, error) {
		out := make(chan client.MessageSet)
		r.dbName = s.(*Session).dbName // set database name
		// ^^ important for good logging
		db := s.(*Session).db
		log.Infof("connection = %v", db)

		go func() {
			defer close(out)
			log.With("db", r.dbName).Infoln("starting Read func")
			// get tables
			tables, err := r.listTables(db, filterFn)
			if err != nil {
				log.With("db", r.dbName).Errorf("unable to list tables, %s", err)
				return
			}
			// iterate tables
			iterationComplete := r.iterateTable(db, tables, out, done)
			func() {
				for {
					select {
					case t, ok := <-iterationComplete:
						if !ok {
							return
						}
						log.With("db", r.dbName).Infof("Table %s done", t)
					case <-done:
						return
					}
				}
			}()
			// end
			log.With("db", r.dbName).Infoln("Read completed")
			return
		}()
		return out, nil
	}
}

// listTables list the tables
func (r *Reader) listTables(db *sql.DB, filterFn func(name string) bool) (<-chan string, error) {
	out := make(chan string)
	// get all tables
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(out)
		defer rows.Close()
		var table string
		for rows.Next() {
			err := rows.Scan(&table)
			if err != nil {
				log.With("db", r.dbName).Errorln(err)
			}

			if filterFn(table) {
				log.Infof("table %s\n", table)
				log.With("db", r.dbName).With("table", table).Infoln("sending for iteration...")
				out <- table
			} else {
				log.With("db", r.dbName).With("table", table).Infoln("skipping iteration...")
			}
		}
		log.With("db", r.dbName).Infoln("done iterating tables")
	}()
	return out, nil
}

// iterateTable takes care of a table
func (r *Reader) iterateTable(db *sql.DB, in <-chan string, out chan<- client.MessageSet, done chan struct{}) <-chan string {
	tableDone := make(chan string)
	go func() {
		defer close(tableDone)
		for {
			select {
			case t, ok := <-in:
				if !ok {
					return
				}
				log.With("db", r.dbName).With("table", t).Infoln("iterating...")
				// read table
				rows, err := db.Query("select * from " + t)
				if err != nil {
					log.With("db", r.dbName).With("table", t).Errorf("Error reading rows %s", err)
					return
				}
				// get columns
				columns, err := rows.Columns()
				if err != nil {
					log.With("db", r.dbName).With("table", t).Errorf("Error reading columns %s", err)
					return
				}
				colCount := len(columns)
				// get data
				for rows.Next() {
					results := make([]interface{}, colCount)
					for i := 0; i < colCount; i++ {
						var temp interface{}
						results[i] = &temp
					}
					err = rows.Scan(results...)
					if err != nil {
						log.With("db", r.dbName).With("table", t).Errorf("Error reading row %s", err)
						return
					}
					data := make(map[string]interface{})
					for i := 0; i < colCount; i++ {
						res := fmt.Sprint(*(results[i].(*interface{})))

						var x interface{} = *(results[i].(*interface{}))
						if value, ok := x.([]uint8); ok {
							res = B2S(value)
						}
						data[columns[i]] = res
					}
					// send data
					// log.Infoln(data)
					out <- client.MessageSet{
						Msg: message.From(ops.Insert, t, data),
					}
					// return early?
					select {
					default:
					case <-done:
						log.With("db", r.dbName).Infoln("Reading table stopped midway")
						return
					}
				}

				tableDone <- t
			case <-done:
				log.With("db", r.dbName).Infoln("iterating no more")
				return
			}
		}
	}()
	return tableDone
}

package mysql

import (
	"database/sql"

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
	rows, err := db.Query("show tables")
	tableCol, err := rows.Columns()
	if err != nil {
		log.Errorf("Error reading columns %s", err)
	}
	tableCount := len(tableCol)
	out := make(chan string, tableCount)
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

				// get column types
				colTypes, err := r.getColumnTypes(db, t, colCount)
				if err != nil {
					log.With("db", r.dbName).With("table", t).Errorf("Error reading types %s", err)
				}

				// get row
				values := make([]sql.RawBytes, colCount)
				results := make([]interface{}, colCount)
				for i := 0; i < colCount; i++ {
					results[i] = &values[i]
				}
				// get data
				for rows.Next() {
					err = rows.Scan(results...)
					if err != nil {
						log.With("db", r.dbName).With("table", t).Errorf("Error reading row %s", err)
						return
					}
					data := make(map[string]interface{})
					for i := 0; i < colCount; i++ {
						data[columns[i]] = castType(colTypes[i], values[i])
						// log.Infoln(data[columns[i]])
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

// getColumnTypes gets types of columns
func (r *Reader) getColumnTypes(db *sql.DB, table string, count int) ([]string, error) {
	log.With("db", r.dbName).With("table", table).Infoln("getting types...")
	types := make([]string, count)

	// read table
	rows, err := db.Query("show columns from " + table)
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		log.With("db", r.dbName).With("table", table).Errorf("Error reading columns %s", err)
		return nil, err
	}
	colCount := len(columns)

	values := make([]sql.RawBytes, colCount)
	results := make([]interface{}, colCount)
	for i := 0; i < colCount; i++ {
		results[i] = &values[i]
	}
	ct := 0

	for rows.Next() {
		err = rows.Scan(results...)
		if err != nil {
			log.Errorln(err)
			return nil, err
		}
		types[ct] = string(values[1])
		ct = ct + 1
	}
	// log.Infoln(types, len(types), count)

	return types, nil
}

// +build !oss

package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func verifyConnections(srcType string, srcURL string, ssl bool) error {
	switch srcType {
	case "postgres":
		if !ssl {
			srcURL = srcURL + "?sslmode=disable"
		}

		conn, err := sql.Open("postgres", srcURL)
		if err != nil {
			return err
		}

		err = conn.Ping()
		if err != nil {
			return err
		}

		conn.Close()

	default:
		return nil
	}
	return nil
}

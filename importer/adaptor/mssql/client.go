package mssql

import (
	"database/sql"
	"net/url"
	"time"

	"github.com/appbaseio/abc/importer/client"

	// "github.com/appbaseio/abc/log"

	_ "github.com/denisenkom/go-mssqldb"
)

const (
	// DefaultURI is the default endpoint for MSSQL on the local machine.
	// Primarily used when initializing a new Client without a specific URI.
	DefaultURI = "sqlserver://sa:mypass@localhost:1234?database=master"

	// DefaultTimeout is the default time.Duration used if one is not provided for options
	// that pertain to timeouts.
	DefaultTimeout = 30 * time.Second
)

var (
	_ client.Client = &Client{}
	_ client.Closer = &Client{}
)

// Client creates and holds the session to RethinkDB
type Client struct {
	uri            string
	db             *sql.DB
	sessionTimeout time.Duration
}

// Session contains an instance of the rethink.Session for use by Readers/Writers
type Session struct {
	db     *sql.DB
	dbName string
}

// ClientOptionFunc It is used in NewClient.
type ClientOptionFunc func(*Client) error

// NewClient creates a new client
func NewClient(options ...ClientOptionFunc) (*Client, error) {
	// Set up the client
	c := &Client{
		uri:            DefaultURI,
		sessionTimeout: DefaultTimeout,
	}

	// Run the options on it
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithURI defines the full connection string of the RethinkDB database.
func WithURI(uri string) ClientOptionFunc {
	return func(c *Client) error {
		_, err := url.Parse(c.uri)
		if err != nil {
			return client.InvalidURIError{URI: uri, Err: err.Error()}
		}
		c.uri = uri
		return nil
	}
}

// Connect wraps the underlying session to the RethinkDB database
func (c *Client) Connect() (client.Session, error) {
	if c.db == nil {
		if err := c.initConnection(); err != nil {
			return nil, err
		}
	}
	// get database name
	pURI, _ := url.Parse(c.uri)
	dbName := pURI.Query().Get("database")
	// create session
	return &Session{c.db, dbName}, nil
}

// Close fulfills the Closer interface and takes care of cleaning up the rethink.Session
func (c *Client) Close() {
	// check for err
	c.db.Close()
}

func (c *Client) initConnection() error {
	db, err := sql.Open("mssql", c.uri)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	c.db = db
	return nil
}

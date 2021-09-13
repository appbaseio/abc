package csv

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/log"
)

const (
	// DefaultPath is the default path for csv file
	// Used as a placeholder
	DefaultPath = "stdout://"
)

var (
	_ client.Client = &Client{}
	_ client.Closer = &Client{}
)

// Client creates and holds the session to RethinkDB
type Client struct {
	uri string
	// db  *sql.DB
	file     *os.File
	reader   *csv.Reader
	typeName string
	// sessionTimeout time.Duration
	deleteFileAfterUsage bool
}

// Session contains an instance of the rethink.Session for use by Readers/Writers
type Session struct {
	reader *csv.Reader
	fName  string
}

// ClientOptionFunc It is used in NewClient.
type ClientOptionFunc func(*Client) error

// NewClient creates a new client
func NewClient(options ...ClientOptionFunc) (*Client, error) {
	// Set up the client
	c := &Client{
		uri: DefaultPath,
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
		c.uri = uri
		return nil
	}
}

// WithType adds type
func WithType(typeName string) ClientOptionFunc {
	return func(c *Client) error {
		if typeName == "" {
			return errors.New("Type name can't be empty")
		}
		c.typeName = typeName
		return nil
	}
}

// Connect wraps the underlying session to the RethinkDB database
func (c *Client) Connect() (client.Session, error) {
	if c.file == nil {
		if err := c.initConnection(); err != nil {
			return nil, err
		}
	}
	// get file name
	index := strings.LastIndex(c.uri, "/")
	fName := c.uri[index+1:]
	// create session
	return &Session{c.reader, fName}, nil
}

// Close fulfills the Closer interface and takes care of cleaning up
func (c *Client) Close() {
	if c.file != nil && c.file != os.Stdout {
		c.file.Close()
		if c.deleteFileAfterUsage {
			common.RemoveFile(c.file.Name())
		}
	}
}

func (c *Client) initConnection() error {
	if strings.HasPrefix(c.uri, "stdout://") {
		c.file = os.Stdout
		return nil
	}

	var name string
	name = strings.Replace(c.uri, "file://", "", 1)

	// read file from remote url
	if url, err := url.ParseRequestURI(c.uri); err == nil {
		name = fmt.Sprintf("%s/%s.csv", common.DefaultDownloadDirectory, fmt.Sprintf("%d", time.Now().Unix()))
		if err := common.DownloadFile(name, url.String()); err != nil {
			return err
		}
		log.Infoln("Download complete for:", name)
		c.deleteFileAfterUsage = true
	}
	f, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		if c.deleteFileAfterUsage {
			common.RemoveFile(name)
		}
		return err
	}
	c.file = f
	c.reader = csv.NewReader(bufio.NewReader(f))
	return nil
}

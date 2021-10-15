package jsonl

import (
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

var (
	_ client.Client = &Client{}
	_ client.Closer = &Client{}
)

// ClientOptionFunc is a function that configures a Client.
type ClientOptionFunc func(*Client) error

// Client represents a client to the underlying File source.
type Client struct {
	uri                  string
	typeName             string
	file                 *os.File
	deleteFileAfterUsage bool
}

// DefaultURI is the default file, outputs to stdout
var (
	DefaultURI = "stdout://"
)

// NewClient creates a default file client
func NewClient(options ...ClientOptionFunc) (*Client, error) {
	// Set up the client
	c := &Client{
		uri: DefaultURI,
	}

	// Run the options on it
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithURI defines the full path to the file, prefixed with file://, or stdout://
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

// Connect initializes the file for IO
func (c *Client) Connect() (client.Session, error) {
	if c.file == nil {
		if err := c.initFile(); err != nil {
			return nil, err
		}
	}
	return &Session{c.file}, nil
}

// Close closes the underlying file
func (c *Client) Close() {
	if c.file != nil && c.file != os.Stdout {
		c.file.Close()
		if c.deleteFileAfterUsage {
			common.RemoveFile(c.file.Name())
		}
	}
}

func (c *Client) initFile() error {
	if strings.HasPrefix(c.uri, "stdout://") {
		c.file = os.Stdout
		return nil
	}
	name := strings.Replace(c.uri, "file://", "", 1)
	proto := strings.Split(c.uri, "//")
	if strings.Contains(proto[0], `http`) {
		if url, err := url.ParseRequestURI(c.uri); err == nil {
			name = fmt.Sprintf("%s/%s.csv", common.DefaultDownloadDirectory, fmt.Sprintf("%d", time.Now().Unix()))
			if err := common.DownloadFile(name, url.String()); err != nil {
				return err
			}
			log.Infoln("Download complete for:", name)
			c.deleteFileAfterUsage = true
		}
	}

	f, err := os.OpenFile(name, os.O_RDWR, 0666)
	if os.IsNotExist(err) {
		f, err = os.Create(name)
		if err != nil {
			if c.deleteFileAfterUsage {
				common.RemoveFile(name)
			}
			return err
		}
	}
	if err != nil {
		if c.deleteFileAfterUsage {
			common.RemoveFile(name)
		}
		return err
	}
	c.file = f
	return nil
}

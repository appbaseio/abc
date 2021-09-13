package redis

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/appbaseio/abc/importer/client"

	"github.com/go-redis/redis"
)

const (
	// DefaultURI is the default endpoint of Postgres on the local machine.
	// Primarily used when initializing a new Client without a specific URI.
	DefaultURI = "redis://127.0.0.1:6379/0"

	// DefaultDB is database used by Redis
	DefaultDB = 0

	// DefaultDialTimeout is the timeout of a Redis connection
	DefaultDialTimeout = 5 * time.Second

	// DefaultReadTimeout is the read timeout of the Redis database
	DefaultReadTimeout = 3 * time.Second

	// DefaultWriteTimeout is the write timeout of the Redis database
	DefaultWriteTimeout = 3 * time.Second
)

var (
	_ client.Client = &Client{}
	_ client.Closer = &Client{}
)

// ClientOptionFunc is a function that configures a Client.
// It is used in NewClient.
type ClientOptionFunc func(*Client) error

// Client represents a client to the underlying Redis source.
type Client struct {
	uri string

	db int

	dialTimeout, readTimeout, writeTimeout time.Duration

	tlsConfig *tls.Config

	conn *redis.Client
}

// Session represents a connection to a Redis client.
type Session struct {
	conn *redis.Client
}

// NewClient creates a new client to work with Redis.
//
// The caller can configure the new client by passing configuration options
// to the func.
//
// Example:
//
//   client, err := NewClient(
//     WithURI("redis://127.0.0.1:6379/0"),
//     WithTimeout("30s"))
//
// If no URI is configured, it uses defaultURI by default.
//
// An error is also returned when some configuration option is invalid
func NewClient(options ...ClientOptionFunc) (*Client, error) {
	// Set up the client
	c := &Client{
		uri:          DefaultURI,
		db:           DefaultDB,
		dialTimeout:  DefaultDialTimeout,
		readTimeout:  DefaultReadTimeout,
		writeTimeout: DefaultWriteTimeout,
		tlsConfig:    nil,
	}

	// Run the options on it
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithURI defines the full connection string of the Redis database.
func WithURI(uri string) ClientOptionFunc {
	return func(c *Client) error {
		if _, err := redis.ParseURL(uri); err != nil {
			return client.InvalidURIError{URI: uri, Err: err.Error()}
		}
		c.uri = uri
		return nil
	}
}

// WithDialTimeout overrides the DefaultDialTimeout and should be parseable by time.ParseDuration
func WithDialTimeout(timeout string) ClientOptionFunc {
	return func(c *Client) error {
		if timeout == "" {
			c.dialTimeout = DefaultDialTimeout
			return nil
		}

		t, err := time.ParseDuration(timeout)
		if err != nil {
			return client.InvalidTimeoutError{Timeout: timeout}
		}
		c.dialTimeout = t
		return nil
	}
}

// WithWriteTimeout overrides the DefaultWriteTimeout and should be parseable by time.ParseDuration
func WithWriteTimeout(timeout string) ClientOptionFunc {
	return func(c *Client) error {
		if timeout == "" {
			c.writeTimeout = DefaultWriteTimeout
			return nil
		}

		t, err := time.ParseDuration(timeout)
		if err != nil {
			return client.InvalidTimeoutError{Timeout: timeout}
		}
		c.writeTimeout = t
		return nil
	}
}

// WithReadTimeout overrides the DefaultReadTimeout and should be parseable by time.ParseDuration
func WithReadTimeout(timeout string) ClientOptionFunc {
	return func(c *Client) error {
		if timeout == "" {
			c.readTimeout = DefaultReadTimeout
			return nil
		}

		t, err := time.ParseDuration(timeout)
		if err != nil {
			return client.InvalidTimeoutError{Timeout: timeout}
		}
		c.readTimeout = t
		return nil
	}
}

// WithCACerts configures the RootCAs for the underlying TLS connection
func WithCACerts(certs []string) ClientOptionFunc {
	return func(c *Client) error {
		if len(certs) > 0 {
			roots := x509.NewCertPool()
			for _, cert := range certs {
				if _, err := os.Stat(cert); err == nil {
					filepath.Abs(cert)
					c, err := ioutil.ReadFile(cert)
					if err != nil {
						return err
					}
					cert = string(c)
				}
				if ok := roots.AppendCertsFromPEM([]byte(cert)); !ok {
					return client.ErrInvalidCert
				}
			}
			if c.tlsConfig != nil {
				c.tlsConfig.RootCAs = roots
			} else {
				c.tlsConfig = &tls.Config{RootCAs: roots}
			}
			c.tlsConfig.InsecureSkipVerify = false
		}
		return nil
	}
}

// Connect tests the mongodb connection and initializes the redis session
func (c *Client) Connect() (client.Session, error) {
	if c.conn == nil {
		if err := c.initConnection(); err != nil {
			return nil, err
		}
	}
	return &Session{c.conn}, nil
}

func (c *Client) initConnection() error {
	opts, _ := redis.ParseURL(c.uri)

	opts.DialTimeout = c.dialTimeout
	opts.ReadTimeout = c.readTimeout
	opts.WriteTimeout = c.writeTimeout
	opts.TLSConfig = c.tlsConfig

	c.conn = redis.NewClient(opts)

	_, err := c.conn.Ping().Result()
	if err != nil {
		return client.ConnectError{Reason: err.Error()}
	}

	return nil
}

// Close implements necessary calls to cleanup the underlying connection.
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

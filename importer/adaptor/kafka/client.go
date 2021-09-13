package kafka

import (
	"net/url"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/log"
)

const (
	// DefaultURI is the default endpoint of Kafka on the local machine.
	// Primarily used when initializing a new Client without a specific URI.
	DefaultURI    = "localhost:9092"
	DefaultOffset = -1
)

var (
	_ client.Client = &Client{}
	_ client.Closer = &Client{}
)

// ClientOptionFunc is a function that configures a Client.
// It is used in NewClient.
type ClientOptionFunc func(*Client) error

// Client represents a client to the underlying File source.
type Client struct {
	uri    []string
	broker *sarama.Broker
	conn   *sarama.Config
	topic  []string
}
type none struct{}

// NewClient creates a new client to work with Kafka.
//
// The caller can configure the new client by passing configuration options
// to the func.
//
// Example:
//
//   client, err := NewClient(
//     WithURI("kafka://localhost:9092"))
//
// If no URI is configured, it uses DefaultURI.
//
// An error is also returned when a configuration option is invalid
func NewClient(options ...ClientOptionFunc) (*Client, error) {
	c := &Client{
		uri:    strings.Split(DefaultURI, ","),
		broker: sarama.NewBroker(DefaultURI),
		topic:  strings.Split(DefaultTopic, ","),
	}

	// Run the options on it
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithURI defines the full connection string for the Kafka connection
func WithURI(uri string) ClientOptionFunc {
	return func(c *Client) error {
		requestedURL, err := url.Parse(uri)
		topic := strings.Trim(requestedURL.Path, "/")
		c.broker = sarama.NewBroker(requestedURL.Host)
		log.Infoln("New Broker created at ", c.broker.Addr())
		c.topic = strings.Split(topic, ",")
		return err
	}
}

// Connect satisfies the client.Client interface.
func (c *Client) Connect() (client.Session, error) {
	connected, _ := c.broker.Connected()
	if connected != true {
		if err := c.initConnection(); err != nil {
			return nil, err
		}
	}
	return &Session{c.broker, c.conn, c.topic}, nil
}

func (c *Client) initConnection() error {
	err := c.broker.Open(c.conn)
	return err
}

// Close implements necessary calls to cleanup the underlying connection.
func (c *Client) Close() {
	c.broker.Close()
	log.Infoln("Broker Closed")
}

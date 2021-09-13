package redis

import (
	"sync"

	adaptor "github.com/appbaseio/abc/importer/adaptor"
	"github.com/appbaseio/abc/importer/client"
)

const (
	sampleConfig = `{
   "uri": "${REDIS_URI}"
  // "timeout": "30s",
  // "tail": false,
  // "ssl": false,
  // "cacerts": ["/path/to/cert.pem"]
}`

	description = "a redis adaptor that functions as both a source and a sink"
)

var _ adaptor.Adaptor = &Redis{}

// Redis is an adaptor that reads data from Redis (http://redis.io/)
type Redis struct {
	adaptor.BaseConfig
	SSL     bool     `json:"ssl"`
	CACerts []string `json:"cacerts"`
}

func init() {
	adaptor.Add(
		"redis",
		func() adaptor.Adaptor {
			return &Redis{}
		},
	)
}

// Client creates an instance of Client to be used for connecting to Redis.
func (r *Redis) Client() (client.Client, error) {
	// TODO: pull db from the URI
	return NewClient(
		WithURI(r.URI),
		WithDialTimeout(r.Timeout),
		WithCACerts(r.CACerts),
	)
}

func (r *Redis) Reader() (client.Reader, error) {
	return newReader(), nil
}

func (r *Redis) Writer(done chan struct{}, wg *sync.WaitGroup) (client.Writer, error) {
	return nil, adaptor.ErrFuncNotSupported{Name: "Writer()", Func: "redis"}
}

func (r *Redis) Verify() error {
	iClient, err := r.Client()
	if err != nil {
		return err
	}

	client := iClient.(*Client)

	_, err = client.Connect()
	if err != nil {
		return err
	}

	client.Close()

	return nil
}

// Description for file adaptor
func (r *Redis) Description() string {
	return description
}

// SampleConfig for file adaptor
func (r *Redis) SampleConfig() string {
	return sampleConfig
}

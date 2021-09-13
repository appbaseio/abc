package kafka

import (
	"strings"
	"sync"

	"github.com/appbaseio/abc/importer/adaptor"
	"github.com/appbaseio/abc/importer/client"
)

const (
	sampleConfig = `{
	host: "http://localhost:9092"
  	"uri": "${KAFKA_URI}",
 	"topic": "",
 	"partition": 0,
	"offset": -1
	}`

	description = "an adaptor that handles publish/subscribe messaging with Kafka"
)

var _ adaptor.Adaptor = &Kafka{}

// KAFKA is an adaptor
type Kafka struct {
	adaptor.BaseConfig
	Topics    []string `json:"topic"`
	Partition int32    `json:"partition"`
	Offset    int64    `json:"offset"`
	SSL       bool     `json:"ssl"`
	CACerts   []string `json:"cacerts"`
}

func init() {
	adaptor.Add(
		"kafka",
		func() adaptor.Adaptor {
			return &Kafka{
				BaseConfig: adaptor.BaseConfig{URI: DefaultURI},
				Topics:     strings.Split(DefaultTopic, ","),
				Partition:  DefaultPartition,
				Offset:     DefaultOffset,
			}
		},
	)
}

// Client creates an instance of Client to be used for connecting to RabbitMQ.
func (r *Kafka) Client() (client.Client, error) {
	return NewClient(WithURI(r.URI))
}

// Reader instantiates a Reader for use with subscribing to one or more topics.
func (r *Kafka) Reader() (client.Reader, error) {
	return &Reader{r.URI, r.Topics}, nil
}

// Writer instantiates a Writer for use with publishing to one or more exchanges.
func (r *Kafka) Writer(done chan struct{}, wg *sync.WaitGroup) (client.Writer, error) {
	return &Writer{r.URI, r.Topics}, nil
}

func (r *Kafka) Verify() error {
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
func (r *Kafka) Description() string {
	return description
}

// SampleConfig for file adaptor
func (r *Kafka) SampleConfig() string {
	return sampleConfig
}

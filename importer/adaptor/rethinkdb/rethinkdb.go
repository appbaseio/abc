package rethinkdb

import (
	"sync"

	adaptor "github.com/appbaseio/abc/importer/adaptor"
	"github.com/appbaseio/abc/importer/client"
)

const (
	sampleConfig = `{
   "uri": "${RETHINKDB_URI}"
  // "timeout": "30s",
  // "tail": false,
  // "ssl": false,
  // "cacerts": ["/path/to/cert.pem"]
}`

	description = "a rethinkdb adaptor that functions as both a source and a sink"
)

var _ adaptor.Adaptor = &RethinkDB{}

// RethinkDB is an adaptor that writes metrics to rethinkdb (http://rethinkdb.com/)
// An open-source distributed database
type RethinkDB struct {
	adaptor.BaseConfig
	Tail    bool     `json:"tail"`
	SSL     bool     `json:"ssl"`
	CACerts []string `json:"cacerts"`
}

func init() {
	adaptor.Add(
		"rethinkdb",
		func() adaptor.Adaptor {
			return &RethinkDB{}
		},
	)
}

func (r *RethinkDB) Client() (client.Client, error) {
	// TODO: pull db from the URI
	return NewClient(
		WithURI(r.URI),
		WithSessionTimeout(r.Timeout),
		WithSSL(r.SSL),
		WithCACerts(r.CACerts),
	)
}

func (r *RethinkDB) Reader() (client.Reader, error) {
	return newReader(r.Tail), nil
}

func (r *RethinkDB) Writer(done chan struct{}, wg *sync.WaitGroup) (client.Writer, error) {
	return newWriter(done, wg), nil
}

func (r *RethinkDB) Verify() error {
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

// Description for rethinkdb adaptor
func (r *RethinkDB) Description() string {
	return description
}

// SampleConfig for rethinkdb adaptor
func (r *RethinkDB) SampleConfig() string {
	return sampleConfig
}

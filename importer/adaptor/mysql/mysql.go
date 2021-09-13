package mysql

import (
	"sync"

	adaptor "github.com/appbaseio/abc/importer/adaptor"
	"github.com/appbaseio/abc/importer/client"
)

const (
	// DefaultDatabase is used when there is not one included in the provided URI.
	DefaultDatabase = "test"

	description = "A mysql source adaptor"

	sampleConfig = `{
  "uri": "${MYSQL_URI}"
  // "timeout": "30s", // defaults to 30s
  // "tail": false // enable tailing
}`
)

var _ adaptor.Adaptor = &MYSQL{}

// MYSQL is an adaptor
type MYSQL struct {
	adaptor.BaseConfig
	Tail bool `json:"tail" doc:"if tail is set, database will be watched for changes"`
}

func init() {
	adaptor.Add(
		"mysql",
		func() adaptor.Adaptor {
			return &MYSQL{}
		},
	)
}

// Client ...
func (ms *MYSQL) Client() (client.Client, error) {
	return NewClient(WithURI(ms.URI))
}

// Reader ...
func (ms *MYSQL) Reader() (client.Reader, error) {
	return newReader(ms.Tail), nil
}

// Writer ...
func (ms *MYSQL) Writer(done chan struct{}, wg *sync.WaitGroup) (client.Writer, error) {
	return nil, adaptor.ErrFuncNotSupported{Name: "Writer()", Func: "MYSQL"}
}

func (ms *MYSQL) Verify() error {
	iClient, err := ms.Client()
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

// Description for the MYSQL adaptor
func (ms *MYSQL) Description() string {
	return description
}

// SampleConfig for MYSQL adaptor
func (ms *MYSQL) SampleConfig() string {
	return sampleConfig
}

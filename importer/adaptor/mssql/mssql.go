package mssql

import (
	"sync"

	adaptor "github.com/appbaseio/abc/importer/adaptor"
	"github.com/appbaseio/abc/importer/client"
)

const (
	// DefaultDatabase is used when there is not one included in the provided URI.
	DefaultDatabase = "Default"

	description = "A mssql source adaptor"

	sampleConfig = `{
  "uri": "${MSSQL_URI}"
  // "timeout": "30s", // defaults to 30s
  // "tail": false // enable tailing
}`
)

var _ adaptor.Adaptor = &MSSQL{}

// MSSQL is an adaptor
type MSSQL struct {
	adaptor.BaseConfig
	Tail bool `json:"tail" doc:"if tail is set, database will be watched for changes"`
}

func init() {
	adaptor.Add(
		"mssql",
		func() adaptor.Adaptor {
			return &MSSQL{}
		},
	)
}

// Client ...
func (ms *MSSQL) Client() (client.Client, error) {
	return NewClient(WithURI(ms.URI))
}

// Reader ...
func (ms *MSSQL) Reader() (client.Reader, error) {
	return newReader(ms.Tail), nil
}

// Writer ...
func (ms *MSSQL) Writer(done chan struct{}, wg *sync.WaitGroup) (client.Writer, error) {
	return nil, adaptor.ErrFuncNotSupported{Name: "Writer()", Func: "mssql"}
}

func (ms *MSSQL) Verify() error {
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

// Description for the MSSQL adaptor
func (ms *MSSQL) Description() string {
	return description
}

// SampleConfig for MSSQL adaptor
func (ms *MSSQL) SampleConfig() string {
	return sampleConfig
}

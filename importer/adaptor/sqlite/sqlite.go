package sqlite

import (
	"sync"

	adaptor "github.com/appbaseio/abc/importer/adaptor"
	"github.com/appbaseio/abc/importer/client"
)

const (
	DefaultDatabase = "data.db"

	description = "A SQLite source adaptor"

	sampleConfig = `{
        "uri" : "${SQLITE_URI}"
        // "timeout": "30s", // default to 30s
        // "tail": false // enable tailing
    }`
)

var _ adaptor.Adaptor = &SQLITE{}

// SQLITE is an adaptor
type SQLITE struct {
	adaptor.BaseConfig
	Tail bool `json:"tail" doc:"if tail is set, database will be watched for changes"`
}

func init() {
	adaptor.Add(
		"sqlite",
		func() adaptor.Adaptor {
			return &SQLITE{}
		},
	)
}

// Client ...
func (sl *SQLITE) Client() (client.Client, error) {
	return NewClient(WithURI(sl.URI))
}

// Reader ...
func (sl *SQLITE) Reader() (client.Reader, error) {
	return newReader(sl.Tail), nil
}

// Writer ...
func (sl *SQLITE) Writer(done chan struct{}, wg *sync.WaitGroup) (client.Writer, error) {
	return nil, adaptor.ErrFuncNotSupported{Name: "Writer()", Func: "sqlite"}
}

func (sl *SQLITE) Verify() error {
	iClient, err := sl.Client()
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
func (sl *SQLITE) Description() string {
	return description
}

// SampleConfig for MSSQL adaptor
func (sl *SQLITE) SampleConfig() string {
	return sampleConfig
}

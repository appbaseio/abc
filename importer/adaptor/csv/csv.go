package csv

import (
	"sync"

	adaptor "github.com/appbaseio/abc/importer/adaptor"
	"github.com/appbaseio/abc/importer/client"
)

const (
	description = "A csv source adaptor"

	sampleConfig = `{
  "uri": "${CSV_PATH}",
  "typeName": "csvType" // type/table name to use
}`
)

var _ adaptor.Adaptor = &CSV{}

// CSV is an adaptor
type CSV struct {
	TypeName string `json:"typeName" doc:"type/table name to use"`
	adaptor.BaseConfig
}

func init() {
	adaptor.Add(
		"csv",
		func() adaptor.Adaptor {
			return &CSV{}
		},
	)
}

// Client ...
func (ms *CSV) Client() (client.Client, error) {
	return NewClient(
		WithURI(ms.URI),
		WithType(ms.TypeName),
	)
}

// Reader ...
func (ms *CSV) Reader() (client.Reader, error) {
	return newReader(ms.TypeName), nil
}

// Writer ...
func (ms *CSV) Writer(done chan struct{}, wg *sync.WaitGroup) (client.Writer, error) {
	return nil, adaptor.ErrFuncNotSupported{Name: "Writer()", Func: "CSV"}
}

func (ms *CSV) Verify() error {
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

// Description for the CSV adaptor
func (ms *CSV) Description() string {
	return description
}

// SampleConfig for CSV adaptor
func (ms *CSV) SampleConfig() string {
	return sampleConfig
}

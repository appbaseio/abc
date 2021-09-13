package jsonl

import (
	"sync"

	"github.com/appbaseio/abc/importer/adaptor"
	"github.com/appbaseio/abc/importer/client"
)

const (
	sampleConfig = `
	{
		"uri": "/full/path/to/json",
  		"typeName": "jsonTypeName"
	}`

	description = "an adaptor that reads / writes json files"
)

// File is an adaptor that can be used as a
// source / sink for file's on disk, as well as a sink to stdout.
type File struct {
	adaptor.BaseConfig
	TypeName string `json:"typeName" doc:"type/table name to use"`
}

func init() {
	adaptor.Add(
		"jsonl",
		func() adaptor.Adaptor {
			return &File{}
		},
	)
}

// Client creates an instance of Client to be used for reading/writing to a file.
func (f *File) Client() (client.Client, error) {
	return NewClient(
		WithURI(f.URI),
		WithType(f.TypeName),
	)
}

// Reader instantiates a Reader for use with working with the file.
func (f *File) Reader() (client.Reader, error) {
	return newReader(f.TypeName), nil
}

// Writer instantiates a Writer for use with working with the file.
func (f *File) Writer(done chan struct{}, wg *sync.WaitGroup) (client.Writer, error) {
	return nil, adaptor.ErrFuncNotSupported{Name: "Writer()", Func: "jsonl"}
}

func (f *File) Verify() error {
	iClient, err := f.Client()
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
func (f *File) Description() string {
	return description
}

// SampleConfig for file adaptor
func (f *File) SampleConfig() string {
	return sampleConfig
}

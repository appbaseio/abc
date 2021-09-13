package firestore

import (
	"sync"

	"github.com/appbaseio/abc/importer/adaptor"
	"github.com/appbaseio/abc/importer/client"
)

const (
	sampleConfig = `{
	"sacPath": "ServiceAccountKey.json",
	"projectId": "sample-project"
	}`

	description = "a firestore source adaptor"
)

type Firestore struct {
	adaptor.BaseConfig
	SACPath   string
	ProjectID string
}

func init() {
	adaptor.Add(
		"firestore",
		func() adaptor.Adaptor {
			return &Firestore{
				SACPath:   DefaultSACPath,
				ProjectID: DefaultProjectID,
			}
		},
	)
}

func (f *Firestore) Client() (client.Client, error) {
	return NewClient(WithConfig(f.SACPath))
}

func (f *Firestore) Reader() (client.Reader, error) {
	return &Reader{}, nil
}

func (f *Firestore) Verify() error {
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

func (f *Firestore) Writer(done chan struct{}, group *sync.WaitGroup) (client.Writer, error) {
	return nil, adaptor.ErrFuncNotSupported{Name: "Writer()", Func: "firestore"}
}

func (f *Firestore) Description() string {
	return description
}

func (f *Firestore) SampleConfig() string {
	return sampleConfig
}

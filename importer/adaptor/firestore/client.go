package firestore

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/log"
	"google.golang.org/api/option"
)

const (
	// Path to the service account credentials key
	DefaultSACPath = "ServiceAccountKey.json"

	// Project-Id designated for the firebase project
	DefaultProjectID = "appbase-firestore-adapter"
)

var (
	_ client.Client = &Client{}
	_ client.Closer = &Client{}
)

type ClientOptionFunc func(*Client) error

type Client struct {
	sacPath   string
	projectID string
	db        string
	fcClient  *firestore.Client
}

type FirebaseProject struct {
	Type string `json:"type"`
	Id   string `json:"project_id"`
}

func NewClient(options ...ClientOptionFunc) (*Client, error) {
	c := &Client{
		db: "firestore",
	}

	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) Connect() (client.Session, error) {
	if c.fcClient == nil {
		sa := option.WithCredentialsFile(c.sacPath)
		c.fcClient, _ = firestore.NewClient(context.Background(), c.projectID, sa)
	}
	return &Session{c.fcClient, c.db}, nil
}

func (c *Client) Close() {
	if c.fcClient != nil {
		c.fcClient.Close()
	}
}

func WithConfig(sacPath string) ClientOptionFunc {
	return func(c *Client) error {
		// read the SAC file to obtain the projectId
		sacFile, err := os.Open(sacPath)
		if err != nil {
			return err
		}
		defer sacFile.Close()
		bytes, _ := ioutil.ReadAll(sacFile)
		var project FirebaseProject
		json.Unmarshal(bytes, &project)
		c.sacPath = sacPath
		c.projectID = project.Id
		log.Infof("Obtaining service account credentials for %v from: %v", project.Id, sacPath)
		return nil
	}
}

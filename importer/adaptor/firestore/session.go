package firestore

import (
	"cloud.google.com/go/firestore"
	"github.com/appbaseio/abc/importer/client"
)

var _ client.Session = &Session{}

type Session struct {
	fc *firestore.Client
	db string
}

package clients

import (
	"net/http"
	"net/url"

	"github.com/appbaseio/abc/importer/client"
	"github.com/hashicorp/go-version"
)

// VersionedClient encapsulates a version.Constraints and Creator func that can be stopred in
// the Clients map.
type VersionedClient struct {
	Constraint version.Constraints
	Creator    Creator
	Reader     Reader
}

// Creator defines the func signature expected for any implementing Writer
type Creator func(*ClientOptions) (client.Writer, error)

// Reader defines the func signature expected for implementing Reader
type Reader func(*ClientOptions) (client.Reader, error)

// Clients contains the map of versioned clients
var Clients = map[string]*VersionedClient{}

// Add exposes the ability for each versioned client to register itself for use
func Add(v string, constraint version.Constraints, creator Creator) {
	Clients[v] = &VersionedClient{constraint, creator, nil}
}

// AddReader ...
func AddReader(v string, constraint version.Constraints, reader Reader) {
	Clients[v+"_reader"] = &VersionedClient{constraint, nil, reader}
}

// ClientOptions defines the available options that can be used to configured the client.Writer
type ClientOptions struct {
	URLs         []string
	UserInfo     *url.Userinfo
	HTTPClient   *http.Client
	Index        string
	RequestSize  int64
	BulkRequests int
	Tail         bool
}

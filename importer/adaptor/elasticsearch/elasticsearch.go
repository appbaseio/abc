package elasticsearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/appbaseio/abc/importer/adaptor"
	"github.com/appbaseio/abc/importer/adaptor/elasticsearch/clients"

	// used to call init function for each client to register itself
	_ "github.com/appbaseio/abc/importer/adaptor/elasticsearch/clients/all"
	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/log"
	"github.com/hashicorp/go-version"
)

const (
	// DefaultIndex is used when there is not one included in the provided URI.
	DefaultIndex = "test"

	DefaultESVersion    = "2.1"
	DefaultRequestSize  = 2 << 19
	DefaultBulkRequests = 1000

	description = "an elasticsearch source/sink adaptor"

	sampleConfig = `{
  "uri": "${ELASTICSEARCH_URI}"
  // "timeout": "10s", // defaults to 30s
  // "aws_access_key": "ABCDEF", // used for signing requests to AWS Elasticsearch service
  // "aws_access_secret": "ABCDEF", // used for signing requests to AWS Elasticsearch service
  // "tail": false, // enable tailing
  // "request_size": 524288,
  // "bulk_requests": 1000
}`
)

var _ adaptor.Adaptor = &Elasticsearch{}

// Elasticsearch is an adaptor to connect a pipeline to
// an elasticsearch cluster.
type Elasticsearch struct {
	adaptor.BaseConfig
	AWSAccessKeyID  string `json:"aws_access_key" doc:"credentials for use with AWS Elasticsearch service"`
	AWSAccessSecret string `json:"aws_access_secret" doc:"credentials for use with AWS Elasticsearch service"`
	Tail            bool   `json:"tail" doc:"if tail is set, ES index will be watched for changes"`
	RequestSize     int64  `json:"request_size"`
	BulkRequests    int    `json:"bulk_requests"`
}

// Description for the Elasticsearcb adaptor
func (e *Elasticsearch) Description() string {
	return description
}

// SampleConfig for elasticsearch adaptor
func (e *Elasticsearch) SampleConfig() string {
	return sampleConfig
}

func init() {
	adaptor.Add(
		"elasticsearch",
		func() adaptor.Adaptor {
			return &Elasticsearch{
				RequestSize:  DefaultRequestSize,
				BulkRequests: DefaultBulkRequests,
			}
		},
	)
}

// Client returns a client that doesn't do anything other than fulfill the client.Client interface.
func (e *Elasticsearch) Client() (client.Client, error) {
	return &client.Mock{}, nil
}

// Reader returns an error because this adaptor is currently not supported as a Source.
func (e *Elasticsearch) Reader() (client.Reader, error) {
	return setupReader(e)
}

func (e *Elasticsearch) Verify() error {
	_, err := setupReader(e)
	if err != nil {
		return err
	}

	return nil
}

// Writer determines the which underlying writer to used based on the cluster's version.
func (e *Elasticsearch) Writer(done chan struct{}, wg *sync.WaitGroup) (client.Writer, error) {
	return setupWriter(e)
}

// setupWriter ...
func setupWriter(conf *Elasticsearch) (client.Writer, error) {
	uri, err := url.Parse(conf.URI)
	if err != nil {
		return nil, client.InvalidURIError{URI: conf.URI, Err: err.Error()}
	}

	if uri.Path == "" {
		uri.Path = fmt.Sprintf("/%s", DefaultIndex)
	}

	hostsAndPorts := strings.Split(uri.Host, ",")
	stringVersion, err := determineVersion(uri, hostsAndPorts[0], uri.User)
	// stringVersion, err := getESVersionFor(uri.String())
	log.Infoln("ES Version: ", stringVersion)
	if err != nil {
		return nil, err
	}

	v, err := version.NewVersion(stringVersion)
	if err != nil {
		return nil, client.VersionError{URI: conf.URI, V: stringVersion, Err: err.Error()}
	}

	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		log.Debugf("failed to parse duration, %s, falling back to default timeout of 30s", conf.Timeout)
		timeout = 300 * time.Second
	}

	httpClient := &http.Client{
		Timeout:   timeout,
		Transport: newTransport(conf.AWSAccessKeyID, conf.AWSAccessSecret),
	}

	for _, vc := range clients.Clients {
		if vc.Constraint.Check(v) && vc.Creator != nil {
			urls := make([]string, len(hostsAndPorts))
			for i, hAndP := range hostsAndPorts {
				urls[i] = fmt.Sprintf("%s://%s", uri.Scheme, hAndP)
			}
			opts := &clients.ClientOptions{
				URLs:         urls,
				UserInfo:     uri.User,
				HTTPClient:   httpClient,
				Index:        uri.Path[1:],
				RequestSize:  conf.RequestSize,
				BulkRequests: conf.BulkRequests,
				Tail:         conf.Tail,
			}
			versionedClient, _ := vc.Creator(opts)
			return versionedClient, nil
		}
	}

	return nil, client.VersionError{URI: conf.URI, V: stringVersion, Err: "unsupported client"}
}

// setupReader ...
func setupReader(conf *Elasticsearch) (client.Reader, error) {
	uri, err := url.Parse(conf.URI)
	if err != nil {
		return nil, client.InvalidURIError{URI: conf.URI, Err: err.Error()}
	}
	// fail if no index defined
	if uri.Path == "" {
		return nil, client.InvalidURIError{URI: conf.URI, Err: "Index not defined in URI"}
	}

	hostsAndPorts := strings.Split(uri.Host, ",")
	stringVersion, err := determineVersion(uri, hostsAndPorts[0], uri.User)
	if err != nil {
		return nil, err
	}

	v, err := version.NewVersion(stringVersion)
	if err != nil {
		return nil, client.VersionError{URI: conf.URI, V: stringVersion, Err: err.Error()}
	}

	timeout, err := time.ParseDuration(conf.Timeout)
	if err != nil {
		log.Debugf("failed to parse duration, %s, falling back to default timeout of 30s", conf.Timeout)
		timeout = 30 * time.Second
	}

	httpClient := &http.Client{
		Timeout:   timeout,
		Transport: newTransport(conf.AWSAccessKeyID, conf.AWSAccessSecret),
	}

	for _, vc := range clients.Clients {
		if vc.Constraint.Check(v) && vc.Reader != nil {
			urls := make([]string, len(hostsAndPorts))
			for i, hAndP := range hostsAndPorts {
				urls[i] = fmt.Sprintf("%s://%s", uri.Scheme, hAndP)
			}
			opts := &clients.ClientOptions{
				URLs:       urls,
				UserInfo:   uri.User,
				HTTPClient: httpClient,
				Index:      uri.Path[1:],
			}
			versionedClient, _ := vc.Reader(opts)
			return versionedClient, nil
		}
	}

	return nil, client.VersionError{URI: conf.URI, V: stringVersion, Err: "unsupported client"}
}

func getESVersionFor(uri string) (string, error) {
	appName := getAppName(uri)
	uri += "/_settings?human"

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return DefaultESVersion, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return DefaultESVersion, client.ConnectError{Reason: uri}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DefaultESVersion, client.VersionError{URI: uri, V: "", Err: "unable to read response body"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DefaultESVersion, client.VersionError{URI: uri, V: "", Err: fmt.Sprintf("bad status code: %d", resp.StatusCode)}
	}

	var jsonResp interface{}
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return DefaultESVersion, client.VersionError{URI: uri, V: "", Err: fmt.Sprintf("malformed JSON: %s", body)}
	} else {
		ver := getVersionStringFrom(jsonResp, appName)
		if ver == "" {
			return DefaultESVersion, client.VersionError{URI: uri, V: "", Err: fmt.Sprintf("missing version: %s", body)}
		} else {
			return ver, nil
		}
	}
}

func getAppName(uri string) string {
	tokens := strings.Split(uri, "/")
	size := len(tokens)
	return tokens[size-1]
}

// TODO: Ugly workaround, find a better alternative
func getVersionStringFrom(jsonResp interface{}, appName string) string {
	jsonObj := jsonResp.(map[string]interface{})
	app := jsonObj[appName].(map[string]interface{})
	settings := app["settings"].(map[string]interface{})
	index := settings["index"].(map[string]interface{})
	ver := index["version"].(map[string]interface{})
	return ver["created_string"].(string)
}

func determineVersion(uri *url.URL, host string, user *url.Userinfo) (string, error) {
	reqURL := fmt.Sprintf("%s://%s", uri.Scheme, host)

	// check if appbase.io
	if strings.Contains(reqURL, "scalr.api.appbase.io") {
		return getESVersionFor(uri.String())
	}

	// normal ES cluster
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "", err
	}
	if user != nil {
		if pwd, ok := user.Password(); ok {
			req.SetBasicAuth(user.Username(), pwd)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", client.ConnectError{Reason: reqURL}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", client.VersionError{URI: reqURL, V: "", Err: "unable to read response body"}
	}
	defer resp.Body.Close()
	var r struct {
		Name    string `json:"name"`
		Version struct {
			Number string `json:"number"`
		} `json:"version"`
		Tagline string `json:"tagline"`
	}
	if resp.StatusCode != http.StatusOK {
		return "", client.VersionError{URI: reqURL, V: "", Err: fmt.Sprintf("bad status code: %d", resp.StatusCode)}
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", client.VersionError{URI: reqURL, V: "", Err: fmt.Sprintf("malformed JSON: %s", body)}
	} else if r.Version.Number == "" {
		return "", client.VersionError{URI: reqURL, V: "", Err: fmt.Sprintf("missing version: %s", body)}
	}

	// If the tagline contains `OpenSearch` in it, that means that this is
	// an OpenSearch cluster so we should always treat it like ES 8.x . Thus we
	// will return a different version that will make abc think this is an ES 8.x
	// cluster instead.
	if strings.Contains(strings.ToLower(r.Tagline), "opensearch") {
		return "8.8.1", nil
	}

	return r.Version.Number, nil
}

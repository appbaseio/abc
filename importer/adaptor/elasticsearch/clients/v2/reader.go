package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/appbaseio/abc/importer/adaptor/elasticsearch/clients"
	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/importer/message"
	"github.com/appbaseio/abc/importer/message/ops"
	"github.com/appbaseio/abc/log"
	version "github.com/hashicorp/go-version"
	elastic "gopkg.in/olivere/elastic.v3"
)

var _ client.Reader = &Reader{}

// Reader fulfills the client.Reader interface
type Reader struct {
	tail      bool
	index     string
	logger    log.Logger
	esClient  *elastic.Client
	isAppbase bool
}

func init() {
	constraint, _ := version.NewConstraint(">= 2.0, < 5.0")
	clients.AddReader("v2", constraint, func(opts *clients.ClientOptions) (client.Reader, error) {
		esOptions := []elastic.ClientOptionFunc{
			elastic.SetURL(opts.URLs...),
			elastic.SetSniff(false),
			elastic.SetHttpClient(opts.HTTPClient),
			elastic.SetMaxRetries(2),
		}
		if opts.UserInfo != nil {
			if pwd, ok := opts.UserInfo.Password(); ok {
				esOptions = append(esOptions, elastic.SetBasicAuth(opts.UserInfo.Username(), pwd))
			}
		}
		esClient, err := elastic.NewClient(esOptions...)
		if err != nil {
			return nil, err
		}
		r := &Reader{}
		r.tail = false // TODO: fix
		r.esClient = esClient
		r.logger = log.With("reader", "elasticsearch").With("version", 2)
		r.index = opts.Index
		// check appbase
		uri := opts.URLs[0]
		if strings.Contains(uri, "scalr.api.appbase.io") {
			r.isAppbase = true
		} else {
			r.isAppbase = false
		}
		return r, nil
	})
}

func (r *Reader) Read(resumeMap map[string]client.MessageSet, filterFn client.NsFilterFunc) client.MessageChanFunc {
	return func(s client.Session, done chan struct{}) (chan client.MessageSet, error) {
		out := make(chan client.MessageSet)
		go func() {
			defer close(out)
			log.With("cluster", r.esClient).Infoln("starting Read func")
			// list mappings
			mappings := r.listMappings(filterFn)
			// fetch data
			tableDone := r.iterateType(mappings, out, done)
			func() {
				for {
					select {
					case i, ok := <-tableDone:
						if !ok {
							return
						}
						r.logger.Infof("Mapping %s done", i)
					case <-done:
						return
					}
				}
			}()

			return
		}()
		return out, nil
	}
}

func (r *Reader) listMappings(filterFn func(name string) bool) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		res, err := r.esClient.PerformRequest(context.Background(), "GET", fmt.Sprintf("/%s/_mapping", r.index), nil, nil)
		if err != nil {
			r.logger.Errorf("Error reading mappings %s", err)
			return // exit
		}
		// convert response
		result, _ := res.Body.MarshalJSON()
		var m map[string]interface{}
		err = json.Unmarshal(result, &m)
		if err != nil {
			r.logger.Errorf("Error reading mappings %s", err)
			return
		}
		// fetch
		typesResp := m[r.index].(map[string]interface{})["mappings"].(map[string]interface{})
		for typeName := range typesResp {
			if r.isAppbase && typeName == "_default_" {
				// appbase reserved mapping
				continue
			}
			if filterFn(typeName) {
				log.With("cluster", r.esClient).With("type", typeName).Infoln("sending for iteration...")
				out <- typeName
			} else {
				log.With("cluster", r.esClient).With("type", typeName).Infoln("skipping iteration...")
			}
		}
	}()
	return out
}

func (r *Reader) iterateType(in <-chan string, out chan<- client.MessageSet, done chan struct{}) <-chan string {
	tableDone := make(chan string)
	go func() {
		defer close(tableDone)
		for {
			select {
			case t, ok := <-in:
				if !ok {
					return
				}
				r.logger.Infoln(t)
				// do scrolling here
				scrollService := r.esClient.Scroll(r.index).Type(t).Size(100).Query(elastic.NewMatchAllQuery())
				// https://godoc.org/gopkg.in/olivere/elastic.v3#ScrollService
				for {
					results, err := scrollService.Do(context.Background())
					if err == io.EOF {
						break // all results retrieved
					}
					if err != nil && r.isAppbase &&
						(strings.Contains(err.Error(), "400") || strings.Contains(err.Error(), "403")) {
						break // another way to show scroll end
					}
					if err != nil {
						r.logger.Errorf("Error scrolling documents %s", err)
						break
					}
					for _, hit := range results.Hits.Hits {
						bytes, _ := hit.Source.MarshalJSON()
						var m map[string]interface{}
						err := json.Unmarshal(bytes, &m)
						if err != nil {
							r.logger.Errorf("Problem unmarshaling document %s", err)
							return
						}
						// r.logger.Infoln(m)
						// send data
						out <- client.MessageSet{
							Msg: message.From(ops.Insert, t, m),
						}
					}
					// return early?
					select {
					default:
					case <-done:
						log.With("cluster", r.esClient).Infoln("Scroll stopped midway")
						return
					}
				}
				// clear scroll
				if !r.isAppbase {
					// no such method in appbase
					scrollService.Clear(context.Background())
				}
				// finish
				tableDone <- t
			case <-done:
				log.With("cluster", r.esClient).Infoln("Done with iterating")
				return
			}
		}
	}()
	return tableDone
}

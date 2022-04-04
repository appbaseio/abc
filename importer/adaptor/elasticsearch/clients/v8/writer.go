package v8

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/appbaseio/abc/importer/adaptor/elasticsearch/clients"
	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/importer/function/mapping"
	"github.com/appbaseio/abc/importer/message"
	"github.com/appbaseio/abc/importer/message/ops"
	"github.com/appbaseio/abc/log"
	"github.com/hashicorp/go-version"
	"github.com/olivere/elastic/v7"
)

var (
	_                client.Writer = &Writer{}
	_                client.Closer = &Writer{}
	isMappingApplied               = false
)

// Writer implements client.Writer and client.Session for sending requests to an elasticsearch
// cluster via its _bulk API.
type Writer struct {
	index      string
	bs         *elastic.BulkService
	indexCount int
	sync.Mutex
	confirmChan  chan struct{}
	esClient     *elastic.Client
	logger       log.Logger
	ticker       *time.Ticker
	requestSize  int64
	bulkRequests int
}

func init() {
	constraint, _ := version.NewConstraint(">=8.0")
	clients.Add("v8", constraint, func(opts *clients.ClientOptions) (client.Writer, error) {
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
		w := &Writer{
			index:  opts.Index,
			logger: log.With("writer", "elasticsearch").With("version", 7),
		}
		// bulk handler
		w.bs = esClient.Bulk().Index(opts.Index)
		w.esClient = esClient
		w.ticker = time.NewTicker(5 * time.Second)
		w.requestSize = opts.RequestSize
		w.bulkRequests = opts.BulkRequests
		if opts.Tail {
			go func() {
				for range w.ticker.C {
					err = w.EsCommit()
					if err != nil {
						w.logger.Errorln(err)
						panic(err)
					}
				}
			}()
		}
		return w, nil
	})
}

func (w *Writer) Write(msg message.Msg) func(client.Session) (message.Msg, error) {
	return func(s client.Session) (message.Msg, error) {
		w.Lock()
		w.confirmChan = msg.Confirms()
		w.Unlock()

		// apply mapping
		if mapping.IsMappingSet {
			isMappingApplied = true
			err := w.setMapping(w.esClient, mapping.CurrentMapping)
			if err != nil {
				return nil, err
			}
		}

		if msg.Data().AsMap() != nil && len(msg.Data().AsMap()) > 0 {
			var id string
			var index string
			if _, ok := msg.Data()["_id"]; ok {
				id = msg.ID()
				msg.Data().Delete("_id")
			}

			// override the default import index if specified in the message data
			if _, ok := msg.Data()["_index"]; ok {
				index = msg.Data()["_index"].(string)
				msg.Data().Delete("_index")
			}

			var br elastic.BulkableRequest
			switch msg.OP() {
			case ops.Delete:
				br = elastic.NewBulkDeleteRequest().Index(index).Id(id)
			case ops.Insert:
				br = elastic.NewBulkIndexRequest().Id(id).Index(index).Doc(msg.Data())
			case ops.Update:
				br = elastic.NewBulkUpdateRequest().Id(id).Index(index).Doc(msg.Data())
			}

			// add a bulk request only if # of requests < --bulk_requests AND size of requests < --request_size switches
			w.Lock()
			if w.bs.NumberOfActions() < w.bulkRequests && w.bs.EstimatedSizeInBytes() < w.requestSize {
				log.Debugln(br.String())
				w.bs.Add(br)
			}
			w.Unlock()
		}

		// clear confirmChan
		if w.confirmChan != nil {
			close(w.confirmChan)
			w.confirmChan = nil
		}

		var err error
		// commit if # requests exceed either constraint
		if w.bs.NumberOfActions() >= w.bulkRequests || w.bs.EstimatedSizeInBytes() >= w.requestSize {
			err = w.EsCommit()
			if err != nil {
				panic(err)
			}
		}
		return msg, err
	}
}

// EsCommit is called to commit changes to ES
func (w *Writer) EsCommit() error {
	defer w.Unlock()
	w.Lock()
	numberOfActions := w.bs.NumberOfActions()
	if numberOfActions > 0 {
		w.logger.Infof("Going through %d", numberOfActions)
		w.indexCount += numberOfActions

		w.logger.Infof("indexing %d data record(s)\n", numberOfActions)
		startTime := time.Now()
		data, err := w.bs.Do(context.Background())
		w.logger.Infof("%d data record(s) indexed in %f seconds\n", numberOfActions, time.Since(startTime).Seconds())
		fmt.Printf("%d total data record(s) indexed\n", w.indexCount)

		if err != nil {
			w.logger.Errorln(err)
		}
		if data != nil && len(data.Failed()) > 0 {
			fl := data.Failed()[0]
			w.logger.Infof("fail %s %s %v %v %v", fl.Id, fl.Index, fl.Type, fl.Error, fl.Status)
			return errors.New(fmt.Sprintf("%v", fl.Error))
		}
		return err
	}
	return nil
}

// Close is called by clients.Close() when it receives on the done channel.
func (w *Writer) Close() {
	err := w.EsCommit() // save changes before exiting
	w.logger.Infoln("closing BulkService")
	w.esClient.Stop()
	w.ticker.Stop()

	if err != nil {
		panic(err)
	}
}

// setMapping sets the index mapping
func (w *Writer) setMapping(esClient *elastic.Client, mapping map[string]interface{}) error {
	log.Debugf("Going to apply mapping %s", mapping)
	_, err := esClient.CreateIndex(w.index).BodyJson(map[string]interface{}{
		"mappings": mapping,
	}).Do(context.Background())
	// BodyJson(mapping).Do(context.Background())
	if err != nil {
		// if above fails, try assuming the index already exists
		_, innerErr := esClient.PutMapping().Index(w.index).BodyJson(mapping).Do(context.Background())

		if innerErr != nil {
			return errors.New("Mapping request failed")
		}
		return nil
	}
	return nil
}

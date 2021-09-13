package kafka

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/importer/commitlog"
	"github.com/appbaseio/abc/importer/message"
	"github.com/appbaseio/abc/importer/message/data"
	"github.com/appbaseio/abc/importer/message/ops"
	"github.com/appbaseio/abc/log"
)

var _ client.Reader = &Reader{}

// Reader implements client.Reader by consuming messages from the cluster based on its configuration.
type Reader struct {
	Uri    string
	Topics []string
}

func (r *Reader) Read(_ map[string]client.MessageSet, filterFn client.NsFilterFunc) client.MessageChanFunc {
	return func(s client.Session, done chan struct{}) (chan client.MessageSet, error) {
		out := make(chan client.MessageSet)
		var topics []string
		config := s.(*Session).conn
		broker := s.(*Session).broker
		topics = s.(*Session).topic
		uri := strings.Split(broker.Addr(), ",")
		Client, err := sarama.NewClient(uri, config)
		if err != nil {
			log.Errorln(err)
			return out, nil
		}
		if topics[0] == "" {
			topics, err = Client.Topics()
			log.Infoln("No topic name given, consuming from topic(s): ", topics)
			if err != nil {
				log.Errorln("Error while consuming: ", err)
			}
		} else {
			log.Infoln("Consuming from topic(s): ", topics)
		}

		// topics, err = client.Topics()
		var filterTopics []string
		for _, q := range topics {
			if filterFn(q) {
				filterTopics = append(filterTopics, q)
			}
		}
		go func(qs []string, session *Session) {
			defer func() {
				broker.Close()
				close(out)
			}()
			master, err := sarama.NewConsumerFromClient(Client)
			if err != nil {
				return
			}
			defer close(out)
			var wg sync.WaitGroup
			for _, topic := range filterTopics {
				wg.Add(1)
				go consumeTopics(master, topic, &wg, done, out)
			}
			wg.Wait()
		}(topics, s.(*Session))
		return out, nil
	}
}

func consumeTopics(master sarama.Consumer, topic string, wg *sync.WaitGroup, done chan struct{}, out chan client.MessageSet) {
	defer func() {
		log.With("topic", topic).Infoln("consuming complete")
		wg.Done()
	}()
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case err := <-consumer.Errors():
			log.With("Error while consuming: ", err)
		case msg := <-consumer.Messages():
			log.With("Sending msg to channel: ", msg)
			var result map[string]interface{}

			if jerr := json.NewDecoder(bytes.NewReader(msg.Value)).Decode(&result); jerr != nil {
				log.Errorf("unable to decode message to JSON, %s", jerr)
				continue
			}
			log.Infoln("Record with offset ", msg.Offset, " indexed")
			out <- client.MessageSet{
				Msg:  message.From(ops.Insert, msg.Topic, data.Data(result)),
				Mode: commitlog.Sync,
			}

		case <-done:

		}
	}
}

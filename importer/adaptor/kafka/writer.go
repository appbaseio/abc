package kafka

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/importer/message"
)

const (

	// DefaultRoutingKey is set to an empty string so all messages published to the exchange will
	DefaultPartition = 0
	// DefaultOffset = -1
	// get routed to whatever queues are bound to it.
	DefaultTopic = "topic"
)

var _ client.Writer = &Writer{}

type Writer struct {
	Uri    string
	Topics []string
}

func (w *Writer) Write(msg message.Msg) func(client.Session) (message.Msg, error) {
	return func(s client.Session) (message.Msg, error) {
		config := s.(*Session).conn
		uri := strings.Split(w.Uri, ",")
		producer, err := sarama.NewAsyncProducer(uri, config)
		if err != nil {
			panic(err) // Should not reach here
		}

		defer func() {
			if err := producer.Close(); err != nil {
				panic(err) // Should not reach here
			}
		}()
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(msg.Data())
		saramaMsg := &sarama.ProducerMessage{
			Topic: w.Topics[0],
			Value: sarama.StringEncoder(b.String()),
		}

		select {
		case producer.Input() <- saramaMsg:
		case err := <-producer.Errors():
			return msg, err
		}

		return msg, err
	}
}

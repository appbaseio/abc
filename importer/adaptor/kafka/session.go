package kafka

import "github.com/Shopify/sarama"

type Session struct {
	broker *sarama.Broker
	conn   *sarama.Config
	topic  []string
}

package kafkaClient

import (
	"log"

	"github.com/Shopify/sarama"
)

type ProducerAsyncClient interface {
	sarama.AsyncProducer
	SendOne(topic string, key string, message string, partition int32) error
}

type producerAsyncClient struct {
	sarama.AsyncProducer
}

// NewProducerAsyncClient 异步生产者
func NewProducerAsyncClient(cfg Config) (ProducerAsyncClient, func(), error) {
	client, err := NewKafkaClient(cfg, true)
	if err != nil {
		log.Fatalln(err)
	}
	// 建立异步生产者连接
	producerClient, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, nil, err
	}

	return &producerAsyncClient{producerClient}, func() {
		defer producerClient.AsyncClose()
	}, nil
}

// SendOne 单条消息发送
func (cli *producerAsyncClient) SendOne(topic string, key string, message string, partition int32) error {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(message),
		Partition: partition,
	}
	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}
	select {
	case cli.Input() <- msg:
		return nil
	case err := <-cli.Errors():
		return err
	}
}

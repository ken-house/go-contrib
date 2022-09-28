package kafkaClient

import (
	"log"

	"github.com/Shopify/sarama"
)

type ProducerSyncClient interface {
	sarama.SyncProducer
	SendOne(topic string, key string, message string, partition int32) (int32, int64, error)
	SendMany(topic string, key string, messageList []string, partition int32) error
}

type producerSyncClient struct {
	sarama.SyncProducer
}

// NewProducerSyncClient 同步生产者
func NewProducerSyncClient(cfg Config) (ProducerSyncClient, func(), error) {
	client, err := NewKafkaClient(cfg, false)
	if err != nil {
		log.Fatalln(err)
	}
	// 建立同步生产者连接
	producerClient, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, nil, err
	}

	return &producerSyncClient{producerClient}, func() {
		defer producerClient.Close()
	}, nil
}

// SendOne 单条消息发送
func (cli *producerSyncClient) SendOne(topic string, key string, message string, partition int32) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(message),
		Partition: partition,
	}
	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}
	return cli.SendMessage(msg)
}

// SendMany 多条消息发送
func (cli *producerSyncClient) SendMany(topic string, key string, messageList []string, partition int32) error {
	msgList := make([]*sarama.ProducerMessage, 0, 100)
	for _, message := range messageList {
		msg := &sarama.ProducerMessage{
			Topic:     topic,
			Value:     sarama.StringEncoder(message),
			Partition: partition,
		}
		if key != "" {
			msg.Key = sarama.StringEncoder(key)
		}
		msgList = append(msgList, msg)
	}
	return cli.SendMessages(msgList)
}

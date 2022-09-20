package kafkaClient

import (
	"time"

	"github.com/Shopify/sarama"
)

type ProducerAsyncClient interface {
	sarama.AsyncProducer
	SendOne(topic string, key string, message string, partition int32)
}

type producerAsyncClient struct {
	sarama.AsyncProducer
}

// NewProducerAsyncClient 异步生产者
func NewProducerAsyncClient(cfg Config) (ProducerAsyncClient, func(), error) {
	config := sarama.NewConfig()
	// 指定应答方式
	config.Producer.RequiredAcks = sarama.RequiredAcks(cfg.ProducerConfig.Ack)
	// 设置达到多少条消息才发送到kafka
	config.Producer.Flush.Messages = cfg.ProducerConfig.FlushMessageNum
	// 设置间隔多少秒才发送到kafka
	if cfg.ProducerConfig.FlushMessageFrequency > 0 {
		config.Producer.Flush.Frequency = time.Duration(cfg.ProducerConfig.FlushMessageFrequency) * time.Second
	}
	// 成功交付的消息将在success channel返回 必须指定为true
	config.Producer.Return.Successes = true
	// 指定分区算法
	setPartition(config, cfg.ProducerConfig.PartitionerType)
	// 建立同步生产者连接
	productClient, err := sarama.NewAsyncProducer(cfg.ServerAddrList, config)
	if err != nil {
		return nil, nil, err
	}

	return &producerAsyncClient{productClient}, func() {
		defer productClient.Close()
	}, nil
}

// SendOne 单条消息发送
func (cli *producerAsyncClient) SendOne(topic string, key string, message string, partition int32) {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(message),
		Partition: partition,
	}
	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}
	cli.Input() <- msg
}

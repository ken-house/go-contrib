package sarama

import (
	"time"

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
	config := sarama.NewConfig()
	//// 指定kafka版本 - 需根据实际kafka版本调整
	config.Version = sarama.V2_8_1_0
	// 指定应答方式
	config.Producer.RequiredAcks = sarama.RequiredAcks(cfg.ProducerConfig.Ack)
	// 设置达到多少条消息才发送到kafka，相当于batch.size(批次大小)
	config.Producer.Flush.Messages = cfg.ProducerConfig.BatchMessageNum
	// 设置间隔多少秒才发送到kafka，相当于linger.ms（等待时间ms）
	if cfg.ProducerConfig.LingerMs > 0 {
		config.Producer.Flush.Frequency = time.Duration(cfg.ProducerConfig.LingerMs) * time.Millisecond
	}
	// 指定数据压缩方式
	config.Producer.Compression = cfg.ProducerConfig.CompressionType
	// 生产者缓冲区大小
	if cfg.ProducerConfig.RecordAccumulator > 0 {
		config.Producer.MaxMessageBytes = cfg.ProducerConfig.RecordAccumulator
	}
	// 指定分区算法
	setProducerPartitionPolicy(config, cfg.ProducerConfig.PartitionerPolicy)
	// 建立同步生产者连接
	producerClient, err := sarama.NewAsyncProducer(cfg.ServerAddrList, config)
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

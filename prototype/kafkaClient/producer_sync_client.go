package kafkaClient

import (
	"time"

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
	config := sarama.NewConfig()
	// 指定kafka版本 - 需根据实际kafka版本调整
	config.Version = sarama.V2_8_1_0
	// 指定应答方式
	config.Producer.RequiredAcks = sarama.RequiredAcks(cfg.ProducerConfig.Ack)
	// 设置达到多少条消息才发送到kafka，相当于batch.size(批次大小)
	config.Producer.Flush.Messages = cfg.ProducerConfig.BatchMessageNum
	// 设置间隔多少秒才发送到kafka，相当于linger.ms（等待时间）
	if cfg.ProducerConfig.LingerMs > 0 {
		config.Producer.Flush.Frequency = time.Duration(cfg.ProducerConfig.LingerMs) * time.Millisecond
	}
	// 指定数据压缩方式
	config.Producer.Compression = cfg.ProducerConfig.CompressionType
	// 生产者缓冲区大小
	if cfg.ProducerConfig.RecordAccumulator > 0 {
		config.Producer.MaxMessageBytes = cfg.ProducerConfig.RecordAccumulator
	}
	// 成功交付的消息将在success channel返回 同步发送必须指定为true
	config.Producer.Return.Successes = true
	// 指定分区算法
	setPartitionPolicy(config, cfg.ProducerConfig.PartitionerPolicy)
	// 建立同步生产者连接
	producerClient, err := sarama.NewSyncProducer(cfg.ServerAddrList, config)
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

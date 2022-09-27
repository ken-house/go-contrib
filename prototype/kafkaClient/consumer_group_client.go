package kafkaClient

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type ConsumerGroupClient interface {
	sarama.ConsumerGroup
	ConsumeTopic(ctx context.Context, topicList []string, consumeHandler func(message *sarama.ConsumerMessage))
}

type consumerGroupClient struct {
	sarama.ConsumerGroup
}

func NewConsumerGroupClient(cfg Config) (ConsumerGroupClient, func(), error) {
	config := sarama.NewConfig()
	// 指定kafka版本 - 需根据实际kafka版本调整
	config.Version = sarama.V2_8_1_0
	// 设置批次抓取最小字节
	if cfg.ConsumerConfig.FetchMinBytes > 0 {
		config.Consumer.Fetch.Min = cfg.ConsumerConfig.FetchMinBytes
	}
	// 设置批次抓取最大字节
	if cfg.ConsumerConfig.FetchMaxBytes > 0 {
		config.Consumer.Fetch.Max = cfg.ConsumerConfig.FetchMaxBytes
	}
	// 设置每批次数据达到的超时时间（毫秒）
	if cfg.ConsumerConfig.MaxWaitTimeMs > 0 {
		config.Consumer.MaxWaitTime = time.Duration(cfg.ConsumerConfig.MaxWaitTimeMs) * time.Millisecond
	}
	// 是否从头消费
	if cfg.ConsumerConfig.FromBeginning {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	// offset是否自动提交
	if !cfg.ConsumerConfig.OffsetAutoCommitEnabled {
		config.Consumer.Offsets.AutoCommit.Enable = cfg.ConsumerConfig.OffsetAutoCommitEnabled
	}
	// offset自动提交时间间隔
	if cfg.ConsumerConfig.OffsetAutoCommitInterval > 0 {
		config.Consumer.Offsets.AutoCommit.Interval = time.Duration(cfg.ConsumerConfig.OffsetAutoCommitInterval) * time.Second
	}
	// 设置消费者分区分配算法
	setConsumerPartitionPolicy(config, cfg.ConsumerConfig.BalanceStrategy)
	// 创建消费者组客户端
	client, err := sarama.NewConsumerGroup(cfg.ServerAddrList, cfg.ConsumerConfig.GroupId, config)
	if err != nil {
		return nil, nil, err
	}

	return &consumerGroupClient{client}, func() {
		defer client.Close()
	}, nil
}

// ConsumeTopic 消费主题
func (cli *consumerGroupClient) ConsumeTopic(ctx context.Context, topicList []string, consumeFunc func(message *sarama.ConsumerMessage)) {
	consumerHandler := consumeHandler{
		handle: consumeFunc,
		ready:  make(chan struct{}),
	}
	go func() {
		for {
			if err := cli.Consume(ctx, topicList, &consumerHandler); err != nil {
				panic(err)
				return
			}
		}
	}()
	<-consumerHandler.ready
}

type consumeHandler struct {
	handle func(*sarama.ConsumerMessage)
	ready  chan struct{}
}

func (handler *consumeHandler) Setup(session sarama.ConsumerGroupSession) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("consumer setup error: %v\n", err)
		}
	}()
	close(handler.ready)
	return nil
}

func (handler *consumeHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (handler *consumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		handler.handle(msg)
		session.MarkMessage(msg, "")
	}
	return nil
}

package kafkaClient

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

var offsetAutoCommitEnabled bool

type ConsumerGroupClient interface {
	sarama.ConsumerGroup
	ConsumeTopic(ctx context.Context, topicList []string, consumeFunc func(message *sarama.ConsumerMessage)) error
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
	// offset是否自动提交，同时设置一个全局变量offsetAutoCommitEnabled
	offsetAutoCommitEnabled = cfg.ConsumerConfig.OffsetAutoCommitEnabled
	if !cfg.ConsumerConfig.OffsetAutoCommitEnabled {
		config.Consumer.Offsets.AutoCommit.Enable = cfg.ConsumerConfig.OffsetAutoCommitEnabled
	}
	// offset自动提交时间间隔
	if cfg.ConsumerConfig.OffsetAutoCommitInterval > 0 {
		config.Consumer.Offsets.AutoCommit.Interval = time.Duration(cfg.ConsumerConfig.OffsetAutoCommitInterval) * time.Second
	}
	// 一次拉取返回消息的最大条数
	if cfg.ConsumerConfig.MaxPollRecords > 0 {
		config.ChannelBufferSize = cfg.ConsumerConfig.MaxPollRecords
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
func (cli *consumerGroupClient) ConsumeTopic(ctx context.Context, topicList []string, consumeFunc func(message *sarama.ConsumerMessage)) error {
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
	return nil
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
	// 指定offset消费
	session.ResetOffset("first", 1, 443, "")
	close(handler.ready)
	return nil
}

func (handler *consumeHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (handler *consumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	i := 0
	for msg := range claim.Messages() {
		// 处理消息
		handler.handle(msg)
		// 标记消息已经被消费
		session.MarkMessage(msg, "")

		// 若设置了手动提交offset，即：offset_auto_commit_enabled: true,需要添加以下代码进行手动提交
		if !offsetAutoCommitEnabled {
			fmt.Println(msg.Offset)
			i++
			// 每10条消息提交一次offset
			if i%10 == 0 {
				session.Commit()
			}
		}
	}
	return nil
}

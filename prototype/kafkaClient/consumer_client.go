package kafkaClient

import (
	"sync"

	"github.com/Shopify/sarama"
)

type ConsumerClient interface {
	sarama.Consumer
	ConsumeTopic(topic string, isNew int64, ConsumerFunc func(message *sarama.ConsumerMessage)) error
}

type consumerClient struct {
	sarama.Consumer
}

func NewConsumerClient(cfg Config) (ConsumerClient, func(), error) {
	config := sarama.NewConfig()
	client, err := sarama.NewConsumer(cfg.ServerAddrList, config)
	if err != nil {
		return nil, nil, err
	}

	return &consumerClient{client}, func() {
		defer client.Close()
	}, nil
}

// ConsumeTopic 消费整个Topic
func (cli *consumerClient) ConsumeTopic(topic string, isNew int64, ConsumerFunc func(message *sarama.ConsumerMessage)) error {
	partitionList, err := cli.Partitions(topic) // 通过topic获取到所有的分区
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, partition := range partitionList {
		partitionConsumer, err := cli.ConsumePartition(topic, partition, isNew)
		if err != nil {
			return err
		}
		defer partitionConsumer.AsyncClose()

		wg.Add(1)
		go func(partitionConsumer sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				defer wg.Done()
				ConsumerFunc(msg)
			}
		}(partitionConsumer)
	}
	wg.Wait()
	return nil
}

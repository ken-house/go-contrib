package kafkaClient

import "github.com/Shopify/sarama"

// 若使用自定义分区器，需完善该方法
func customPartitionerFunc(topic string, message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	return 0, nil
}

func NewCustomPartitioner(topic string) sarama.Partitioner {
	return &CustomPartitioner{
		Topic:                 topic,
		CustomPartitionerFunc: customPartitionerFunc,
	}
}

// CustomPartitioner 自定义生产者分区策略
type CustomPartitioner struct {
	Topic                 string
	CustomPartitionerFunc func(topic string, message *sarama.ProducerMessage, numPartitions int32) (int32, error)
}

// Partition 根据自定义方法计算出消息在总分区数分配到哪个分区
func (p *CustomPartitioner) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	return p.CustomPartitionerFunc(p.Topic, message, numPartitions)
}

// RequiresConsistency 是否要求分区程序需要一致性
func (p *CustomPartitioner) RequiresConsistency() bool {
	return false
}

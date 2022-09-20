package kafkaClient

import "github.com/Shopify/sarama"

// Config kafka连接及配置信息
type Config struct {
	ServerAddrList []string       `json:"server_addr_list" mapstructure:"server_addr_list"` // kafka地址
	ProducerConfig ProducerConfig `json:"producer_config"  mapstructure:"producer_config"`  // 生产者配置
	ConsumerConfig ConsumerConfig `json:"consumer_config" mapstructure:"consumer_config"`   // 消费者配置
}

// ProducerConfig kafka生产者配置参数
type ProducerConfig struct {
	Ack                   int `json:"acks" mapstructure:"acks"`                                       // 应答类型 0  1 -1
	PartitionerType       int `json:"partitioner_type" mapstructure:"partitioner_type"`               // 分区算法
	FlushMessageNum       int `json:"flush_message_num" mapstructure:"flush_message_num"`             // 达到多少条消息才发送
	FlushMessageFrequency int `json:"flush_message_frequency" mapstructure:"flush_message_frequency"` // 达到多少秒消息才发送
}

// ConsumerConfig kafka消费者配置参数
type ConsumerConfig struct {
}

// 指定分区算法
func setPartition(config *sarama.Config, partitionerType int) {
	switch partitionerType {
	case 1: // 随机算法
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	case 2: // robin算法
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	case 3: // 按消息内容计算分区
		config.Producer.Partitioner = sarama.NewManualPartitioner
	default: // 按key计算分区
		config.Producer.Partitioner = sarama.NewHashPartitioner
	}
}

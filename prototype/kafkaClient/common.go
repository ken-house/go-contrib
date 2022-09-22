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
	Ack               int                     `json:"ack" mapstructure:"ack"`                               // 应答类型 0  1 -1
	PartitionerPolicy int                     `json:"partitioner_policy" mapstructure:"partitioner_policy"` // 分区算法
	BatchMessageNum   int                     `json:"batch_message_num" mapstructure:"batch_message_num"`   // 达到多少条消息才发送
	LingerMs          int                     `json:"linger_ms" mapstructure:"linger_ms"`                   // 达到多少秒消息才发送
	CompressionType   sarama.CompressionCodec `json:"compression_type" mapstructure:"compression_type"`     // 压缩方式
	RecordAccumulator int                     `json:"record_accumulator" mapstructure:"record_accumulator"` // 生产区缓冲区大小，单位为字节
}

// ConsumerConfig kafka消费者配置参数
type ConsumerConfig struct {
}

// 指定分区算法
func setPartitionPolicy(config *sarama.Config, partitionerPolicy int) {
	switch partitionerPolicy {
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

package sarama

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
	PartitionerPolicy string                  `json:"partitioner_policy" mapstructure:"partitioner_policy"` // 分区算法
	BatchMessageNum   int                     `json:"batch_message_num" mapstructure:"batch_message_num"`   // 达到多少条消息才发送
	LingerMs          int                     `json:"linger_ms" mapstructure:"linger_ms"`                   // 达到多少秒消息才发送（毫秒）
	CompressionType   sarama.CompressionCodec `json:"compression_type" mapstructure:"compression_type"`     // 压缩方式
	RecordAccumulator int                     `json:"record_accumulator" mapstructure:"record_accumulator"` // 生产区缓冲区大小，单位为字节
}

// ConsumerConfig kafka消费者配置参数
type ConsumerConfig struct {
	GroupId                  string `json:"group_id" mapstructure:"group_id"`                                       // 消费者组id
	BalanceStrategy          string `json:"balance_strategy" mapstructure:"balance_strategy"`                       // 分区分配算法
	FetchMinBytes            int32  `json:"fetch_min_bytes" mapstructure:"fetch_min_bytes"`                         // 每批次最小抓取字节数
	FetchMaxBytes            int32  `json:"fetch_max_bytes" mapstructure:"fetch_max_bytes"`                         // 每批次最大抓取字节数
	MaxWaitTimeMs            int64  `json:"max_wait_time_ms" mapstructure:"max_wait_time_ms"`                       // 一批数据发送的超时时间(毫秒)
	FromBeginning            bool   `json:"from_beginning" mapstructure:"from_beginning"`                           // 是否从头开始消费
	OffsetAutoCommitEnabled  bool   `json:"offset_auto_commit_enabled" mapstructure:"offset_auto_commit_enabled"`   // offset是否自动提交
	OffsetAutoCommitInterval int    `json:"offset_auto_commit_interval" mapstructure:"offset_auto_commit_interval"` // 自动提交offset的时间间隔(秒)
}

// 指定生产者分区算法
func setProducerPartitionPolicy(config *sarama.Config, partitionerPolicy string) {
	switch partitionerPolicy {
	case "random": // 随机算法
		config.Producer.Partitioner = sarama.NewRandomPartitioner
	case "robin": // robin算法
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	case "manual": // 按消息内容计算分区
		config.Producer.Partitioner = sarama.NewManualPartitioner
	case "hash": // 按key的hashcode计算分区
		config.Producer.Partitioner = sarama.NewHashPartitioner
	default: // 按key的hashcode计算分区
		config.Producer.Partitioner = sarama.NewHashPartitioner
	}
}

// 指定消费者分区算法
func setConsumerPartitionPolicy(config *sarama.Config, balanceStrategy string) {
	switch balanceStrategy {
	case sarama.RangeBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	case sarama.RoundRobinBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case sarama.StickyBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	default:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	}
}

package kafkaClient

import (
	"time"

	"github.com/Shopify/sarama"
)

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
	MessageMaxBytes   int                     `json:"message_max_bytes" mapstructure:"message_max_bytes"`   // 一条消息最大字节数
	CompressionType   sarama.CompressionCodec `json:"compression_type" mapstructure:"compression_type"`     // 压缩方式
	IdempotentEnabled bool                    `json:"idempotent_enabled" mapstructure:"idempotent_enabled"` // 是否开启事务幂等
	MaxOpenRequests   int                     `json:"max_open_requests" mapstructure:"max_open_requests"`   // 生产者sender线程最大缓存请求数
	RetryMax          int                     `json:"retry_max" mapstructure:"retry_max"`                   // 重试次数
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

// NewKafkaClient 创建一个kafka客户端
func NewKafkaClient(cfg Config, asyncProducer bool) (sarama.Client, error) {
	config := sarama.NewConfig()
	// 指定kafka版本 - 需根据实际kafka版本调整
	config.Version = sarama.V2_8_1_0
	// 开启幂等，保证数据不重复
	config.Producer.Idempotent = cfg.ProducerConfig.IdempotentEnabled
	// 生产者sender线程最大缓存请求数
	if cfg.ProducerConfig.MaxOpenRequests > 0 {
		config.Net.MaxOpenRequests = cfg.ProducerConfig.MaxOpenRequests
	}
	// 开启幂等需要设置重试次数
	if cfg.ProducerConfig.RetryMax > 0 {
		config.Producer.Retry.Max = cfg.ProducerConfig.RetryMax
	}
	// 指定应答方式
	config.Producer.RequiredAcks = sarama.RequiredAcks(cfg.ProducerConfig.Ack)
	// 设置达到多少条消息才发送到kafka，相当于batch.size(批次大小)
	config.Producer.Flush.Messages = cfg.ProducerConfig.BatchMessageNum
	// 设置间隔多少秒才发送到kafka，相当于linger.ms（等待时间ms）
	if cfg.ProducerConfig.LingerMs > 0 {
		config.Producer.Flush.Frequency = time.Duration(cfg.ProducerConfig.LingerMs) * time.Millisecond
	}
	// 一条消息的最大字节数
	if cfg.ProducerConfig.MessageMaxBytes > 0 {
		config.Producer.MaxMessageBytes = cfg.ProducerConfig.MessageMaxBytes
	}
	// 指定数据压缩方式
	config.Producer.Compression = cfg.ProducerConfig.CompressionType
	// 指定分区算法
	setProducerPartitionPolicy(config, cfg.ProducerConfig.PartitionerPolicy)
	// 成功交付的消息将在success channel返回 同步发送必须指定为true
	if !asyncProducer {
		config.Producer.Return.Successes = true
	}
	return sarama.NewClient(cfg.ServerAddrList, config)
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
	case "custom":
		config.Producer.Partitioner = NewCustomPartitioner
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

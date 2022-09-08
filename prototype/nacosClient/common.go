package nacosClient

import "github.com/nacos-group/nacos-sdk-go/v2/common/constant"

// Config Nacos连接配置
type Config struct {
	ServerIpList   []string `json:"server_ip_list" mapstructure:"server_ip_list"`
	ServerHttpPort uint64   `json:"server_http_port" mapstructure:"server_http_port"`
	ServerGrpcPort uint64   `json:"server_grpc_port" mapstructure:"server_grpc_port"`
	NamespaceId    string   `json:"namespace_id" mapstructure:"namespace_id"`
	Timeout        uint64   `json:"timeout" mapstructure:"timeout"`
	LogLevel       string   `json:"log_level" mapstructure:"log_level"`
	LogPath        string   `json:"log_path" mapstructure:"log_path"`
	CachePath      string   `json:"cache_path" mapstructure:"cache_path"`
	Group          string   `json:"group" mapstructure:"group"`
	DataId         string   `json:"data_id" mapstructure:"data_id"`
}

// 获取服务端配置
func getServerConfig(cfg Config) []constant.ServerConfig {
	serverConfigList := make([]constant.ServerConfig, 0, 10)
	for _, ip := range cfg.ServerIpList {
		serverConfigList = append(serverConfigList, *constant.NewServerConfig(ip, cfg.ServerHttpPort, constant.WithGrpcPort(cfg.ServerGrpcPort)))
	}
	return serverConfigList
}

// 获取客户端配置
func getClientConfig(cfg Config) constant.ClientConfig {
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(cfg.NamespaceId),
		constant.WithTimeoutMs(cfg.Timeout),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(cfg.LogPath),
		constant.WithCacheDir(cfg.CachePath),
		constant.WithLogLevel(cfg.LogLevel),
	)
	return clientConfig
}

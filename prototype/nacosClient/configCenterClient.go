package nacosClient

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type ConfigCenterClient interface {
	config_client.IConfigClient
}

type configCenterClient struct {
	config_client.IConfigClient
}

func NewConfigClient(cfg Config) (ConfigCenterClient, func(), error) {
	// Nacos服务端配置
	serverConfigList := getServerConfig(cfg)
	// Nacos客户端配置
	clientConfig := getClientConfig(cfg)

	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ServerConfigs: serverConfigList,
		ClientConfig:  &clientConfig,
	})

	if err != nil {
		panic(err)
	}

	return &configCenterClient{client}, func() {
		client.CloseClient()
	}, nil
}

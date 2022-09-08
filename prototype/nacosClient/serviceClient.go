package nacosClient

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type ServiceClient interface {
	naming_client.INamingClient
}

type serviceClient struct {
	naming_client.INamingClient
}

func NewServiceClient(cfg Config) (ServiceClient, func(), error) {
	// Nacos服务端配置
	serverConfigList := getServerConfig(cfg)
	// Nacos客户端配置
	clientConfig := getClientConfig(cfg)

	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ServerConfigs: serverConfigList,
		ClientConfig:  &clientConfig,
	})

	if err != nil {
		panic(err)
	}

	return &serviceClient{client}, func() {
		client.CloseClient()
	}, nil
}

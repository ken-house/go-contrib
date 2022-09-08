package nacosClient

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type ServiceClient interface {
	naming_client.INamingClient
	FindHealthInstanceAddress(clusters []string, serviceName string, groupName string) (string, error)
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

// FindHealthInstanceAddress 根据serviceId查找到一个健康的服务实例，获取其地址
func (cli *serviceClient) FindHealthInstanceAddress(clusters []string, serviceName string, groupName string) (string, error) {
	serviceInfo, err := cli.INamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		Clusters:    clusters,
		ServiceName: serviceName,
		GroupName:   groupName,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", serviceInfo.Ip, serviceInfo.Port), nil
}

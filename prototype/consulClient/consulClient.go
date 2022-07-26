package consulClient

import (
	"fmt"

	"errors"

	"github.com/hashicorp/consul/api/watch"

	consulApi "github.com/hashicorp/consul/api"
)

type ConsulClient interface {
	RegisterService(serviceName string, ip string, port int) error
	DeregisterService(serviceId string) error
	GetConfig(consulPath string) ([]byte, error)
	WatchConfig(addr string, consulPath string, OnChange func([]byte)) error
	FindHealthInstanceAddress(serviceName string) (string, error)
}

type consulClient struct {
	Client *consulApi.Client
}

type ConsulConfig struct {
	Host string `json:"host" mapstructure:"host"`
	Port int    `json:"port" mapstructure:"port"`
}

func NewClient(config ConsulConfig) (ConsulClient, error) {
	cfg := consulApi.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", config.Host, config.Port)
	cli, err := consulApi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &consulClient{
		Client: cli,
	}, nil
}

// RegisterService 注册服务
func (cli *consulClient) RegisterService(serviceName string, ip string, port int) error {
	// 服务健康检查
	check := &consulApi.AgentServiceCheck{
		Interval:                       "2s",
		Timeout:                        "10s",
		GRPC:                           fmt.Sprintf("%s:%d", ip, port),
		GRPCUseTLS:                     true,
		TLSSkipVerify:                  true,
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务
	}
	return cli.Client.Agent().ServiceRegister(&consulApi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, ip, port),
		Name:    serviceName,
		Tags:    []string{"my_grpc"},
		Port:    port,
		Address: ip,
		Check:   check,
	})
}

// DeregisterService 注销服务
func (cli *consulClient) DeregisterService(serviceId string) error {
	return cli.Client.Agent().ServiceDeregister(serviceId)
}

// FindHealthInstanceAddress 根据serviceId查找到一个健康的服务实例，获取其地址
func (cli *consulClient) FindHealthInstanceAddress(serviceName string) (string, error) {
	_, serviceList, err := cli.Client.Agent().AgentHealthServiceByName(serviceName)
	if err != nil {
		return "", err
	}
	for _, v := range serviceList {
		return v.Service.Address, nil
	}
	return "", errors.New("no found")
}

// GetConfig 获取配置
func (cli *consulClient) GetConfig(consulPath string) ([]byte, error) {
	kv, _, err := cli.Client.KV().Get(consulPath, nil)
	return kv.Value, err
}

// WatchConfig 监听配置
func (cli *consulClient) WatchConfig(addr string, consulPath string, OnChange func([]byte)) error {
	params := make(map[string]interface{})
	params["type"] = "key"
	params["key"] = consulPath
	w, err := watch.Parse(params)
	if err != nil {
		return err
	}
	w.Handler = func(u uint64, i interface{}) {
		kv := i.(*consulApi.KVPair)
		OnChange(kv.Value)
	}
	if err = w.Run(addr); err != nil {
		return err
	}
	return nil
}

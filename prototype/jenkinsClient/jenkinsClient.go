package jenkinsClient

import (
	"context"
	"github.com/bndr/gojenkins"
)

type JenkinsClient interface {
	GetQueue(ctx context.Context) (*gojenkins.Queue, error)
}

type jenkinsClient struct {
	*gojenkins.Jenkins
}

// JenkinsConfig Jenkins连接配置
type JenkinsConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

func NewJenkinsClient(cfg JenkinsConfig) (JenkinsClient, error) {
	ctx := context.Background()
	client, err := gojenkins.CreateJenkins(nil, cfg.Host, cfg.Username, cfg.Password).Init(ctx)
	if err != nil {
		return nil, err
	}
	return &jenkinsClient{
		Jenkins: client,
	}, nil
}

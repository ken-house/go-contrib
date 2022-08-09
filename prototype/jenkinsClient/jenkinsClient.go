package jenkinsClient

import (
	"context"
	"github.com/bndr/gojenkins"
	"github.com/pkg/errors"
)

type JenkinsClient interface {
	GetQueueTaskIdList() (taskIdList []int64, err error)
}

type jenkinsClient struct {
	Ctx    context.Context
	Client *gojenkins.Jenkins
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
		Ctx:    ctx,
		Client: client,
	}, nil
}

// GetQueueTaskIdList 获取队列中的任务ID
func (cli *jenkinsClient) GetQueueTaskIdList() (taskIdList []int64, err error) {
	taskIdList = make([]int64, 0, 100)
	queueList, err := cli.Client.GetQueue(cli.Ctx)
	if err != nil {
		return taskIdList, err
	}
	if len(queueList.Raw.Items) == 0 {
		return taskIdList, errors.New("queue empty")
	}

	for _, v := range queueList.Raw.Items {
		taskIdList = append(taskIdList, v.ID)
	}
	return taskIdList, nil
}

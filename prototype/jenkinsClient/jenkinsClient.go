package jenkinsClient

import (
	"context"
	"github.com/bndr/gojenkins"
)

type JenkinsClient interface {
	Info(ctx context.Context) (*gojenkins.ExecutorResponse, error)
	SafeRestart(ctx context.Context) error
	GetQueue(ctx context.Context) (*gojenkins.Queue, error)
	GetQueueUrl() string
	GetQueueItem(ctx context.Context, id int64) (*gojenkins.Task, error)
	GetAllJobs(ctx context.Context) ([]*gojenkins.Job, error)
	GetAllJobNames(ctx context.Context) ([]gojenkins.InnerJob, error)
	GetJob(ctx context.Context, id string, parentIDs ...string) (*gojenkins.Job, error)
	GetSubJob(ctx context.Context, parentId string, childId string) (*gojenkins.Job, error)
	CreateJob(ctx context.Context, config string, options ...interface{}) (*gojenkins.Job, error)
	UpdateJob(ctx context.Context, job string, config string) *gojenkins.Job
	RenameJob(ctx context.Context, job string, name string) *gojenkins.Job
	CopyJob(ctx context.Context, copyFrom string, newName string) (*gojenkins.Job, error)
	DeleteJob(ctx context.Context, name string) (bool, error)
	BuildJob(ctx context.Context, name string, params map[string]string) (int64, error)
	GetBuildFromQueueID(ctx context.Context, queueid int64) (*gojenkins.Build, error)
	GetAllNodes(ctx context.Context) ([]*gojenkins.Node, error)
	GetNode(ctx context.Context, name string) (*gojenkins.Node, error)
	CreateNode(ctx context.Context, name string, numExecutors int, description string, remoteFS string, label string, options ...interface{}) (*gojenkins.Node, error)
	DeleteNode(ctx context.Context, name string) (bool, error)
	GetFolder(ctx context.Context, id string, parents ...string) (*gojenkins.Folder, error)
	CreateFolder(ctx context.Context, name string, parents ...string) (*gojenkins.Folder, error)
	CreateJobInFolder(ctx context.Context, config string, jobName string, parentIDs ...string) (*gojenkins.Job, error)
	GetLabel(ctx context.Context, name string) (*gojenkins.Label, error)
	GetAllBuildIds(ctx context.Context, job string) ([]gojenkins.JobBuild, error)
	GetBuild(ctx context.Context, jobName string, number int64) (*gojenkins.Build, error)
	GetArtifactData(ctx context.Context, id string) (*gojenkins.FingerPrintResponse, error)
	GetPlugins(ctx context.Context, depth int) (*gojenkins.Plugins, error)
	UninstallPlugin(ctx context.Context, name string) error
	HasPlugin(ctx context.Context, name string) (*gojenkins.Plugin, error)
	InstallPlugin(ctx context.Context, name string, version string) error
	ValidateFingerPrint(ctx context.Context, id string) (bool, error)
	GetAllViews(ctx context.Context) ([]*gojenkins.View, error)
	GetView(ctx context.Context, name string) (*gojenkins.View, error)
	CreateView(ctx context.Context, name string, viewType string) (*gojenkins.View, error)
	Poll(ctx context.Context) (int, error)
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

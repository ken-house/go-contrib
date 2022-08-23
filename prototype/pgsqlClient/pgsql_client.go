package pgsqlClient

import (
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

// PgsqlClient 单实例客户端
type PgsqlClient interface {
	xorm.EngineInterface
	Transaction(f func(*xorm.Session) (interface{}, error)) (interface{}, error)
	GetEngine() *xorm.Engine
}

type pgsqlClient struct {
	*xorm.Engine
}

func NewClient(cfg PgsqlConfig) (PgsqlClient, func(), error) {
	engine, err := newEngine(cfg)
	if err != nil {
		return nil, nil, err
	}
	client := &pgsqlClient{Engine: engine}
	return client, func() {
		_ = client.Close()
	}, nil
}

func (cli *pgsqlClient) GetEngine() *xorm.Engine {
	return cli.Engine
}

func (cli *pgsqlClient) Transaction(f func(*xorm.Session) (interface{}, error)) (interface{}, error) {
	return cli.Engine.Transaction(f)
}

type PgsqlConfig struct {
	MaxIdle     int    `json:"max_idle" mapstructure:"max_idle"`
	MaxOpen     int    `json:"max_open" mapstructure:"max_open"`
	MaxLifetime int    `json:"max_lifetime" mapstructure:"max_lifetime"`
	Dsn         string `json:"dsn" mapstructure:"dsn"`
}

func newEngine(cfg PgsqlConfig) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("postgres", cfg.Dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := engine.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}
	engine.ShowSQL(true)
	if cfg.MaxIdle > 0 {
		engine.SetMaxIdleConns(cfg.MaxIdle)
	}
	if cfg.MaxOpen > 0 {
		engine.SetMaxOpenConns(cfg.MaxOpen)
	}
	if cfg.MaxLifetime > 0 {
		engine.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)
	}
	return engine, nil
}

package redisClient

import (
	"context"

	"github.com/pkg/errors"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	redis.UniversalClient
}

type redisClient struct {
	redis.UniversalClient
}

func NewClient(cfg RedisConfig) (RedisClient, func(), error) {
	client, err := NewEngine(cfg)
	if err != nil {
		return nil, nil, err
	}
	sc := &redisClient{UniversalClient: client}
	return sc, func() {
		client.Close()
	}, nil
}

type RedisConfig struct {
	Addr     string `json:"addr" mapstructure:"addr"`
	Password string `json:"password" mapstructure:"password"`
	DB       int    `json:"db" mapstructure:"db"`
	PoolSize int    `json:"pool_size" mapstructure:"pool_size"`
}

func NewEngine(cfg RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client, err
}

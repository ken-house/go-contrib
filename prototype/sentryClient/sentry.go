package sentryClient

import (
	"github.com/getsentry/sentry-go"
	"github.com/ken-house/go-contrib/utils/env"
)

func InitSentry(cfg SentryConfig) error {
	err := sentry.Init(
		sentry.ClientOptions{
			Dsn:              cfg.Dsn,
			Debug:            !env.IsReleasing(),            // 线上环境为false 其他环境为true
			Transport:        sentry.NewHTTPSyncTransport(), // 同步发送到sentry
			SampleRate:       cfg.SampleRate,
			TracesSampleRate: cfg.TracesSampleRate,
			AttachStacktrace: cfg.AttachStacktrace,
			IgnoreErrors:     cfg.IgnoreErrors,
			ServerName:       cfg.ServerName,
			Environment:      env.Mode(),
		},
	)

	return err
}

type SentryConfig struct {
	Dsn              string   `json:"dsn" mapstructure:"dsn"`
	ServerName       string   `json:"server_name" mapstructure:"server_name"`
	SampleRate       float64  `json:"sample_rate" mapstructure:"sample_rate"`
	AttachStacktrace bool     `json:"attach_stacktrace" mapstructure:"attach_stacktrace"`
	TracesSampleRate float64  `json:"traces_sample_rate" mapstructure:"traces_sample_rate"`
	IgnoreErrors     []string `json:"ignore_errors" mapstructure:"ignore_errors"`
}

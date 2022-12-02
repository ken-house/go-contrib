package sentryClient

import (
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/ken-house/go-contrib/utils/env"
)

type SentryClient interface {
	CaptureException(err error)
	CaptureMessage(message string)
	Flush(duration time.Duration)
	CaptureExceptionForGin(ctx *gin.Context, err error)
	CaptureMessageForGin(ctx *gin.Context, message string)
	SentryMiddlewareForGin() gin.HandlerFunc
}

type sentryClient struct {
}

func NewSentryClient(cfg SentryConfig) (SentryClient, func(), error) {
	err := initSentry(cfg)
	return &sentryClient{}, func() {
		sentry.Flush(time.Second)
	}, err
}

// 初始化Sentry
func initSentry(cfg SentryConfig) error {
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

// CaptureException 捕获异常错误
func (cli *sentryClient) CaptureException(err error) {
	sentry.CaptureException(err)
}

// CaptureMessage 捕获异常信息
func (cli *sentryClient) CaptureMessage(message string) {
	sentry.CaptureMessage(message)
}

// Flush 刷新sentry缓存
func (cli *sentryClient) Flush(duration time.Duration) {
	sentry.Flush(duration)
}

// CaptureExceptionForGin Gin捕获自定义错误
func (cli *sentryClient) CaptureExceptionForGin(ctx *gin.Context, err error) {
	if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
		hub.CaptureException(err)
	}
}

// CaptureMessageForGin  Gin捕获自定义信息
func (cli *sentryClient) CaptureMessageForGin(ctx *gin.Context, message string) {
	if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
		hub.CaptureMessage(message)
	}
}

// SentryMiddlewareForGin Gin中间件
func (cli *sentryClient) SentryMiddlewareForGin() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{Repanic: true})
}

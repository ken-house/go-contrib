package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ken-house/go-contrib/prototype/sentryClient"
)

func main() {
	cfg := sentryClient.SentryConfig{
		Dsn:              "https://789e8b4d389e40c5994f6b09bd89d519@o435470.ingest.sentry.io/4504257001422848",
		ServerName:       "go_example",
		SampleRate:       1.0,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
		IgnoreErrors:     nil,
	}

	client, clean, err := sentryClient.NewSentryClient(cfg)
	if err != nil {
		panic(err)
	}
	defer clean()

	app := gin.Default()

	app.Use(client.SentryMiddlewareForGin())

	app.GET("/foo", func(ctx *gin.Context) {
		// painc捕获
		//panic("y tho222")

		// 自定义错误捕获
		client.CaptureMessage("自定义错误捕获")
	})

	_ = app.Run(":3000")
}

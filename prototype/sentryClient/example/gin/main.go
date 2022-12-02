package main

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	_ = sentry.Init(sentry.ClientOptions{
		Dsn:              "https://789e8b4d389e40c5994f6b09bd89d519@o435470.ingest.sentry.io/4504257001422848",
		Debug:            true,
		AttachStacktrace: true,
	})

	app := gin.Default()

	app.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	app.GET("/foo", func(ctx *gin.Context) {
		// painc捕获
		panic("y tho222")

		// 自定义错误捕获
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.CaptureMessage("自定义错误捕获")
		}
	})

	_ = app.Run(":3000")
}

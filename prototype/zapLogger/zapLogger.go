package zapLogger

import (
	"github.com/ken-house/go-contrib/utils/env"
	"go.uber.org/zap"
	"log"
)

func init() {
	var logger *zap.Logger
	var err error
	if env.IsReleasing() {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalln(err)
	}

	// 注册全局的单例的logger
	zap.ReplaceGlobals(logger)
	// 改变全局的标准库的log的输出，将其通过zap.Logger来输出
	zap.RedirectStdLog(logger)
}

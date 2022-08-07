package zapLogger

import (
	"log"

	"github.com/ken-house/go-contrib/utils/tools"

	"github.com/ken-house/go-contrib/utils/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SimpleLogger 使用zap包自带的配置文件
func SimpleLogger(outPutPaths []string) {
	var logger *zap.Logger
	var err error

	config := zap.NewProductionConfig()

	// 增加自定义日志记录位置
	if len(outPutPaths) > 0 {
		outPutPathsArr := config.OutputPaths
		for _, filePath := range outPutPaths {
			_, err = tools.FileNotExistAndCreate(filePath)
			if err != nil {
				log.Fatalln(err)
			}
			outPutPathsArr = append(outPutPathsArr, filePath)
		}

		config.OutputPaths = outPutPathsArr
		config.ErrorOutputPaths = outPutPathsArr
	}

	// 更改时间编码
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig = encoderConfig

	// 保证生产环境和其他环境日志存储格式一致，仅日志等级不同
	if env.IsReleasing() {
		logger, err = config.Build()
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.Development = true

		logger, err = config.Build()
	}
	if err != nil {
		log.Fatalln(err)
	}
	defer logger.Sync()

	// 注册全局的单例的logger
	zap.ReplaceGlobals(logger)
	// 改变全局的标准库的log的输出，将其通过zap.Logger来输出
	zap.RedirectStdLog(logger)
}

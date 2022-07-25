package zapLogger

import (
	"github.com/ken-house/go-contrib/utils/env"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

// CustomLogger 自定义zap日志，支持日志切割归档
func CustomLogger(lumberjackLogger lumberjack.Logger, outPutFile string) {
	encoder := getEncoder()
	writeSyncer := getWriteSyncer(&lumberjackLogger, outPutFile)

	logLevel := zapcore.DebugLevel
	if env.IsReleasing() {
		logLevel = zapcore.InfoLevel
	}
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	logger := zap.New(core, zap.AddCaller())

	defer logger.Sync()

	// 注册全局的单例的logger
	zap.ReplaceGlobals(logger)
	// 改变全局的标准库的log的输出，将其通过zap.Logger来输出
	zap.RedirectStdLog(logger)
}

// 获取编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 日志写入目标，使用lumberjack进行日志切割
func getWriteSyncer(lumberjackLogger *lumberjack.Logger, outPutFile string) zapcore.WriteSyncer {
	if lumberjackLogger != nil {
		return zapcore.AddSync(lumberjackLogger)
	} else {
		// todo 确保目录存在，不存在则创建目录
		if outPutFile == "" {
			outPutFile = "./log/test.log"
		}
		file, err := os.Create(outPutFile)
		if err != nil {
			log.Fatalln(err)
		}
		return zapcore.AddSync(file)
	}
}

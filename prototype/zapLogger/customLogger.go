package zapLogger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ken-house/go-contrib/utils/tools"

	"github.com/ken-house/go-contrib/utils/env"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CustomLogger 自定义zap日志，支持日志切割归档
func CustomLogger(lumberjackLogger *lumberjack.Logger, outPutFile string) {
	encoder := getEncoder()
	writeSyncer := getWriteSyncer(lumberjackLogger, outPutFile)

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
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 日志写入目标，使用lumberjack进行日志切割
func getWriteSyncer(lumberjackLogger *lumberjack.Logger, outPutFile string) zapcore.WriteSyncer {
	if lumberjackLogger != nil { // 使用lumberjack进行日志切割
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger))
	} else { // 不使用切割
		// 确保目录存在，不存在则创建目录
		if outPutFile == "" {
			outPutFile = fmt.Sprintf("./logs/log_%s.log", time.Now().Format("20060102"))
		}
		file, err := tools.FileNotExistAndCreate(outPutFile)
		if err != nil {
			log.Fatalln(err)
		}
		return zapcore.AddSync(file)
	}
}

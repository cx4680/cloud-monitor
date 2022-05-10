package logger

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	logger           *zap.SugaredLogger
	once             sync.Once
	defaultDirectory = "../logs"
)

func Logger() *zap.SugaredLogger {
	once.Do(func() {
		InitLogger(config.GetConfig().LogConfig)
	})
	return logger
}
func InitLogger(cfg *config.LogConfig) {
	encoder := getEncoder()

	logWriter, logLevel := getLogLevelAndWriter(cfg)

	errorWriter := getLogWriter("error", cfg)
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, getWriteSyncer(cfg, logWriter), logLevel),
		zapcore.NewCore(encoder, getWriteSyncer(cfg, errorWriter), errorLevel),
	)

	sugarLogger := zap.New(core, zap.AddCaller())
	logger = sugarLogger.Sugar()
	zap.ReplaceGlobals(sugarLogger)
}

func getWriteSyncer(cfg *config.LogConfig, logWriter zapcore.WriteSyncer) zapcore.WriteSyncer {
	if cfg.Stdout {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logWriter))
	} else {
		return zapcore.AddSync(logWriter)
	}
}

func getLogLevelAndWriter(cfg *config.LogConfig) (zapcore.WriteSyncer, zap.LevelEnablerFunc) {
	if cfg.Debug {
		logLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.DebugLevel && lvl < zapcore.ErrorLevel
		})
		logWriter := getLogWriter("debug", cfg)
		return logWriter, logLevel
	} else {
		logLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
		})
		logWriter := getLogWriter("info", cfg)
		return logWriter, logLevel
	}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(level string, cfg *config.LogConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		MaxSize:    100,          // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,           // 日志文件最多保存多少个备份
		MaxAge:     7,            // 文件最多保存多少天
		Compress:   cfg.Compress, // 是否压缩
	}

	fileName := "/application"
	if level == "error" {
		fileName = "/error"
	}

	directory := defaultDirectory
	if cfg.Directory != "" {
		directory = cfg.Directory
	}

	lumberJackLogger.Filename = directory + fileName + ".log"

	if cfg.MaxSize > 0 {
		lumberJackLogger.MaxSize = cfg.MaxSize
	}
	if cfg.MaxBackups > 0 {
		lumberJackLogger.MaxBackups = cfg.MaxBackups
	}
	if cfg.MaxAge > 0 {
		lumberJackLogger.MaxAge = cfg.MaxAge
	}
	return zapcore.AddSync(lumberJackLogger)
}

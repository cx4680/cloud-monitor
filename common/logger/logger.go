package logger

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	logger               *zap.SugaredLogger
	once                 sync.Once
	defaultApp           = "ucmp"
	defaultGroup         = "cec"
	defaultDataLogPrefix = "/data/logs/"
)

func Logger() *zap.SugaredLogger {
	once.Do(func() {
		InitLogger(&config.GetConfig().Logger)
	})
	return logger
}
func InitLogger(cfg *config.LogConfig) {
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	infoWriter := getLogWriter("info", cfg)
	errorWriter := getLogWriter("error", cfg)

	encoder := getEncoder()

	var core zapcore.Core

	if cfg.Debug {
		debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.DebugLevel
		})
		debugWriter := getLogWriter("debug", cfg)

		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
		)
	}

	sugarLogger := zap.New(core, zap.AddCaller())
	logger = sugarLogger.Sugar()
	zap.ReplaceGlobals(sugarLogger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(level string, cfg *config.LogConfig) zapcore.WriteSyncer {
	if !cfg.Stdout {
		lumberJackLogger := &lumberjack.Logger{
			MaxSize:    100,          // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: 10,           // 日志文件最多保存多少个备份
			MaxAge:     2,            // 文件最多保存多少天
			Compress:   cfg.Compress, // 是否压缩
		}

		if level == "error" {
			level = "common-error"
		} else {
			level = "digest_info"
		}

		fileName := defaultDataLogPrefix
		if cfg.DataLogPrefix != "" {
			fileName = cfg.DataLogPrefix
		} else {
			fileName = defaultDataLogPrefix
		}
		if cfg.Group != "" {
			fileName = fileName + cfg.Group + "_"
		} else {
			fileName = fileName + defaultGroup + "_"
		}
		if cfg.App != "" {
			fileName = fileName + cfg.App + "_"
		} else {
			fileName = fileName + defaultApp + "_"
		}
		fileName = fileName + level + ".log"

		lumberJackLogger.Filename = fileName

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

	return zapcore.AddSync(os.Stdout)
}

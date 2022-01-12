package logger

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

var (
	logger               *zap.SugaredLogger
	TrailLogger          *zap.Logger
	once                 sync.Once
	defaultDataLogPrefix = "/data/logs/"
	TimeFormatLayout     = "2006/01/02 15:04:05"
)

func Logger() *zap.SugaredLogger {
	once.Do(func() {
		InitLogger(config.Cfg.Logger)
	})
	return logger
}
func GetTrailLogger() *zap.Logger {
	once.Do(func() {
		InitLogger(config.Cfg.Logger)
	})
	return TrailLogger
}
func InitLogger(cfg config.LogConfig) {
	var core zapcore.Core
	core = newCore(cfg, core)

	sugarLogger := zap.New(core, zap.AddCaller())
	logger = sugarLogger.Sugar()

	//增加log tail
	coreTrail := newTrail(cfg)
	TrailLogger = zap.New(coreTrail)
}

func newCore(cfg config.LogConfig, core zapcore.Core) zapcore.Core {
	if cfg.Debug {
		debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.DebugLevel
		})
		debugWriter := getLogWriter("debug", cfg)

		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(getEncoderConfig("File")), zapcore.NewMultiWriteSyncer(zapcore.AddSync(debugWriter)), debugLevel),
			zapcore.NewCore(zapcore.NewConsoleEncoder(getEncoderConfig("Stdout")), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), debugLevel),
		)
	} else {
		infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.InfoLevel
		})
		infoWriter := getLogWriter("info", cfg)
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(getEncoderConfig("File")), zapcore.NewMultiWriteSyncer(zapcore.AddSync(infoWriter)), infoLevel),
			zapcore.NewCore(zapcore.NewConsoleEncoder(getEncoderConfig("Stdout")), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), infoLevel),
		)
	}
	return core
}

func newTrail(cfg config.LogConfig) zapcore.Core {
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	coreTrail := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(getEncoderConfig("File")), zapcore.NewMultiWriteSyncer(zapcore.AddSync(getLogWriter("operation_trail", cfg))), infoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(getEncoderConfig("Stdout")), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), infoLevel),
	)
	return coreTrail
}

func getEncoderConfig(writerType string) zapcore.EncoderConfig {
	EncodeLevel := zapcore.CapitalLevelEncoder
	if writerType == "Stdout" {
		EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return zapcore.EncoderConfig{
		FunctionKey:      "func",
		StacktraceKey:    "stack",
		NameKey:          "name",
		MessageKey:       "msg",
		LevelKey:         "level",
		ConsoleSeparator: " | ",
		EncodeLevel:      EncodeLevel,
		TimeKey:          "s",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(TimeFormatLayout))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeName: func(n string, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(n)
		},
	}
}

func getLogWriter(level string, cfg config.LogConfig) zapcore.WriteSyncer {
	if !cfg.Stdout {
		lumberJackLogger := &lumberjack.Logger{
			MaxSize:    100,          // 每个日志文件保存的最大尺寸 单位：M
			MaxBackups: 10,           // 日志文件最多保存多少个备份
			MaxAge:     15,           // 文件最多保存多少天
			Compress:   cfg.Compress, // 是否压缩
		}

		//if level == "error" {
		//	level = "common-error"
		//} else {
		//	level = "digest_info"
		//}

		fileName := defaultDataLogPrefix
		if cfg.DataLogPrefix != "" {
			fileName = cfg.DataLogPrefix
		} else {
			fileName = defaultDataLogPrefix
		}
		fileName = fileName + cfg.ServiceName + "/" + level + ".log"

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

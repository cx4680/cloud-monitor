package logger

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/config"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		Logger().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					Logger().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					Logger().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					Logger().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

package web

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/web/middleware"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/task"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var router = gin.New()

func GetRouter() *gin.Engine {
	return router
}

// Start /
func Start(cfg *config.Config) error {
	//加载中间件
	loadPlugin(cfg)
	//加载路由
	loadRouters()
	//加载定时任务
	cronFunc()
	//启动服务
	return doStart(cfg)
}

func doStart(cfg *config.Config) error {

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "path not found", nil)
	})

	port := "8080"
	if cfg.Serve.Port > 0 {
		port = strconv.Itoa(cfg.Serve.Port)
	}

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func loadPlugin(cfg *config.Config) {
	//加载全局
	if cfg.Serve.Debug {
		//router.Use(middleware.GinLogger())
	}
	//自定义组件
	router.Use(middleware.Cors())
	router.Use(middleware.Auth())
}

func cronFunc() {
	task.CronClear()
	task.CronInstanceJob()
}

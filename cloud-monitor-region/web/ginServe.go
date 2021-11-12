package web

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/web/middleware"
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
func Start(cfg config.Serve) error {
	//加载中间件
	loadPlugin(cfg)
	//加载路由
	loadRouters()
	//启动服务
	return doStart(cfg)
}

func doStart(cfg config.Serve) error {

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "path not found", nil)
	})

	port := "8081"
	if cfg.Port > 0 {
		port = strconv.Itoa(cfg.Port)
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

func loadPlugin(cfg config.Serve) {
	//加载全局
	if cfg.Debug {
		//router.Use(middleware.GinLogger())
	}
	//自定义组件
	router.Use(middleware.Cors())
	router.Use(middleware.Auth())
}

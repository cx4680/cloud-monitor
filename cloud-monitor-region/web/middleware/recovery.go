package middleware

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/openapi"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

// 提取当前运行时文件信息 名称 行数
func trace() []map[string]string {
	var pcs [32]uintptr
	n := runtime.Callers(5, pcs[:])
	var traceData []map[string]string
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		lineInfo := map[string]string{"file": file, "line": strconv.Itoa(line)}
		traceData = append(traceData, lineInfo)
		// @todo 只记录打印错误最近的一行信息 break 了
		//break
	}
	return traceData
}

// Recovery 请求异常处理
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				if openapi.OpenApiRouter(c) {
					c.JSON(http.StatusInternalServerError, openapi.NewRespError(openapi.SystemError, c))
					return
				}
				c.JSON(http.StatusInternalServerError, global.NewError("系统异常"))
				return
			}
		}()

		c.Next()
	}
}

package middleware

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/global"
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
				// 异常处理
				//switch err.(type) {
				//// 接管自定义的一些异常抛出
				//case http.StatusInternalServerError:
				log.Println(err)
				c.JSON(http.StatusInternalServerError, global.NewError("系统异常"))
				return
				//default:
				//	message := fmt.Sprintf("%s", err)
				//	// 提取整理一下
				//	traceData := trace()
				//	position := traceData[0]["file"] + ":" + traceData[0]["line"]
				//	ResponseData := gin.H{
				//		"code": 5000,
				//		"msg":  "Internal Server Error ",
				//		"result": map[string]interface{}{
				//			"error_info": message,
				//			"position":   position,
				//			"detail":     trace(),
				//		},
				//	}
				//	c.Set("ResponseData", ResponseData)
				//	c.JSON(500, ResponseData)
				//}
			}
			if c.Writer.Status() == 404 {
				c.JSON(http.StatusNotFound, "path not found")
			}

		}()

		c.Next()
	}
}

package logs

import (
	"bytes"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var userTypeMap = map[string]string{
			"0":"root-account",
			"1":"root-account",
			"2":"root-account",
			"3":"iam-user",
			"4":"assumed-role",
			"5":"system",
            }

func GinTrailzap(utc bool, requestType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName := "CloudMonitor"
		var data []byte
		if http.MethodGet == c.Request.Method {
			params := c.Request.URL.RawQuery
			requestParams, err := json.Marshal(params)
			if err != nil {
				c.Abort()
				logger.Logger().Error(err.Error())
				return
			}
			data = requestParams
		} else {
			body, err := c.GetRawData()
			if err != nil {
				c.Abort()
				logger.Logger().Error(err.Error())
				return
			}
			data = body
		}
		formatData := strings.Replace(strings.Replace(string(data), "\r\n", "", -1), " ", "", -1)
		requestParameters := make(map[string]string, 0)
		requestParameters["request"] = formatData
		eventRegion := "*"
		if strings.Contains(formatData, "regionCode") {
			arr1 := strings.Split(formatData, "regionCode")
			arr2 := strings.Split(arr1[1], "\"")
			eventRegion = arr2[2]
		}
		requestParamJson, _ := json.Marshal(requestParameters)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		requestID := c.GetHeader("X-Request-ID")
		ctx := context.WithValue(context.Background(), "X-Request-ID", requestID)
		c.Set("ctx", ctx)
		c.Next()
		defer func() {
			accountId := c.GetString(global.TenantId)
			userType := c.GetString(global.UserType)
			loginId := c.GetString(global.UserId)
			userName := c.GetString(global.UserName)
			resourceName := c.GetString("ResourceId")
			eventName := c.GetString("Action")
			source := c.Request.Header["Origin"]
			eventSource := ""
			if len(source) > 0 {
				eventSource = source[0]
			}

			end := time.Now()
			if utc {
				end = end.UTC()
			}
			resourceType := "CCS::CloudMonitor::Manager"

			logger.GetTrailLogger().Info("[ACTION_TRAIL_LOG]",
				zap.String("event_id", requestID),
				zap.String("event_version", "1"),
				zap.String("event_source", eventSource),
				zap.String("source_ip_address", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.String("service_name", serviceName),
				zap.String("event_name", eventName),
				zap.String("request_type", requestType),
				zap.String("api_version", "1.0"),
				zap.String("request_id", requestID),
				zap.String("event_time", end.Format(util.FullTimeFmt)),
				zap.String("event_region", eventRegion),
				zap.String("resource_type", resourceType),
				zap.String("resource_name", resourceName),
				zap.String("request_parameters", string(requestParamJson)),
				zap.Int("error_code", c.Writer.Status()),
				zap.Namespace("user_info"),
				zap.String("account_id", accountId),
				zap.String("type", userTypeMap[userType]),
				zap.String("user_name", userName),
				zap.String("principal_id", loginId),
			)
		}()
	}

}

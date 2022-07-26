package logs

import (
	"bufio"
	"bytes"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var userTypeMap = map[string]string{
	"0": "root-account",
	"1": "root-account",
	"2": "root-account",
	"3": "iam-user",
	"4": "assumed-role",
	"5": "system",
}

type EventLevel = string

const (
	INFO  EventLevel = "Info"
	Warn  EventLevel = "Warn"
	Fatal EventLevel = "Fatal"
)

type ResourceType = string

const (
	MonitorProduct    ResourceType = "MonitorProduct"
	MonitorItem       ResourceType = "MonitorItem"
	AlertContact      ResourceType = "AlertContact"
	AlertContactGroup ResourceType = "AlertContactGroup"
	AlertRule         ResourceType = "AlertRule"
	MonitorChart      ResourceType = "MonitorChart"
	MonitorReportForm ResourceType = "MonitorReportForm"
	ConfigItem        ResourceType = "ConfigItem"
	Instance          ResourceType = "Instance"
	Notice            ResourceType = "Notice"
	AlertRecord       ResourceType = "AlertRecord"
	Resource          ResourceType = "Resource"
	AlertRuleTemplate ResourceType = "AlertRuleTemplate"
)

func GinTrailzap(utc bool, requestType string, eventLevel EventLevel, resourceType ResourceType) gin.HandlerFunc {
	return func(c *gin.Context) {
		replaceResponseWriter(c)
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
		eventRegion := config.Cfg.Common.RegionName
		requestParamJson, _ := json.Marshal(requestParameters)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		requestID := ""
		if openapi.OpenApiRouter(c) {
			requestID = c.GetHeader("RequestId")
		} else {
			requestID = c.GetHeader("X-Request-ID")
		}
		if len(requestID) == 0 {
			if newUUID, err := uuid.NewUUID(); err == nil {
				requestID = newUUID.String()
			}
		}
		ctx := context.WithValue(context.Background(), "X-Request-ID", requestID)
		c.Set("ctx", ctx)
		c.Next()
		defer func() {
			accountId := c.GetString(global.TenantId)
			userType := c.GetString(global.UserType)
			loginId := c.GetString(global.UserId)
			userName := c.GetString(global.UserName)
			resourceName := c.GetString(global.ResourceName)
			eventName := c.GetString("Action")
			source := c.Request.Header["Origin"]
			eventSource := ""
			if len(source) > 0 {
				eventSource = source[0]
			} else {
				eventSource = c.Request.Host
			}
			if len(userName) == 0 {
				userName = getUserNameFromRemote(loginId)
			}

			end := time.Now()
			if utc {
				end = end.UTC()
			}
			var response string
			var errs map[string]interface{}
			if writer, ok := c.Writer.(*responseWriter); ok {
				response = string(writer.response)
				json.Unmarshal(writer.response, &errs)
			}
			resourceTypeNew := "CCS::CloudMonitor::" + resourceType
			result := getResult(c.Writer.Status())

			errMessage := ""
			errCode := ""
			if openapi.OpenApiRouter(c) {
				var errMap map[string]string
				err := jsonutil.ToObjectWithError(jsonutil.ToString(errs["Error"]), &errMap)
				if err != nil {
					errMessage = "error"
				} else {
					errMessage = errMap["Message"]

				}
				errCode = strconv.Itoa(c.Writer.Status())
			} else {
				errMessage = jsonutil.ToString(errs["errorMsg"])
				errCode = errs["errorCode"].(string)
			}

			var resError interface{}
			if err := recover(); err != nil {
				result = "Fail"
				errMessage = fmt.Sprint(err)
				resError = err
			}
			var eventId string
			if newUUID, err := uuid.NewUUID(); err == nil {
				eventId = newUUID.String()
			}
			logger.GetTrailLogger().Info("[ACTION_TRAIL_LOG]",
				zap.String("event_id", eventId),
				zap.String("event_version", "1"),
				zap.String("event_source", eventSource),
				zap.String("source_ip_address", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.String("service_name", serviceName),
				zap.String("event_name", eventName),
				zap.String("request_type", requestType),
				zap.String("api_version", "1.0"),
				zap.String("event_level", eventLevel),
				zap.String("request_id", requestID),
				zap.String("event_time", end.Format(util.FullTimeFmt)),
				zap.String("event_region", eventRegion),
				zap.String("resource_type", resourceTypeNew),
				zap.String("resource_name", resourceName),
				zap.String("request_parameters", string(requestParamJson)),
				zap.String("result", result),
				zap.String("response_elements", response),
				zap.String("error_code", errCode),
				zap.String("error_message", errMessage),
				zap.Namespace("user_info"),
				zap.String("account_id", accountId),
				zap.String("type", userTypeMap[userType]),
				zap.String("user_name", userName),
				zap.String("principal_id", loginId),
				zap.String("access_key_id", "-"),
			)
			if resError != nil {
				panic(resError)
			}
		}()
	}

}

func getUserNameFromRemote(loginId string) string {
	params := struct {
		LoginId string `json:"loginId"`
	}{
		LoginId: loginId,
	}

	resp, err := httputil.HttpPostJson(config.Cfg.Common.AccountApiHost+"/api/outer/userinfo/login-info", params, nil)
	if err != nil {
		logger.Logger().Errorf("getUserNameFromRemote error, %v", err)
		return ""
	}

	type userInfo struct {
		LoginCode string `json:"loginCode"`
		LoginId   string `json:"loginId"`
	}

	type respObj struct {
		ErrorMsg  string    `json:"errorMsg"`
		ErrorCode string    `json:"errorCode"`
		Success   bool      `json:"success"`
		Module    *userInfo `json:"module"`
	}

	var respMap respObj
	err = jsonutil.ToObjectWithError(resp, &respMap)
	if err != nil {
		logger.Logger().Errorf("getUserNameFromRemote error, serialization, %v", err)
		return ""
	}
	if respMap.Module == nil {
		return ""
	}
	return respMap.Module.LoginCode
}

func getResult(status int) string {
	if status != 200 {
		return "Fail"
	}
	return "Success"
}

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

func replaceResponseWriter(c *gin.Context) {
	writer := &responseWriter{
		ResponseWriter: c.Writer,
		size:           noWritten,
		status:         defaultStatus,
	}
	c.Writer = writer
}

type responseWriter struct {
	http.ResponseWriter
	size     int
	status   int
	response []byte
}

func (w *responseWriter) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.size = noWritten
	w.status = defaultStatus
}

func (w *responseWriter) WriteHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			print("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code)
		}
		w.status = code
	}
}

func (w *responseWriter) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.response = data
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

func (w *responseWriter) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.size += n
	return
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Written() bool {
	return w.size != noWritten
}

// Hijack implements the http.Hijacker interface.
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotify interface.
func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Flush implements the http.Flush interface.
func (w *responseWriter) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *responseWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}

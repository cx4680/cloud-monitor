package openapi

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

const RequestId = "RequestId"
const locale = "locale"

type ResCommon struct {
	RequestId string
}
type RespError struct {
	ResCommon
	Error Error
}

type Error struct {
	Code    string
	Message string
}

type ResCommonPage struct {
	ResCommon
	TotalCount int
	PageSize   int
	PageNumber int
}

func NewRespError(error *ErrorCode, ctx *gin.Context) *RespError {
	c2 := ctx.Copy()
	requestId := c2.GetHeader(RequestId)
	language := c2.DefaultQuery(locale, "zh")
	var message string
	if language == "zh" {
		message = error.MessageCn
	}
	return &RespError{
		ResCommon: ResCommon{
			RequestId: requestId,
		},
		Error: Error{
			Code:    error.Code,
			Message: message,
		},
	}
}

type ResSuccess struct {
	ResCommon
}

func NewResSuccess(ctx *gin.Context) *ResSuccess {
	return &ResSuccess{
		ResCommon: ResCommon{
			RequestId: ctx.GetHeader(RequestId),
		},
	}
}

func GetErrorCode(err error) *ErrorCode {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range errs {
			if fieldError.ActualTag() == "required" {
				return MissingParameter
			}
		}
	}
	return InvalidParameter
}

func NewResCommonPage(c *gin.Context, pageVo *vo.PageVO) *ResCommonPage {
	return &ResCommonPage{
		ResCommon: ResCommon{
			RequestId: GetRequestId(c),
		},
		TotalCount: pageVo.Total,
		PageSize:   pageVo.Size,
		PageNumber: pageVo.Current,
	}
}

func OpenApiRouter(c *gin.Context) bool {
	uri := c.Request.RequestURI
	return !strings.HasPrefix(uri, "/hawkeye/")
}

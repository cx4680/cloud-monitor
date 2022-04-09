package openapi

import (
	"github.com/gin-gonic/gin"
)

type PageQuery struct {
	PageNumber int `binding:"min=1"`
	PageSize   int `binding:"min=1,max=100"`
}

func NewPageQuery() *PageQuery {
	return &PageQuery{PageNumber: 1, PageSize: 10}
}

func GetRequestId(ctx *gin.Context) string {
	return ctx.GetHeader(RequestId)
}

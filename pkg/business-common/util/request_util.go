package util

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestId = "requestId"

func GenerateRequest(ctx *gin.Context) *context.Context {
	c := context.WithValue(ctx, requestId, uuid.New().String())
	return &c
}

func GetRequestId(c *context.Context) string {
	return (*c).Value(requestId).(string)
}

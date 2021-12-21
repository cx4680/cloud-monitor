package errortypes

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/constants/autherrorenum"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/logger"
	"fmt"
)

type AuthRuntimeError struct {
	Code    string `json:"errorCode"`
	Message string `json:"errorMsg"`
}

func (e AuthRuntimeError) Error() string {
	return fmt.Sprintf("code: %s msg:%s", e.Code, e.Message)
}

func IAMSDKError(code string) error {
	msg := autherrorenum.ErrorText(code)
	logger.Logger().Errorf("【IAM SDK】 IAMSDKError code:%s Message:%s", code, msg)
	return &AuthRuntimeError{Code: code, Message: msg}
}

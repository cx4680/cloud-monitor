package handler

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/config"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/constants/autherrorenum"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/domain"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/errortypes"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetToken(crn string, ssoSid string) (string, error) {
	sidMap := &map[string]string{"sid": ssoSid, "crn": crn}
	requestObject, err := json.Marshal(sidMap)
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 GetToken object to json failed, err:%v", err)
		return "", errortypes.IAMSDKError(autherrorenum.JsonFormatExceptionGetTokenRequestParam)
	}
	logger.Logger().Infof("【IAM SDK】 GetToken 获得令牌请求:%s", string(requestObject))

	resp, err := http.Post(getTokenUrl(config.GetConfig().AuthSdkConfig.AuthRequestSite), "application/json", strings.NewReader(string(requestObject)))
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 GetToken 获得令牌响应失败,err:%v", err)
		return "", errortypes.IAMSDKError(autherrorenum.RequestFailGetToken)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 GetToken 获得令牌响应读取失败,%v", err)
		return "", errortypes.IAMSDKError(autherrorenum.IoReadExceptionGetToken)
	}

	logger.Logger().Infof("【IAM SDK】 GetToken 获得令牌响应:%s", string(b))

	result := &RoleToken{}
	err = json.Unmarshal(b, result)
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 GetToken 获得令牌响应解析失败, err:%v", err)
		return "", errortypes.IAMSDKError(autherrorenum.JsonFormatExceptionGetTokenResponse)
	}

	if !result.Success || result.Module == nil {
		logger.Logger().Errorf("【IAM SDK】 GetToken 获得令牌失败:%s", result.ErrorCode)
		if autherrorenum.IamRoleStsTokenInvalid == result.ErrorCode {
			return "", errortypes.IAMSDKError(autherrorenum.StsInvokeError)
		} else {
			return "", errortypes.IAMSDKError(autherrorenum.ActionNotAllowedGetToken)
		}
	}

	return result.Module.(string), nil
}

type RoleToken struct {
	domain.Response
}

func getTokenUrl(url string) string {
	builder := strings.Builder{}
	builder.WriteString(url)
	builder.WriteString("/role/token")
	return builder.String()
}

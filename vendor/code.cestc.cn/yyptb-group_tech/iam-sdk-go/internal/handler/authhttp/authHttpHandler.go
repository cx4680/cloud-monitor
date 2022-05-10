package authhttp

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/config"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/constants/autherrorenum"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/errortypes"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type UserAuthRequest struct {
	AccountUserId string   `json:"accountUserId"`
	Product       string   `json:"product"`
	ActionName    string   `json:"actionName"`
	Resource      []string `json:"resource"`
}

type RoleAuthRequest struct {
	RoleCrn    string   `json:"roleCrn"`
	Product    string   `json:"product"`
	Token      string   `json:"token"`
	ActionName string   `json:"actionName"`
	Resource   []string `json:"resource"`
}

// UserAuth IAM 用户鉴权
func (authRequest *UserAuthRequest) UserAuth() (*AuthResponse, error) {
	requestObject, err := json.Marshal(*authRequest)
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 UserAuth  object to json failed, err:%v", err)
		return nil, errortypes.IAMSDKError(autherrorenum.JsonFormatExceptionUserAuthRequestParam)
	}
	logger.Logger().Infof("【IAM SDK】 UserAuth IAM用户鉴权请求:%s", string(requestObject))

	resp, err := http.Post(getUserIdentityUrl(config.GetConfig().AuthSdkConfig.AuthRequestSite), "application/json", strings.NewReader(string(requestObject)))
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 UserAuth  IAM用户鉴权响应失败, err:%v", err)
		return nil, errortypes.IAMSDKError(autherrorenum.RequestFailUserAuth)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 UserAuth  IAM用户鉴权响应体读取失败, err:%v", err)
		return nil, errortypes.IAMSDKError(autherrorenum.IoReadExceptionUserAuth)
	}

	logger.Logger().Infof("【IAM SDK】 UserAuth IAM用户鉴权响应:%s", string(b))

	result := &AuthResponse{}
	err = json.Unmarshal(b, result)
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 UserAuth IAM用户鉴权响应体解析失败, err:%v", err)
		return nil, errortypes.IAMSDKError(autherrorenum.JsonFormatExceptionUserAuthResponse)
	}

	return result, nil
}

// RoleAuth IAM角色鉴权
func (authRequest *RoleAuthRequest) RoleAuth() (*AuthResponse, error) {
	requestObject, err := json.Marshal(*authRequest)
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 RoleAuth  object to json failed, err:%v", err)
		return nil, errortypes.IAMSDKError(autherrorenum.JsonFormatExceptionRoleAuthRequestParam)
	}
	logger.Logger().Infof("【IAM SDK】 RoleAuth 角色鉴权请求:%s", string(requestObject))

	resp, err := http.Post(getRoleIdentityUrl(config.GetConfig().AuthSdkConfig.AuthRequestSite), "application/json", strings.NewReader(string(requestObject)))
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 RoleAuth  角色鉴权响应失败, err:%v", err)
		return nil, errortypes.IAMSDKError(autherrorenum.RequestFailRoleAuth)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 RoleAuth  角色鉴权响应体读取失败, err:%v", err)
		return nil, errortypes.IAMSDKError(autherrorenum.IoReadExceptionRoleAuth)
	}

	logger.Logger().Infof("【IAM SDK】 RoleAuth 角色鉴权响应:%s", string(b))

	result := &AuthResponse{}
	err = json.Unmarshal(b, result)
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 RoleAuth 角色鉴权响应体解析失败, err:%v", err)
		return nil, errortypes.IAMSDKError(autherrorenum.JsonFormatExceptionRoleAuthResponse)
	}

	return result, nil
}

func getUserIdentityUrl(url string) string {
	builder := strings.Builder{}
	builder.WriteString(url)
	builder.WriteString("/v2/identity")
	logger.Logger().Infof("【IAM SDK】 getUserIdentityUrl url:%s", builder.String())
	return builder.String()
}

func getRoleIdentityUrl(url string) string {
	builder := strings.Builder{}
	builder.WriteString(url)
	builder.WriteString("/v2/role/identity")
	logger.Logger().Infof("【IAM SDK】 getRoleIdentityUrl url:%s", builder.String())
	return builder.String()
}

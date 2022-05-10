package roleswitchrecordhandler

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/config"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/constants/autherrorenum"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/errortypes"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/internal/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetRoleRecord(sid string) (*IamRoleSwitchRecord, error) {
	sidMap := &map[string]string{"sid": sid}
	requestObj, err := json.Marshal(sidMap)
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 GetRoleRecord  JsonFormatExceptionGetRoleRecordRequestParam")
		return nil, errortypes.IAMSDKError(autherrorenum.JsonFormatExceptionGetRoleRecordRequestParam)
	}
	logger.Logger().Infof("【IAM SDK】 GetRoleRecord 获得角色切换记录请求:{sid:%s}", sid)

	resp, err := http.Post(getRoleRecordUrl(config.GetConfig().AuthSdkConfig.AuthRequestSite), "application/json", strings.NewReader(string(requestObj)))
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 GetRoleRecord  获得角色切换记录失败, err:%v", err)
		return nil, errortypes.IAMSDKError(autherrorenum.RequestFailGetRoleRecord)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger().Errorf("【IAM SDK】 GetRoleRecord  获得角色切换记录读取响应体失败, err:%v", err)
		return nil, errortypes.IAMSDKError(autherrorenum.IoReadExceptionGetRoleRecord)
	}

	logger.Logger().Infof("【IAM SDK】 GetRoleRecord 获得角色切换记录响应:%s", string(b))

	result := &IamRoleSwitchResponse{}
	errParseResult := json.Unmarshal(b, result)
	if errParseResult != nil {
		logger.Logger().Errorf("【IAM SDK】 GetRoleRecord 获得角色切换记录响应体解析失败, err:%v", errParseResult)
		return nil, errortypes.IAMSDKError(autherrorenum.JsonFormatExceptionGetRoleRecordResponse)
	}

	if !result.Success || result.Module == nil {
		logger.Logger().Errorf("【IAM SDK】 GetRoleRecord 获得角色切换记录失败:%s", result.ErrorCode)
		if autherrorenum.IamRoleStsTokenInvalid == result.ErrorCode {
			return nil, errortypes.IAMSDKError(autherrorenum.IamRoleStsTokenInvalid)
		} else {
			return nil, errortypes.IAMSDKError(autherrorenum.ActionNotAllowedGetRoleRecord)
		}
	}

	return result.Module, nil
}

func getRoleRecordUrl(url string) string {
	builder := strings.Builder{}
	builder.WriteString(url)
	builder.WriteString("/role/record")
	logger.Logger().Infof("【IAM SDK】 getRoleRecordUrl url:%s", builder.String())
	return builder.String()
}

type IamRoleSwitchResponse struct {
	ErrorMsg   string               `json:"errorMsg"`
	ErrorCode  string               `json:"errorCode"`
	Success    bool                 `json:"success"`
	Module     *IamRoleSwitchRecord `json:"module"`
	AllowRetry bool                 `json:"allowRetry"`
	ErrorArgs  interface{}          `json:"errorArgs"`
}

type IamRoleSwitchRecord struct {

	// 主键
	id int64

	/**
	 * 主键
	 */
	sid string

	/**
	 * 角色ID
	 */
	roleId int64

	//角色名称
	roleName string

	/**
	 * 策略类型
	 */
	strategy int64 `json:"type""`

	/**
	 * 最大会话时间
	 */
	durationSecond int64

	/**
	 * 扮演者归属云账号ID
	 */
	cloudAccount string

	/**
	 * 租户ID
	 */
	accountId int64

	/**
	 * 归属云账号用户名
	 */
	loginCode string

	/**
	 * 角色CRN
	 */
	RoleCrn string

	/**
	 * securityToken
	 */
	SecurityToken string

	/**
	 * 失效时间
	 */
	expiration string

	/**
	 * AK的id
	 */
	accessKeyId string

	/**
	 * 最后刷新时间
	 */
	lastRefreshTime time.Time
}

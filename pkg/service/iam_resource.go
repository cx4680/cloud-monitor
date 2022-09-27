package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"fmt"
	"time"
)

func GetIamUserByAccountUid(accountUid string) (string, error) {
	response, err := httputil.HttpPostJson(config.Cfg.Iam.IamUser, map[string]string{"accountUserId": accountUid}, nil)
	logger.Logger().Infof("IamUser:%v", response)
	if err != nil {
		logger.Logger().Errorf("获取iam部门错误：%v", err)
		return "", errors.NewBusinessError("获取iam部门错误")
	}
	var iamUser form.IamUser
	jsonutil.ToObject(response, &iamUser)
	return iamUser.Module.DirectoryId, nil
}

func GetIamResourcesIdList(iamUserId string) ([]string, error) {
	directoryId, err := GetIamUserByAccountUid(iamUserId)
	if err != nil {
		return nil, err
	}
	param := form.InstanceRequest{
		DirectoryIds: []string{directoryId},
		CurrPage:     "1",
		PageSize:     "99999",
	}
	response, err := httputil.HttpPostJson(config.Cfg.Common.Rc, param, nil)
	logger.Logger().Infof("IamResources:%v", response)
	if err != nil {
		logger.Logger().Errorf("获取实例列表错误：%v", err)
		return nil, errors.NewBusinessError("获取实例列表错误")
	}
	var result form.InstanceResponse
	jsonutil.ToObject(response, &result)
	var resourcesIdList []string
	for _, v := range result.Data.List {
		resourcesIdList = append(resourcesIdList, v.ResourceId)
	}
	return resourcesIdList, nil
}

func CheckIamDirectory(loginId string) (bool, error) {
	key := fmt.Sprintf(constant.TenantDirectoryKey, loginId)
	value, err := sys_redis.Get(key)
	if err != nil {
		logger.Logger().Error("key=" + loginId + ", error:" + err.Error())
	}
	var iamLoginStart form.IamLoginStart
	if strutil.IsNotBlank(value) {
		jsonutil.ToObject(value, &iamLoginStart)
		return iamLoginStart.Module, nil
	}
	response, err := httputil.HttpPostJson(config.Cfg.Iam.LoginModel, map[string]string{"loginId": loginId}, nil)
	if err != nil {
		logger.Logger().Errorf("判断是否开启云组织管理接口错误：%v", err)
		return false, errors.NewBusinessError("判断是否开启云组织管理接口错误")
	}
	var result form.IamLoginStart
	jsonutil.ToObject(response, &result)
	if e := sys_redis.SetByTimeOut(key, jsonutil.ToString(result), 5*time.Minute); e != nil {
		logger.Logger().Error("设置监控项缓存错误, key=" + loginId)
	}
	return result.Module, nil
}

func CheckIamLogin(tenantId, iamUserId string) bool {
	isOpen, err := CheckIamDirectory(tenantId)
	if err != nil {
		logger.Logger().Errorf("IamLogin接口错误：%v", err)
		return false
	}
	return isOpen && strutil.IsNotBlank(iamUserId) && iamUserId != tenantId
}

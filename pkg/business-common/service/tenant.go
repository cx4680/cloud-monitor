package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"fmt"
	"time"
)

const RedisTimeOutOneHour = 60 * 60

type TenantService struct {
}

func NewTenantService() *TenantService {
	return &TenantService{}
}

type tenantResponse struct {
	Result tenantResult `json:"result"`
	Module tenantResult `json:"module"`
}

type tenantResult struct {
	Login tenantLogin `json:"login"`
}

type tenantLogin struct {
	LoginCode    string `json:"loginCode"`
	SerialNumber string `json:"serialNumber"`
}

func (s *TenantService) GetTenantInfo(tenantId string) dto.TenantDTO {
	var tenant dto.TenantDTO
	key := fmt.Sprintf(constant.TenantInfoKey, tenantId)

	value, err := sys_redis.Get(key)
	if err != nil {
		logger.Logger().Info("key=" + key + ", error:" + err.Error())
	}
	if strutil.IsNotBlank(value) {
		jsonutil.ToObject(value, &tenant)
		return tenant
	}
	tenantFromRemote := getTenantFromServer(tenantId)
	if tenantFromRemote == nil {
		logger.Logger().Info("获取租户信息为空")
		return tenant
	}
	if e := sys_redis.SetByTimeOut(key, jsonutil.ToString(tenantFromRemote), RedisTimeOutOneHour*time.Second); e != nil {
		logger.Logger().Error("设置redis出错, key=" + key)
	}
	return *tenantFromRemote
}

func getTenantFromServer(tenantId string) *dto.TenantDTO {
	logger.Logger().Info("tenantId:", tenantId)
	m := make(map[string]string, 1)
	m["loginId"] = tenantId
	resp, err := httputil.HttpPostJson(config.Cfg.Common.TenantUrl, m, nil)
	if err != nil {
		logger.Logger().Error("查询租户信息失败:", resp, err)
		return nil
	}
	var result tenantResponse
	var loginName, serialNumber string
	jsonutil.ToObject(resp, &result)
	if result.Result != (tenantResult{}) {
		loginName = result.Result.Login.LoginCode
		serialNumber = result.Result.Login.SerialNumber
	} else {
		loginName = result.Module.Login.LoginCode
		serialNumber = result.Module.Login.SerialNumber
	}
	return &dto.TenantDTO{
		Name:  loginName,
		Phone: serialNumber,
	}
}

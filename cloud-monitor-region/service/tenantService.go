package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"log"
	"time"
)

const RedisTimeOutOneHour = 60 * 60

type TenantService struct {
}

func NewTenantService() *TenantService {
	return &TenantService{}
}

func (s *TenantService) GetTenantInfo(tenantId string) dtos.TenantDTO {
	var tenant dtos.TenantDTO
	key := "userInfo:" + tenantId
	value, err := redis.Get(key)
	if err != nil {
		log.Fatalln("获取缓存出错, key=" + key)
	}
	if tools.IsNotBlank(value) {
		//
		tools.ToObject(value, &tenant)
	} else {
		tenant := getTenant(tenantId)
		if tenant != nil {
			err := redis.SetByTimeOut(key, tools.ToString(tenant), RedisTimeOutOneHour*time.Second)
			if err != nil {
				log.Fatalln("设置redis出错, key=" + key)
			}
		} else {
			log.Fatalln("获取租户信息为空")
		}
	}

	return tenant
}

func getTenant(tenantId string) *dtos.TenantDTO {
	m := make(map[string]string, 1)
	m["loginId"] = tenantId
	resp, err := tools.HttpPostJson(config.GetCommonConfig().TenantUrl, m)
	if err != nil {
		log.Fatalln("查询租户信息失败")
		return nil
	}
	var result map[string]map[string]map[string]string
	var loginName, serialNumber string
	tools.ToObject(resp, &result)
	if result["module"] != nil {
		loginName = result["module"]["login"]["loginCode"]
		serialNumber = result["module"]["login"]["serialNumber"]
	} else {
		loginName = result["result"]["login"]["loginCode"]
		serialNumber = result["result"]["login"]["serialNumber"]
	}
	return &dtos.TenantDTO{
		Name:  loginName,
		Phone: serialNumber,
	}
}

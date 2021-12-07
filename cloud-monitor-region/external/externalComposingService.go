package external

import commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"

var ProductInstanceServiceMap = map[string]commonService.InstanceService{
	"cbr": &CbrInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	"ecs": &EcsInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	"eip": &EipInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	"nat": &NatInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	"slb": &SlbInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
}

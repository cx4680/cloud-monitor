package external

import commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"

//已接入的产品简称
const (
	ECS = "ecs"
	CBR = "cbr"
	EIP = "eip"
	NAT = "nat"
	SLB = "slb"
)

var ProductInstanceServiceMap = map[string]commonService.InstanceService{
	CBR: &CbrInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	ECS: &EcsInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	EIP: &EipInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	NAT: &NatInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
	SLB: &SlbInstanceService{
		InstanceServiceImpl: commonService.InstanceServiceImpl{},
	},
}

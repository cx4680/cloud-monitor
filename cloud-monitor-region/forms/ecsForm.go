package forms

type EcsQueryPageForm struct {
	TenantId     string `json:"tenantId"`
	Current      int    `json:"current"`
	PageSize     int    `json:"pageSize"`
	InstanceName string `json:"instanceName"`
	InstanceId   string `json:"instanceId"`
	Status       int    `json:"status"`
	StatusList   []int  `json:"statusList"`
}

type EcsQueryPageVO struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    EcsPageVO `json:"data"`
}

type EcsPageVO struct {
	Total int      `json:"total"`
	Rows  []*ECSVO `json:"rows"`
}

type ECSVO struct {
	InstanceId   string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	Region       string `json:"region"`
	Ip           string `json:"ip"`
	Status       int    `json:"status"`
	OsType       string `json:"osType"`
}

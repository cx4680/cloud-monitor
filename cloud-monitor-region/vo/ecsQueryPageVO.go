package vo

type EcsQueryPageVO struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    EcsPageVO `json:"data"`
}

type EcsPageVO struct {
	Total int     `json:"total"`
	Rows  []ECSVO `json:"rows"`
}

type ECSVO struct {
	InstanceId   string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	Region       string `json:"region"`
	Ip           string `yaml:"ip"`
	Status       int    `yaml:"status"`
}

package form

type BmsQueryPageRequest struct {
	PageNumber       int    `json:"pageNumber,omitempty"`
	PageSize         int    `json:"pageSize,omitempty"`
	TenantId         string `json:"tenantId,omitempty"`
	AllTenants       string `json:"all_tenants,omitempty"`
	Id               string `json:"id,omitempty"`
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

type BmsResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Data    BmsQueryPageResult `json:"data"`
}

type BmsQueryPageResult struct {
	Servers    []BmsServers `json:"servers"`
	TotalCount int          `json:"total_count"`
}

type BmsServers struct {
	FlavorId string `json:"flavor_id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Id       string `json:"id"`
}

type SerialResponse struct {
	Code    string                `json:"code"`
	Message string                `json:"message"`
	Data    SerialQueryPageResult `json:"data"`
}

type SerialQueryPageResult struct {
	Bmcnodes   []Bmcnodes `json:"bmcnodes"`
	TotalCount int        `json:"total_count"`
}

type Bmcnodes struct {
	SerialNumber      string `json:"serial_number"`
	EcsId             string `json:"ecs_id"`
	Manufactory       string `json:"manufactory"`
	Productor         string `json:"productor"`
	Cpu               int    `json:"cpu"`
	CpuManufactory    string `json:"cpu_manufactory"`
	CpuProductor      string `json:"cpu_productor"`
	CpuArch           string `json:"cpu_arch"`
	Memory            int    `json:"memory"`
	SystemDisk        int    `json:"system_disk"`
	DataDisk          int    `json:"data_disk"`
	Nic               string `json:"nic"`
	GpuType           string `json:"gpu_type"`
	FirmwareVersion   string `json:"firmware_version"`
	IdcRoom           string `json:"idc_room"`
	CabinetNum        string `json:"cabinet_num"`
	RankNum           string `json:"rank_num"`
	AvailabilityZone  string `json:"availability_zone"`
	IpmiBmcIp         string `json:"ipmi_bmc_ip"`
	IpmiBmcUsername   string `json:"ipmi_bmc_username"`
	IpmiBmcPassword   string `json:"ipmi_bmc_password"`
	SmartnicManagerIp string `json:"smartnic_manager_ip"`
	Smartniciqn       string `json:"smartniciqn"`
	HealthyStatus     bool   `json:"healthy_status"`
	UnhealthyReason   string `json:"unhealthy_reason"`
	UserSystemStatus  string `json:"user_system_status"`
}

package form

type InstanceRequest struct {
	CloudProductCode string   `json:"cloudProductCode"`
	ResourceTypeCode string   `json:"resourceTypeCode"`
	ResourceId       string   `json:"resourceId"`
	StatusList       []string `json:"statusList"`
	RegionCode       string   `json:"regionCode"`
	Name             string   `json:"name"`
	TenantId         string   `json:"tenantId"`
	PageSize         string   `json:"pageSize"`
	CurrPage         string   `json:"currPage"`
	DirectoryIds     []string `json:"directoryIds"`
}

type InstanceResponse struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	TraceId string `json:"traceId"`
	Data    *struct {
		Total int `json:"total"`
		List  []*struct {
			Id               int    `json:"id"`
			UuidStr          string `json:"uuidStr"`
			RegionCode       string `json:"regionCode"`
			RegionName       string `json:"regionName"`
			ResourceTypeCode string `json:"resourceTypeCode"`
			CloudProductCode string `json:"cloudProductCode"`
			TenantId         string `json:"tenantId"`
			TenantName       string `json:"tenantName"`
			ResourceId       string `json:"resourceId"`
			ResourceName     string `json:"resourceName"`
			OrderId          string `json:"orderId"`
			ResourceUrl      string `json:"resourceUrl"`
			AvailabilityZone string `json:"availabilityZone"`
			Status           int    `json:"status"`
			StatusDesc       string `json:"statusDesc"`
			Deleted          int    `json:"deleted"`
			Additional       string `json:"additional"`
			Modifier         string `json:"modifier"`
		} `json:"list"`
	} `json:"data"`
}

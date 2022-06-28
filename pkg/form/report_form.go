package form

type ReportFormParam struct {
	Name         string         `form:"name" json:"name"`
	ItemList     []string       `form:"itemList" json:"itemList"`
	InstanceList []InstanceForm `form:"instanceList" json:"instanceList"`
	/**
	 * 区间数据查询参数 时间戳
	 */
	Start int `form:"start" json:"start"`
	End   int `form:"end" json:"end"`
	Step  int `form:"step" json:"step"`
	/**
	 * 统计方式
	 * 聚合函数 sum(求和)  min(最小值)  max (最大值)  avg (平均值)  stddev (标准差)  stdvar (标准差异)  count (计数)
	 */
	Statistics []string `form:"statistics" json:"statistics"`
	Current    int      `form:"current" json:"current"`
	PageSize   int      `form:"pageSize" json:"pageSize"`
	RegionCode string   `form:"regionCode" json:"regionCode"`
	FileId     string   `form:"fileId" json:"fileId"`
}

type InstanceForm struct {
	InstanceName string `form:"instanceName" json:"instanceName"`
	InstanceId   string `form:"instanceId" json:"instanceId"`
	Status       string `form:"status" json:"status"`
}

type ReportForm struct {
	Region       string `json:"region"`
	InstanceName string `json:"instanceName"`
	InstanceId   string `json:"instanceId"`
	Status       string `json:"status"`
	ItemName     string `json:"itemName"`
	Time         string `json:"time"`
	Timestamp    int64  `json:"timestamp"`
	Value        string `json:"value"`
	MaxValue     string `json:"maxValue"`
	MinValue     string `json:"minValue"`
	AvgValue     string `json:"avgValue"`
}

type AsyncExportRequest struct {
	TemplateId string             `json:"templateId"`
	Params     []AsyncExportParam `json:"params"`
}

type AsyncExportParam struct {
	SheetSeq   int    `json:"sheetSeq"`
	SheetName  string `json:"sheetName"`
	SheetParam string `json:"sheetParam"`
}

type ExportRecordsRequest struct {
	SheetSeq   int    `json:"sheetSeq"`
	SheetName  string `json:"sheetName"`
	SheetParam string `json:"sheetParam"`
}

type ExportRecordsResponse struct {
	PageCount int                   `json:"pageCount"`
	PageSize  int                   `json:"pageSize"`
	Current   int                   `json:"current"`
	Result    []ExportRecordsResult `json:"result"`
}

type ExportRecordsResult struct {
	BusiId     int64  `json:"busiId"`
	Module     string `json:"module"`
	FileName   string `json:"fileName"`
	CreateDate string `json:"createDate"`
	FileSuffix string `json:"fileSuffix"`
	FileStatus string `json:"fileStatus"`
	FileId     string `json:"fileId"`
}

type CallbackReportForm struct {
	PageSize int    `json:"pageSize"`
	Current  int    `json:"current"`
	Param    string `json:"param"`
}

type ExportRecords struct {
	PageSize  int `json:"pageSize"`
	Current   int `json:"current"`
	PageCount int `json:"pageCount"`
	Result    []struct {
		BusiId     int64  `json:"busiId"`
		Module     string `json:"module"`
		FileName   string `json:"fileName"`
		CreateDate string `json:"createDate"`
		FileSuffix string `json:"fileSuffix"`
		FileStatus string `json:"fileStatus"`
		FileId     string `json:"fileId"`
	} `json:"result"`
}

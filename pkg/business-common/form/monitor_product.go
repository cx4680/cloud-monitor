package form

type MonitorProductParam struct {
	ProductCodeList []string `form:"productCodeList"`
	Status          uint8    `form:"status"`
}

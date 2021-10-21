package global

const (
	// SuccessServer HTTP code 正常
	SuccessServer = "200"
	// ErrorServer 服务错误
	ErrorServer = "500"
)

const TenantId = "tenantId"
const UserId = "userId"

const (
	ALL      = "ALL"
	INSTANCE = "INSTANCE"
)

var ResourceScopeText = map[string]int{
	ALL:      1,
	INSTANCE: 2,
}

func GetResourceScopeInt(code string) int {
	return ResourceScopeText[code]
}

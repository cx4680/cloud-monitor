package constant

const (
	INSTANCE    = "instance"
	FILTER      = "device!='tmpfs'"
	MetricLabel = "$INSTANCE"
	TopExpr     = "topk(5,(%s))"

	TenantId = "accountId"
)

var ProductMap = map[string]string{
	"1": "ecs",
	"2": "eip",
	"3": "slb",
	"4": "mysql",
	"5": "cbr",
	"6": "nat",
	"7": "bms",
}

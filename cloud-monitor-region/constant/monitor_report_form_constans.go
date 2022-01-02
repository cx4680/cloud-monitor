package constant

const (
	INSTANCE           = "instance"
	FILTER             = "device!='tmpfs'"
	MetricLabel        = "$INSTANCE"
	EcsCpuUsage        = "ecs_cpu_usage"
	EcsCpuUsageTopExpr = "topk(5,(100 - (100 * (sum by(instance) (irate(ecs_cpu_seconds_total{mode='idle',$INSTANCE}[3m])) / sum by(instance) (irate(ecs_cpu_seconds_total{$INSTANCE}[3m]))))))"

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

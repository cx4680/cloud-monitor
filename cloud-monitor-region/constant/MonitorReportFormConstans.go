package constant

const (
	INSTANCE           = "instance"
	FILTER             = "device!='tmpfs'"
	MetricLabel        = "$INSTANCE"
	EcsCpuUsage        = "ecs_cpu_usage"
	EcsCpuUsageTopExpr = "topk(5,(100 - (100 * (sum by(instance) (irate(ecs_cpu_seconds_total{mode='idle',$INSTANCE}[3m])) / sum by(instance) (irate(ecs_cpu_seconds_total{$INSTANCE}[3m]))))))"

	Ecs   = "ecs"
	Eip   = "eip"
	Slb   = "slb"
	MySql = "mysql"
	Cbr   = "cbr"
	Nat   = "nat"

	TenantId = "accountId"
)

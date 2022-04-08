package constant

const (
	INSTANCE    = "instance"
	FILTER      = "device!='tmpfs'"
	MetricLabel = "$INSTANCE"
	TopExpr     = "topk(5,(%s))"

	TenantId = "accountId"
)

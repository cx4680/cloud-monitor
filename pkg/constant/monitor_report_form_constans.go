package constant

const (
	INSTANCE    = "instance"
	FILTER      = "device!='tmpfs'"
	MetricLabel = "$INSTANCE"
	TopExpr     = "topk(%s,(%s))"
	PId         = "pid='%s'"
)

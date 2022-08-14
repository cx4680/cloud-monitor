package constant

const (
	TenantInfoKey      = "cec-hawkeye:userInfo:%s"       //缓存租户信息
	TenantRuleKey      = "cec-hawkeye:rule:%s"           //租户获取规则锁key
	TenantDirectoryKey = "cloud-monitor:iamDirectory:%s" //租户获取规则锁key

	SyncTaskKey          = "task_%s_started" //同步器MasterId
	SyncTaskStartLockKey = "task_%s_start_lock"

	MonitorItemKey = "cloud-monitor:monitorItem:%s" //监控项redis缓存key
)

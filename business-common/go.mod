module code.cestc.cn/ccos-ops/cloud-monitor/business-common

go 1.15

require (
	code.cestc.cn/ccos-ops/cloud-monitor/common v0.0.0-20211028062752-e559c17fe0f2
	github.com/go-redis/redis/v8 v8.11.4
	github.com/google/uuid v1.3.0
	github.com/pkg/errors v0.8.1
	github.com/robfig/cron v1.2.0
	gorm.io/gorm v1.22.1
)

replace code.cestc.cn/ccos-ops/cloud-monitor/common => ../common

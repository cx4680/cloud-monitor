module code.cestc.cn/ccos-ops/cloud-monitor/business-common

go 1.15

require (
	code.cestc.cn/ccos-ops/cloud-monitor/common v0.0.0-20211028062752-e559c17fe0f2
	github.com/go-redis/redis/v8 v8.11.4
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/pkg/errors v0.8.1
	gorm.io/gorm v1.22.1
)

replace code.cestc.cn/ccos-ops/cloud-monitor/common => ../common

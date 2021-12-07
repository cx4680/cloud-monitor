module code.cestc.cn/ccos-ops/cloud-monitor/business-common

go 1.15

require (
	code.cestc.cn/ccos-ops/cloud-monitor/common v0.0.0-20211028062752-e559c17fe0f2
	github.com/apache/rocketmq-client-go/v2 v2.1.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/google/uuid v1.3.0
	github.com/json-iterator/go v1.1.12
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/pkg/errors v0.8.1
	github.com/robfig/cron v1.2.0
	github.com/tidwall/gjson v1.2.1
	go.uber.org/atomic v1.9.0
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.22.1
    code.cestc.cn/yyptb-group_tech/iam-sdk-go v1.0.3
)

replace code.cestc.cn/ccos-ops/cloud-monitor/common => ../common

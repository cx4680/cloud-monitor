package sys_redis

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"context"
	"log"
	"sync"
	"testing"
)

func TestRedisLock(t *testing.T) {
	config.InitConfig("config.local.yml")
	InitClient(config.Cfg.Redis)
	var waitGroup = &sync.WaitGroup{}
	waitGroup.Add(10)
	for i := 0; i < 10; i++ {
		go test1(waitGroup, i)
	}
	waitGroup.Wait()

}

func test1(wg *sync.WaitGroup, i int) {
	defer wg.Done()
	background := context.Background()
	err := Lock(background, "test", 10000000000, true)
	if err != nil {
		log.Printf("xxx get lock fail %+v", i)
	}
	log.Printf("get lock  success %+v", i)
	Unlock(ctx, "test")
}

package sysRedis

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"time"
)

var (
	ctx   = context.Background()
	rdb   *redis.Client
	mutex sync.Mutex
)

func InitClient(config config.RedisConfig) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		DB:           0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("redis connection error", err)
		return err
	}
	return nil
}

func Set(key, value string) error {
	return rdb.Set(ctx, key, value, 0).Err()
}

func SetByTimeOut(key, value string, timeout time.Duration) error {
	return rdb.Set(ctx, key, value, timeout).Err()
}

func Get(key string) (string, error) {
	cmd := rdb.Get(ctx, key)
	return cmd.Result()
}

func Lock(key string) error {
	mutex.Lock()
	defer mutex.Unlock()
	return rdb.SetNX(ctx, key, 1, 10*time.Second).Err()
}

func UnLock(key string) error {
	return rdb.Del(ctx, key).Err()
}

func GetClient() *redis.Client {
	return rdb
}

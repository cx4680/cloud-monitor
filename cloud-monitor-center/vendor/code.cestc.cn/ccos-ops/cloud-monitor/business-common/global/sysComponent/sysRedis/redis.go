package sysRedis

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"os"
	"time"
)

var (
	ctx            = context.Background()
	rdb            *redis.Client
	MaxRetry       = 200
	DefaultLease   = 10 * time.Second
	ErrLockByOther = errors.New("Lock by others")
)

func InitClient(config config.RedisConfig) error {
	redisPwd := os.Getenv("REDIS_PWD")
	logger.Logger().Info("redis init pwd:", redisPwd)
	rdb = redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     redisPwd,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		DB:           0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Logger().Errorf("redis connection error", err)
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

func GetClient() *redis.Client {
	return rdb
}

func Lock(ctx context.Context, key string, lease time.Duration, wait bool) error {
	retry := 0
	select {
	case <-ctx.Done():
		return ErrLockByOther
	default:
	BEGIN:
		isSet, err := rdb.SetNX(ctx, key, 1, lease).Result()
		if err != nil {
			return err
		}
		if isSet {
			return nil
		} else if wait {
			retry++
			if retry > MaxRetry {
				return ErrLockByOther
			}
			// max wait 10 seconds
			time.Sleep(100 * time.Millisecond)
			goto BEGIN
		}
	}
	return ErrLockByOther
}

func Unlock(ctx context.Context, key string) error {
	err := rdb.Del(ctx, key).Err()
	return err
}

package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func initEngine(name string, cfg redisConfig) {
	rdb, err := newRedis(cfg)
	if err != nil {
		panic(err)
	}
	setConnection(name, rdb)
}

func newRedis(cfg redisConfig) (redis.UniversalClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // no password set
		DB:       cfg.Db,       // use default DB
	})
	// 加日志钩子
	rdb.AddHook(newLoggerHook())
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

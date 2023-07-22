package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

func Default() redis.UniversalClient {
	return Connection("default")
}

var (
	_mapConn map[string]redis.UniversalClient
)

func Connection(name string) redis.UniversalClient {
	if e, ok := _mapConn[name]; ok {
		return e
	}
	return nil
}

func setConnection(name string, rdb redis.UniversalClient) {
	fmt.Printf("redis connect success! name=[%s]\n", name)
	_mapConn[name] = rdb
}

func init() {
	_mapConn = make(map[string]redis.UniversalClient, 1)
}

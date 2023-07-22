package redis

import (
	"errors"
	"github.com/redis/go-redis/v9"
)

const (
	Nil = redis.Nil
)

/**
 * IsRedisNil
 * @Description: Redis Nil
 * @author lijunpeng<lijunpeng@weimiao.cn>
 * @date: 2021-08-06 18:42:34
 * @param err
 * @return bool
 */
func IsRedisNil(err error) bool {
	return errors.Is(err, Nil)
}

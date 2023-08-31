/**
 * @Description:
 * @FilePath: /bull-golang/internel/redisAction/init.go
 * @Author: liyibing liyibing@lixiang.com
 * @Date: 2023-07-26 10:13:03
 */
package redisAction

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

func InitRedisClient(ip string, passwd string) (redis.Cmdable, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     ip,
		Password: passwd,
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New("redis init failed")
	}
	return rdb, nil
}

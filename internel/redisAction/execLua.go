/**
 * @Description:
 * @FilePath: /bull-golang/internel/redisAction/execLua.go
 * @Author: liyibing liyibing@lixiang.com
 * @Date: 2023-07-26 17:08:14
 */
package redisAction

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func ExecLua(luaScript string, rdb redis.Cmdable, keys []string, args []interface{}) error {
	ctx := context.Background()
	_, err := rdb.Eval(ctx, luaScript, keys, args).Result()
	if err != nil {
		return err
	}
	return nil
}

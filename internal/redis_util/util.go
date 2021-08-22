/**
 * @File: util
 * @Date: 2021/8/3 下午9:57
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package redis_util

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func newRedisPool(host, port, password, db string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host+":"+port)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				_ = c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", db); err != nil {
				_ = c.Close()
				return nil, err
			}
			return c, nil
		},
		MaxIdle:         20,
		IdleTimeout:     time.Minute * 3,
		Wait:            true,
		MaxConnLifetime: 0,
	}
}

/**
 * @File: redis_utils
 * @Date: 2021/8/3 下午7:41
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package redis_util

import (
	"errors"
	"github.com/gomodule/redigo/redis"
)

type RedisDAO struct {
	pool *redis.Pool
}

//var (
//	rao *RedisDAO
//)
//
//func Init(host, port, password, db string) error {
//	rao = &RedisDAO{
//		pool: newRedisPool(host, port, password, db),
//	}
//	if err := rao.TestConnection(); err != nil {
//		return err
//	}
//
//	return nil
//}

func New(host, port, password, db string) (*RedisDAO, error) {
	r := &RedisDAO{
		pool: newRedisPool(host, port, password, db),
	}
	if err := r.TestConnection(); err != nil {
		return nil, err
	}

	return r, nil
}

func (t *RedisDAO) TestConnection() error {
	c := t.pool.Get()
	defer c.Close()
	v, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}
	if v != "PONG" {
		return errors.New("redis internal_connect_types fail")
	}

	return nil
}

func (t *RedisDAO) Reset() error {
	c := t.pool.Get()
	defer c.Close()
	_, err := c.Do("FLUSHDB")
	return err
}

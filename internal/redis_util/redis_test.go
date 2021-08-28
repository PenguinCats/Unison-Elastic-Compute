package redis_util

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestRedisSetAndDel(t *testing.T) {
	pool := newRedisPool("223.3.84.194", "19500", "tnFA40IR", "0")
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "1", "1")
	if err != nil {
		fmt.Println(err.Error())
		t.Fatalf("test set and del fail")
	}
	_, err = conn.Do("SET", "2", "2")
	if err != nil {
		fmt.Println(err.Error())
		t.Fatalf("test set and del fail")
	}

	keys := []string{"1", "2"}
	v, err := redis.Int(conn.Do("DEL", redis.Args{}.AddFlat(keys)...))
	fmt.Println(v)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatalf("test set and del fail")
	}
	if v != 2 {
		t.Fatalf("test set and del fail with wrong number")
	}
}

func TestRedisSetArray(t *testing.T) {
	pool := newRedisPool("223.3.84.194", "19500", "tnFA40IR", "0")
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", "ll")
	if err != nil {
		fmt.Println(err.Error())
		t.Fatalf("test array fail")
	}

	_, err = conn.Do("RPUSH", redis.Args{}.Add("ll").AddFlat([]string{"1", "2", "3"})...)
	if err != nil {
		fmt.Println(err.Error())
		t.Fatalf("test array fail")
	}

	v, err := redis.Strings(conn.Do("LRANGE", "ll", "0", "-1"))
	if err != nil {
		fmt.Println(err.Error())

	}

	if len(v) != 3 {
		t.Fatalf("array element number wrong")
	}
}

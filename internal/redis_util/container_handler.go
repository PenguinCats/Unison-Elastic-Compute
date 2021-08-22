package redis_util

import "fmt"

func (t *RedisDAO) ContainerSet(slaveID, containerID, key string, value interface{}) error {
	containerKey := fmt.Sprintf("ues:container:%s.%s:%s", slaveID, containerID, key)
	conn := t.pool.Get()

	_, err := conn.Do("SET", containerKey, value)
	if err != nil {
		return err
	}

	return nil
}

func (t *RedisDAO) ContainerSetWithTime(slaveID, containerID, key string, value interface{}, seconds int) error {
	containerKey := fmt.Sprintf("ues:container:%s.%s:%s", slaveID, containerID, key)
	conn := t.pool.Get()

	_, err := conn.Do("SET", containerKey, value, "EX", seconds)
	if err != nil {
		return err
	}

	return nil
}

package redis_util

import (
	"fmt"
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/container"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

func (t *RedisDAO) ContainerSetBusy(containerID string) bool {
	conn := t.pool.Get()
	defer conn.Close()
	containerKey := fmt.Sprintf("uec:container:%s:busy", containerID)
	res, err := conn.Do("SET", containerKey, "1", "NX", "EX", 300)
	if err != nil || res == nil {
		return false
	}
	return true
}

func (t *RedisDAO) ContainerReleaseBusy(containerID string) {
	conn := t.pool.Get()
	defer conn.Close()
	containerKey := fmt.Sprintf("uec:container:%s:busy", containerID)
	_, _ = conn.Do("DEL", containerKey)
}

func (t *RedisDAO) ContainerDelAll(containerID string) {
	keys := []string{"busy", "slave_ID", "status", "profile.dict",
		"profile.exposed_tcp_ports", "profile.exposed_tcp_mapping_ports",
		"profile.exposed_udp_ports", "profile.exposed_udp_mapping_ports"}
	var containerKeys []string
	for _, key := range keys {
		containerKey := fmt.Sprintf("uec:container:%s:%s", containerID, key)
		containerKeys = append(containerKeys, containerKey)
	}

	conn := t.pool.Get()
	defer conn.Close()
	_, _ = conn.Do("DEL", redis.Args{}.AddFlat(containerKeys)...)
}

func (t *RedisDAO) ContainerSetProfile(containerID string, profile container.ContainerProfile) error {
	conn := t.pool.Get()
	defer conn.Close()

	_ = conn.Send("MULTI")

	_ = conn.Send("DEL", fmt.Sprintf("uec:container:%s:profile.dict", containerID))
	_ = conn.Send("DEL", fmt.Sprintf("uec:container:%s:profile.exposed_tcp_ports", containerID))
	_ = conn.Send("DEL", fmt.Sprintf("uec:container:%s:profile.exposed_tcp_mapping_ports", containerID))
	_ = conn.Send("DEL", fmt.Sprintf("uec:container:%s:profile.exposed_udp_ports", containerID))
	_ = conn.Send("DEL", fmt.Sprintf("uec:container:%s:profile.exposed_udp_mapping_ports", containerID))

	_ = conn.Send("HMSET", fmt.Sprintf("uec:container:%s:profile.dict", containerID),
		"ext_container_id", profile.ExtContainerID,
		"image_name", profile.ImageName,
		"core_request", profile.CoreRequest,
		"memory_request", profile.MemoryRequest,
		"storage_request", profile.StorageRequest)

	if len(profile.ExposedTCPPorts) > 0 {
		_ = conn.Send("RPUSH", redis.Args{}.
			Add(fmt.Sprintf("uec:container:%s:profile.exposed_tcp_ports", containerID)).
			AddFlat(profile.ExposedTCPPorts)...)

		_ = conn.Send("RPUSH", redis.Args{}.
			Add(fmt.Sprintf("uec:container:%s:profile.exposed_tcp_mapping_ports", containerID)).
			AddFlat(profile.ExposedTCPMappingPorts)...)
	}
	if len(profile.ExposedUDPPorts) > 0 {
		_ = conn.Send("RPUSH", redis.Args{}.
			Add(fmt.Sprintf("uec:container:%s:profile.exposed_udp_ports", containerID)).
			AddFlat(profile.ExposedUDPPorts)...)
		_ = conn.Send("RPUSH", redis.Args{}.
			Add(fmt.Sprintf("uec:container:%s:profile.exposed_udp_mapping_ports", containerID)).
			AddFlat(profile.ExposedUDPMappingPorts)...)
	}

	_, err := conn.Do("EXEC")

	if err != nil {
		logrus.Warning(err.Error())
	}
	return err
}

func (t *RedisDAO) ContainerSet(containerID, key string, value interface{}) error {
	containerKey := fmt.Sprintf("uec:container:%s:%s", containerID, key)
	conn := t.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", containerKey, value)
	if err != nil {
		return err
	}

	return nil
}

func (t *RedisDAO) ContainerSetWithTime(slaveID, containerID, key string, value interface{}, seconds int) error {
	containerKey := fmt.Sprintf("ues:container:%s.%s:%s", slaveID, containerID, key)
	conn := t.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", containerKey, value, "EX", seconds)
	if err != nil {
		return err
	}

	return nil
}

func (t *RedisDAO) ContainerHSet(containerID, key, field string, value interface{}) error {
	containerKey := fmt.Sprintf("uec:container:%s:%s", containerID, key)
	conn := t.pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", containerKey, field, value)
	if err != nil {
		return err
	}

	return nil
}

func (t *RedisDAO) ContainerHSetWithTime(containerID, key, field string, value interface{}, second int) error {
	containerKey := fmt.Sprintf("uec:container:%s:%s", containerID, key)
	conn := t.pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", containerKey, field, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", containerKey, second)
	if err != nil {
		return err
	}

	return nil
}

func (t *RedisDAO) ContainerHMSETWithTime(containerID, key string) error {
	containerKey := fmt.Sprintf("uec:container:%s:%s", containerID, key)
	conn := t.pool.Get()
	defer conn.Close()

	_, err := conn.Do("HMSET", containerKey)
	return err
}

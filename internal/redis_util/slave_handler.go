package redis_util

import (
	"errors"
	"fmt"
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/hosts"
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/resource"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"regexp"
)

//func (t *RedisDAO) SlaveResetProfile(slaveID string) error {
//
//}

func (t *RedisDAO) SlaveUUIDList() ([]string, error) {
	conn := t.pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "uec:slave:*:hostinfo"))

	if err != nil {
		logrus.Warning(err.Error())
		return nil, err

	}

	uuidRegexp := regexp.MustCompile(`^uec:slave:([\w-]+):hostinfo$`)

	var uuids []string
	for _, s := range keys {
		params := uuidRegexp.FindStringSubmatch(s)
		if len(params) < 2 {
			return nil, errors.New("regexp error")
		}
		uuids = append(uuids, params[1])
	}

	return uuids, nil
}

func (t *RedisDAO) SlaveProfile(slaveID string) (map[string]string, error) {
	slaveHostInfoKey := fmt.Sprintf("uec:slave:%s:hostinfo", slaveID)

	conn := t.pool.Get()
	defer conn.Close()

	mp, err := redis.StringMap(conn.Do("HGETALL", slaveHostInfoKey))
	if err != nil {
		return nil, err
	}

	profileKey := []string{"platform", "platform_family", "platform_version", "mem_total",
		"cpu_model_name", "logical_cpu_num", "physical_cpu_num"}
	for _, key := range profileKey {
		_, ok := mp[key]
		if !ok {
			err = fmt.Errorf("profile not enough for slaveID: %s", slaveID)
			logrus.Warning(err.Error())
			return nil, err
		}
	}

	return mp, nil
}

func (t *RedisDAO) SlaveStats(slaveID string) (string, error) {
	slaveStatsKey := fmt.Sprintf("uec:slave:%s:stats", slaveID)

	conn := t.pool.Get()
	defer conn.Close()

	stats, err := redis.String(conn.Do("GET", slaveStatsKey))
	if err != nil {
		return "", err
	}

	return stats, nil
}

func (t *RedisDAO) SlaveStatus(slaveID string) (map[string]string, error) {
	slaveStatusKey := fmt.Sprintf("uec:slave:%s:status", slaveID)

	conn := t.pool.Get()
	defer conn.Close()

	mp, err := redis.StringMap(conn.Do("HGETALL", slaveStatusKey))
	if err != nil {
		return nil, err
	}

	profileKey := []string{"core_available", "mem_available", "storage_available"}
	for _, key := range profileKey {
		_, ok := mp[key]
		if !ok {
			return nil, errors.New("status not enough")
		}
	}

	return mp, nil
}

func (t *RedisDAO) SlaveResetHostInfo(slaveID string, hostInfo hosts.HostInfo) error {
	slaveHostInfoKey := fmt.Sprintf("uec:slave:%s:hostinfo", slaveID)

	conn := t.pool.Get()
	defer conn.Close()

	_ = conn.Send("MULTI")

	_ = conn.Send("HMSET", slaveHostInfoKey,
		"platform", hostInfo.Platform,
		"platform_family", hostInfo.PlatformFamily,
		"platform_version", hostInfo.PlatformVersion,
		"mem_total", hostInfo.MemoryTotalSize,
		"cpu_model_name", hostInfo.CpuModelName,
		"logical_cpu_num", hostInfo.LogicalCoreCnt,
		"physical_cpu_num", hostInfo.PhysicalCoreCnt)

	_, err := conn.Do("EXEC")

	if err != nil {
		logrus.Warning(err.Error())
	}
	return err
}

func (t *RedisDAO) SlaveUpdateStats(slaveID string, stats types.StatsSlave) error {
	statsString := types.GetSlaveStatsString(stats)
	slaveKey := fmt.Sprintf("uec:slave:%s:stats", slaveID)

	conn := t.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", slaveKey, statsString, "EX", 45)
	if err != nil {
		logrus.Warning(err.Error())
	}

	return err
}

func (t *RedisDAO) SlaveUpdateStatus(slaveID string, stats types.StatsSlave, resource resource.ResourceAvailable) error {
	statsString := types.GetSlaveStatsString(stats)
	slaveKey := fmt.Sprintf("uec:slave:%s:stats", slaveID)

	conn := t.pool.Get()
	defer conn.Close()

	_ = conn.Send("MULTI")

	_ = conn.Send("SET", slaveKey, statsString, "EX", 45)

	slaveStatusKey := fmt.Sprintf("uec:slave:%s:status", slaveID)

	_ = conn.Send("HMSET", slaveStatusKey,
		"core_available", resource.CoreAvailable,
		"mem_available", resource.MemoryAvailable,
		"storage_available", resource.StorageAvailable)

	_ = conn.Send("EXPIRE", slaveStatusKey, 45)

	_, err := conn.Do("EXEC")

	if err != nil {
		logrus.Warning(err.Error())
	}
	return err
}

func (t *RedisDAO) SlaveUpdateAddToken(token string) error {
	key := "uec:slave_add_token"

	conn := t.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, token, "EX", 300)
	if err != nil {
		logrus.Warning(err.Error())
	}

	return err
}

func (t *RedisDAO) SlaveGetAddToken() (string, error) {
	key := "uec:slave_add_token"

	conn := t.pool.Get()
	defer conn.Close()

	token, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return token, nil
}

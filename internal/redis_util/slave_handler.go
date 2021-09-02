package redis_util

import (
	"fmt"
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/hosts"
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/resource"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/sirupsen/logrus"
)

//func (t *RedisDAO) SlaveResetProfile(slaveID string) error {
//
//}

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

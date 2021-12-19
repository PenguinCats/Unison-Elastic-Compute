/**
 * @File: setting
 * @Date: 2021/7/20 下午8:33
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package main

import (
	"github.com/go-ini/ini"
	"log"
)

type Host struct {
	PortBias int
}

var HostSetting = &Host{}

type Connect struct {
	MasterIP        string
	MasterPort      string
	MasterSecretKey string
}

var ConnectSetting = &Connect{}

type DockerController struct {
	MemoryReserveRatio          int64
	StorageReserveRatioForImage int64
	StoragePoolName             string
	CoreAvailableList           []string
	HostPortRange               string
	ContainerStopTimeout        int
}

var DockerControllerSetting = &DockerController{}

type System struct {
	Reload bool
}

var SystemSetting = &System{}

var cfg *ini.File

func LoadGlobalSetting(configPath string) (err error) {
	cfg, err = ini.Load(configPath)
	if err != nil {
		return
	}

	mapTo("connect", ConnectSetting)
	mapTo("docker_controller", DockerControllerSetting)
	mapTo("host", HostSetting)
	mapTo("system", SystemSetting)

	return err
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

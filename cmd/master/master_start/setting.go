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

type General struct {
	Recovery string
}

var GeneralSetting = &General{}

type Connect struct {
	SlaveControlListenerPort string
}

var ConnectSetting = &Connect{}

type Api struct {
	MasterAPIPort string
}

var ApiSetting = &Api{}

type Redis struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
}

var RedisSetting = &Redis{}

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
	mapTo("api", ApiSetting)
	mapTo("redis", RedisSetting)
	mapTo("general", GeneralSetting)
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

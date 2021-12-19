/**
 * @File: master_start
 * @Date: 2021/7/20 下午8:58
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package main

import (
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&nested.Formatter{
		ShowFullLevel: true,
	})

	if err := LoadGlobalSetting("cmd/master/master_start/master.ini"); err != nil {
		panic(err.Error())
	}

	cmb := types.CreatMasterBody{
		Recovery:                 GeneralSetting.Recovery,
		SlaveControlListenerPort: ConnectSetting.SlaveControlListenerPort,
		APIPort:                  ApiSetting.MasterAPIPort,
		RedisHost:                RedisSetting.RedisHost,
		RedisPort:                RedisSetting.RedisPort,
		RedisPassword:            RedisSetting.RedisPassword,
		RedisDB:                  RedisSetting.RedisDB,
		Reload:                   SystemSetting.Reload,
	}
	m, err := master.New(cmb)
	if err != nil {
		panic(err.Error())
	}

	if err := m.Start(); err != nil {
		panic(err.Error())
	}

	ch := make(chan bool, 1)
	<-ch
}

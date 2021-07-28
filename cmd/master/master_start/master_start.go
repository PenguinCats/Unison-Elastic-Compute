/**
 * @File: master_start
 * @Date: 2021/7/20 下午8:58
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package main

import (
	master2 "Unison-Elastic-Compute/api/types/control/master"
	"Unison-Elastic-Compute/cmd/master/internal/settings"
	"Unison-Elastic-Compute/pkg/master"
	"time"
)

func main() {
	ch := make(chan bool, 1)

	if err := settings.LoadGlobalSetting("cmd/master/conf.ini"); err != nil {

	}

	cmb := master2.CreatMasterBody{
		SlaveControlListenerPort: settings.ConnectSetting.SlaveControlListenerPort,
		MasterAPIPort:            settings.ApiSetting.MasterAPIPort,
	}
	m := master.New(cmb)

	if err := m.Start(); err != nil {
		panic(err.Error())
	}

	time.Sleep(time.Second * 5)
	<-ch
}

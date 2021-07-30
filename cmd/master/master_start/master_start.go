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
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	logrus.SetFormatter(&nested.Formatter{
		ShowFullLevel: true,
	})

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

	ch := make(chan bool, 1)
	<-ch
}

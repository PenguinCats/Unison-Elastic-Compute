package main

import (
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/docker_controller"
	slave3 "github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	slave2 "github.com/PenguinCats/Unison-Elastic-Compute/pkg/slave"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&nested.Formatter{
		ShowFullLevel: true,
	})

	err := LoadGlobalSetting("./slave.ini")
	if err != nil {
		return
	}

	slave, err := slave2.NewSlave(slave3.CreatSlaveBody{
		MasterIP:        ConnectSetting.MasterIP,
		MasterPort:      ConnectSetting.MasterPort,
		MasterSecretKey: ConnectSetting.MasterSecretKey,
		HostPortBias:    HostSetting.PortBias,
		Reload:          SystemSetting.Reload,
	}, docker_controller.DockerControllerCreatBody{
		MemoryReserveRatio:          DockerControllerSetting.MemoryReserveRatio,
		StorageReserveRatioForImage: DockerControllerSetting.StorageReserveRatioForImage,
		StoragePoolName:             DockerControllerSetting.StoragePoolName,
		CoreAvailableList:           DockerControllerSetting.CoreAvailableList,
		HostPortRange:               DockerControllerSetting.HostPortRange,
		ContainerStopTimeout:        DockerControllerSetting.ContainerStopTimeout,
		Reload:                      SystemSetting.Reload,
	})
	if err != nil {
		panic(err.Error())
	}

	slave.Start()

	ch := make(chan bool, 1)
	<-ch
}

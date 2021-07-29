package master

import (
	"Unison-Elastic-Compute/api/types/control/master"
	"Unison-Elastic-Compute/pkg/master/internal/slave-controller"
)

type Master struct {
	slaveController          *slave_controller.SlaveController
	SlaveControlListenerPort string

	masterAPIPort string
}

func New(cmb master.CreatMasterBody) *Master {
	m := &Master{
		SlaveControlListenerPort: cmb.SlaveControlListenerPort,
		masterAPIPort:            cmb.MasterAPIPort,
	}
	return m
}

func (m *Master) Start() error {
	slaveController, err := slave_controller.NewSlaveController(slave_controller.CreateSlaveControllerBody{
		SlaveControlListenerPort: m.SlaveControlListenerPort,
	})
	if err != nil {
		return slave_controller.ErrSlaveControllerCreat
	}
	m.slaveController = slaveController

	m.slaveController.Start()

	return nil
}

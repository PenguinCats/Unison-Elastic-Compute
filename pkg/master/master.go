package master

import (
	"Unison-Elastic-Compute/api/types/control/master"
	"Unison-Elastic-Compute/pkg/master/internal/slave-control"
)

type Master struct {
	slaveController          *slave_control.SlaveController
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
	slaveController, err := slave_control.NewSlaveController(slave_control.CreateSlaveControllerBody{
		SlaveControlListenerPort: m.SlaveControlListenerPort,
	})
	if err != nil {
		return slave_control.ErrSlaveControllerCreat
	}
	m.slaveController = slaveController

	m.slaveController.Start()

	return nil
}

package master

import (
	"Unison-Elastic-Compute/api/types/control/master"
	"Unison-Elastic-Compute/pkg/master/internal/slave-control"
)

type Master struct {
	slaveController *slave_control.SlaveController
}

func New(cmb master.CreatMasterBody) (*Master, error) {
	slaveController, err := slave_control.NewSlaveController(slave_control.CreateSlaveControllerBody{
		SlaveControlListenerPort: cmb.SlaveControlListenerPort,
	})
	if err != nil {
		return nil, slave_control.ErrSlaveControllerCreat
	}
	m := &Master{
		slaveController: slaveController,
	}
	return m, nil
}

func (m *Master) Start() {
	m.slaveController.Start()
}

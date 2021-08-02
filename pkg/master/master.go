package master

import (
	"Unison-Elastic-Compute/api/types/control/master"
	"Unison-Elastic-Compute/pkg/master/internal/slave-controller"
)

type Master struct {
	slaveController          *slave_controller.SlaveController
	SlaveControlListenerPort string

	httpApiController *HttpApiController
	APIPort           string
}

func New(cmb master.CreatMasterBody) *Master {
	m := &Master{
		SlaveControlListenerPort: cmb.SlaveControlListenerPort,
		APIPort:                  cmb.APIPort,
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

	httpApiController := newHttpApiController(m)
	m.httpApiController = httpApiController

	m.slaveController.Start()
	m.httpApiController.startHttpApiServe()

	return nil
}

func (m *Master) Stop() {

}

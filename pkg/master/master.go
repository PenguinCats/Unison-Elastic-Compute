package master

import (
	"errors"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/redis_util"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/slave-controller"
)

type Master struct {
	slaveController          *slave_controller.SlaveController
	slaveControlListenerPort string

	httpApiController *HttpApiController
	apiPort           string

	redisDAO *redis_util.RedisDAO
}

func New(cmb types.CreatMasterBody) (*Master, error) {
	rdao, err := redis_util.New(cmb.RedisHost, cmb.RedisPort, cmb.RedisPassword, cmb.RedisDB)
	if err != nil {
		panic(err.Error())
	}

	m := &Master{
		slaveControlListenerPort: cmb.SlaveControlListenerPort,
		apiPort:                  cmb.APIPort,
		redisDAO:                 rdao,
	}

	if err := m.init(cmb.Recovery); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Master) init(recovery string) (err error) {
	if recovery == "true" {

	} else if recovery == "false" {
		err = m.redisDAO.Reset()
	} else {
		err = errors.New("invalid parameter: recovery")
	}

	return
}

func (m *Master) Start() error {
	slaveController, err := slave_controller.NewSlaveController(slave_controller.CreateSlaveControllerBody{
		SlaveControlListenerPort: m.slaveControlListenerPort,
		RedisDAO:                 m.redisDAO,
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

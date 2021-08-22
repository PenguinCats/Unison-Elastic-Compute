package master

import (
	"errors"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/redis_util"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/http-controller"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/operation"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/slave-controller"
)

type Master struct {
	slaveController          *slave_controller.SlaveController
	slaveControlListenerPort string

	httpApiController *http_controller.HttpApiController
	apiPort           string

	operationTaskChan     chan *operation.OperationTask
	operationResponseChan chan *operation.OperationResponse

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
		operationTaskChan:        make(chan *operation.OperationTask, 100),
		operationResponseChan:    make(chan *operation.OperationResponse, 100),
	}

	slaveController, err := slave_controller.NewSlaveController(slave_controller.CreateSlaveControllerBody{
		SlaveControlListenerPort: m.slaveControlListenerPort,
		RedisDAO:                 m.redisDAO,
		OperationResponseChan:    m.operationResponseChan,
	})
	if err != nil {
		return nil, slave_controller.ErrSlaveControllerCreat
	}
	m.slaveController = slaveController

	httpApiController := http_controller.NewHttpApiController(m.apiPort, m.operationTaskChan)
	m.httpApiController = httpApiController

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
	m.slaveController.Start()

	m.handleOperationTask()
	m.handleOperationResponse()

	err := m.httpApiController.StartHttpApiServe()
	if err != nil {
		return err
	}

	return nil
}

func (m *Master) Stop() {

}

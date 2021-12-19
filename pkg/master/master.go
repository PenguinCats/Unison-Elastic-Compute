package master

import (
	"errors"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/redis_util"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/util"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/http-controller"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/operation"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/slave-controller"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

type Master struct {
	slaveController          *slave_controller.SlaveController
	slaveControlListenerPort string

	httpApiController *http_controller.HttpApiController
	apiPort           string

	operationTaskChan     chan *operation.OperationTask
	operationResponseChan chan *operation.OperationResponse

	redisDAO *redis_util.RedisDAO
	db       *leveldb.DB
}

func New(cmb types.CreatMasterBody) (*Master, error) {
	rdao, err := redis_util.New(cmb.RedisHost, cmb.RedisPort, cmb.RedisPassword, cmb.RedisDB)
	if err != nil {
		panic(err.Error())
	}

	dbPath := "/var/opt/uec/master.db"
	if !cmb.Reload {
		exist, err := util.IsPathExists(dbPath)
		if err != nil {
			return nil, err
		}
		if exist {
			err := os.RemoveAll(dbPath)
			if err != nil {
				return nil, err
			}
		}
	}
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}

	m := &Master{
		slaveControlListenerPort: cmb.SlaveControlListenerPort,
		apiPort:                  cmb.APIPort,
		redisDAO:                 rdao,
		db:                       db,
		operationTaskChan:        make(chan *operation.OperationTask, 100),
		operationResponseChan:    make(chan *operation.OperationResponse, 100),
	}

	slaveController, err := slave_controller.NewSlaveController(slave_controller.CreateSlaveControllerBody{
		SlaveControlListenerPort: m.slaveControlListenerPort,
		RedisDAO:                 m.redisDAO,
		Db:                       m.db,
		OperationResponseChan:    m.operationResponseChan,
	})
	if err != nil {
		return nil, slave_controller.ErrSlaveControllerCreat
	}
	m.slaveController = slaveController

	if cmb.Reload {
		err := m.slaveController.Reload()
		if err != nil {
			return nil, slave_controller.ErrSlaveControllerCreat
		}
	}

	httpApiController := http_controller.NewHttpApiController(m.apiPort, m.operationTaskChan, rdao)
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

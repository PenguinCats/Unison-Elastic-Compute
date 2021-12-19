package slave

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/docker_controller"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/util"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"net"
	"os"
	"sync"
	"time"

	"github.com/PenguinCats/Unison-Docker-Controller/pkg/controller"
)

type Slave struct {
	status types.StatsSlave

	masterIP      string
	masterPort    string
	joinSecretKey string

	uuid  string
	token string

	ctrlConn    *net.TCPConn
	ctrlEncoder *json.Encoder
	ctrlDecoder *json.Decoder
	dataConn    *net.TCPConn
	dataEncoder *json.Encoder
	dataDecoder *json.Decoder

	scbStopFunc context.CancelFunc

	mu sync.RWMutex

	lastHeartbeatTime     time.Time
	lastHeartbeatTimeLock sync.RWMutex

	dc           *controller.DockerController
	hostPortBias int

	db *leveldb.DB
}

func NewSlave(cb types.CreatSlaveBody, dccb docker_controller.DockerControllerCreatBody) (*Slave, error) {
	dbPath := "/var/opt/uec/slave.db"

	if !cb.Reload {
		exist, err := util.IsPathExists(dbPath)
		if err != nil {
			return nil, err
		}
		if exist {
			e := os.RemoveAll(dbPath)
			if e != nil {
				logrus.Warning(e.Error())
				return nil, err
			}
		}
	}

	db, err := leveldb.OpenFile("/var/opt/uec/slave.db", nil)
	if err != nil {
		return nil, err
	}

	dc, err := controller.NewDockerController(&dccb)
	if err != nil {
		return nil, err
	}

	slave := &Slave{
		masterIP:      cb.MasterIP,
		masterPort:    cb.MasterPort,
		joinSecretKey: cb.MasterSecretKey,
		dc:            dc,
		hostPortBias:  cb.HostPortBias,
		db:            db,
	}

	if cb.Reload {
		exist_token, err_token := db.Has([]byte("uec:token"), nil)
		if err_token != nil {
			return nil, err
		}

		exist_uuid, err_uuid := db.Has([]byte("uec:uuid"), nil)
		if err_uuid != nil {
			return nil, err
		}

		if !exist_uuid || !exist_token {
			logrus.Warning("No existing connection record, perform new registration")
		} else {
			token, err := db.Get([]byte("uec:token"), nil)
			if err_uuid != nil {
				return nil, err
			}

			uuid, err := db.Get([]byte("uec:uuid"), nil)
			if err != nil {
				return nil, err
			}

			slave.token = string(token)
			slave.uuid = string(uuid)
		}
	}

	return slave, nil
}

func (s *Slave) Start() {
	if s.uuid != "" && s.token != "" {
		// 执行恢复
		err := s.reconnect()
		if err != nil {
			panic(fmt.Sprintf("Slave Start Reconnect Error with [%s]", err.Error()))
		}
	} else {
		// 新注册
		err := s.register()
		if err != nil {
			panic(fmt.Sprintf("Slave Connect Error with [%s]", err.Error()))
		}
	}

	logrus.Warning("register success")

	ctx, cancel := context.WithCancel(context.Background())
	s.scbStopFunc = cancel

	s.startHandleCtrlMessage(ctx)
	s.startHandleDataMessage(ctx)

	time.Sleep(time.Second)

	s.startSendHeartbeat(ctx)
	//slave.startHeartbeatCheck(ctx)

}

func (s *Slave) StopWork() {
	s.stopActivity()
	s.mu.Lock()
	s.status = types.StatsStopped
	s.mu.Unlock()
}

func (s *Slave) GetStatus() types.StatsSlave {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.status
}

func (s *Slave) stopActivity() {
	s.mu.Lock()
	s.scbStopFunc()
	_ = s.ctrlConn.Close()
	_ = s.dataConn.Close()
	s.mu.Unlock()
}

func (s *Slave) offline() {
	s.stopActivity()
	s.mu.Lock()
	s.status = types.StatsOffline
	s.mu.Unlock()
}

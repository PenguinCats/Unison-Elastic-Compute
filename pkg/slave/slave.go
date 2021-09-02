package slave

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/docker_controller"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"

	"github.com/PenguinCats/Unison-Docker-Controller/pkg/controller"
)

type Slave struct {
	status types.StatsSlave

	masterIP   string
	masterPort string
	secretKey  string

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
}

func NewSlave(cb types.CreatSlaveBody, dccb docker_controller.DockerControllerCreatBody) (*Slave, error) {
	dc, err := controller.NewDockerController(&dccb)
	if err != nil {
		return nil, err
	}

	return &Slave{
		masterIP:     cb.MasterIP,
		masterPort:   cb.MasterPort,
		secretKey:    cb.MasterSecretKey,
		dc:           dc,
		hostPortBias: cb.HostPortBias,
	}, nil
}

func (s *Slave) Start() {
	err := s.register()
	if err != nil {
		panic(fmt.Sprintf("Slave Start Error with [%s]", err.Error()))
	}
	logrus.Warning("register success")

	ctx, cancel := context.WithCancel(context.Background())
	s.scbStopFunc = cancel

	s.startHandleCtrlMessage(ctx)
	s.startHandleDataMessage(ctx)

	time.Sleep(time.Second)

	// TODO do heartbeat check
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

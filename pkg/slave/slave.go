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
	status types.StatusSlave

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

func (slave *Slave) Start() {
	err := slave.register()
	if err != nil {
		panic(fmt.Sprintf("Slave Start Error with [%s]", err.Error()))
	}
	logrus.Warning("register success")

	ctx, cancel := context.WithCancel(context.Background())
	slave.scbStopFunc = cancel

	slave.startHandleCtrlMessage(ctx)
	slave.startSendHeartbeat(ctx)
	//slave.startHeartbeatCheck(ctx)
}

func (slave *Slave) StopWork() {
	slave.stopActivity()
	slave.mu.Lock()
	slave.status = types.StatusStopped
	slave.mu.Unlock()
}

func (slave *Slave) GetStatus() types.StatusSlave {
	slave.mu.RLock()
	defer slave.mu.RUnlock()

	return slave.status
}

func (slave *Slave) stopActivity() {
	slave.mu.Lock()
	slave.scbStopFunc()
	_ = slave.ctrlConn.Close()
	_ = slave.dataConn.Close()
	slave.mu.Unlock()
}

func (slave *Slave) offline() {
	slave.stopActivity()
	slave.mu.Lock()
	slave.status = types.StatusOffline
	slave.mu.Unlock()
}

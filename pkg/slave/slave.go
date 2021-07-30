package slave

import (
	slave2 "Unison-Elastic-Compute/api/types/control/slave"
	"Unison-Elastic-Compute/pkg/internal/communication/api/control"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

type Slave struct {
	status slave2.StatusSlave

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

	mu sync.Mutex

	lastHeartbeatTime     time.Time
	lastHeartbeatTimeLock sync.RWMutex
}

func New(cb slave2.CreatSlaveBody) *Slave {
	return &Slave{
		masterIP:   cb.MasterIP,
		masterPort: cb.MasterPort,
		secretKey:  cb.MasterSecretKey,
	}
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
	slave.status = slave2.StatusOffline
	slave.mu.Unlock()
}

func (slave *Slave) StopWork() {
	slave.stopActivity()
	slave.mu.Lock()
	slave.status = slave2.StatusStopped
	slave.mu.Unlock()
}

func (slave *Slave) startHandleCtrlMessage(ctx context.Context) {
	go func() {
		var err error = nil
		defer func() {
			if err != nil {
				logrus.Warning(err.Error())
				slave.offline()
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				message := control.Message{}
				err = slave.ctrlDecoder.Decode(&message)
				if err != nil {
					err = control.ErrControlInvalidMessage
					return
				}

				switch message.MessageType {
				case control.MessageCtrlTypeHeartbeat:
					slave.handleHeartbeatMessage(message.Value)
				default:
					err = control.ErrControlInvalidMessage
					return
				}
			}
		}
	}()
}

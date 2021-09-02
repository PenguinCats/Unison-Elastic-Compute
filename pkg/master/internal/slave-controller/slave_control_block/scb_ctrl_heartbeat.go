package slave_control_block

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/container"
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/resource"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_control_types"
	"github.com/sirupsen/logrus"
	"time"
)

func (scb *SlaveControlBlock) SetLastHeartbeatTime(t time.Time) {
	scb.lastHeartbeatTimeLock.Lock()
	scb.lastHeartbeatTime = t
	scb.lastHeartbeatTimeLock.Unlock()
}

func (scb *SlaveControlBlock) GetLastHeartbeatTime() (t time.Time) {
	scb.lastHeartbeatTimeLock.RLock()
	t = scb.lastHeartbeatTime
	scb.lastHeartbeatTimeLock.RUnlock()
	return
}

func (scb *SlaveControlBlock) sendHeartbeatACK() {
	m := internal_control_types.HeartBeatMessageAck{}
	v, err := json.Marshal(&m)
	if err != nil {
		logrus.Warning(err.Error())
	}

	err = scb.ctrlEncoder.Encode(&internal_control_types.Message{
		MessageType: internal_control_types.MessageCtrlTypeHeartbeat,
		Value:       v,
	})
	if err != nil {
		logrus.Warning(err.Error())
	}
}

func (scb *SlaveControlBlock) handleHeartbeatMessage(v []byte) {
	m := internal_control_types.HeartBeatMessageReport{}
	err := json.Unmarshal(v, &m)
	if err != nil {
		logrus.Warning(internal_control_types.ErrControlInvalidHeartbeat.Error())
		return
	}

	fmt.Println(m)

	go scb.heartbeatContainerStatusUpdate(m.ContainerStatus)
	go scb.heartbeatSlaveStatusUpdate(scb.uuid, m.Stats, m.Resource)

	scb.sendHeartbeatACK()

	scb.SetLastHeartbeatTime(time.Now())
}

func (scb *SlaveControlBlock) startHeartbeatCheck(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Second * 15)
				lastHeartBeatTime := scb.GetLastHeartbeatTime()
				if time.Now().Sub(lastHeartBeatTime) > time.Minute*3 {
					scb.offline()
					return
				}
			}
		}
	}()
}

func (scb *SlaveControlBlock) heartbeatContainerStatusUpdate(containerStatus map[string]container.ContainerStatus) {
	for k, v := range containerStatus {
		_ = scb.RedisDAO.ContainerUpdateStatus(k, v)
	}
}

func (scb *SlaveControlBlock) heartbeatSlaveStatusUpdate(slaveID string, stats types.StatsSlave, resource resource.ResourceAvailable) {
	_ = scb.RedisDAO.SlaveUpdateStatus(slaveID, stats, resource)
}

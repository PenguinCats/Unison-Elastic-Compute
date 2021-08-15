package slave_control_block

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/control"
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
	m := control.HeartBeatMessageAck{}
	v, err := json.Marshal(&m)
	if err != nil {
		logrus.Warning(err.Error())
	}

	err = scb.ctrlEncoder.Encode(&control.Message{
		MessageType: control.MessageCtrlTypeHeartbeat,
		Value:       v,
	})
	if err != nil {
		logrus.Warning(err.Error())
	}
}

func (scb *SlaveControlBlock) handleHeartbeatMessage(v []byte) {
	m := control.HeartBeatMessageReport{}
	err := json.Unmarshal(v, &m)
	if err != nil {
		logrus.Warning(control.ErrControlInvalidHeartbeat.Error())
		return
	}

	fmt.Println(m)

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
				time.Sleep(time.Second * 30)
				lastHeartBeatTime := scb.GetLastHeartbeatTime()
				if time.Now().Sub(lastHeartBeatTime) > time.Minute*3 {
					scb.offline()
					return
				}
			}
		}
	}()
}

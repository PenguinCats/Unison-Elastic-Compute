/**
 * @File: slave_heart_beat
 * @Date: 2021/7/22 下午8:59
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave

import (
	"context"
	"encoding/json"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/control"
	"github.com/sirupsen/logrus"
	"time"
)

func (slave *Slave) setLastHeartbeatTime(t time.Time) {
	slave.lastHeartbeatTimeLock.Lock()
	slave.lastHeartbeatTime = t
	slave.lastHeartbeatTimeLock.Unlock()
}

func (slave *Slave) GetLastHeartbeatTime() (t time.Time) {
	slave.lastHeartbeatTimeLock.RLock()
	t = slave.lastHeartbeatTime
	slave.lastHeartbeatTimeLock.RUnlock()
	return
}

func (slave *Slave) handleHeartbeatMessage(v []byte) {
	m := control.HeartBeatMessageAck{}
	err := json.Unmarshal(v, &m)
	if err != nil {
		logrus.Warning(control.ErrControlInvalidHeartbeat.Error())
		return
	}

	slave.setLastHeartbeatTime(time.Now())
}

func (slave *Slave) startSendHeartbeat(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m := control.HeartBeatMessageReport{
					Status:   slave.GetStatus(),
					Resource: *slave.dc.GetResourceAvailable(),
				}
				v, err := json.Marshal(&m)
				if err != nil {
					logrus.Warning(err.Error())
				}

				err = slave.ctrlEncoder.Encode(&control.Message{
					MessageType: control.MessageCtrlTypeHeartbeat,
					Value:       v,
				})
				if err != nil {
					logrus.Warning(err.Error())
				}
			}
			time.Sleep(time.Second * 30)
		}
	}()
}

func (slave *Slave) startHeartbeatCheck(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Second * 30)
				lastHeartBeatTime := slave.GetLastHeartbeatTime()
				if time.Now().Sub(lastHeartBeatTime) > time.Minute*3 {
					slave.offline()
					return
				}
			}
		}
	}()
}

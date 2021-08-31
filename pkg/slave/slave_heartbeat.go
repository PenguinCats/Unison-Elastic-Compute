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
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_control_types"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *Slave) setLastHeartbeatTime(t time.Time) {
	s.lastHeartbeatTimeLock.Lock()
	s.lastHeartbeatTime = t
	s.lastHeartbeatTimeLock.Unlock()
}

func (s *Slave) GetLastHeartbeatTime() (t time.Time) {
	s.lastHeartbeatTimeLock.RLock()
	t = s.lastHeartbeatTime
	s.lastHeartbeatTimeLock.RUnlock()
	return
}

func (s *Slave) handleHeartbeatMessage(v []byte) {
	m := internal_control_types.HeartBeatMessageAck{}
	err := json.Unmarshal(v, &m)
	if err != nil {
		logrus.Warning(internal_control_types.ErrControlInvalidHeartbeat.Error())
		return
	}

	s.setLastHeartbeatTime(time.Now())
}

func (s *Slave) startSendHeartbeat(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m := internal_control_types.HeartBeatMessageReport{
					Status:          s.GetStatus(),
					Resource:        *s.dc.GetResourceAvailable(),
					ContainerStatus: s.dc.ContainerAllStats(),
				}
				v, err := json.Marshal(&m)
				if err != nil {
					logrus.Warning(err.Error())
				}

				err = s.ctrlEncoder.Encode(&internal_control_types.Message{
					MessageType: internal_control_types.MessageCtrlTypeHeartbeat,
					Value:       v,
				})
				if err != nil {
					logrus.Warning(err.Error())
				}
			}
			time.Sleep(time.Second * 15)
		}
	}()
}

func (s *Slave) startHeartbeatCheck(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Second * 15)
				lastHeartBeatTime := s.GetLastHeartbeatTime()
				if time.Now().Sub(lastHeartBeatTime) > time.Minute*3 {
					s.offline()
					return
				}
			}
		}
	}()
}

/**
 * @File: slave_control_register
 * @Date: 2021/7/15 上午10:49
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_controller

import (
	"encoding/json"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	register2 "github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_connect_types"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

func (sc *SlaveController) reconnectCtrlConnection(c net.Conn, d *json.Decoder) {
	var err error = nil
	defer func() {
		if err != nil {
			_ = c.Close()
			logrus.Info(err.Error())
		}
	}()

	e := json.NewEncoder(c)

	// Reconnect Ctrl Communication Handshake Step 1
	hs1b := register2.ReconnectCtrlConnectionHandshakeStep1Body{}
	err = d.Decode(&hs1b)
	if err != nil {
		err = ErrReconnectCtrlConnInvalidRequest
		return
	}

	hs2b := register2.ReconnectCtrlConnectionHandshakeStep2Body{
		Ack:   hs1b.Seq + 1,
		Agree: true,
	}

	scb, ok := sc.GetSlaveCtrlBlk(hs1b.UUID)
	if !ok || scb.GetToken() != hs1b.Token {
		err = ErrReconnectCtrlConnStepFail
		hs2b.Agree = false
	} else {
		scb.ResetCtrl(c, e, d)
		scb.SetLastHeartbeatTime(time.Now())
		scb.SetStatus(types.StatsWaitingEstablishDataConnection)
	}

	// Establish Ctrl Communication Handshake Step 2
	err = e.Encode(&hs2b)
	if err != nil {
		err = ErrReconnectCtrlConnStepFail
		return
	}
}

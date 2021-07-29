/**
 * @File: slave_control_register
 * @Date: 2021/7/15 上午10:49
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_controller

import (
	"Unison-Elastic-Compute/api/types/control/slave"
	"Unison-Elastic-Compute/internal/auth"
	"Unison-Elastic-Compute/pkg/internal/communication/connect/register"
	"Unison-Elastic-Compute/pkg/master/internal/slave-controller/slave_control_block"
	"encoding/json"
	"log"
	"net"
	"time"
)

func (sc *SlaveController) establishCtrlConnection(c net.Conn, d *json.Decoder) {
	var err error = nil
	defer func() {
		if err != nil {
			_ = c.Close()
			log.Println(err.Error())
		}
	}()

	e := json.NewEncoder(c)

	// Establish Ctrl Communication Hand shake Step 1
	hs1b := register.EstablishCtrlConnectionHandshakeStep1Body{}
	err = d.Decode(&hs1b)
	if err != nil {
		err = ErrEstablishCtrlConnInvalidRequest
		return
	}

	// TODO: check Secret Key

	token := auth.GenerateRandomUUID()
	uuid := auth.GenerateRandomUUID()

	scb := slave_control_block.NewWithCtrl(slave.SlaveStatusWaitingEstablishControlConnection, uuid, token, c, d)
	scb.UpdateLastHeartbeatTime(time.Now())

	sc.slaveCtrBlkMutex.Lock()
	sc.slaveCtrBlk[uuid] = scb
	sc.slaveCtrBlkMutex.Unlock()

	// Establish Ctrl Communication Hand shake Step 2
	localSeq := auth.GenerateRandomInt()
	hs2b := register.EstablishCtrlConnectionHandshakeStep2Body{
		Ack:   hs1b.Seq + 1,
		Seq:   localSeq,
		Token: token,
		UUID:  uuid,
	}

	err = e.Encode(&hs2b)
	if err != nil {
		err = ErrEstablishCtrlConnStepFail
		return
	}

	// Establish Ctrl Communication Hand shake Step 3
	hs3b := register.EstablishCtrlConnectionHandshakeStep3Body{}
	err = d.Decode(&hs3b)
	if err != nil {
		err = ErrEstablishCtrlConnInvalidRequest
		return
	}

	if hs3b.Ack != localSeq+1 {
		err = ErrEstablishCtrlConnInvalidRequest
		return
	}

	scb.SetStatus(slave.SlaveStatusWaitingEstablishDataConnection)
	scb.UpdateLastHeartbeatTime(time.Now())
}

func (sc *SlaveController) establishDataConnection(c net.Conn, d *json.Decoder) {
	var err error = nil
	defer func() {
		if err != nil {
			_ = c.Close()
			log.Println(err.Error())
		}
	}()

	e := json.NewEncoder(c)

	// Establish Data Communication Step 1
	hs1b := register.EstablishDataConnectionHandShakeStep1Body{}
	err = d.Decode(&hs1b)
	if err != nil {
		err = ErrEstablishDataConnInvalidRequest
		return
	}

	sc.slaveCtrBlkMutex.RLock()
	scb, ok := sc.slaveCtrBlk[hs1b.UUID]
	sc.slaveCtrBlkMutex.RUnlock()
	if ok != true {
		err = ErrEstablishDataConnInvalidRequest
		return
	}

	ok = scb.GetToken() == hs1b.Token
	if ok != true {
		err = ErrEstablishDataConnInvalidRequest
		return
	}

	scb.SetDataConn(c)
	scb.SetDataDecoder(d)
	scb.UpdateLastHeartbeatTime(time.Now())

	// Establish Data Communication Step 2
	edcs2 := register.EstablishDataConnectionHandShakeStep2Body{}

	err = e.Encode(&edcs2)
	if err != nil {
		err = ErrEstablishDataConnStepFail
		return
	}

	scb.SetStatus(slave.SlaveStatusNormal)
}

// TODO: 巡检

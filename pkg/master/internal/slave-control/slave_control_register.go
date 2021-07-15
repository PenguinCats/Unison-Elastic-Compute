/**
 * @File: slave_control_register
 * @Date: 2021/7/15 上午10:49
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_control

import (
	"Unison-Elastic-Compute/api/types/control/register"
	"Unison-Elastic-Compute/internal/auth"
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"
)

func (sc *SlaveController) register(c net.Conn, hs1b register.HandshakeStep1Body) {
	var err error = nil
	defer func() {
		if err != nil {
			_ = c.Close()
			log.Println(err.Error())
		}
	}()

	// Handshake Step 1
	// TODO: check Secret Key

	e := json.NewEncoder(c)
	adds := strings.Split(c.RemoteAddr().String(), ":")
	if len(adds) != 2 {
		err = ErrConnectionRequestWrongAddress
		return
	}
	ip, port := adds[0], adds[1]

	// Handshake Step 2
	localSeq := auth.GenerateRandomInt()
	token := auth.GenerateRandomUUID()
	uuid := auth.GenerateRandomUUID()
	hs2b := register.HandshakeStep2Body{
		Ack:   hs1b.Seq + 1,
		Seq:   localSeq,
		Token: token,
		UUID:  uuid,
	}

	err = e.Encode(&hs2b)
	if err != nil {
		err = ErrRegisterInvalidBody
		return
	}

	// Handshake Step 3
	d := json.NewDecoder(c)
	hs3b := register.HandshakeStep3Body{}
	err = d.Decode(&hs3b)
	if err != nil {
		err = ErrRegisterInvalidBody
		return
	}

	if hs3b.Ack != localSeq+1 {
		err = ErrRegisterInvalidInfo
		return
	}

	sc.slaveCtrBlkMutex.Lock()
	scb := &SlaveControlBlock{
		UUID:              uuid,
		IP:                ip,
		Port:              port,
		Token:             token,
		ctrConn:           c,
		LastHeartBeatTime: time.Now(),
	}
	sc.slaveCtrBlk[uuid] = scb
	sc.slaveCtrBlkMutex.Unlock()
	go scb.start()
}

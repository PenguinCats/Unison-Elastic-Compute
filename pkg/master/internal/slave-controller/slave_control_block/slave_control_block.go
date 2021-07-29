/**
 * @File: slave-controller
 * @Date: 2021/7/15 上午9:48
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_control_block

import (
	"Unison-Elastic-Compute/api/types/control/slave"
	"Unison-Elastic-Compute/pkg/internal/communication/control"
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

type SlaveControlBlock struct {
	status slave.StatusSlave

	uuid  string
	token string

	ctrlConn    net.Conn
	ctrlDecoder *json.Decoder
	dataConn    net.Conn
	dataDecoder *json.Decoder

	lastHeartBeatTime     time.Time
	lastHeartBeatTimeLock sync.RWMutex

	mu sync.RWMutex
}

func NewWithCtrl(status slave.StatusSlave, uuid, token string,
	ctrlConn net.Conn, ctrlDecoder *json.Decoder) *SlaveControlBlock {
	return &SlaveControlBlock{
		status:      status,
		uuid:        uuid,
		token:       token,
		ctrlConn:    ctrlConn,
		ctrlDecoder: ctrlDecoder,
	}
}

func (scb *SlaveControlBlock) SetStatus(statusSlave slave.StatusSlave) {
	scb.mu.Lock()
	scb.status = statusSlave
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) GetUUID() (uuid string) {
	scb.mu.RLock()
	uuid = scb.uuid
	scb.mu.RUnlock()
	return
}

func (scb *SlaveControlBlock) GetToken() (token string) {
	scb.mu.RLock()
	token = scb.token
	scb.mu.RUnlock()
	return
}

func (scb *SlaveControlBlock) SetCtrlConn(c net.Conn) {
	scb.mu.Lock()
	scb.ctrlConn = c
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) SetCtrlDecoder(d *json.Decoder) {
	scb.mu.Lock()
	scb.ctrlDecoder = d
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) SetDataConn(c net.Conn) {
	scb.mu.Lock()
	scb.dataConn = c
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) SetDataDecoder(d *json.Decoder) {
	scb.mu.Lock()
	scb.dataDecoder = d
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) start() {
	scb.startHandleCtrlMessage()
	scb.startHeartbeatCheck()
}

func (scb *SlaveControlBlock) stop() {

}

func (scb *SlaveControlBlock) startHandleCtrlMessage() {
	go func() {
		var err error = nil
		defer func() {
			if err != nil {
				log.Println(err.Error())
				scb.stop()
			}
		}()

		for {
			message := control.Message{}
			err = scb.ctrlDecoder.Decode(&message)
			if err != nil {
				err = ErrControlInvalidMessage
				return
			}

			switch message.MessageType {
			case control.MessageCtrlTypeHeartbeat:
				scb.HandleHeartbeatMessage(message.Value)
			default:
				err = ErrControlInvalidMessage
				return
			}
		}
	}()
}

/**
 * @File: slave-controller
 * @Date: 2021/7/15 上午9:48
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_control_block

import (
	"Unison-Elastic-Compute/api/types/control/slave"
	control2 "Unison-Elastic-Compute/pkg/internal/communication/api/control"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

type SlaveControlBlock struct {
	status slave.StatusSlave

	uuid  string
	token string

	ctrlConn    net.Conn
	ctrlEncoder *json.Encoder
	ctrlDecoder *json.Decoder
	dataConn    net.Conn
	dataEncoder *json.Encoder
	dataDecoder *json.Decoder

	scbStopFunc context.CancelFunc

	mu sync.RWMutex

	lastHeartbeatTime     time.Time
	lastHeartbeatTimeLock sync.RWMutex
}

func NewWithCtrl(status slave.StatusSlave, uuid, token string,
	ctrlConn net.Conn, ctrlEncoder *json.Encoder, ctrlDecoder *json.Decoder) *SlaveControlBlock {
	return &SlaveControlBlock{
		status:      status,
		uuid:        uuid,
		token:       token,
		ctrlConn:    ctrlConn,
		ctrlEncoder: ctrlEncoder,
		ctrlDecoder: ctrlDecoder,
	}
}

func (scb *SlaveControlBlock) SetStatus(statusSlave slave.StatusSlave) {
	scb.mu.Lock()
	scb.status = statusSlave
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) GetStatus() slave.StatusSlave {
	scb.mu.RLock()
	defer scb.mu.RUnlock()
	return scb.status
}

func (scb *SlaveControlBlock) GetUUID() string {
	scb.mu.RLock()
	defer scb.mu.RUnlock()
	return scb.uuid
}

func (scb *SlaveControlBlock) GetToken() string {
	scb.mu.RLock()
	defer scb.mu.RUnlock()
	return scb.token
}

func (scb *SlaveControlBlock) SetCtrlConn(c net.Conn) {
	scb.mu.Lock()
	scb.ctrlConn = c
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) SetCtrlEncoderDecoder(e *json.Encoder, d *json.Decoder) {
	scb.mu.Lock()
	scb.ctrlEncoder = e
	scb.ctrlDecoder = d
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) SetDataConn(c net.Conn) {
	scb.mu.Lock()
	scb.dataConn = c
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) SetDataEncoderDecoder(e *json.Encoder, d *json.Decoder) {
	scb.mu.Lock()
	scb.dataEncoder = e
	scb.dataDecoder = d
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) Start() {
	logrus.Warning("new slave joined")

	ctx, cancel := context.WithCancel(context.Background())

	scb.mu.Lock()
	scb.scbStopFunc = cancel
	scb.mu.Unlock()

	scb.startHandleCtrlMessage(ctx)
	//scb.startHeartbeatCheck(ctx)
}

func (scb *SlaveControlBlock) stopActivity() {
	scb.mu.Lock()
	scb.scbStopFunc()
	_ = scb.ctrlConn.Close()
	_ = scb.dataConn.Close()
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) offline() {
	scb.stopActivity()
	scb.mu.Lock()
	scb.status = slave.StatusOffline
	scb.mu.Unlock()
	logrus.Warningf("slave [%s] offline", scb.GetUUID())
}

func (scb *SlaveControlBlock) StopWork() {
	scb.stopActivity()
	scb.mu.Lock()
	scb.status = slave.StatusStopped
	scb.mu.Unlock()
	logrus.Warningf("slave [%s] stop", scb.GetUUID())
}

func (scb *SlaveControlBlock) startHandleCtrlMessage(ctx context.Context) {
	go func() {
		var err error = nil
		defer func() {
			if err != nil {
				logrus.Warning(err.Error())
				scb.offline()
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				message := control2.Message{}
				err = scb.ctrlDecoder.Decode(&message)
				if err != nil {
					err = control2.ErrControlInvalidMessage
					return
				}

				switch message.MessageType {
				case control2.MessageCtrlTypeHeartbeat:
					scb.handleHeartbeatMessage(message.Value)
				default:
					err = control2.ErrControlInvalidMessage
					return
				}
			}
		}
	}()
}

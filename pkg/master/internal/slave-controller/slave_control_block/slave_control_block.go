/**
 * @File: slave-controller
 * @Date: 2021/7/15 上午9:48
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_control_block

import (
	"context"
	"encoding/json"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/redis_util"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/operation"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

type SlaveControlBlock struct {
	status types.StatsSlave

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

	operationResponseChan chan *operation.OperationResponse
	RedisDAO              *redis_util.RedisDAO
}

func NewWithCtrl(uuid, token string,
	ctrlConn net.Conn, ctrlEncoder *json.Encoder, ctrlDecoder *json.Decoder,
	operationResponseChan chan *operation.OperationResponse, redisDAO *redis_util.RedisDAO) *SlaveControlBlock {

	scb := &SlaveControlBlock{
		uuid:                  uuid,
		token:                 token,
		ctrlConn:              ctrlConn,
		ctrlEncoder:           ctrlEncoder,
		ctrlDecoder:           ctrlDecoder,
		operationResponseChan: operationResponseChan,
		RedisDAO:              redisDAO,
	}

	scb.SetStatus(types.StatsWaitingEstablishControlConnection)

	return scb
}

func NewWithReload(uuid, token string,
	operationResponseChan chan *operation.OperationResponse, redisDAO *redis_util.RedisDAO) *SlaveControlBlock {

	scb := &SlaveControlBlock{
		uuid:                  uuid,
		token:                 token,
		operationResponseChan: operationResponseChan,
		RedisDAO:              redisDAO,
	}

	scb.SetStatus(types.StatsOffline)

	return scb
}

func (scb *SlaveControlBlock) ResetCtrl(ctrlConn net.Conn, ctrlEncoder *json.Encoder, ctrlDecoder *json.Decoder) {
	scb.mu.Lock()
	defer scb.mu.Unlock()

	scb.ctrlConn = ctrlConn
	scb.ctrlEncoder = ctrlEncoder
	scb.ctrlDecoder = ctrlDecoder
}

func (scb *SlaveControlBlock) SetStatus(statusSlave types.StatsSlave) {
	scb.status = statusSlave
	go func() {
		_ = scb.RedisDAO.SlaveUpdateStats(scb.uuid, statusSlave)
	}()
}

func (scb *SlaveControlBlock) GetStatus() types.StatsSlave {
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
	logrus.Warningf("slave [%s] joined", scb.GetUUID())

	ctx, cancel := context.WithCancel(context.Background())

	scb.mu.Lock()
	scb.scbStopFunc = cancel
	scb.mu.Unlock()

	scb.SetStatus(types.StatsNormal)
	scb.startHandleCtrlMessage(ctx)
	scb.startHandleDataMessage(ctx)
	scb.startHeartbeatCheck(ctx)
}

func (scb *SlaveControlBlock) StopWork() {
	scb.stopActivity()
	scb.mu.Lock()
	scb.SetStatus(types.StatsStopped)
	scb.mu.Unlock()
	logrus.Warningf("slave [%s] stop", scb.GetUUID())
}

func (scb *SlaveControlBlock) stopActivity() {
	scb.mu.Lock()
	if scb.scbStopFunc != nil {
		scb.scbStopFunc()
	}
	if scb.ctrlConn != nil {
		_ = scb.ctrlConn.Close()
	}
	if scb.dataConn != nil {
		_ = scb.dataConn.Close()
	}
	scb.mu.Unlock()
}

func (scb *SlaveControlBlock) offline() {
	if scb.GetStatus() != types.StatsOffline {
		scb.stopActivity()
		scb.mu.Lock()
		scb.SetStatus(types.StatsOffline)
		scb.mu.Unlock()
		logrus.Warningf("slave [%s] offline", scb.GetUUID())
	}
}

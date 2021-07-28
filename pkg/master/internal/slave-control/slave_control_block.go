/**
 * @File: slave-control
 * @Date: 2021/7/15 上午9:48
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_control

import (
	"Unison-Elastic-Compute/api/types/control/slave"
	"encoding/json"
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

	lastHeartBeatTime time.Time

	mu sync.RWMutex
}

func (scb *SlaveControlBlock) start() {

}

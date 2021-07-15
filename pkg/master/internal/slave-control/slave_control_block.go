/**
 * @File: slave-control
 * @Date: 2021/7/15 上午9:48
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_control

import (
	"net"
	"time"
)

type SlaveControlBlock struct {
	UUID  string
	IP    string
	Port  string
	Token string

	ctrConn net.Conn

	LastHeartBeatTime time.Time
}

func (scb *SlaveControlBlock) start() {

}

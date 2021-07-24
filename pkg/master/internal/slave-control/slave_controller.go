/**
 * @File: slave_controller
 * @Date: 2021/7/15 上午9:54
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_control

import (
	"net"
	"sync"
)

type CreateSlaveControllerBody struct {
	SlaveControlListenerPort string
}

type SlaveController struct {
	ctrlLn net.Listener

	slaveCtrBlk      map[string]*SlaveControlBlock
	slaveCtrBlkMutex sync.RWMutex
}

func NewSlaveController(cscb CreateSlaveControllerBody) (*SlaveController, error) {
	ln, err := net.Listen("tcp", ":"+cscb.SlaveControlListenerPort)
	if err != nil {
		return nil, ErrListenerCreat
	}
	defer func() {
		if err != nil {
			_ = ln.Close()
		}
	}()

	sc := &SlaveController{
		ctrlLn:           ln,
		slaveCtrBlk:      make(map[string]*SlaveControlBlock),
		slaveCtrBlkMutex: sync.RWMutex{},
	}

	return sc, nil
}

func (sc *SlaveController) Start() {
	sc.startControlListen()
}

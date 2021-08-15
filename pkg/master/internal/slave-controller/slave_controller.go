/**
 * @File: slave_control
 * @Date: 2021/7/15 上午9:54
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_controller

import (
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/redis_util"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/slave-controller/slave_control_block"
	"net"
	"sync"
)

type CreateSlaveControllerBody struct {
	SlaveControlListenerPort string
	RedisDAO                 *redis_util.RedisDAO
}

type SlaveController struct {
	ctrlLn net.Listener

	slaveCtrBlk      map[string]*slave_control_block.SlaveControlBlock
	slaveCtrBlkMutex sync.RWMutex

	redisDAO *redis_util.RedisDAO
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
		slaveCtrBlk:      make(map[string]*slave_control_block.SlaveControlBlock),
		slaveCtrBlkMutex: sync.RWMutex{},
		redisDAO:         cscb.RedisDAO,
	}

	return sc, nil
}

func (sc *SlaveController) Start() {
	sc.startControlListen()
}

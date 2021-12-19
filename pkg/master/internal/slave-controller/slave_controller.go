/**
 * @File: slave_control
 * @Date: 2021/7/15 上午9:54
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_controller

import (
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/redis_util"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/operation"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/slave-controller/slave_control_block"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"net"
	"sync"
)

type CreateSlaveControllerBody struct {
	SlaveControlListenerPort string
	RedisDAO                 *redis_util.RedisDAO
	Db                       *leveldb.DB
	OperationResponseChan    chan *operation.OperationResponse
}

type SlaveController struct {
	ctrlLn net.Listener

	slaveCtrBlk      map[string]*slave_control_block.SlaveControlBlock
	slaveCtrBlkMutex sync.RWMutex

	redisDAO              *redis_util.RedisDAO
	db                    *leveldb.DB
	operationResponseChan chan *operation.OperationResponse
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
		ctrlLn:                ln,
		slaveCtrBlk:           make(map[string]*slave_control_block.SlaveControlBlock),
		slaveCtrBlkMutex:      sync.RWMutex{},
		redisDAO:              cscb.RedisDAO,
		db:                    cscb.Db,
		operationResponseChan: cscb.OperationResponseChan,
	}

	return sc, nil
}

func (sc *SlaveController) Reload() error {
	iter := sc.db.NewIterator(util.BytesPrefix([]byte("slave:token:")), nil)
	for iter.Next() {
		key := string(iter.Key())
		uuid := key[12:]
		token := string(iter.Value())

		scb := slave_control_block.NewWithReload(uuid, token, sc.operationResponseChan, sc.redisDAO)

		sc.slaveCtrBlkMutex.Lock()
		sc.slaveCtrBlk[uuid] = scb
		sc.slaveCtrBlkMutex.Unlock()
	}

	return nil
}

func (sc *SlaveController) Start() {
	sc.startControlListen()
}

func (sc *SlaveController) GetSlaveCtrlBlk(slaveID string) (*slave_control_block.SlaveControlBlock, bool) {
	sc.slaveCtrBlkMutex.RLock()
	defer sc.slaveCtrBlkMutex.RUnlock()

	scb, ok := sc.slaveCtrBlk[slaveID]
	if !ok {
		return nil, false
	}

	return scb, true
}

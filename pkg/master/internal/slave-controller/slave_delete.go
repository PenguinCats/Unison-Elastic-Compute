/**
 * @File: slave_control
 * @Date: 2021/7/15 上午9:54
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_controller

import (
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func (sc *SlaveController) SlaveDelete(uuid string) error {
	scb, ok := sc.GetSlaveCtrlBlk(uuid)
	if !ok {
		// 不存在就当已经删掉了
		return nil
	}

	scb.StopWork()
	scb.Delete()

	err := sc.db.Delete([]byte(uuid), &opt.WriteOptions{
		Sync: true,
	})
	if err != nil {
		return err
	}

	sc.slaveCtrBlkMutex.Lock()
	delete(sc.slaveCtrBlk, uuid)
	sc.slaveCtrBlkMutex.Unlock()

	return nil
}

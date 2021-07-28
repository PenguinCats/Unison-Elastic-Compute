/**
 * @File: slave_register_test
 * @Date: 2021/7/20 下午6:45
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave

import (
	"log"
	"sync"
	"testing"
	"time"
)

var (
	cnt = 0
	mu  sync.Mutex
)

func TestSlaveRegister(t *testing.T) {
	slave := CreateDefaultTestSlave(t)

	err := slave.register()

	mu.Lock()
	cnt += 1
	mu.Unlock()

	if err != nil {
		t.Fatalf("slave register failed, err message: %s", err.Error())
	}
}

func TestMultipleSlaveRegister(t *testing.T) {
	cnt = 0

	times := 1000

	for i := 0; i < times; i++ {
		log.Println(i)
		go TestSlaveRegister(t)
	}

	time.Sleep(time.Second * 3)
	if cnt != cnt {
		t.Fatalf("message loss: [%d : %d]", cnt, times)
	}
}

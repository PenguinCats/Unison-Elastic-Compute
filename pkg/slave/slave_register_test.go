/**
 * @File: slave_register_test
 * @Date: 2021/7/20 下午6:45
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave

import (
	"testing"
)

func TestSlaveRegister(t *testing.T) {
	slave := CreateDefaultTestSlave(t)

	err := slave.register()

	if err != nil {
		t.Fatalf("slave register failed, err message: %s", err.Error())
	}
}

func TestMultipleSlaveRegister(t *testing.T) {
	for i := 0; i < 100; i++ {
		TestSlaveRegister(t)
	}
}

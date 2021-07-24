/**
 * @File: slave_test
 * @Date: 2021/7/20 下午6:47
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave

import (
	slave2 "Unison-Elastic-Compute/api/types/control/slave"
	"testing"
)

func CreateDefaultTestSlave(t *testing.T) *Slave {
	t.Helper()
	cb := slave2.CreatSlaveBody{
		MasterIP:        "127.0.0.1",
		MasterPort:      "9700",
		MasterSecretKey: "1234567890abcde",
	}
	return New(cb)
}

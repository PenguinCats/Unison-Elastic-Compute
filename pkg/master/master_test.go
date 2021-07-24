/**
 * @File: master_test
 * @Date: 2021/7/20 下午7:05
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package master

import (
	"Unison-Elastic-Compute/api/types/control/master"
	"testing"
)

func CreateDefaultTestMaster(t *testing.T) *Master {
	t.Helper()
	return New(master.CreatMasterBody{
		SlaveControlListenerPort: "9700",
		MasterAPIPort:            "9600",
	})
}

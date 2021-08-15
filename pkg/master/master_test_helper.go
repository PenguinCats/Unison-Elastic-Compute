/**
 * @File: master_test
 * @Date: 2021/7/20 下午7:05
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package master

import (
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"testing"
)

func CreateDefaultTestMaster(t *testing.T) *Master {
	t.Helper()

	cmb := types.CreatMasterBody{
		Recovery:                 "false",
		SlaveControlListenerPort: "9700",
		APIPort:                  "9600",
		RedisHost:                "223.3.84.194",
		RedisPort:                "19500",
		RedisPassword:            "tnFA40IR",
		RedisDB:                  "0",
	}
	m, err := New(cmb)
	if err != nil {
		t.Fatalf("unexpected error: [%s]", err.Error())
	}

	return m
}

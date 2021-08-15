/**
 * @File: slave_test
 * @Date: 2021/7/20 下午6:47
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave

import (
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/docker_controller"
	slave2 "github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"testing"
)

func CreateDefaultTestSlave(t *testing.T) *Slave {
	t.Helper()
	cb := slave2.CreatSlaveBody{
		MasterIP:        "127.0.0.1",
		MasterPort:      "9700",
		MasterSecretKey: "1234567890abcde",
	}

	slave, err := NewSlave(cb, docker_controller.DockerControllerCreatBody{
		MemoryReserveRatio:          5,
		StorageReserveRatioForImage: 10,
		StoragePoolName:             "docker-thinpool",
		CoreAvailableList:           []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"},
		HostPortRange:               "14000-15000",
		ContainerStopTimeout:        5,
	})
	if err != nil {
		t.Fatalf("unexpected error: [%s]", err.Error())
	}

	return slave
}

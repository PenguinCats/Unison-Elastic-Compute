package slave

import (
	"testing"
	"time"
)

func TestHeartBeAT(t *testing.T) {
	slave := CreateDefaultTestSlave(t)
	slave.Start()

	for i := 0; i < 3; i += 1 {
		time.Sleep(time.Second * 35)
		if time.Now().Sub(slave.GetLastHeartbeatTime()) >= time.Second*30 {
			t.Fatal("heartbeat message error")
		}
	}

}

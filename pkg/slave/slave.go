package slave

import (
	slave2 "Unison-Elastic-Compute/api/types/control/slave"
	"fmt"
	"net"
)

type Slave struct {
	status slave2.StatusSlave

	masterIP   string
	masterPort string
	secretKey  string

	uuid  string
	token string

	ctrlConnWithMaster *net.TCPConn
	dataConnWithMaster *net.TCPConn
}

func New(cb slave2.CreatSlaveBody) *Slave {
	return &Slave{
		masterIP:   cb.MasterIP,
		masterPort: cb.MasterPort,
		secretKey:  cb.MasterSecretKey,
	}
}

func (slave *Slave) Start() {
	err := slave.register()
	if err != nil {
		panic(fmt.Sprintf("Slave Start Error with [%s]", err.Error()))
	}
}

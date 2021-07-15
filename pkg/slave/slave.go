package slave

import (
	"net"
)

type Slave struct {
	connWithMaster *net.TCPConn
}

func New() *Slave {
	return &Slave{}
}

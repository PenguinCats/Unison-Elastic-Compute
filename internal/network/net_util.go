package network

import (
	"fmt"
	"net"
)

func CreateConn(ip string, port string) (*net.TCPConn, error) {
	host := ip + ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, fmt.Errorf("create connection error: %s", err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("create connection error: %s", err.Error())
	}

	return conn, nil
}

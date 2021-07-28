/**
 * @File: start_control_listen
 * @Date: 2021/7/15 上午10:11
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_control

import (
	"Unison-Elastic-Compute/pkg/internal/communication/connect"
	"encoding/json"
	"log"
	"net"
)

func (sc *SlaveController) startControlListen() {
	for {
		conn, err := sc.ctrlLn.Accept()
		if err != nil {
			continue
		}
		go sc.handleControlConnection(conn)
	}
}

func (sc *SlaveController) handleControlConnection(c net.Conn) {
	var err error = nil
	defer func() {
		if err != nil {
			log.Println(err.Error())
			_ = c.Close()
		}
	}()

	d := json.NewDecoder(c)
	connectionHead := connect.ConnectionHead{}
	err = d.Decode(&connectionHead)
	if err != nil {
		return
	}

	switch connectionHead.ConnectionType {
	case connect.ConnectionTypeEstablishCtrlConnection:
		sc.establishCtrlConnection(c, d)
	case connect.ConnectionTypeEstablishDataConnection:
		sc.establishDataConnection(c, d)
	case connect.ConnectionTypeReconnect:
	case connect.ConnectionTypeError:
	default:
		err = ErrInvalidConnectionRequest
	}
}

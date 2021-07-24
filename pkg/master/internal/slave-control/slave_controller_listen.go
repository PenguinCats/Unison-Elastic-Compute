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
	go func() {
		for {
			conn, err := sc.ctrlLn.Accept()
			if err != nil {
				continue
			}
			go sc.handleControlConnection(conn)
		}
	}()
}

func (sc *SlaveController) handleControlConnection(c net.Conn) {
	d := json.NewDecoder(c)
	connectionHead := connect.ConnectionHead{}
	err := d.Decode(&connectionHead)
	defer func() {
		if err != nil {
			log.Println(err.Error())
			_ = c.Close()
		}
	}()
	if err != nil {
		return
	}

	switch connectionHead.ConnectionType {
	case connect.ConnectionTypeEstablishCtrlConnection:
		sc.establishCtrlConnection(c)
	case connect.ConnectionTypeEstablishDataConnection:
		sc.establishDataConnection(c)
	case connect.ConnectionTypeReconnect:
	case connect.ConnectionTypeError:
	default:
		err = ErrInvalidConnectionRequest
	}
}

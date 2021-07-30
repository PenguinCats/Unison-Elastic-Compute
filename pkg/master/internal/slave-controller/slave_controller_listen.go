/**
 * @File: start_control_listen
 * @Date: 2021/7/15 上午10:11
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_controller

import (
	connect2 "Unison-Elastic-Compute/pkg/internal/communication/api/connect"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net"
)

func (sc *SlaveController) startControlListen() {
	go func() {
		for {
			conn, err := sc.ctrlLn.Accept()
			if err != nil {
				logrus.Info(err.Error())
				continue
			}
			go sc.handleControlConnection(conn)
		}
	}()
}

func (sc *SlaveController) handleControlConnection(c net.Conn) {
	var err error = nil
	defer func() {
		if err != nil {
			logrus.Info(err.Error())
			_ = c.Close()
		}
	}()

	d := json.NewDecoder(c)
	connectionHead := connect2.ConnectionHead{}
	err = d.Decode(&connectionHead)
	if err != nil {
		return
	}

	switch connectionHead.ConnectionType {
	case connect2.ConnectionTypeEstablishCtrlConnection:
		sc.establishCtrlConnection(c, d)
	case connect2.ConnectionTypeEstablishDataConnection:
		sc.establishDataConnection(c, d)
	case connect2.ConnectionTypeReconnect:
	case connect2.ConnectionTypeError:
	default:
		err = ErrInvalidConnectionRequest
	}
}

/**
 * @File: start_control_listen
 * @Date: 2021/7/15 上午10:11
 * @Author: Binjie Zhang (bj_zhang@seu.edu.cn)
 * @Description: nil
 */

package slave_control

import (
	"Unison-Elastic-Compute/api/types/control"
	"Unison-Elastic-Compute/api/types/control/register"
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
	message := control.Message{}
	err := d.Decode(&message)
	defer func() {
		if err != nil {
			log.Println(err.Error())
			_ = c.Close()
		}
	}()
	if err != nil {
		return
	}

	switch message.MessageType {
	case control.MessageTypeRegister:
		if hs1b, ok := message.Value.(register.HandshakeStep1Body); ok {
			sc.register(c, hs1b)
		} else {
			err = ErrConnectionRequestInvalid
		}

	case control.MessageTypeReconnect:
	default:
		err = ErrConnectionRequestInvalid
	}
}

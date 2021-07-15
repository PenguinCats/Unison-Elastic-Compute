package slave

import (
	"Unison-Elastic-Compute/api/types/control"
	"Unison-Elastic-Compute/api/types/control/register"
	"Unison-Elastic-Compute/api/types/control/slave"
	"Unison-Elastic-Compute/internal/auth"
	"Unison-Elastic-Compute/internal/network"
	"encoding/json"
	"fmt"
)

func (slave *Slave) Register(rb slave.RegisterBody) error {
	conn, err := network.CreateConn(rb.IP, rb.Port)
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()
	if err != nil {
		return fmt.Errorf("register to : %s:%s failed with [%s]", rb.IP, rb.Port, err.Error())
	}

	e := json.NewEncoder(conn)
	d := json.NewDecoder(conn)

	// Handshake Step 1
	hs1b := register.HandshakeStep1Body{
		SecretKey: rb.SecretKey,
		Seq:       auth.GenerateRandomInt(),
	}
	err = e.Encode(&control.Message{
		MessageType: control.MessageTypeRegister,
		Value:       hs1b,
	})
	if err != nil {
		return fmt.Errorf("register to : %s:%s failed in step 1 with [%s]", rb.IP, rb.Port, err.Error())
	}

	// Handshake Step 2
	hs2b := register.HandshakeStep2Body{}
	err = d.Decode(&hs2b)
	if err != nil {
		return fmt.Errorf("register to : %s:%s failed in step 2 with [%s]", rb.IP, rb.Port, err.Error())
	}

	// Handshake Step 3
	hs3b := register.HandshakeStep3Body{
		Ack: hs2b.Seq + 1,
	}
	err = e.Encode(&control.Message{
		MessageType: control.MessageTypeRegister,
		Value:       hs3b,
	})
	if err != nil {
		return fmt.Errorf("register to : %s:%s failed in step 3 with [%s]", rb.IP, rb.Port, err.Error())
	}

	slave.connWithMaster = conn
	return nil
}

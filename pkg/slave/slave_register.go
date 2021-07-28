package slave

import (
	slave2 "Unison-Elastic-Compute/api/types/control/slave"
	"Unison-Elastic-Compute/internal/auth"
	"Unison-Elastic-Compute/internal/network"
	"Unison-Elastic-Compute/pkg/internal/communication/connect"
	"Unison-Elastic-Compute/pkg/internal/communication/connect/register"
	"encoding/json"
	"fmt"
)

func (slave *Slave) register() error {
	if err := slave.establishCtrlConn(slave.masterIP, slave.masterPort, slave.secretKey); err != nil {
		return err
	}

	if err := slave.establishDataConn(slave.masterIP, slave.masterPort, slave.uuid, slave.token); err != nil {
		return err
	}

	slave.status = slave2.SlaveStatusNormal
	return nil
}

func (slave *Slave) establishCtrlConn(ip, port, secretKey string) error {
	conn, err := network.CreateConn(ip, port)
	if err != nil {
		return fmt.Errorf("connect to : %s:%s failed with [%s]", ip, port, err.Error())
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	e := json.NewEncoder(conn)
	d := json.NewDecoder(conn)

	// Mark the purpose of the connect
	err = e.Encode(&connect.ConnectionHead{ConnectionType: connect.ConnectionTypeEstablishCtrlConnection})
	if err != nil {
		return fmt.Errorf("establish ctrl connect failed with [%s]", err.Error())
	}

	// Handshake Step 1
	localSeq := auth.GenerateRandomInt()
	hs1b := register.EstablishCtrlConnectionHandshakeStep1Body{
		SecretKey: secretKey,
		Seq:       localSeq,
	}
	err = e.Encode(&hs1b)
	if err != nil {
		return fmt.Errorf("establish ctrl connect failed in step 1 with [%s]", err.Error())
	}

	// Handshake Step 2
	hs2b := register.EstablishCtrlConnectionHandshakeStep2Body{}
	err = d.Decode(&hs2b)
	if err != nil {
		return fmt.Errorf("establish ctrl connect failed in step 2 with [%s]", err.Error())
	}
	if localSeq+1 != hs2b.Ack {
		return fmt.Errorf("establish ctrl connect failed with [wrong ack]")
	}
	slave.uuid = hs2b.UUID
	slave.token = hs2b.Token

	// Handshake Step 3
	hs3b := register.EstablishCtrlConnectionHandshakeStep3Body{
		Ack: hs2b.Seq + 1,
	}
	err = e.Encode(&hs3b)
	if err != nil {
		return fmt.Errorf("establish ctrl connect failed in step 3 with [%s]", err.Error())
	}

	slave.ctrlConnWithMaster = conn
	return nil
}

func (slave *Slave) establishDataConn(ip, port, uuid, token string) error {
	conn, err := network.CreateConn(ip, port)
	if err != nil {
		return fmt.Errorf("connect to : %s:%s failed with [%s]", ip, port, err.Error())
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	e := json.NewEncoder(conn)
	d := json.NewDecoder(conn)

	// Mark the purpose of the connect
	err = e.Encode(&connect.ConnectionHead{ConnectionType: connect.ConnectionTypeEstablishDataConnection})
	if err != nil {
		return fmt.Errorf("establish data connect failed with [%s]", err.Error())
	}

	// Establish Data Connection Handshake Step 1
	hs1b := register.EstablishDataConnectionHandShakeStep1Body{
		UUID:  uuid,
		Token: token,
	}
	err = e.Encode(&hs1b)
	if err != nil {
		return fmt.Errorf("establish data connect failed in step 1 with [%s]", err.Error())
	}

	// Establish Data Connection Handshake Step 2
	hs2b := register.EstablishDataConnectionHandShakeStep2Body{}
	err = d.Decode(&hs2b)
	if err != nil {
		return fmt.Errorf("establish data connect failed in step 2 with [%s]", err.Error())
	}

	slave.dataConnWithMaster = conn
	return nil
}

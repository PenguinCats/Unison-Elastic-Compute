package slave

import (
	"encoding/json"
	"fmt"
	slave2 "github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/auth"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/network"
	connect2 "github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_connect_types"
)

func (s *Slave) register() error {
	if err := s.establishCtrlConn(s.masterIP, s.masterPort, s.secretKey); err != nil {
		return err
	}

	if err := s.establishDataConn(s.masterIP, s.masterPort, s.uuid, s.token); err != nil {
		return err
	}

	s.status = slave2.StatusNormal
	return nil
}

func (s *Slave) establishCtrlConn(ip, port, secretKey string) error {
	conn, err := network.CreateConn(ip, port)
	if err != nil {
		return fmt.Errorf("internal_connect_types to : %s:%s failed with [%s]", ip, port, err.Error())
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	e := json.NewEncoder(conn)
	d := json.NewDecoder(conn)

	// Mark the purpose of the internal_connect_types
	err = e.Encode(&connect2.ConnectionHead{ConnectionType: connect2.ConnectionTypeEstablishCtrlConnection})
	if err != nil {
		return fmt.Errorf("establish ctrl internal_connect_types failed with [%s]", err.Error())
	}

	// Handshake Step 1
	localSeq := auth.GenerateRandomInt()
	hs1b := connect2.EstablishCtrlConnectionHandshakeStep1Body{
		SecretKey: secretKey,
		Seq:       localSeq,
	}
	err = e.Encode(&hs1b)
	if err != nil {
		return fmt.Errorf("establish ctrl internal_connect_types failed in step 1 with [%s]", err.Error())
	}

	// Handshake Step 2
	hs2b := connect2.EstablishCtrlConnectionHandshakeStep2Body{}
	err = d.Decode(&hs2b)
	if err != nil {
		return fmt.Errorf("establish ctrl internal_connect_types failed in step 2 with [%s]", err.Error())
	}
	if localSeq+1 != hs2b.Ack {
		return fmt.Errorf("establish ctrl internal_connect_types failed with [wrong ack]")
	}
	s.uuid = hs2b.UUID
	s.token = hs2b.Token

	// Handshake Step 3
	hs3b := connect2.EstablishCtrlConnectionHandshakeStep3Body{
		Ack: hs2b.Seq + 1,
	}
	err = e.Encode(&hs3b)
	if err != nil {
		return fmt.Errorf("establish ctrl internal_connect_types failed in step 3 with [%s]", err.Error())
	}

	s.ctrlConn = conn
	s.ctrlDecoder = d
	s.ctrlEncoder = e
	return nil
}

func (s *Slave) establishDataConn(ip, port, uuid, token string) error {
	conn, err := network.CreateConn(ip, port)
	if err != nil {
		return fmt.Errorf("internal_connect_types to : %s:%s failed with [%s]", ip, port, err.Error())
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	e := json.NewEncoder(conn)
	d := json.NewDecoder(conn)

	// Mark the purpose of the internal_connect_types
	err = e.Encode(&connect2.ConnectionHead{ConnectionType: connect2.ConnectionTypeEstablishDataConnection})
	if err != nil {
		return fmt.Errorf("establish internal_data_types internal_connect_types failed with [%s]", err.Error())
	}

	// Establish Data Connection Handshake Step 1
	hs1b := connect2.EstablishDataConnectionHandShakeStep1Body{
		UUID:  uuid,
		Token: token,
	}
	err = e.Encode(&hs1b)
	if err != nil {
		return fmt.Errorf("establish internal_data_types internal_connect_types failed in step 1 with [%s]", err.Error())
	}

	// Establish Data Connection Handshake Step 2
	hs2b := connect2.EstablishDataConnectionHandShakeStep2Body{}
	err = d.Decode(&hs2b)
	if err != nil {
		return fmt.Errorf("establish internal_data_types internal_connect_types failed in step 2 with [%s]", err.Error())
	}

	s.dataConn = conn
	s.dataDecoder = d
	s.dataEncoder = e
	return nil
}

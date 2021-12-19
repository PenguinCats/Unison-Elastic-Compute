package slave

import (
	"encoding/json"
	"fmt"
	slave2 "github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/auth"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/network"
	connect2 "github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_connect_types"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func (s *Slave) register() error {
	if err := s.establishCtrlConn(); err != nil {
		return err
	}

	if err := s.establishDataConn(); err != nil {
		return err
	}

	s.status = slave2.StatsNormal
	return nil
}

func (s *Slave) establishCtrlConn() error {
	conn, err := network.CreateConn(s.masterIP, s.masterPort)
	if err != nil {
		return fmt.Errorf("internal_connect_types to : %s:%s failed with [%s]", s.masterIP, s.masterPort, err.Error())
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	ctrlEncoder := json.NewEncoder(conn)
	ctrlDecoder := json.NewDecoder(conn)

	// Mark the purpose of the internal_connect_types
	err = ctrlEncoder.Encode(&connect2.ConnectionHead{ConnectionType: connect2.ConnectionTypeEstablishCtrlConnection})
	if err != nil {
		err = fmt.Errorf("establish ctrl internal_connect_types failed with [%s]", err.Error())
		return err
	}

	// Handshake Step 1
	localSeq := auth.GenerateRandomInt()
	hs1b := connect2.EstablishCtrlConnectionHandshakeStep1Body{
		SecretKey: s.joinSecretKey,
		Seq:       localSeq,
	}
	err = ctrlEncoder.Encode(&hs1b)
	if err != nil {
		return fmt.Errorf("establish ctrl internal_connect_types failed in step 1 with [%s]", err.Error())
	}

	// Handshake Step 2
	hs2b := connect2.EstablishCtrlConnectionHandshakeStep2Body{}
	err = ctrlDecoder.Decode(&hs2b)
	if err != nil {
		return fmt.Errorf("establish ctrl internal_connect_types failed in step 2 with [%s]", err.Error())
	}
	if localSeq+1 != hs2b.Ack {
		err = fmt.Errorf("establish ctrl internal_connect_types failed with [wrong ack]")
		return err
	}
	s.uuid = hs2b.UUID
	s.token = hs2b.Token

	// Handshake Step 3
	hs3b := connect2.EstablishCtrlConnectionHandshakeStep3Body{
		Ack: hs2b.Seq + 1,
	}
	err = ctrlEncoder.Encode(&hs3b)
	if err != nil {
		return fmt.Errorf("establish ctrl internal_connect_types failed in step 3 with [%s]", err.Error())
	}

	s.ctrlConn = conn
	s.ctrlDecoder = ctrlDecoder
	s.ctrlEncoder = ctrlEncoder

	batch := new(leveldb.Batch)
	batch.Put([]byte("uec:token"), []byte(s.token))
	batch.Put([]byte("uec:uuid"), []byte(s.uuid))
	err = s.db.Write(batch, &opt.WriteOptions{
		Sync: true,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Slave) establishDataConn() error {
	conn, err := network.CreateConn(s.masterIP, s.masterPort)
	if err != nil {
		return fmt.Errorf("internal_connect_types to : %s:%s failed with [%s]", s.masterIP, s.masterPort, err.Error())
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	dataEncoder := json.NewEncoder(conn)
	dataDecoder := json.NewDecoder(conn)

	// Mark the purpose of the internal_connect_types
	err = dataEncoder.Encode(&connect2.ConnectionHead{ConnectionType: connect2.ConnectionTypeEstablishDataConnection})
	if err != nil {
		return fmt.Errorf("establish internal_data_types internal_connect_types failed with [%s]", err.Error())
	}

	// Establish Data Connection Handshake Step 1
	hs1b := connect2.EstablishDataConnectionHandShakeStep1Body{
		UUID:     s.uuid,
		Token:    s.token,
		HostInfo: s.dc.GetHostInfo(),
	}
	err = dataEncoder.Encode(&hs1b)
	if err != nil {
		return fmt.Errorf("establish internal_data_types internal_connect_types failed in step 1 with [%s]", err.Error())
	}

	// Establish Data Connection Handshake Step 2
	hs2b := connect2.EstablishDataConnectionHandShakeStep2Body{}
	err = dataDecoder.Decode(&hs2b)
	if err != nil {
		return fmt.Errorf("establish internal_data_types internal_connect_types failed in step 2 with [%s]", err.Error())
	}

	s.dataConn = conn
	s.dataDecoder = dataDecoder
	s.dataEncoder = dataEncoder
	return nil
}

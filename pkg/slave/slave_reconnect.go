package slave

import (
	"encoding/json"
	"errors"
	"fmt"
	slave2 "github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/auth"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/network"
	connect2 "github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_connect_types"
)

func (s *Slave) reconnect() error {
	if err := s.reconnectCtrlConn(); err != nil {
		return err
	}

	if err := s.establishDataConn(); err != nil {
		return err
	}

	s.status = slave2.StatsNormal
	return nil
}

func (s *Slave) reconnectCtrlConn() error {
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
	err = ctrlEncoder.Encode(&connect2.ConnectionHead{ConnectionType: connect2.ConnectionTypeReconnectCtrlConnection})
	if err != nil {
		return fmt.Errorf("reconnect ctrl internal_connect_types failed with [%s]", err.Error())
	}

	// Handshake Step 1
	localSeq := auth.GenerateRandomInt()
	hs1b := connect2.ReconnectCtrlConnectionHandshakeStep1Body{
		Seq:   localSeq,
		Token: s.token,
		UUID:  s.uuid,
	}
	err = ctrlEncoder.Encode(&hs1b)
	if err != nil {
		return fmt.Errorf("reconnect ctrl internal_connect_types failed in step 1 with [%s]", err.Error())
	}

	// Handshake Step 2
	hs2b := connect2.ReconnectCtrlConnectionHandshakeStep2Body{}
	err = ctrlDecoder.Decode(&hs2b)
	if err != nil {
		return fmt.Errorf("establish ctrl internal_connect_types failed in step 2 with [%s]", err.Error())
	}
	if localSeq+1 != hs2b.Ack {
		err = fmt.Errorf("establish ctrl internal_connect_types failed with [wrong ack]")
		return err
	}
	if !hs2b.Agree {
		err = errors.New("master node refuse to reconnect")
		return err
	}

	s.ctrlConn = conn
	s.ctrlDecoder = ctrlDecoder
	s.ctrlEncoder = ctrlEncoder
	return nil
}

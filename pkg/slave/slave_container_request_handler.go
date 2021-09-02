package slave

import (
	"encoding/json"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_data_types"
	"github.com/sirupsen/logrus"
)

func (s *Slave) handleContainerCreateMessage(v []byte) {
	m := internal_data_types.ContainerCreateMessage{}
	err := json.Unmarshal(v, &m)
	if err != nil {
		return
	}

	resp := internal_data_types.ContainerCreateResponse{
		OperationID:    m.OperationID,
		ExtContainerID: m.CCB.ExtContainerID,
		Error:          nil,
	}
	defer func() {
		resV, err := json.Marshal(&resp)
		if err != nil {
			logrus.Warning(err.Error())
		} else {
			s.sendDataMessage(internal_data_types.MessageDataTypeContainerCreate, resV)
		}
	}()

	_, err = s.dc.ContainerCreate(m.CCB)
	if err != nil {
		resp.Error = err
		return
	}

	profile, err := s.dc.ContainerProfile(m.CCB.ExtContainerID)
	if err != nil {
		resp.Error = err
		return
	}

	resp.Profile = profile
}

func (s *Slave) handleContainerStartMessage(v []byte) {
	m := internal_data_types.ContainerStartMessage{}
	err := json.Unmarshal(v, &m)
	if err != nil {
		return
	}

	resp := internal_data_types.ContainerStartResponse{
		OperationID:    m.OperationID,
		ExtContainerID: m.ExtContainerID,
		Error:          nil,
	}
	defer func() {
		resV, err := json.Marshal(&resp)
		if err != nil {
			logrus.Warning(err.Error())
		} else {
			s.sendDataMessage(internal_data_types.MessageDataTypeContainerStart, resV)
		}
	}()

	err = s.dc.ContainerStart(m.ExtContainerID)
	if err != nil {
		resp.Error = err
		return
	}
}

func (s *Slave) handleContainerStopMessage(v []byte) {
	m := internal_data_types.ContainerStopMessage{}
	err := json.Unmarshal(v, &m)
	if err != nil {
		return
	}

	resp := internal_data_types.ContainerStopResponse{
		OperationID:    m.OperationID,
		Error:          nil,
		ExtContainerID: m.ExtContainerID,
	}
	defer func() {
		resV, err := json.Marshal(&resp)
		if err != nil {
			logrus.Warning(err.Error())
		}
		s.sendDataMessage(internal_data_types.MessageDataTypeContainerStop, resV)
	}()

	err = s.dc.ContainerStop(m.ExtContainerID)
	if err != nil {
		resp.Error = err
		return
	}
}

func (s *Slave) handleContainerRemoveMessage(v []byte) {
	m := internal_data_types.ContainerRemoveMessage{}
	err := json.Unmarshal(v, &m)
	if err != nil {
		return
	}

	resp := internal_data_types.ContainerRemoveResponse{
		OperationID:    m.OperationID,
		Error:          nil,
		ExtContainerID: m.ExtContainerID,
	}
	defer func() {
		resV, err := json.Marshal(&resp)
		if err != nil {
			logrus.Warning(err.Error())
		}
		s.sendDataMessage(internal_data_types.MessageDataTypeContainerRemove, resV)
	}()

	err = s.dc.ContainerRemove(m.ExtContainerID)
	if err != nil {
		resp.Error = err
		return
	}
}

//func (s *Slave) handleContainerProfileMessage(v []byte) {
//	m := internal_data_types.ContainerProfileMessage{}
//	err := json.Unmarshal(v, &m)
//	if err != nil {
//		return
//	}
//
//	resp := internal_data_types.ContainerProfileResponse{
//		OperationUUID: m.OperationUUID,
//		ContainerID: m.ContainerID,
//		Profile:     container.ContainerProfile{},
//		Flag:        internal_data_types.ContainerProfileFail,
//	}
//	defer func() {
//		resV, err := json.Marshal(&resp)
//		if err != nil {
//			logrus.Warning(err.Error())
//		}
//		s.sendDataMessage(internal_data_types.MessageDataTypeContainerProfile, resV)
//	}()
//
//	profile, err := s.dc.ContainerProfile(m.ContainerID)
//	if err != nil {
//		resp.Profile = profile
//		resp.Flag = internal_data_types.ContainerProfileSuccess
//	}
//}
//
//func (s *Slave) handleContainerStatusMessage(v []byte) {
//	m := internal_data_types.ContainerStatusMessage{}
//	err := json.Unmarshal(v, &m)
//	if err != nil {
//		return
//	}
//
//	resp := internal_data_types.ContainerStatusResponse{
//		OperationUUID: m.OperationUUID,
//		ContainerID: m.ContainerID,
//		Stats:      container.ContainerStatus{},
//		Flag:        internal_data_types.ContainerStatusFail,
//	}
//	defer func() {
//		resV, err := json.Marshal(&resp)
//		if err != nil {
//			logrus.Warning(err.Error())
//		}
//		s.sendDataMessage(internal_data_types.MessageDataTypeContainerStatus, resV)
//	}()
//
//	status, err := s.dc.ContainerStats(m.ContainerID)
//	if err != nil {
//		resp.Stats = status
//		resp.Flag = internal_data_types.ContainerStatusSuccess
//	}
//}

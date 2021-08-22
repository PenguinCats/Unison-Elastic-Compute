package master

import (
	"encoding/json"
	types2 "github.com/PenguinCats/Unison-Docker-Controller/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/http-controller"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/operation"
	"github.com/sirupsen/logrus"
)

func (m *Master) handleOperationContainerCreateTask(task operation.OperationContainerCreateTask) {
	var errCode = types.SUCCESS
	defer func() {
		if errCode != types.SUCCESS {
			oprInfo := operation.OprInfoUtil.GetOptInfo(task.OperationID)
			resp := types.APIContainerCreateResponse{
				APIResponseBase: types.APIResponseBase{
					OperationID: task.OperationID,
					Code:        errCode,
					Msg:         types.GetMsg(errCode),
				},
				ExposedTCPPorts:        nil,
				ExposedTCPMappingPorts: nil,
				ExposedUDPPorts:        nil,
				ExposedUDPMappingPorts: nil,
			}
			sendDataByte, jsonErr := json.Marshal(resp)
			if jsonErr != nil {
				logrus.Warning(jsonErr.Error())
				return
			}
			http_controller.SendCallback(oprInfo.CallbackURL, sendDataByte)
		}
	}()

	scb, ok := m.slaveController.GetSlaveCtrlBlk(task.SlaveID)
	if !ok {
		errCode = types.SLAVE_INVALID
		return
	}

	err := scb.SendDataContainerCreateMsg(task.ContainerCreateMessage)
	if err != nil {
		errCode = types.ERROR
		return
	}
}

func (m *Master) handleOperationContainerCreateResponse(resp operation.OperationContainerCreateResponse) {
	oprInfo := operation.OprInfoUtil.GetOptInfo(resp.OperationID)

	var errCode = types.SUCCESS
	if resp.Error != nil {
		switch resp.Error {
		case types2.ErrInternalError:
			errCode = types.ERROR
		case types2.ErrInsufficientResource:
			errCode = types.SLAVE_INSUFFICIENT_RESOURCE
		default:
			errCode = types.UNKNOWN_ERROR
		}
	}

	response := types.APIContainerCreateResponse{
		APIResponseBase: types.APIResponseBase{
			OperationID: resp.OperationID,
			Code:        errCode,
			Msg:         types.GetMsg(errCode),
		},
		ExposedTCPPorts:        nil,
		ExposedTCPMappingPorts: nil,
		ExposedUDPPorts:        nil,
		ExposedUDPMappingPorts: nil,
	}

	if resp.Error == nil {
		response.ExposedTCPPorts = resp.Profile.ExposedTCPPorts
		response.ExposedTCPMappingPorts = resp.Profile.ExposedTCPMappingPorts
		response.ExposedUDPPorts = resp.Profile.ExposedUDPPorts
		response.ExposedUDPMappingPorts = resp.Profile.ExposedUDPMappingPorts
	}

	sendDataByte, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		logrus.Warning(jsonErr.Error())
		return
	}
	http_controller.SendCallback(oprInfo.CallbackURL, sendDataByte)
}

func (m *Master) handleOperationContainerStartTask(task operation.OperationContainerStartTask) {
	var errCode = types.SUCCESS
	defer func() {
		if errCode != types.SUCCESS {
			oprInfo := operation.OprInfoUtil.GetOptInfo(task.OperationID)
			resp := types.APIContainerStartResponse{
				APIResponseBase: types.APIResponseBase{
					OperationID: task.OperationID,
					Code:        errCode,
					Msg:         types.GetMsg(errCode),
				},
			}
			sendDataByte, jsonErr := json.Marshal(resp)
			if jsonErr != nil {
				logrus.Warning(jsonErr.Error())
				return
			}
			http_controller.SendCallback(oprInfo.CallbackURL, sendDataByte)
		}
	}()

	scb, ok := m.slaveController.GetSlaveCtrlBlk(task.SlaveID)
	if !ok {
		errCode = types.SLAVE_INVALID
		return
	}

	err := scb.SendDataContainerStartMsg(task.ContainerStartMessage)
	if err != nil {
		errCode = types.ERROR
		return
	}
}

func (m *Master) handleOperationContainerStartResponse(resp operation.OperationContainerStartResponse) {
	oprInfo := operation.OprInfoUtil.GetOptInfo(resp.OperationID)

	var errCode = types.SUCCESS
	if resp.Error != nil {
		switch resp.Error {
		case types2.ErrInternalError:
			errCode = types.ERROR
		case types2.ErrInsufficientResource:
			errCode = types.SLAVE_INSUFFICIENT_RESOURCE
		default:
			errCode = types.UNKNOWN_ERROR
		}
	}

	response := types.APIContainerStartResponse{
		APIResponseBase: types.APIResponseBase{
			OperationID: resp.OperationID,
			Code:        errCode,
			Msg:         types.GetMsg(errCode),
		},
	}

	sendDataByte, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		logrus.Warning(jsonErr.Error())
		return
	}
	http_controller.SendCallback(oprInfo.CallbackURL, sendDataByte)
}

func (m *Master) handleOperationContainerStopTask(task operation.OperationContainerStopTask) {
	var errCode = types.SUCCESS
	defer func() {
		if errCode != types.SUCCESS {
			oprInfo := operation.OprInfoUtil.GetOptInfo(task.OperationID)
			resp := types.APIContainerStopResponse{
				APIResponseBase: types.APIResponseBase{
					OperationID: task.OperationID,
					Code:        errCode,
					Msg:         types.GetMsg(errCode),
				},
			}
			sendDataByte, jsonErr := json.Marshal(resp)
			if jsonErr != nil {
				logrus.Warning(jsonErr.Error())
				return
			}
			http_controller.SendCallback(oprInfo.CallbackURL, sendDataByte)
		}
	}()

	scb, ok := m.slaveController.GetSlaveCtrlBlk(task.SlaveID)
	if !ok {
		errCode = types.SLAVE_INVALID
		return
	}

	err := scb.SendDataContainerStopMsg(task.ContainerStopMessage)
	if err != nil {
		errCode = types.ERROR
		return
	}
}

func (m *Master) handleOperationContainerStopResponse(resp operation.OperationContainerStopResponse) {
	oprInfo := operation.OprInfoUtil.GetOptInfo(resp.OperationID)

	var errCode = types.SUCCESS
	if resp.Error != nil {
		switch resp.Error {
		case types2.ErrInternalError:
			errCode = types.ERROR
		default:
			errCode = types.UNKNOWN_ERROR
		}
	}

	response := types.APIContainerStopResponse{
		APIResponseBase: types.APIResponseBase{
			OperationID: resp.OperationID,
			Code:        errCode,
			Msg:         types.GetMsg(errCode),
		},
	}

	sendDataByte, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		logrus.Warning(jsonErr.Error())
		return
	}
	http_controller.SendCallback(oprInfo.CallbackURL, sendDataByte)
}

func (m *Master) handleOperationContainerRemoveTask(task operation.OperationContainerRemoveTask) {
	var errCode = types.SUCCESS
	defer func() {
		if errCode != types.SUCCESS {
			oprInfo := operation.OprInfoUtil.GetOptInfo(task.OperationID)
			resp := types.APIContainerRemoveResponse{
				APIResponseBase: types.APIResponseBase{
					OperationID: task.OperationID,
					Code:        errCode,
					Msg:         types.GetMsg(errCode),
				},
			}
			sendDataByte, jsonErr := json.Marshal(resp)
			if jsonErr != nil {
				logrus.Warning(jsonErr.Error())
				return
			}
			http_controller.SendCallback(oprInfo.CallbackURL, sendDataByte)
		}
	}()

	scb, ok := m.slaveController.GetSlaveCtrlBlk(task.SlaveID)
	if !ok {
		errCode = types.SLAVE_INVALID
		return
	}

	err := scb.SendDataContainerRemoveMsg(task.ContainerRemoveMessage)
	if err != nil {
		errCode = types.ERROR
		return
	}
}

func (m *Master) handleOperationContainerRemoveResponse(resp operation.OperationContainerRemoveResponse) {
	oprInfo := operation.OprInfoUtil.GetOptInfo(resp.OperationID)

	var errCode = types.SUCCESS
	if resp.Error != nil {
		switch resp.Error {
		case types2.ErrInternalError:
			errCode = types.ERROR
		default:
			errCode = types.UNKNOWN_ERROR
		}
	}

	response := types.APIContainerRemoveResponse{
		APIResponseBase: types.APIResponseBase{
			OperationID: resp.OperationID,
			Code:        errCode,
			Msg:         types.GetMsg(errCode),
		},
	}

	sendDataByte, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		logrus.Warning(jsonErr.Error())
		return
	}
	http_controller.SendCallback(oprInfo.CallbackURL, sendDataByte)
}

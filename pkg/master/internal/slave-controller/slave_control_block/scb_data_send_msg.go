package slave_control_block

import (
	"encoding/json"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_data_types"
	"github.com/sirupsen/logrus"
)

func (scb *SlaveControlBlock) SendDataContainerCreateMsg(msg internal_data_types.ContainerCreateMessage) error {
	v, err := json.Marshal(&msg)
	if err != nil {
		logrus.Warning(err.Error())
		return ErrJsonEncode
	}

	err = scb.dataEncoder.Encode(&internal_data_types.Message{
		MessageType: internal_data_types.MessageDataTypeContainerCreate,
		Value:       v,
	})
	if err != nil {
		logrus.Warning(err.Error())
		return ErrSendDataMsg
	}

	return nil
}

func (scb *SlaveControlBlock) SendDataContainerStartMsg(msg internal_data_types.ContainerStartMessage) error {
	v, err := json.Marshal(&msg)
	if err != nil {
		logrus.Warning(err.Error())
		return ErrJsonEncode
	}

	err = scb.dataEncoder.Encode(&internal_data_types.Message{
		MessageType: internal_data_types.MessageDataTypeContainerStart,
		Value:       v,
	})
	if err != nil {
		logrus.Warning(err.Error())
		return ErrSendDataMsg
	}

	return nil
}

func (scb *SlaveControlBlock) SendDataContainerStopMsg(msg internal_data_types.ContainerStopMessage) error {
	v, err := json.Marshal(&msg)
	if err != nil {
		logrus.Warning(err.Error())
		return ErrJsonEncode
	}

	err = scb.dataEncoder.Encode(&internal_data_types.Message{
		MessageType: internal_data_types.MessageDataTypeContainerStop,
		Value:       v,
	})
	if err != nil {
		logrus.Warning(err.Error())
		return ErrSendDataMsg
	}

	return nil
}

func (scb *SlaveControlBlock) SendDataContainerRemoveMsg(msg internal_data_types.ContainerRemoveMessage) error {
	v, err := json.Marshal(&msg)
	if err != nil {
		logrus.Warning(err.Error())
		return ErrJsonEncode
	}

	err = scb.dataEncoder.Encode(&internal_data_types.Message{
		MessageType: internal_data_types.MessageDataTypeContainerRemove,
		Value:       v,
	})
	if err != nil {
		logrus.Warning(err.Error())
		return ErrSendDataMsg
	}

	return nil
}

package slave_control_block

import (
	"context"
	"encoding/json"
	"github.com/PenguinCats/Unison-Docker-Controller/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_data_types"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/operation"
	"github.com/sirupsen/logrus"
)

func (scb *SlaveControlBlock) startHandleDataMessage(ctx context.Context) {
	go func() {
		var err error = nil
		defer func() {
			if err != nil {
				logrus.Warning(err.Error())
				scb.offline()
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				message := internal_data_types.Message{}
				err = scb.ctrlDecoder.Decode(&message)
				if err != nil {
					logrus.Warning(internal_data_types.ErrDataInvalidMessage.Error())
					err = types.ErrInternalError
					return
				}

				switch message.MessageType {
				case internal_data_types.MessageDataTypeContainerCreate:
					go scb.handleContainerCreateResponse(message.Value)
				case internal_data_types.MessageDataTypeContainerStart:
					go scb.handleContainerStartResponse(message.Value)
				default:
					logrus.Warning(internal_data_types.ErrDataInvalidMessage.Error())
					err = types.ErrInternalError
					return
				}
			}
		}
	}()
}

func (scb *SlaveControlBlock) handleContainerCreateResponse(v []byte) {
	m := internal_data_types.ContainerCreateResponse{}
	err := json.Unmarshal(v, &m)
	if err != nil {
		logrus.Warning(internal_data_types.ErrDataInvalidContainerCreateMessage.Error())
		return
	}

	scb.operationResponseChan <- &operation.OperationResponse{
		OperationID: m.OperationID,
		OperationResponseBody: operation.OperationContainerCreateResponse{
			OperationID:    m.OperationID,
			Error:          m.Error,
			UECContainerID: m.ExtContainerID,
			Profile:        m.Profile,
		},
	}
}

func (scb *SlaveControlBlock) handleContainerStartResponse(v []byte) {
	m := internal_data_types.ContainerStartResponse{}
	err := json.Unmarshal(v, &m)
	if err != nil {
		logrus.Warning(internal_data_types.ErrDataInvalidContainerStartMessage.Error())
		return
	}

	scb.operationResponseChan <- &operation.OperationResponse{
		OperationID: m.OperationID,
		OperationResponseBody: operation.OperationContainerStartResponse{
			OperationID:    m.OperationID,
			Error:          m.Error,
			UECContainerID: m.ExtContainerID,
		},
	}
}

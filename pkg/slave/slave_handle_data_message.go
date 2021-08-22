package slave

import (
	"context"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_data_types"
	"github.com/sirupsen/logrus"
)

func (s *Slave) startHandleDataMessage(ctx context.Context) {
	go func() {
		var err error = nil
		defer func() {
			if err != nil {
				logrus.Warning(err.Error())
				s.offline()
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				message := internal_data_types.Message{}
				err = s.dataDecoder.Decode(&message)
				if err != nil {
					err = internal_data_types.ErrDataInvalidMessage
					return
				}

				switch message.MessageType {
				case internal_data_types.MessageDataTypeContainerCreate:
					s.handleContainerCreateMessage(message.Value)
				case internal_data_types.MessageDataTypeContainerStart:
					s.handleContainerStartMessage(message.Value)
				case internal_data_types.MessageDataTypeContainerStop:
					s.handleContainerStopMessage(message.Value)
				case internal_data_types.MessageDataTypeContainerRemove:
					s.handleContainerRemoveMessage(message.Value)
				//case internal_data_types.MessageDataTypeContainerProfile:
				//	s.handleContainerProfileMessage(message.Value)
				//case internal_data_types.MessageDataTypeContainerStatus:
				//	s.handleContainerStatusMessage(message.Value)

				default:
					err = internal_data_types.ErrDataInvalidMessage
					return
				}
			}
		}
	}()
}

func (s *Slave) sendDataMessage(dataType internal_data_types.MessageDataType, v []byte) {
	message := internal_data_types.Message{
		MessageType: dataType,
		Value:       v,
	}
	err := s.dataEncoder.Encode(&message)
	if err != nil {
		logrus.Warning(err.Error())
	}
}

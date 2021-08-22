package slave

import (
	"context"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_control_types"
	"github.com/sirupsen/logrus"
)

func (s *Slave) startHandleCtrlMessage(ctx context.Context) {
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
				message := internal_control_types.Message{}
				err = s.ctrlDecoder.Decode(&message)
				if err != nil {
					err = internal_control_types.ErrControlInvalidMessage
					return
				}

				switch message.MessageType {
				case internal_control_types.MessageCtrlTypeHeartbeat:
					s.handleHeartbeatMessage(message.Value)
				default:
					err = internal_control_types.ErrControlInvalidMessage
					return
				}
			}
		}
	}()
}

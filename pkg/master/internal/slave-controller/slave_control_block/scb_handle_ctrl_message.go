package slave_control_block

import (
	"context"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/control"
	"github.com/sirupsen/logrus"
)

func (scb *SlaveControlBlock) startHandleCtrlMessage(ctx context.Context) {
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
				message := control.Message{}
				err = scb.ctrlDecoder.Decode(&message)
				if err != nil {
					err = control.ErrControlInvalidMessage
					return
				}

				switch message.MessageType {
				case control.MessageCtrlTypeHeartbeat:
					scb.handleHeartbeatMessage(message.Value)
				default:
					err = control.ErrControlInvalidMessage
					return
				}
			}
		}
	}()
}

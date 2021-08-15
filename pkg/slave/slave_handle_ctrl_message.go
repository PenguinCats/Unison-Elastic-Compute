package slave

import (
	"context"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/control"
	"github.com/sirupsen/logrus"
)

func (slave *Slave) startHandleCtrlMessage(ctx context.Context) {
	go func() {
		var err error = nil
		defer func() {
			if err != nil {
				logrus.Warning(err.Error())
				slave.offline()
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				message := control.Message{}
				err = slave.ctrlDecoder.Decode(&message)
				if err != nil {
					err = control.ErrControlInvalidMessage
					return
				}

				switch message.MessageType {
				case control.MessageCtrlTypeHeartbeat:
					slave.handleHeartbeatMessage(message.Value)
				default:
					err = control.ErrControlInvalidMessage
					return
				}
			}
		}
	}()
}

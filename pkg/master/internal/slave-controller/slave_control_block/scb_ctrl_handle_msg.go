package slave_control_block

import (
	"context"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_control_types"
	"github.com/sirupsen/logrus"
	"io"
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
				message := internal_control_types.Message{}
				err = scb.ctrlDecoder.Decode(&message)
				if err != nil {
					if err == io.EOF {

					} else {
						err = internal_control_types.ErrControlInvalidMessage
					}
					return
				}

				switch message.MessageType {
				case internal_control_types.MessageCtrlTypeHeartbeat:
					go scb.handleHeartbeatMessage(message.Value)
				default:
					err = internal_control_types.ErrControlInvalidMessage
					return
				}
			}
		}
	}()
}

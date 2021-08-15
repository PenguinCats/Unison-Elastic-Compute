package slave

import (
	"context"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/data"
	"github.com/sirupsen/logrus"
)

func (slave *Slave) startHandleDataMessage(ctx context.Context) {
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
				message := data.Message{}
				err = slave.dataDecoder.Decode(&message)
				if err != nil {
					err = data.ErrDataInvalidMessage
					return
				}

				switch message.MessageType {

				default:
					err = data.ErrDataInvalidMessage
					return
				}
			}
		}
	}()
}

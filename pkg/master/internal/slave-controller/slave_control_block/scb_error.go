package slave_control_block

import "errors"

var (
	ErrJsonEncode = errors.New("json encode error")

	ErrSendDataMsg = errors.New("send data message fail")

	ErrSendCtrlMsg = errors.New("send control message fail")
)

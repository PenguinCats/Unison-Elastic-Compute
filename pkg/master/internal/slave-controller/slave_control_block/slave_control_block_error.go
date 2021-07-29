package slave_control_block

import "errors"

var (
	ErrControlInvalidMessage   = errors.New("invalid control message")
	ErrControlInvalidHeartbeat = errors.New("invalid heartbeat message")
)

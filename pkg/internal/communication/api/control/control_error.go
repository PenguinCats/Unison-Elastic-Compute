package control

import "errors"

var (
	ErrControlInvalidMessage   = errors.New("invalid control_types message")
	ErrControlInvalidHeartbeat = errors.New("invalid heartbeat message")
)

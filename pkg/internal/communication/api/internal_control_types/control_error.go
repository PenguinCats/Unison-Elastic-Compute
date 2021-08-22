package internal_control_types

import "errors"

var (
	ErrControlInvalidMessage   = errors.New("invalid control_types message")
	ErrControlInvalidHeartbeat = errors.New("invalid heartbeat message")
)

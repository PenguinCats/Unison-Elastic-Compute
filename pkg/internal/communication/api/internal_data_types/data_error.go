package internal_data_types

import "errors"

var (
	ErrDataInvalidMessage                = errors.New("invalid data type message")
	ErrDataInvalidContainerCreateMessage = errors.New("invalid container create message")
	ErrDataInvalidContainerStartMessage  = errors.New("invalid container start message")
)

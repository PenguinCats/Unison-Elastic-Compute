package operation

import "errors"

var (
	ErrOperationTask     = errors.New("invalid operation task")
	ErrOperationResponse = errors.New("invalid operation response")
)

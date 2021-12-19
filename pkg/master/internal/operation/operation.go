package operation

type OperationType int

type OperationTask struct {
	OperationID int64
	CallbackURL string

	OperationTaskBody interface{}
}

type OperationResponse struct {
	OperationID           int64
	OperationResponseBody interface{}
}

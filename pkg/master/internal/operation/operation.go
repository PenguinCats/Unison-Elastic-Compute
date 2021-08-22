package operation

type OperationType int

type OperationTask struct {
	OperationID string
	CallbackURL string

	OperationTaskBody interface{}
}

type OperationResponse struct {
	OperationID           string
	OperationResponseBody interface{}
}

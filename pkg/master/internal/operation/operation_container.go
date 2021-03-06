package operation

import (
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/container"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_data_types"
)

type OperationContainerCreateTask struct {
	OperationID int64

	SlaveID string

	ContainerCreateMessage internal_data_types.ContainerCreateMessage
}

type OperationContainerCreateResponse struct {
	OperationID int64
	Error       error

	UECContainerID string
	Profile        container.ContainerProfile
}

type OperationContainerStartTask struct {
	OperationID int64

	SlaveID string

	internal_data_types.ContainerStartMessage
}

type OperationContainerStartResponse struct {
	OperationID    int64
	Error          error
	UECContainerID string
}

type OperationContainerStopTask struct {
	OperationID int64

	SlaveID string

	internal_data_types.ContainerStopMessage
}

type OperationContainerStopResponse struct {
	OperationID    int64
	Error          error
	UECContainerID string
}

type OperationContainerRemoveTask struct {
	OperationID int64

	SlaveID string

	internal_data_types.ContainerRemoveMessage
}

type OperationContainerRemoveResponse struct {
	OperationID    int64
	Error          error
	UECContainerID string
}

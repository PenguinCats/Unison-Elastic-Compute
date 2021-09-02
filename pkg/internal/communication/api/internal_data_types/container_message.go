package internal_data_types

import "github.com/PenguinCats/Unison-Docker-Controller/api/types/container"

type ContainerCreateMessage struct {
	OperationID string

	CCB container.ContainerCreateBody
}

type ContainerCreateResponse struct {
	OperationID string
	Error       error

	ExtContainerID string

	Profile container.ContainerProfile
}

type ContainerStartMessage struct {
	OperationID string

	ExtContainerID string
}

type ContainerStartResponse struct {
	OperationID string
	Error       error

	ExtContainerID string
}

type ContainerStopMessage struct {
	OperationID string

	ExtContainerID string
}

type ContainerStopResponse struct {
	OperationID string
	Error       error

	ExtContainerID string
}

type ContainerRemoveMessage struct {
	OperationID string

	ExtContainerID string
}

type ContainerRemoveResponse struct {
	OperationID string
	Error       error

	ExtContainerID string
}

//type ContainerProfileMessage struct {
//	OperationID string
//	ContainerID string
//}
//
//type ContainerProfileResultFlag int
//
//const (
//	ContainerProfileSuccess ContainerProfileResultFlag = iota
//	ContainerProfileFail
//)
//
//type ContainerProfileResponse struct {
//	OperationID string
//	ContainerID string
//	Profile     container.ContainerProfile
//	Flag        ContainerProfileResultFlag
//}
//
//type ContainerStatusMessage struct {
//	OperationID string
//	ContainerID string
//}
//
//type ContainerStatusResultFlag int
//
//const (
//	ContainerStatusSuccess ContainerStatusResultFlag = iota
//	ContainerStatusFail
//)
//
//type ContainerStatusResponse struct {
//	OperationID string
//	ContainerID string
//	Stats      container.ContainerStatus
//	Flag        ContainerStatusResultFlag
//}

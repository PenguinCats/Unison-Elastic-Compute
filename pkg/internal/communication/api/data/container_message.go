package data

import "github.com/PenguinCats/Unison-Docker-Controller/api/types/container"

type ContainerCreateMessage struct {
	containerUUID string // this UUID is an external information for the outer system, not the docker Container ID.

	ccb container.ContainerCreateBody
}

type ContainerCreateResultFlag int

const (
	ContainerCreateSuccess ContainerCreateResultFlag = iota
	ContainerCreateFail
)

type ContainerCreateResponse struct {
	containerUUID string // this UUID is an external information for the outer system, not the docker Container ID.

	flag ContainerCreateResultFlag
}

type ContainerStartMessage struct {
	containerID string
}

type ContainerStartResultFlag int

const (
	ContainerStartSuccess ContainerStartResultFlag = iota
	ContainerStartFail
)

type ContainerStartResponse struct {
	containerID string
	flag        ContainerStartResultFlag
}

type ContainerStopMessage struct {
	containerID string
}

type ContainerStopResultFlag int

const (
	ContainerStopSuccess ContainerStopResultFlag = iota
	ContainerStopFail
)

type ContainerStopResponse struct {
	containerID string
	flag        ContainerStopResultFlag
}

type ContainerRemoveMessage struct {
	containerID string
}

type ContainerRemoveResultFlag int

const (
	ContainerRemoveSuccess ContainerRemoveResultFlag = iota
	ContainerRemoveFail
)

type ContainerRemoveResponse struct {
	containerID string
	flag        ContainerRemoveResultFlag
}

type ContainerProfileMessage struct {
	containerID string
}

type ContainerProfileResultFlag int

const (
	ContainerProfileSuccess ContainerProfileResultFlag = iota
	ContainerProfileFail
)

type ContainerProfileResponse struct {
	containerID string
	profile     container.ContainerProfile
	flag        ContainerProfileResultFlag
}

type ContainerStatusMessage struct {
	containerID string
}

type ContainerStatusResultFlag int

const (
	ContainerStatusSuccess ContainerStatusResultFlag = iota
	ContainerStatusFail
)

type ContainerStatusResponse struct {
	containerID string
	profile     container.ContainerStatus
	flag        ContainerStatusResultFlag
}

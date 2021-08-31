package internal_control_types

import (
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/container"
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/resource"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
)

type HeartBeatMessageReport struct {
	Status   types.StatusSlave
	Resource resource.ResourceAvailable

	ContainerStatus map[string]container.ContainerStatus
}

type HeartBeatMessageAck struct {
}

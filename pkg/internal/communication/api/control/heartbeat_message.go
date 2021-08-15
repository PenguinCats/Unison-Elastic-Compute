package control

import (
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/resource"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
)

type HeartBeatMessageReport struct {
	Status   types.StatusSlave
	Resource resource.ResourceAvailable
}

type HeartBeatMessageAck struct {
}

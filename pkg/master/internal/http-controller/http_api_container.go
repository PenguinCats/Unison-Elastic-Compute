package http_controller

import (
	"github.com/PenguinCats/Unison-Docker-Controller/api/types/container"
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/http_wrapper"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/internal/communication/api/internal_data_types"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/operation"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (hac *HttpApiController) containerCreate(c *gin.Context) {
	var (
		appG = http_wrapper.Gin{C: c}
		form types.APIContainerCreateRequest
	)

	httpCode, errCode := http_wrapper.BindAndValid(c, &form)
	if errCode != types.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	hac.operationTaskChan <- &operation.OperationTask{
		OperationID: form.OperationID,
		CallbackURL: form.CallbackURL,
		OperationTaskBody: operation.OperationContainerCreateTask{
			OperationID: form.OperationID,
			SlaveID:     form.SlaveID,
			ContainerCreateMessage: internal_data_types.ContainerCreateMessage{
				OperationID: form.OperationID,
				CCB: container.ContainerCreateBody{
					ExtContainerID:  form.UECContainerID,
					ImageName:       form.ImageName,
					ExposedTCPPorts: form.ExposedTCPPorts,
					ExposedUDPPorts: form.ExposedUDPPorts,
					Mounts:          form.Mounts,
					CoreCnt:         form.CoreCnt,
					MemorySize:      form.MemorySize,
					StorageSize:     form.StorageSize,
				},
			},
		},
	}

	appG.Response(http.StatusOK, types.SUCCESS, nil)
}

func (hac *HttpApiController) containerStart(c *gin.Context) {
	var (
		appG = http_wrapper.Gin{C: c}
		form types.APIContainerStartRequest
	)

	httpCode, errCode := http_wrapper.BindAndValid(c, &form)
	if errCode != types.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	hac.operationTaskChan <- &operation.OperationTask{
		OperationID: form.OperationID,
		CallbackURL: form.CallbackURL,
		OperationTaskBody: operation.OperationContainerStartTask{
			OperationID: form.OperationID,
			SlaveID:     form.SlaveID,
			ContainerStartMessage: internal_data_types.ContainerStartMessage{
				OperationID:    form.OperationID,
				ExtContainerID: form.UECContainerID,
			},
		},
	}

	appG.Response(http.StatusOK, types.SUCCESS, nil)
}

func (hac *HttpApiController) containerStop(c *gin.Context) {
	var (
		appG = http_wrapper.Gin{C: c}
		form types.APIContainerStopRequest
	)

	httpCode, errCode := http_wrapper.BindAndValid(c, &form)
	if errCode != types.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	hac.operationTaskChan <- &operation.OperationTask{
		OperationID: form.OperationID,
		CallbackURL: form.CallbackURL,
		OperationTaskBody: operation.OperationContainerStopTask{
			OperationID: form.OperationID,
			SlaveID:     form.SlaveID,
			ContainerStopMessage: internal_data_types.ContainerStopMessage{
				OperationID:    form.OperationID,
				ExtContainerID: form.UECContainerID,
			},
		},
	}

	appG.Response(http.StatusOK, types.SUCCESS, nil)
}

func (hac *HttpApiController) containerRemove(c *gin.Context) {
	var (
		appG = http_wrapper.Gin{C: c}
		form types.APIContainerRemoveRequest
	)

	httpCode, errCode := http_wrapper.BindAndValid(c, &form)
	if errCode != types.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	hac.operationTaskChan <- &operation.OperationTask{
		OperationID: form.OperationID,
		CallbackURL: form.CallbackURL,
		OperationTaskBody: operation.OperationContainerRemoveTask{
			OperationID: form.OperationID,
			SlaveID:     form.SlaveID,
			ContainerRemoveMessage: internal_data_types.ContainerRemoveMessage{
				OperationID:    form.OperationID,
				ExtContainerID: form.UECContainerID,
			},
		},
	}

	appG.Response(http.StatusOK, types.SUCCESS, nil)
}

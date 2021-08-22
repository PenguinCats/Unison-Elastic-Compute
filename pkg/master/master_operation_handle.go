package master

import (
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/operation"
	"github.com/sirupsen/logrus"
)

func (m *Master) handleOperationTask() {
	go func() {
		for {
			task, ok := <-m.operationTaskChan
			if !ok {
				break
			}

			operation.OprInfoUtil.SetOptInfo(task.OperationID,
				operation.OprInfo{
					CallbackURL: task.CallbackURL,
				})

			switch body := task.OperationTaskBody.(type) {
			case operation.OperationContainerCreateTask:
				go m.handleOperationContainerCreateTask(body)
			case operation.OperationContainerStartTask:
				go m.handleOperationContainerStartTask(body)
			case operation.OperationContainerStopTask:
				go m.handleOperationContainerStopTask(body)
			case operation.OperationContainerRemoveTask:
				go m.handleOperationContainerRemoveTask(body)
			default:
				logrus.Warning(operation.ErrOperationTask.Error())
			}
		}
	}()
}

func (m *Master) handleOperationResponse() {
	go func() {
		for {
			resp, ok := <-m.operationResponseChan
			if !ok {
				break
			}

			switch body := resp.OperationResponseBody.(type) {
			case operation.OperationContainerCreateResponse:
				go m.handleOperationContainerCreateResponse(body)
			case operation.OperationContainerStartResponse:
				go m.handleOperationContainerStartResponse(body)
			case operation.OperationContainerStopResponse:
				go m.handleOperationContainerStopResponse(body)
			case operation.OperationContainerRemoveResponse:
				go m.handleOperationContainerRemoveResponse(body)
			default:
				logrus.Warning(operation.ErrOperationResponse.Error())
			}
		}
	}()
}

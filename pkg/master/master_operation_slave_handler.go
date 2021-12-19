package master

import (
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	http_controller "github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/http-controller"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/operation"
	"github.com/sirupsen/logrus"
)

// slave 删除不需要 slave 的 response
func (m *Master) handleOperationSlaveDeleteTask(task operation.OperationSlaveDeleteTask) {
	var errCode = types.SUCCESS

	oprInfo, ok := operation.OprInfoUtil.GetOptInfo(task.OperationID)
	if !ok {
		return
	}
	operation.OprInfoUtil.DelOptInfo(task.OperationID)

	err := m.slaveController.SlaveDelete(task.SlaveID)
	if err != nil {
		errCode = types.SLAVE_INVALID
	}

	resp := types.APISlaveDeleteResponse{
		APICallBackResponseBase: types.APICallBackResponseBase{
			OperationID: task.OperationID,
			Code:        errCode,
			Msg:         types.GetMsg(errCode),
		},
		SlaveUUID: task.SlaveID,
	}

	err = http_controller.SendCallbackPostWithoutResponse(oprInfo.CallbackURL, &resp)
	if err != nil {
		logrus.Warning(err.Error())
		return
	}
}

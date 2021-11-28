package http_controller

import (
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/http_wrapper"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (hac *HttpApiController) slaveList(c *gin.Context) {
	var appG = http_wrapper.Gin{C: c}

	response := types.APISlaveListResponse{Slaves: []types.SlaveProfile{}}
	code := types.SUCCESS

	defer func() {
		if code != types.SUCCESS {
			appG.Response(http.StatusOK, code, nil)
		}
	}()

	uuidList, err := hac.redisDAO.SlaveUUIDList()
	if err != nil {
		code = types.ERROR
		return
	}

	for _, uuid := range uuidList {
		profile, err := hac.redisDAO.SlaveProfile(uuid)
		if err != nil {
			continue
		}

		memTotalSize, err := strconv.ParseUint(profile["mem_total"], 10, 64)
		if err != nil {
			continue
		}
		logicalCoreCnt, err := strconv.Atoi(profile["logical_cpu_num"])
		if err != nil {
			continue
		}
		physicalCoreCnt, err := strconv.Atoi(profile["physical_cpu_num"])
		if err != nil {
			continue
		}

		response.Slaves = append(response.Slaves, types.SlaveProfile{
			SlaveUUId:       uuid,
			Platform:        profile["platform"],
			PlatformFamily:  profile["platform_family"],
			PlatformVersion: profile["platform_version"],
			MemoryTotalSize: memTotalSize,
			CpuModelName:    profile["cpu_model_name"],
			LogicalCoreCnt:  logicalCoreCnt,
			PhysicalCoreCnt: physicalCoreCnt,
		})
	}

	appG.Response(http.StatusOK, code, response)
}

func (hac *HttpApiController) slaveStatus(c *gin.Context) {
	var (
		appG = http_wrapper.Gin{C: c}
		form types.APISlaveStatusRequest
	)

	httpCode, errCode := http_wrapper.BindAndValid(c, &form)
	if errCode != types.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	defer func() {
		if errCode != types.SUCCESS {
			appG.Response(http.StatusOK, errCode, nil)
		}
	}()

	response := types.APISlaveStatusResponse{Status: []types.APISlaveStatusItem{}}

	for _, uuid := range form.SlaveUUID {
		item := types.APISlaveStatusItem{
			UUID:  uuid,
			Stats: "offline",
		}

		func() {
			stats, err := hac.redisDAO.SlaveStats(uuid)
			if err != nil || stats == "" || stats == "offline" {
				return
			}
			item.Stats = stats

			status, err := hac.redisDAO.SlaveStatus(uuid)
			if err != nil {
				stats = "offline"
				return
			}
			coreAvailable, err := strconv.Atoi(status["core_available"])
			if err != nil {
				return
			}
			item.CoreAvailable = coreAvailable
			memAvailable, err := strconv.ParseUint(status["mem_available"], 10, 64)
			if err != nil {
				return
			}
			item.MemAvailable = memAvailable
			storageAvailable, err := strconv.ParseUint(status["storage_available"], 10, 64)
			if err != nil {
				return
			}
			item.StorageAvailable = storageAvailable
		}()

		response.Status = append(response.Status, item)
	}

	appG.Response(http.StatusOK, errCode, response)
}

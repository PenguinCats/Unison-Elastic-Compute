package http_controller

import (
	"github.com/PenguinCats/Unison-Elastic-Compute/internal/redis_util"
	"github.com/PenguinCats/Unison-Elastic-Compute/pkg/master/internal/operation"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HttpApiController struct {
	r                 *gin.Engine
	apiPort           string
	operationTaskChan chan *operation.OperationTask

	redisDAO *redis_util.RedisDAO
}

func NewHttpApiController(apiPort string, operationTaskChan chan *operation.OperationTask,
	redisDAO *redis_util.RedisDAO) *HttpApiController {
	hac := &HttpApiController{
		r:                 gin.Default(),
		apiPort:           apiPort,
		operationTaskChan: operationTaskChan,
		redisDAO:          redisDAO,
	}

	hac.initRouter()

	return hac
}

func (hac *HttpApiController) StartHttpApiServe() error {
	err := hac.r.Run(":" + hac.apiPort)
	if err != nil {
		logrus.Error("HTTP api server are down")
		return err
	}

	return nil
}

func (hac *HttpApiController) initRouter() {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	api := r.Group("/api")
	{
		apiContainer := api.Group("/container")
		{
			// change state
			apiContainer.POST("/create", hac.containerCreate)
			apiContainer.POST("/start", hac.containerStart)
			apiContainer.POST("/stop", hac.containerStop)
			apiContainer.POST("/remove", hac.containerRemove)
			// read-only
		}

		apiSlave := api.Group("/slave")
		{
			// read-only
			apiSlave.POST("/list", hac.slaveList)
			apiSlave.POST("/status", hac.slaveStatus)
		}

		////获取标签列表
		//api.GET("/tags", v1.GetTags)
		////新建标签
		//api.POST("/tags", v1.AddTag)
		////更新指定标签
		//api.PUT("/tags/:id", v1.EditTag)
		////删除指定标签
		//api.DELETE("/tags/:id", v1.DeleteTag)
	}

	hac.r = r
}

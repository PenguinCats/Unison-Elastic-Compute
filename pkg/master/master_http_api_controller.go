package master

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HttpApiController struct {
	master *Master
	r      *gin.Engine
}

func newHttpApiController(master *Master) *HttpApiController {
	hac := &HttpApiController{
		master: master,
		r:      gin.Default(),
	}
	return hac
}

func (hac *HttpApiController) startHttpApiServe() {
	err := hac.r.Run(":" + hac.master.apiPort)
	if err != nil {
		logrus.Error("HTTP api server are down")
		hac.master.Stop()
	}
}

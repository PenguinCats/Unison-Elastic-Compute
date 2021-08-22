package http_wrapper

import (
	"github.com/PenguinCats/Unison-Elastic-Compute/api/types"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, types.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, types.ERROR
	}
	if !check {
		return http.StatusBadRequest, types.INVALID_PARAMS
	}

	return http.StatusOK, types.SUCCESS
}

package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadStats struct {
	FuncUnitTodo    int `json:"funcUnitTodo"`
	FuncCtxUnitTodo int `json:"funcCtxUnitTodo"`
}

// @Summary upload status query
// @Produce json
// @Success 200
// @Router  /api/v1/monitor/upload [get]
func HandleStatusUpload(c *gin.Context) {
	stat := &UploadStats{
		FuncUnitTodo:    len(funcUnitQueue),
		FuncCtxUnitTodo: len(funcCtxUnitQueue),
	}
	c.JSON(http.StatusOK, stat)
}

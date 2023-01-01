package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/pkg/server/worker"
)

type UploadStats struct {
	FuncUnitTodo    int `json:"funcUnitTodo"`
	FuncCtxUnitTodo int `json:"funcCtxUnitTodo"`
	ClazzUnitTodo   int `json:"clazzUnitTodo"`
}

// @BasePath /
// @Summary upload status query
// @Produce json
// @Success 200
// @Router  /ops/monitor/upload [get]
// @Tags OPS
func HandleStatusUpload(c *gin.Context) {
	stat := &UploadStats{
		FuncUnitTodo:    worker.GetFuncQueueTodoCount(),
		FuncCtxUnitTodo: worker.GetFuncCtxQueueTodoCount(),
		ClazzUnitTodo:   worker.GetClazzQueueTodoCount(),
	}
	c.JSON(http.StatusOK, stat)
}

// @BasePath /
// @Summary  ping example
// @Produce  json
// @Success  200
// @Router   /ops/ping [get]
// @Tags OPS
func HandlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

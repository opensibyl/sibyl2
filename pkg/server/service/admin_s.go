package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamfzc/sibyl2/pkg/server/worker"
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
		FuncUnitTodo:    worker.GetFuncQueueTodoCount(),
		FuncCtxUnitTodo: worker.GetFuncCtxQueueTodoCount(),
	}
	c.JSON(http.StatusOK, stat)
}

// @BasePath /
// @Summary  ping example
// @Produce  json
// @Success  200
// @Router   /ping [get]
func HandlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

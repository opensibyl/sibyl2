package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/server/object"
	"github.com/williamfzc/sibyl2/pkg/server/worker"
)

func HandleRepoFuncUpload(c *gin.Context) {
	result := &object.FunctionUploadUnit{}
	err := c.BindJSON(result)
	if err != nil {
		core.Log.Errorf("error when parse: %v\n", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("parse json error: %v", err))
		return
	}
	if err := result.WorkspaceConfig.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	worker.SubmitFunc(result)
	c.JSON(http.StatusOK, "received")
}

func HandleFunctionContextUpload(c *gin.Context) {
	result := &object.FuncContextUploadUnit{}
	err := c.BindJSON(result)
	if err != nil {
		core.Log.Errorf("error when parse: %v\n", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("parse json error: %v", err))
		return
	}
	if err := result.WorkspaceConfig.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	worker.SubmitFuncCtx(result)
	c.JSON(http.StatusOK, "received")
}

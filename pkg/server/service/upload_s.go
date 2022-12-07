package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/server/object"
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
	sharedQueue.SubmitFunc(result)
	c.JSON(http.StatusOK, "received")
}

func HandleFunctionContextUpload(c *gin.Context) {
	result := &object.FunctionContextUploadUnit{}
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
	sharedQueue.SubmitFuncCtx(result)
	c.JSON(http.StatusOK, "received")
}

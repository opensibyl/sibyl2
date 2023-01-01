package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

func HandleFunctionUpload(c *gin.Context) {
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
	go sharedQueue.SubmitFunc(result)
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
	go sharedQueue.SubmitFuncCtx(result)
	c.JSON(http.StatusOK, "received")
}

func HandleClazzUpload(c *gin.Context) {
	result := &object.ClazzUploadUnit{}
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
	go sharedQueue.SubmitClazz(result)
	c.JSON(http.StatusOK, "received")
}

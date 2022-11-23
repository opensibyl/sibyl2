package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamfzc/sibyl2/pkg/core"
)

func HandleRepoFuncUpload(c *gin.Context) {
	result := &FunctionUploadUnit{}
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

	go func() {
		ctx := context.Background()

		err = sharedDriver.CreateWorkspace(result.WorkspaceConfig, ctx)
		if err != nil {
			core.Log.Warnf("error when init: %v\n", err)
		}

		err := sharedDriver.CreateFuncFile(result.WorkspaceConfig, result.FunctionResult, ctx)
		if err != nil {
			core.Log.Errorf("error when upload: %v\n", err)
		}
	}()

	c.JSON(http.StatusOK, "received")
}

func HandleFunctionContextUpload(c *gin.Context) {
	result := &FuncContextUploadUnit{}
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

	go func() {
		ctx := context.Background()
		err = sharedDriver.CreateWorkspace(result.WorkspaceConfig, ctx)
		if err != nil {
			// allow existed
			core.Log.Warnf("error when init: %v\n", err)
		}

		for _, each := range result.FunctionContexts {
			err := sharedDriver.CreateFuncContext(result.WorkspaceConfig, each, ctx)
			if err != nil {
				core.Log.Warnf("error when upload: %v\n", err)
			}
		}
	}()

	c.JSON(http.StatusOK, "received")
}

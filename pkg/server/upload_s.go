package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williamfzc/sibyl2/pkg/core"
)

var funcUnitQueue chan *FunctionUploadUnit
var funcCtxUnitQueue chan *FuncContextUploadUnit

// default neo4j db may be very slow in I/O
var workerCount = 64

// tiny mq, will block request when it is full
// so make it large enough in your scenario
// or add a real mq
// for example,
// each build for repo which contains 3000 files = 3000 jobs in seconds
var workerQueueSize = 10240

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
	funcUnitQueue <- result
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
	funcCtxUnitQueue <- result
	c.JSON(http.StatusOK, "received")
}

func initUpload(config ExecuteConfig) {
	workerCount = config.UploadWorkerCount
	workerQueueSize = config.UploadQueueSize

	funcUnitQueue = make(chan *FunctionUploadUnit, workerQueueSize)
	funcCtxUnitQueue = make(chan *FuncContextUploadUnit, workerQueueSize)
	initWorkers(LifecycleContext)
}

func initWorkers(ctx context.Context) {
	for i := 0; i < workerCount; i++ {
		go func() {
			startWorker(ctx)
		}()
	}
}

func startWorker(ctx context.Context) {
	for {
		select {
		case result := <-funcUnitQueue:
			err := sharedDriver.CreateWorkspace(result.WorkspaceConfig, ctx)
			if err != nil {
				// allow existed
				core.Log.Warnf("error when init: %v\n", err)
			}

			err = sharedDriver.CreateFuncFile(result.WorkspaceConfig, result.FunctionResult, ctx)
			if err != nil {
				core.Log.Errorf("error when upload: %v\n", err)
			}

		case result := <-funcCtxUnitQueue:
			err := sharedDriver.CreateWorkspace(result.WorkspaceConfig, ctx)
			if err != nil {
				// allow existed
				core.Log.Warnf("error when init: %v\n", err)
			}

			for _, each := range result.FunctionContexts {
				err = sharedDriver.CreateFuncContext(result.WorkspaceConfig, each, ctx)
				if err != nil {
					// deadlock easily happen in neo4j when creating complex edges
					// append to the queue
					// should replace with dead message queue
					core.Log.Warnf("err when create ctx for: %v", each.GetSignature())
					funcCtxUnitQueue <- result
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

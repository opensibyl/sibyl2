package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server/binding"
	_ "github.com/opensibyl/sibyl2/pkg/server/docs"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/opensibyl/sibyl2/pkg/server/queue"
	"github.com/opensibyl/sibyl2/pkg/server/service"
	"github.com/opensibyl/sibyl2/pkg/server/worker"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/*
server mode: everything
worker mode:
- no http api
- read only queue
- full worker threads

proxy mode:
- full http api
- write only queue
- no worker threads
*/

// @title openapi for sibyl2 server
// @version         1.0
// @termsOfService  http://swagger.io/terms/
// @contact.name   williamfzc
// @contact.url    https://github.com/williamfzc
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func Execute(config object.ExecuteConfig) {
	configStr, err := config.ToJson()
	if err != nil {
		core.Log.Errorf("parse config to string failed: %v", err)
		return
	}

	core.Log.Infof("started with config: %s", configStr)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sharedDriver, err := binding.InitDriver(config, ctx)
	if err != nil {
		core.Log.Errorf("failed to create binding: %v", err)
		return
	}

	defer sharedDriver.DeferDriver()
	mq := queue.InitQueue(config, ctx)
	defer mq.Defer()

	service.InitService(config, ctx, sharedDriver, mq)
	worker.InitWorker(config, ctx, sharedDriver, mq)

	engine := gin.Default()

	v1group := engine.Group("/api/v1")
	// query
	v1group.Handle(http.MethodGet, "/repo", service.HandleRepoQuery)
	v1group.Handle(http.MethodGet, "/rev", service.HandleRevQuery)
	v1group.Handle(http.MethodGet, "/file", service.HandleFileQuery)
	v1group.Handle(http.MethodGet, "/func", service.HandleFunctionsQuery)
	v1group.Handle(http.MethodGet, "/funcctx", service.HandleFunctionCtxQuery)
	// upload
	v1group.Handle(http.MethodPost, "/func", service.HandleRepoFuncUpload)
	v1group.Handle(http.MethodPost, "/funcctx", service.HandleFunctionContextUpload)
	// admin
	v1group.Handle(http.MethodGet, "/monitor/upload", service.HandleStatusUpload)
	engine.Handle(http.MethodGet, "/ping", service.HandlePing)
	engine.Handle(http.MethodGet, "/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = engine.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		core.Log.Errorf("failed to start repoctor_receiver: %s", err.Error())
	}
}

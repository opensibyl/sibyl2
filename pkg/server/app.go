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

	/*
		server mode: everything
		worker mode:
		- no http api
		- read only queue
		- full worker threads

		gateway mode:
		- full http api
		- write only queue
		- no worker threads
	*/
	needWorker := config.Mode == object.ServerTypeAll || config.Mode == object.ServerTypeWorker
	needGateway := config.Mode == object.ServerTypeAll || config.Mode == object.ServerTypeGateway
	core.Log.Infof("current mode: %s, worker: %v, gateway: %v", config.Mode, needWorker, needGateway)

	// middleware start up
	// data driver is always required for query
	sharedDriver, err := binding.InitDriver(config, ctx)
	if err != nil {
		core.Log.Errorf("failed to create binding: %v", err)
		return
	}
	defer sharedDriver.DeferDriver()
	// queue is always required for submit
	mq := queue.InitQueue(config, ctx)
	defer mq.Defer()
	service.InitService(config, ctx, sharedDriver, mq)

	// worker start up
	if needWorker {
		worker.InitWorker(config, ctx, sharedDriver, mq)
	}

	// webserver start up
	engine := gin.Default()
	v1group := engine.Group("/api/v1")

	// for CRUD
	if needGateway {
		// query
		v1group.Handle(http.MethodGet, "/repo", service.HandleRepoQuery)
		v1group.Handle(http.MethodGet, "/rev", service.HandleRevQuery)
		v1group.Handle(http.MethodGet, "/file", service.HandleFileQuery)
		v1group.Handle(http.MethodGet, "/func", service.HandleFunctionsQuery)
		v1group.Handle(http.MethodGet, "/funcctx", service.HandleFunctionCtxQuery)
		v1group.Handle(http.MethodGet, "/clazz", service.HandleClazzQuery)
		// upload
		v1group.Handle(http.MethodPost, "/func", service.HandleFunctionUpload)
		v1group.Handle(http.MethodPost, "/funcctx", service.HandleFunctionContextUpload)
		v1group.Handle(http.MethodPost, "/clazz", service.HandleClazzUpload)
		// swagger
		engine.Handle(http.MethodGet, "/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	// for ops
	engine.Handle(http.MethodGet, "/ops/ping", service.HandlePing)
	engine.Handle(http.MethodGet, "/ops/monitor/upload", service.HandleStatusUpload)

	err = engine.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		core.Log.Errorf("failed to start repoctor_receiver: %s", err.Error())
	}
}

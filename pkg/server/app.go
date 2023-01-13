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
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		// scope query
		v1group.Handle(http.MethodGet, "/repo", service.HandleRepoQuery)
		v1group.Handle(http.MethodDelete, "/repo", service.HandleRepoDelete)
		v1group.Handle(http.MethodGet, "/rev", service.HandleRevQuery)
		v1group.Handle(http.MethodDelete, "/rev", service.HandleRevDelete)
		v1group.Handle(http.MethodGet, "/file", service.HandleFileQuery)
		// normal upload
		v1group.Handle(http.MethodPost, "/func", service.HandleFunctionUpload)
		v1group.Handle(http.MethodPost, "/funcctx", service.HandleFunctionContextUpload)
		v1group.Handle(http.MethodPost, "/clazz", service.HandleClazzUpload)
		// normal query
		v1group.Handle(http.MethodGet, "/func", service.HandleFunctionsQuery)
		v1group.Handle(http.MethodGet, "/funcctx", service.HandleFunctionContextsQuery)
		v1group.Handle(http.MethodGet, "/clazz", service.HandleClazzesQuery)

		// EXPERIMENTAL
		// global query (e.g. Method name is known, but do not know where it is)
		v1group.Handle(http.MethodGet, "/func/signature", service.HandleFunctionSignaturesQuery)
		v1group.Handle(http.MethodGet, "/func/with/signature", service.HandleFunctionQueryWithSignature)
		v1group.Handle(http.MethodGet, "/func/with/regex", service.HandleFunctionQueryWithRegex)
		v1group.Handle(http.MethodGet, "/clazz/with/regex", service.HandleClazzQueryWithRegex)
		v1group.Handle(http.MethodGet, "/funcctx/with/regex", service.HandleFuncCtxQueryWithRegex)
		v1group.Handle(http.MethodGet, "/funcctx/with/reference/count", service.HandleFunctionContextQueryWithReferenceCount)
		v1group.Handle(http.MethodGet, "/funcctx/with/referenced/count", service.HandleFunctionContextQueryWithReferencedCount)
		v1group.Handle(http.MethodGet, "/rev/stat", service.HandleRevStatQuery)
	}
	// for ops
	opsGroup := engine.Group("/ops")
	opsGroup.Handle(http.MethodGet, "/ping", service.HandlePing)
	opsGroup.Handle(http.MethodGet, "/monitor/upload", service.HandleStatusUpload)
	opsGroup.Handle(http.MethodGet, "/version", service.HandleVersion)
	opsGroup.Handle(http.MethodGet, "/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	opsGroup.Handle(http.MethodGet, "/metrics", gin.WrapH(promhttp.Handler()))

	err = engine.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		core.Log.Errorf("failed to start repoctor_receiver: %s", err.Error())
	}
}

package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/frontend"
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
func Execute(config object.ExecuteConfig, ctx context.Context) error {
	configStr, err := config.ToJson()
	if err != nil {
		core.Log.Errorf("parse config to string failed: %v", err)
		return err
	}

	defer core.Log.Infof("sibyl everything down")
	core.Log.Infof("started with config: %s", configStr)

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
		return err
	}
	defer func() {
		sharedDriver.DeferDriver()
		core.Log.Infof("shared driver down")
	}()

	// queue is always required for submit
	mq := queue.InitQueue(config, ctx)
	defer func() {
		mq.Defer()
		core.Log.Infof("mq down")
	}()

	service.InitService(config, ctx, sharedDriver, mq)

	// worker start up
	if needWorker {
		worker.InitWorker(config, ctx, sharedDriver, mq)
	}

	// webserver start up
	engine := gin.Default()

	// for CRUD
	if needGateway {
		// NOTICE:
		// I did not maintain a controller layer like spring for clear arch.
		// I will add it when we really need v2 api.
		v1group := engine.Group("/api/v1")
		injectV1Group(v1group)
	}
	// for ops
	opsGroup := engine.Group("/ops")
	injectOpsGroup(opsGroup)

	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: engine,
	}

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			core.Log.Errorf("sibyl server down: %s", err.Error())
		}
	}()
	<-ctx.Done()
	err = srv.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func ExecuteFrontend(port int, ctx context.Context) error {
	engine := gin.Default()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: engine,
	}

	engine.StaticFS("/", http.FS(frontend.Static))
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			core.Log.Errorf("sibyl server down: %s", err.Error())
		}
	}()
	<-ctx.Done()
	err := srv.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func injectV1Group(v1group *gin.RouterGroup) {
	// scope
	scopeGroup := v1group.Group("/")
	scopeGroup.Handle(http.MethodGet, "/repo", service.HandleRepoQuery)
	scopeGroup.Handle(http.MethodDelete, "/repo", service.HandleRepoDelete)
	scopeGroup.Handle(http.MethodGet, "/rev", service.HandleRevQuery)
	scopeGroup.Handle(http.MethodDelete, "/rev", service.HandleRevDelete)
	scopeGroup.Handle(http.MethodGet, "/file", service.HandleFileQuery)
	// upload
	uploadGroup := v1group.Group("/")
	uploadGroup.Handle(http.MethodPost, "/func", service.HandleFunctionUpload)
	uploadGroup.Handle(http.MethodPost, "/funcctx", service.HandleFunctionContextUpload)
	uploadGroup.Handle(http.MethodPost, "/clazz", service.HandleClazzUpload)
	// basic
	basicGroup := v1group.Group("/")
	basicGroup.Handle(http.MethodGet, "/func", service.HandleFunctionsQuery)
	basicGroup.Handle(http.MethodGet, "/funcctx", service.HandleFunctionContextsQuery)
	basicGroup.Handle(http.MethodGet, "/clazz", service.HandleClazzesQuery)

	// query by signature
	signatureGroup := v1group.Group("signature")
	signatureGroup.Handle(http.MethodGet, "/regex/func", service.HandleSignatureRegexFunc)
	signatureGroup.Handle(http.MethodGet, "/func", service.HandleSignatureFunc)
	signatureGroup.Handle(http.MethodGet, "/funcctx", service.HandleSignatureFuncctx)
	signatureGroup.Handle(http.MethodGet, "/funcctx/chain", service.HandleSignatureFuncctxChain)
	signatureGroup.Handle(http.MethodGet, "/funcctx/rchain", service.HandleSignatureFuncctxReverseChain)
	// query by regex
	regexGroup := v1group.Group("regex")
	regexGroup.Handle(http.MethodGet, "/func", service.HandleRegexFunc)
	regexGroup.Handle(http.MethodGet, "/clazz", service.HandleRegexClazz)
	regexGroup.Handle(http.MethodGet, "/funcctx", service.HandleRegexFuncctx)
	// query by reference
	referenceGroup := v1group.Group("reference")
	countGroup := referenceGroup.Group("count")
	countGroup.Handle(http.MethodGet, "/funcctx", service.HandleReferenceCountFuncctx)
	countGroup.Handle(http.MethodGet, "/funcctx/reverse", service.HandleReferenceCountFuncctxReverse)
	// query by stat
	v1group.Handle(http.MethodGet, "/rev/stat", service.HandleRevStatQuery)
}

func injectOpsGroup(opsGroup *gin.RouterGroup) {
	opsGroup.Handle(http.MethodGet, "/ping", service.HandlePing)
	opsGroup.Handle(http.MethodGet, "/monitor/upload", service.HandleStatusUpload)
	opsGroup.Handle(http.MethodGet, "/version", service.HandleVersion)
	opsGroup.Handle(http.MethodGet, "/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	opsGroup.Handle(http.MethodGet, "/metrics", gin.WrapH(promhttp.Handler()))
}

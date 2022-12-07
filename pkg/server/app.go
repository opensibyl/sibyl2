package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/williamfzc/sibyl2/pkg/server/binding"
	_ "github.com/williamfzc/sibyl2/pkg/server/docs"
	"github.com/williamfzc/sibyl2/pkg/server/object"
	"github.com/williamfzc/sibyl2/pkg/server/service"
	"github.com/williamfzc/sibyl2/pkg/server/worker"
)

// @title swagger doc for sibyl2 server
func Execute(config object.ExecuteConfig) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sharedDriver := initDriver(config)
	service.InitServices(config, ctx, sharedDriver)
	worker.InitWorker(config, ctx, sharedDriver)

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

	err := engine.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		fmt.Printf("failed to start repoctor_receiver: %s", err.Error())
	}
}

func initDriver(config object.ExecuteConfig) binding.Driver {
	var driver binding.Driver
	switch config.DbType {
	case binding.DtInMemory:
		driver = initMemDriver()
	case binding.DtNeo4j:
		driver = initNeo4jDriver(config)
	default:
		driver = initMemDriver()
	}
	err := driver.InitDriver()
	if err != nil {
		panic(err)
	}
	return driver
}

func initMemDriver() binding.Driver {
	driver, err := binding.NewInMemoryDriver()
	if err != nil {
		panic(err)
	}
	return driver
}

func initNeo4jDriver(config object.ExecuteConfig) binding.Driver {
	var authToken = neo4j.BasicAuth(config.Neo4jUserName, config.Neo4jPassword, "")
	driver, err := neo4j.NewDriverWithContext(config.Neo4jUri, authToken)
	if err != nil {
		panic(err)
	}
	final, err := binding.NewNeo4jDriver(driver)
	if err != nil {
		panic(err)
	}
	return final
}

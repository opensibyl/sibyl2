package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"github.com/williamfzc/sibyl2/pkg/server/binding"
)

var sharedDriver binding.Driver

type FunctionWithSignature struct {
	*sibyl2.FunctionWithPath
	Signature string `json:"signature"`
}

type FunctionUploadUnit struct {
	WorkspaceConfig *binding.WorkspaceConfig      `json:"workspace"`
	FunctionResult  *extractor.FunctionFileResult `json:"funcResult"`
}

type FuncContextUploadUnit struct {
	WorkspaceConfig  *binding.WorkspaceConfig  `json:"workspace"`
	FunctionContexts []*sibyl2.FunctionContext `json:"functionContext"`
}

func HandlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

type ExecuteConfig struct {
	DbType        binding.DriverType
	Neo4jUri      string
	Neo4jUserName string
	Neo4jPassword string
}

func DefaultExecuteConfig() ExecuteConfig {
	return ExecuteConfig{
		binding.DtInMemory,
		"bolt://localhost:7687",
		"neo4j",
		"neo4j",
	}
}

func Execute(config ExecuteConfig) {
	initDriver(config)

	engine := gin.Default()
	engine.Handle(http.MethodGet, "/ping", HandlePing)

	v1group := engine.Group("/api/v1")
	v1group.Handle(http.MethodGet, "/repo", HandleRepoQuery)
	v1group.Handle(http.MethodGet, "/rev", HandleRevQuery)
	v1group.Handle(http.MethodGet, "/file", HandleFileQuery)
	v1group.Handle(http.MethodGet, "/func", HandleFunctionsQuery)
	v1group.Handle(http.MethodGet, "/funcctx", HandleFunctionCtxQuery)

	v1group.Handle(http.MethodPost, "/func", HandleRepoFuncUpload)
	v1group.Handle(http.MethodPost, "/funcctx", HandleFunctionContextUpload)

	err := engine.Run(fmt.Sprintf(":%d", 9876))
	if err != nil {
		fmt.Printf("failed to start repoctor_receiver: %s", err.Error())
	}
}

func initDriver(config ExecuteConfig) {
	switch config.DbType {
	case binding.DtInMemory:
		initMemDriver()
	case binding.DtNeo4j:
		initNeo4jDriver(config)
	default:
		initMemDriver()
	}
	err := sharedDriver.InitDriver()
	if err != nil {
		panic(err)
	}
}

func initMemDriver() {
	driver, err := binding.NewInMemoryDriver()
	if err != nil {
		panic(err)
	}
	sharedDriver = driver
}

func initNeo4jDriver(config ExecuteConfig) {
	var authToken = neo4j.BasicAuth(config.Neo4jUserName, config.Neo4jPassword, "")
	driver, err := neo4j.NewDriverWithContext(config.Neo4jUri, authToken)
	if err != nil {
		panic(err)
	}
	sharedDriver, err = binding.NewNeo4jDriver(driver)
	if err != nil {
		panic(err)
	}
}

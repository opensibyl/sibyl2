package binding

import (
	"context"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/pkg/errors"
)

// binding to backend databases
// mainly designed for graph db
// such as neo4j/nebula

/*
About how to insert a func node to graph db

- create func node itself with all the properties
- check and create nodes:
	- file node, create if absent
	- rev node, create if absent
	- repo node, create if absent
- create links
	- file INCLUDE func
	- rev INCLUDE file
	- repo INCLUDE rev

About how to create link between functions

- check:
	- func 1 existed
	- func 2 existed
- link
	- func1 CALL func2

cypher:

create nodes:
	MERGE (func:Func {signature: "abcde:fdeglkb"})
	MERGE (f:File {path: 'abcde'})
	MERGE (rev:Rev {hash: '123456F'})
	MERGE (repo:Repo {id: 1234, name: "haha"})

	MERGE (f)-[:INCLUDE]->(func)
	MERGE (rev)-[:INCLUDE]->(f)
	MERGE (repo)-[:INCLUDE]->(rev)
	RETURN *

create func links:
	MATCH (src:Func {signature:"abcde:fdeglkb"})
	MATCH (tar:Func {signature:"eytjkdgfhs"})
	MERGE (src)-[r:CALL]->(tar)
	RETURN *
*/

type driverBase interface {
	GetType() object.DriverType
	InitDriver(ctx context.Context) error
}

type driverCreate interface {
	CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error
	CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContext, ctx context.Context) error
	CreateWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error
}

type driverRead interface {
	ReadRepos(ctx context.Context) ([]string, error)
	ReadRevs(repoId string, ctx context.Context) ([]string, error)
	ReadFiles(wc *object.WorkspaceConfig, ctx context.Context) ([]string, error)
	ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.FunctionWithPath, error)
	ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionWithPath, error)
	ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionWithPath, error)
	ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionContext, error)
}

type driverUpdate interface {
	UpdateRevProperties(wc *object.WorkspaceConfig, k string, v any, ctx context.Context) error
	UpdateFileProperties(wc *object.WorkspaceConfig, path string, k string, v any, ctx context.Context) error
	UpdateFuncProperties(wc *object.WorkspaceConfig, signature string, k string, v any, ctx context.Context) error
}

type driverDelete interface {
	DeleteWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error
}

type Driver interface {
	driverBase
	driverCreate
	driverRead
	driverUpdate
	driverDelete
}

func NewNeo4jDriver(dwc neo4j.DriverWithContext) (Driver, error) {
	return &neo4jDriver{dwc}, nil
}

func NewInMemoryDriver() (Driver, error) {
	return newMemDriver(), nil
}

func WorkspaceConfigFromKey(key string) (*object.WorkspaceConfig, error) {
	parts := strings.Split(key, object.FlagWcKeySplit)
	if len(parts) < 2 {
		return nil, errors.New("invalid workspace repr: " + key)
	}
	ret := &object.WorkspaceConfig{
		RepoId:  parts[0],
		RevHash: parts[1],
	}
	return ret, nil
}

func InitDriver(config object.ExecuteConfig, ctx context.Context) Driver {
	var driver Driver
	switch config.DbType {
	case object.DtInMemory:
		driver = initMemDriver()
	case object.DtNeo4j:
		driver = initNeo4jDriver(config)
	default:
		driver = initMemDriver()
	}
	err := driver.InitDriver(ctx)
	if err != nil {
		panic(err)
	}
	return driver
}

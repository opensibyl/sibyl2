package binding

import (
	"context"
	"strings"

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
	DeferDriver() error
}

/*
Rule

Rule is a query structure implemented with regex and gjson syntax.
- key: gjson path syntax
- value: regex
*/
type Rule = map[string]string

type driverCreate interface {
	CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error
	CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContext, ctx context.Context) error
	CreateClazzFile(wc *object.WorkspaceConfig, c *extractor.ClazzFileResult, ctx context.Context) error
	CreateWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error
}

type driverRead interface {
	ReadRepos(ctx context.Context) ([]string, error)
	ReadRevs(repoId string, ctx context.Context) ([]string, error)
	ReadFiles(wc *object.WorkspaceConfig, ctx context.Context) ([]string, error)
	ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.FunctionWithPath, error)
	ReadFunctionsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*sibyl2.FunctionWithPath, error)
	ReadFunctionSignaturesWithRegex(wc *object.WorkspaceConfig, regex string, ctx context.Context) ([]string, error)
	ReadClasses(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.ClazzWithPath, error)
	ReadClassesWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.ClazzWithPath, error)
	ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionWithPath, error)
	ReadFunctionContextsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionContext, error)
	ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionWithPath, error)
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

func InitDriver(config object.ExecuteConfig, ctx context.Context) (Driver, error) {
	var driver Driver

	// create driver obj, do some settings
	switch config.DbType {
	case object.DriverTypeInMemory:
		// now in memory driver handled by badger
		driver = initBadgerDriver(config)
	case object.DriverTypeNeo4j:
		driver = initNeo4jDriver(config)
	case object.DriverTypeBadger:
		driver = initBadgerDriver(config)
	case object.DriverTypeTikv:
		driver = initTikvDriver(config)

	default:
		return nil, errors.New("invalid driver: " + string(config.DbType))
	}

	// init driver instance (maybe pre connection, etc.)
	err := driver.InitDriver(ctx)
	if err != nil {
		return nil, err
	}
	return driver, nil
}

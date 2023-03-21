package binding

import (
	"context"
	"strings"

	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/pkg/errors"
)

/*
binding to backend databases

Mainly designed for k/v databases.
You can implement Driver interface to adapt different backends.

Such as:
- tikv
- redis (WIP)
- badger
- ...

In the past days, we have tried neo4j but removed because of performance and distribution.
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
- value: verify function

full serialization is expensive.
*/
type Rule = map[string]func(string) bool

type driverCreate interface {
	CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error
	CreateFuncTag(wc *object.WorkspaceConfig, signature string, tag string, ctx context.Context) error
	CreateFuncContext(wc *object.WorkspaceConfig, f *object.FunctionContextSlim, ctx context.Context) error
	CreateClazzFile(wc *object.WorkspaceConfig, c *extractor.ClazzFileResult, ctx context.Context) error
	CreateWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error
}

type driverRead interface {
	ReadRepos(ctx context.Context) ([]string, error)
	ReadRevs(repoId string, ctx context.Context) ([]string, error)
	ReadRevInfo(wc *object.WorkspaceConfig, ctx context.Context) (*object.RevInfo, error)
	ReadFiles(wc *object.WorkspaceConfig, ctx context.Context) ([]string, error)

	ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*object.FunctionServiceDTO, error)
	ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.FunctionServiceDTO, error)
	ReadFunctionsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*object.FunctionServiceDTO, error)
	ReadFunctionSignaturesWithRegex(wc *object.WorkspaceConfig, regex string, ctx context.Context) ([]string, error)
	ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*object.FunctionServiceDTO, error)
	ReadFunctionsWithTag(wc *object.WorkspaceConfig, tag object.FuncTag, ctx context.Context) ([]string, error)

	ReadClasses(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*object.ClazzServiceDTO, error)
	ReadClassesWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.ClazzServiceDTO, error)
	ReadClassesWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*object.ClazzServiceDTO, error)

	ReadFunctionContextsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.FunctionContextSlim, error)
	ReadFunctionContextsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*object.FunctionContextSlim, error)
	ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*object.FunctionContextSlim, error)
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
	case object.DriverTypeBadger:
		driver = initBadgerDriver(config)
	case object.DriverTypeTikv:
		driver = initTikvDriver(config)
	case object.DriverTypeMongoDB:
		driver = initMongoDriver(config)

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

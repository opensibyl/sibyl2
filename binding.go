package sibyl2

import (
	"context"

	"github.com/williamfzc/sibyl2/pkg/extractor"
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

type Driver interface {
	UploadFileResult(wc *WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error
	UploadFuncContext(wc *WorkspaceConfig, f *FunctionContext, ctx context.Context) error
	QueryFiles(wc *WorkspaceConfig, ctx context.Context) ([]string, error)
	QueryFunctions(wc *WorkspaceConfig, path string, ctx context.Context) ([]*FunctionWithPath, error)
	QueryFunctionWithSignature(wc *WorkspaceConfig, signature string, ctx context.Context) (*FunctionWithPath, error)
	QueryFunctionsWithLines(wc *WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*FunctionWithPath, error)
	QueryFunctionContextWithSignature(wc *WorkspaceConfig, signature string, ctx context.Context) (*FunctionContext, error)
}

/*
WorkspaceConfig

as an infra lib, it will not assume what kind of repo you used.

just two fields:
- repoId: unique id of your repo, no matter git or svn, even appId.
- revHash: unique id of your version.
*/
type WorkspaceConfig struct {
	RepoId  string `json:"repoId"`
	RevHash string `json:"revHash"`
}

func (wc *WorkspaceConfig) Verify() bool {
	// all the fields should be filled
	return wc.RepoId != "" && wc.RevHash != ""
}

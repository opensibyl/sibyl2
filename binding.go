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
	UploadFileResultWithContext(wc *WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error
}

type RepoConfig struct {
	RepoId   int    `json:"repoId"`
	RepoName string `json:"repoName"`
	RepoType string `json:"repoType"`
}

type WorkspaceConfig struct {
	*RepoConfig
	RevHash string `json:"revHash"`
}

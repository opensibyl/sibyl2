package sibyl2

import (
	"context"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestUploadFile(t *testing.T) {
	t.Skip("always skip in CI")
	wc := &WorkspaceConfig{
		RepoConfig: &RepoConfig{
			RepoId:   102994,
			RepoName: "sibyl2",
			RepoType: "github",
		},
		RevHash: "79068a0e21f095c3b7f35aff28f44db74173f3fe",
	}

	dbUri := "bolt://localhost:7687"
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("neo4j", "williamfzc", ""))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer driver.Close(ctx)
	newDriver := &Neo4jDriver{driver}
	functions, _ := ExtractFunction(".", DefaultConfig())
	for _, each := range functions {
		newDriver.UploadFileResultWithContext(wc, each, ctx)
	}
}

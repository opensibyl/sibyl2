package sibyl2

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/williamfzc/sibyl2/pkg/extractor"
)

const (
	TemplateMergeFunc = "MERGE (func%d:Func {" +
		"name: $func_name, " +
		"receiver: $func_receiver, " +
		"parameters: $func_parameters, " +
		"returns: $func_returns, " +
		"span: $func_span, " +
		"extras: $func_extras})"
	TemplateMergeFile = "MERGE (file:File {" +
		"path: $file_path, " +
		"lang: $file_lang})"
	TemplateMergeRev = "MERGE (rev:Rev {" +
		"hash: $rev_hash})"
	TemplateMergeRepo = "MERGE (repo:Repo {" +
		"name: $repo_name, " +
		"type: $repo_type, " +
		"id: $repo_id})"
	TemplateMergeLinkInclude   = "MERGE (%s)-[:INCLUDE]->(%s)"
	TemplateMergeLinkReference = "MERGE (%s)-[:REFERENCE]->(%s)"
)

type Neo4jDriver struct {
	neo4j.DriverWithContext
}

func (d *Neo4jDriver) UploadFileResultWithContext(wc *WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) {
	driver := d.DriverWithContext
	err := insert(wc, f, ctx, driver)
	if err != nil {
		panic(err)
	}
}

func insert(wc *WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context, driver neo4j.DriverWithContext) error {
	// session is cheap to create
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, createItemFn(wc, f, ctx))
	if err != nil {
		return err
	}
	return nil
}

func createItemFn(wc *WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		merged := []string{
			TemplateMergeRepo,
			TemplateMergeRev,
			TemplateMergeFile,
			fmt.Sprintf(TemplateMergeLinkInclude, "repo", "rev"),
			fmt.Sprintf(TemplateMergeLinkInclude, "rev", "file"),
		}

		for i, each := range f.Units {
			newMerged := append(merged, []string{
				// create func node
				fmt.Sprintf(TemplateMergeFunc, i),
				// link them
				fmt.Sprintf(TemplateMergeLinkInclude, "file", fmt.Sprintf("func%d", i)),
			}...)

			recv, _ := json.Marshal(each.Receiver)
			params, _ := json.Marshal(each.Parameters)
			returns, _ := json.Marshal(each.Returns)
			extras, _ := json.Marshal(each.Extras)
			_, err := tx.Run(ctx, strings.Join(newMerged, " "), map[string]any{
				"repo_id":         wc.RepoId,
				"repo_name":       wc.RepoName,
				"repo_type":       wc.RepoType,
				"rev_hash":        wc.RevHash,
				"file_path":       f.Path,
				"file_lang":       f.Language,
				"func_name":       each.Name,
				"func_receiver":   string(recv),
				"func_parameters": string(params),
				"func_returns":    string(returns),
				"func_span":       each.Span.String(),
				"func_extras":     string(extras),
			})
			_, err = tx.Run(ctx, fmt.Sprintf(TemplateMergeLinkInclude, "file", "func"), nil)
			if err != nil {
				return nil, err
			}
		}

		// temp ignore return values
		return nil, nil
	}
}

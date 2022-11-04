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
	TemplateMergeFunc = "CREATE (func%d:Func {" +
		"name: $func_name, " +
		"receiver: $func_receiver, " +
		"parameters: $func_parameters, " +
		"returns: $func_returns, " +
		"span: $func_span, " +
		"extras: $func_extras," +
		"signature: $func_signature })"
	TemplateMergeFile = "MERGE (file:File {" +
		"path: $file_path, " +
		"lang: $file_lang})"
	TemplateMergeRev = "MERGE (rev:Rev {" +
		"hash: $rev_hash})"
	TemplateMergeRepo = "MERGE (repo:Repo {id: $repo_id})"

	TemplateMatchFuncFull = "MATCH " +
		"(repo:Repo {name: $repo_name, type: $repo_type, id: $repo_id})" +
		"-[:INCLUDE]->" +
		"(rev:Rev {hash: $rev_hash})" +
		"-[:INCLUDE]->" +
		"(%s:File {path: %s, lang: $file_lang})" +
		"-[:INCLUDE]->" +
		"(func%d:Func {signature: $func%d_signature})"

	TemplateMergeLinkInclude       = "MERGE (%s)-[:INCLUDE]->(%s)"
	TemplateMergeLinkFileReference = "MERGE (%s)-[:FILE_REFERENCE]->(%s)"
	TemplateMergeLinkFuncReference = "MERGE (%s)-[:FUNC_REFERENCE]->(%s)"
)

type Neo4jDriver struct {
	neo4j.DriverWithContext
}

func (d *Neo4jDriver) UploadFileResultWithContext(wc *WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
	// session is cheap to create
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, createFunctionFileTransaction(wc, f, ctx))
	if err != nil {
		return err
	}
	return nil
}

func (d *Neo4jDriver) UploadFuncContextWithContext(wc *WorkspaceConfig, f *FunctionContext, ctx context.Context) error {
	// session is cheap to create
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, createFuncGraphTransaction(wc, f, ctx))
	if err != nil {
		return err
	}
	return nil
}

func createFunctionFileTransaction(wc *WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) neo4j.ManagedTransactionWork {
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

			receiver, _ := json.Marshal(each.Receiver)
			params, _ := json.Marshal(each.Parameters)
			returns, _ := json.Marshal(each.Returns)
			extras, _ := json.Marshal(each.Extras)
			spanLines := each.Span.Lines()

			// todo: merge these run together
			_, err := tx.Run(ctx, strings.Join(newMerged, " "), map[string]any{
				"repo_id":         wc.RepoId,
				"rev_hash":        wc.RevHash,
				"file_path":       f.Path,
				"file_lang":       f.Language,
				"func_name":       each.Name,
				"func_receiver":   string(receiver),
				"func_parameters": string(params),
				"func_returns":    string(returns),
				"func_span":       []int{spanLines[0], spanLines[len(spanLines)-1]},
				"func_extras":     string(extras),
				"func_signature":  each.GetSignature(),
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

func createFuncGraphTransaction(wc *WorkspaceConfig, f *FunctionContext, ctx context.Context) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		for i, each := range f.Calls {
			id := i + 1
			eachFuncName := fmt.Sprintf("func%d", id)
			eachMerged := []string{
				fmt.Sprintf(TemplateMatchFuncFull, "srcfile", "$srcPath", 0, 0),
				fmt.Sprintf(TemplateMatchFuncFull, "targetfile", "$targetPath", id, id),
			}
			if each.Path != f.Path {
				eachMerged = append(eachMerged, fmt.Sprintf(TemplateMergeLinkFileReference, "srcfile", "targetfile"))
			}
			eachMerged = append(eachMerged,
				// create link
				fmt.Sprintf(TemplateMergeLinkFuncReference, "func0", eachFuncName),
				"RETURN *")

			// todo: merge these run together
			_, err := tx.Run(ctx, strings.Join(eachMerged, " "), map[string]any{
				"repo_id":                           wc.RepoId,
				"rev_hash":                          wc.RevHash,
				"srcPath":                           f.Path,
				"targetPath":                        each.Path,
				"file_lang":                         each.Language,
				"func0_signature":                   f.GetSignature(),
				fmt.Sprintf("func%d_signature", id): each.GetSignature(),
			})
			if err != nil {
				return nil, err
			}
		}

		// temp ignore return values
		return nil, nil
	}
}

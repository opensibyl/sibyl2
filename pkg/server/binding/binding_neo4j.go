package binding

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"github.com/williamfzc/sibyl2/pkg/server/object"
)

// NOTICE:
// I am not pretty sure that neo4j is a proper database
// so my implementation here may be a little casual

// functions can be reused

/*
indexes and constraint (neo4j 5.x
avoid some conflicts
*/
const (
	TemplateMergeFunc = `
MATCH (rev:Rev {id: $repo_id, hash: $rev_hash})
WITH rev
MERGE (rev)-[:INCLUDE]->(file:File {path: $file_path, lang: $file_lang})
MERGE (func:Func {name: $func_name, receiver: $func_receiver, parameters: $func_parameters, returns: $func_returns, span: $func_span, extras: $func_extras, signature: $func_signature })
MERGE (file)-[:INCLUDE]->(func)
`
	TemplateMatchFuncFull = `
MATCH (:Rev {id: $repo_id, hash: $rev_hash})
-[:INCLUDE]->(%s:File {path: %s, lang: $file_lang})
-[:INCLUDE]->(func%d:Func {signature: $func%d_signature})
`

	TemplateMergeLinkInclude       = "MERGE (%s)-[:INCLUDE]->(%s)"
	TemplateMergeLinkFuncReference = "WITH func%d, func%d MERGE (%s)-[:FUNC_REFERENCE]->(%s)"
)

type neo4jDriver struct {
	neo4j.DriverWithContext
}

func (d *neo4jDriver) GetType() object.DriverType {
	return object.DtNeo4j
}

func (d *neo4jDriver) InitDriver(ctx context.Context) error {
	// create indexes
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	// unique index for rev
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		create := `CREATE CONSTRAINT IF NOT EXISTS FOR (r:Rev) REQUIRE (r.hash, r.id) IS UNIQUE`
		_, err := tx.Run(ctx, create, nil)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	// index for lookup and write:
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		create := `CREATE INDEX ON :File(path)`
		_, err := tx.Run(ctx, create, nil)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		create := `CREATE INDEX ON :Func(signature)`
		_, err := tx.Run(ctx, create, nil)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	return nil
}

func (d *neo4jDriver) CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
	if err := wc.Verify(); err != nil {
		return err
	}
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, createFunctionFileTransaction(wc, f, ctx))
	if err != nil {
		return err
	}
	return nil
}

func (d *neo4jDriver) CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContext, ctx context.Context) error {
	if err := wc.Verify(); err != nil {
		return err
	}
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, createFuncGraphTransaction(wc, f, ctx))
	if err != nil {
		return err
	}
	return nil
}

func (d *neo4jDriver) ReadFiles(wc *object.WorkspaceConfig, ctx context.Context) ([]string, error) {
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)
	ret, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (:Rev {id: $repoId, hash: $revHash})-[:INCLUDE]->(f:File) RETURN f.path`
		results, err := tx.Run(ctx, query, map[string]any{
			"repoId":  wc.RepoId,
			"revHash": wc.RevHash,
		})
		if err != nil {
			return nil, err
		}
		nodes, err := results.Collect(ctx)
		if err != nil {
			return nil, err
		}
		ret := make([]string, 0, len(nodes))
		for _, each := range nodes {
			ret = append(ret, each.Values[0].(string))
		}
		return ret, nil
	})
	if err != nil {
		return nil, err
	}
	return ret.([]string), nil
}

func (d *neo4jDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)
	ret, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (:Rev {id: $repoId, hash: $revHash})-[:INCLUDE]->(file:File {path: $path})-[:INCLUDE]->(f) RETURN f, file`
		results, err := tx.Run(ctx, query, map[string]any{
			"repoId":  wc.RepoId,
			"revHash": wc.RevHash,
			"path":    path,
		})
		if err != nil {
			return nil, err
		}
		nodes, err := results.Collect(ctx)
		if err != nil {
			return nil, err
		}
		ret := make([]*sibyl2.FunctionWithPath, 0, len(nodes))
		for _, each := range nodes {
			rawMap := each.Values[0].(dbtype.Node).Props
			file := each.Values[1].(dbtype.Node).Props

			// special handlers for neo4j :)
			rawMap = funcMapAdapter(rawMap)

			f, err := extractor.FromMap(rawMap)
			if err != nil {
				return nil, err
			}
			fwp := &sibyl2.FunctionWithPath{
				Function: f,
				Path:     file["lang"].(string),
				Language: core.LangType(file["path"].(string)),
			}
			ret = append(ret, fwp)
		}
		return ret, nil
	})
	if err != nil {
		return nil, err
	}
	return ret.([]*sibyl2.FunctionWithPath), nil
}

func (d *neo4jDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionWithPath, error) {
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)
	ret, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (:Rev {id: $repoId, hash: $revHash})-[:INCLUDE]->(file:File)-[:INCLUDE]->(f:Func {signature: $signature}) RETURN f, file`
		results, err := tx.Run(ctx, query, map[string]any{
			"repoId":    wc.RepoId,
			"revHash":   wc.RevHash,
			"signature": signature,
		})
		if err != nil {
			return nil, err
		}
		nodes, err := results.Collect(ctx)
		if err != nil {
			return nil, err
		}
		ret := make([]*sibyl2.FunctionWithPath, 0, len(nodes))
		for _, each := range nodes {
			rawMap := each.Values[0].(dbtype.Node).Props
			file := each.Values[1].(dbtype.Node).Props

			// special handlers for neo4j :)
			rawMap = funcMapAdapter(rawMap)

			f, err := extractor.FromMap(rawMap)
			if err != nil {
				return nil, err
			}
			fwp := &sibyl2.FunctionWithPath{
				Function: f,
				Path:     file["lang"].(string),
				Language: core.LangType(file["path"].(string)),
			}
			ret = append(ret, fwp)
		}
		// normally, it will contain only one element
		if len(ret) == 0 {
			return nil, nil
		}
		return ret[0], nil
	})
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	return ret.(*sibyl2.FunctionWithPath), nil
}

func (d *neo4jDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	functions, err := d.ReadFunctions(wc, path, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*sibyl2.FunctionWithPath, 0)
	for _, each := range functions {
		if each.GetSpan().ContainAnyLine(lines...) {
			ret = append(ret, each)
		}
	}
	return ret, nil
}

func (d *neo4jDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionContext, error) {
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)
	ret, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
MATCH (rev:Rev {id: $repoId, hash: $revHash})-[:INCLUDE]->(file:File)-[:INCLUDE]->(f:Func {signature: $signature}) 
MATCH (rev)-[:INCLUDE]->(srcFile:File)-[:INCLUDE]->(srcFunc)-[:FUNC_REFERENCE]->(f)
MATCH (f)-[:FUNC_REFERENCE]->(targetFunc)
MATCH (targetFile:File)-[:INCLUDE]->(targetFunc)
RETURN f, file, srcFunc, srcFile, targetFunc, targetFile
`
		results, err := tx.Run(ctx, query, map[string]any{
			"repoId":    wc.RepoId,
			"revHash":   wc.RevHash,
			"signature": signature,
		})
		if err != nil {
			return nil, err
		}
		nodes, err := results.Collect(ctx)
		if err != nil {
			return nil, err
		}
		fc := &sibyl2.FunctionContext{
			FunctionWithPath: nil,
			Calls:            nil,
			ReverseCalls:     nil,
		}
		srcCache := make(map[string]*sibyl2.FunctionWithPath)
		targetCache := make(map[string]*sibyl2.FunctionWithPath)
		for _, each := range nodes {
			rawMap := each.Values[0].(dbtype.Node).Props
			file := each.Values[1].(dbtype.Node).Props
			srcMap := each.Values[2].(dbtype.Node).Props
			srcFile := each.Values[3].(dbtype.Node).Props
			targetMap := each.Values[4].(dbtype.Node).Props
			targetFile := each.Values[5].(dbtype.Node).Props

			// special handlers for neo4j :)
			rawMap = funcMapAdapter(rawMap)
			srcMap = funcMapAdapter(srcMap)
			targetMap = funcMapAdapter(targetMap)

			f, err := extractor.FromMap(rawMap)
			if err != nil {
				return nil, err
			}
			srcf, err := extractor.FromMap(srcMap)
			if err != nil {
				return nil, err
			}
			srcfwp := &sibyl2.FunctionWithPath{
				Function: srcf,
				Path:     srcFile["lang"].(string),
				Language: core.LangType(srcFile["path"].(string)),
			}
			targetf, err := extractor.FromMap(targetMap)
			if err != nil {
				return nil, err
			}
			targetfwp := &sibyl2.FunctionWithPath{
				Function: targetf,
				Path:     targetFile["lang"].(string),
				Language: core.LangType(targetFile["path"].(string)),
			}

			fwp := &sibyl2.FunctionWithPath{
				Function: f,
				Path:     file["lang"].(string),
				Language: core.LangType(file["path"].(string)),
			}
			fc.FunctionWithPath = fwp
			srcCache[srcfwp.GetSignature()] = srcfwp
			targetCache[targetfwp.GetSignature()] = targetfwp
		}
		core.Log.Infof("srccache: %v", srcCache)

		srcFinal := make([]*sibyl2.FunctionWithPath, 0, len(srcCache))
		targetFinal := make([]*sibyl2.FunctionWithPath, 0, len(targetCache))
		for _, each := range srcCache {
			srcFinal = append(srcFinal, each)
		}
		for _, each := range targetCache {
			targetFinal = append(targetFinal, each)
		}
		fc.ReverseCalls = srcFinal
		fc.Calls = targetFinal

		return fc, nil
	})
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	return ret.(*sibyl2.FunctionContext), nil
}

func (d *neo4jDriver) DeleteWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	// todo: this will remove functions shared by multi revs.
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
MATCH (rev:Rev {id: $repoId, hash: $revHash})-[:INCLUDE]->(file:File)-[:INCLUDE]->(func:Func) DETACH DELETE rev, file, func
`
		_, err := tx.Run(ctx, query, map[string]any{
			"repoId":  wc.RepoId,
			"revHash": wc.RevHash,
		})
		return nil, err
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *neo4jDriver) CreateWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	if err := wc.Verify(); err != nil {
		return err
	}

	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MERGE (:Rev {id: $repoId, hash: $revHash})`
		_, err := tx.Run(ctx, query, map[string]any{
			"repoId":  wc.RepoId,
			"revHash": wc.RevHash,
		})
		return nil, err
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *neo4jDriver) ReadRepos(ctx context.Context) ([]string, error) {
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)
	ret, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (r:Rev) RETURN r.id`
		ret, err := tx.Run(ctx, query, map[string]any{})
		if err != nil {
			return nil, err
		}
		revs, err := ret.Collect(ctx)
		if err != nil {
			return nil, err
		}

		returns := make([]string, 0, len(revs))
		for _, each := range revs {
			returns = append(returns, each.Values[0].(string))
		}
		return returns, nil
	})

	if err != nil {
		return nil, err
	}
	return ret.([]string), err
}

func (d *neo4jDriver) ReadRevs(repoId string, ctx context.Context) ([]string, error) {
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)
	ret, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (r:Rev {id: $repoId}) RETURN r.hash`
		ret, err := tx.Run(ctx, query, map[string]any{
			"repoId": repoId,
		})
		if err != nil {
			return nil, err
		}
		revs, err := ret.Collect(ctx)
		if err != nil {
			return nil, err
		}

		returns := make([]string, 0, len(revs))
		for _, each := range revs {
			returns = append(returns, each.Values[0].(string))
		}
		return returns, err
	})

	if err != nil {
		return nil, err
	}
	return ret.([]string), nil
}

func (d *neo4jDriver) UpdateRevProperties(wc *object.WorkspaceConfig, k string, v any, ctx context.Context) error {
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
MATCH (r:Rev {id: $repoId, hash: $revHash})
SET r.%s = $v
RETURN r`
		query = fmt.Sprintf(query, k)
		_, err := tx.Run(ctx, query, map[string]any{
			"repoId":  wc.RepoId,
			"revHash": wc.RevHash,
			"v":       v,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *neo4jDriver) UpdateFileProperties(wc *object.WorkspaceConfig, path string, k string, v any, ctx context.Context) error {
	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
MATCH (:Rev {id: $repoId, hash: $revHash})-[:INCLUDE]->(file:File {path: $path})
SET file.%s = $v
RETURN file`
		query = fmt.Sprintf(query, k)
		_, err := tx.Run(ctx, query, map[string]any{
			"repoId":  wc.RepoId,
			"revHash": wc.RevHash,
			"path":    path,
			"v":       v,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *neo4jDriver) UpdateFuncProperties(wc *object.WorkspaceConfig, signature string, k string, v any, ctx context.Context) error {
	if err := wc.Verify(); err != nil {
		return err
	}

	session := d.DriverWithContext.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
MATCH (:Rev {id: $repoId, hash: $revHash})-[:INCLUDE]->(file:File)-[:INCLUDE]->(f:Func {signature: $signature}) 
SET f.%s = $v
RETURN f`
		query = fmt.Sprintf(query, k)
		_, err := tx.Run(ctx, query, map[string]any{
			"repoId":    wc.RepoId,
			"revHash":   wc.RevHash,
			"signature": signature,
			"v":         v,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func createFunctionFileTransaction(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		for _, each := range f.Units {
			funcMap, _ := each.ToMap()

			// neo4j can not handle nested properties and null
			// awful support and I really hate this shit
			for k, v := range funcMap {
				kind := reflect.ValueOf(v).Kind()
				// nested map and slice -> string
				if kind == reflect.Map || kind == reflect.Slice {
					pureContent, _ := json.Marshal(v)
					funcMap[k] = string(pureContent)
					continue
				}
				// nil -> empty string
				if v == nil {
					funcMap[k] = ""
				}
			}

			// todo: merge these run together
			_, err := tx.Run(ctx, TemplateMergeFunc, map[string]any{
				"repo_id":         wc.RepoId,
				"rev_hash":        wc.RevHash,
				"file_path":       f.Path,
				"file_lang":       f.Language,
				"func_name":       each.Name,
				"func_receiver":   funcMap["receiver"],
				"func_parameters": funcMap["parameters"],
				"func_returns":    funcMap["returns"],
				"func_span":       funcMap["span"],
				"func_extras":     funcMap["extras"],
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

func createFuncGraphTransaction(wc *object.WorkspaceConfig, f *sibyl2.FunctionContext, ctx context.Context) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		for i, each := range f.Calls {
			id := i + 1
			eachFuncName := fmt.Sprintf("func%d", id)
			eachMerged := []string{
				fmt.Sprintf(TemplateMatchFuncFull, "srcfile", "$srcPath", 0, 0),
				fmt.Sprintf(TemplateMatchFuncFull, "targetfile", "$targetPath", id, id),
			}
			eachMerged = append(eachMerged,
				// create link
				fmt.Sprintf(TemplateMergeLinkFuncReference, 0, id, "func0", eachFuncName),
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

func funcMapAdapter(origin map[string]any) map[string]any {
	var params []any
	var returns []any
	var span map[string]any
	_ = json.Unmarshal([]byte(origin["parameters"].(string)), &params)
	_ = json.Unmarshal([]byte(origin["returns"].(string)), &returns)
	_ = json.Unmarshal([]byte(origin["span"].(string)), &span)
	origin["parameters"] = params
	origin["returns"] = returns
	origin["span"] = span
	return origin
}
